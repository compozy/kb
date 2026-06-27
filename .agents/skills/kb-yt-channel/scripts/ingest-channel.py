#!/usr/bin/env python3
"""Create or update a KB topic from a YouTube channel."""

from __future__ import annotations

import argparse
import datetime as dt
import json
import os
import re
import shutil
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any
from urllib.parse import parse_qs, urlparse, urlunparse


SLUG_RE = re.compile(r"^[a-z0-9]+(-[a-z0-9]+)*$")
VIDEO_ID_RE = re.compile(r"^[A-Za-z0-9_-]{11}$")


@dataclass(frozen=True)
class Video:
    video_id: str
    title: str
    url: str


class CommandError(RuntimeError):
    def __init__(self, args: list[str], returncode: int, stdout: str, stderr: str) -> None:
        self.args_list = args
        self.returncode = returncode
        self.stdout = stdout
        self.stderr = stderr
        super().__init__(f"command failed ({returncode}): {' '.join(args)}")


def eprint(message: str) -> None:
    print(message, file=sys.stderr)


def run(args: list[str], cwd: Path, check: bool = True) -> subprocess.CompletedProcess[str]:
    eprint("$ " + " ".join(args))
    completed = subprocess.run(
        args,
        cwd=str(cwd),
        text=True,
        capture_output=True,
        env=os.environ.copy(),
    )
    if check and completed.returncode != 0:
        raise CommandError(args, completed.returncode, completed.stdout, completed.stderr)
    return completed


def normalize_channel_url(raw_url: str) -> str:
    value = raw_url.strip()
    if not value:
        raise ValueError("channel URL is required")
    if "://" not in value:
        value = "https://" + value
    parsed = urlparse(value)
    if "youtube.com" not in parsed.netloc.lower():
        raise ValueError(f"expected a youtube.com channel URL, got {raw_url!r}")
    if parsed.path == "/watch" or parse_qs(parsed.query).get("v"):
        raise ValueError("expected a channel URL, not a video URL")
    path = parsed.path.rstrip("/")
    if not path:
        raise ValueError("channel URL path is required")
    if path.endswith("/videos"):
        normalized_path = path
    elif path.endswith("/shorts") or path.endswith("/streams"):
        normalized_path = path.rsplit("/", 1)[0] + "/videos"
    else:
        normalized_path = path + "/videos"
    return urlunparse((parsed.scheme or "https", parsed.netloc, normalized_path, "", "", ""))


def resolve_videos(channel_url: str, limit: int | None, yt_dlp: str, vault: Path) -> list[Video]:
    command = [
        yt_dlp,
        "--flat-playlist",
        "--no-warnings",
        "--print",
        "%(id)s\t%(title)s\t%(webpage_url)s",
    ]
    if limit is not None:
        command.extend(["--playlist-end", str(limit)])
    command.append(channel_url)
    completed = run(command, cwd=vault)
    videos: list[Video] = []
    seen: set[str] = set()
    for line in completed.stdout.splitlines():
        parts = line.split("\t")
        if len(parts) < 2:
            continue
        video_id = parts[0].strip()
        title = parts[1].strip()
        url = parts[2].strip() if len(parts) >= 3 else ""
        if not VIDEO_ID_RE.match(video_id):
            continue
        if not url or url == "NA":
            url = f"https://www.youtube.com/watch?v={video_id}"
        if video_id in seen:
            continue
        seen.add(video_id)
        videos.append(Video(video_id=video_id, title=title, url=url))
    if not videos:
        raise RuntimeError("yt-dlp returned no videos for the channel")
    if limit is not None and len(videos) < limit:
        eprint(f"warning: requested {limit} videos, resolved {len(videos)}")
    return videos


def validate_inputs(args: argparse.Namespace) -> None:
    if not SLUG_RE.match(args.topic_slug):
        raise ValueError("topic slug must use lowercase alphanumerics separated by single hyphens")
    if not args.title.strip():
        raise ValueError("title is required")
    if not args.domain.strip():
        raise ValueError("domain is required")


def topic_paths(vault: Path, slug: str) -> tuple[Path, Path]:
    return vault / slug, vault / "yt-channels" / slug


def scaffold_topic(vault: Path, slug: str, title: str, domain: str, kb_path: str) -> Path:
    root_topic, category_topic = topic_paths(vault, slug)
    category_root = vault / "yt-channels"
    category_root.mkdir(parents=True, exist_ok=True)
    if category_topic.exists():
        if not (category_topic / "CLAUDE.md").exists():
            raise RuntimeError(f"existing topic is missing CLAUDE.md: {category_topic}")
        eprint(f"reusing existing topic {category_topic}")
        return category_topic
    if root_topic.exists():
        raise RuntimeError(f"cannot scaffold {slug}: root topic already exists at {root_topic}")
    run([kb_path, "topic", "new", slug, title, domain], cwd=vault)
    shutil.move(str(root_topic), str(category_topic))
    return category_topic


def write_topic_metadata(topic_dir: Path, slug: str, title: str, domain: str) -> None:
    content = "\n".join(
        [
            f"slug: {slug}",
            f"title: {title}",
            f"domain: {domain}",
            "category: yt-channels",
            f"path: yt-channels/{slug}",
            f"qmd_collection: {slug}",
            "",
        ]
    )
    (topic_dir / "topic.yaml").write_text(content, encoding="utf-8")


def ensure_agents_symlink(topic_dir: Path) -> None:
    agents_path = topic_dir / "AGENTS.md"
    if agents_path.exists() or agents_path.is_symlink():
        return
    try:
        agents_path.symlink_to("CLAUDE.md")
    except OSError:
        agents_path.write_text((topic_dir / "CLAUDE.md").read_text(encoding="utf-8"), encoding="utf-8")


def read_simple_yaml(path: Path) -> dict[str, str]:
    values: dict[str, str] = {}
    if not path.exists():
        return values
    for line in path.read_text(encoding="utf-8").splitlines():
        if ":" not in line:
            continue
        key, value = line.split(":", 1)
        values[key.strip()] = value.strip().strip('"')
    return values


def category_topics(category_root: Path) -> list[dict[str, str]]:
    topics: list[dict[str, str]] = []
    if not category_root.exists():
        return topics
    for child in sorted(category_root.iterdir()):
        if not child.is_dir():
            continue
        metadata = read_simple_yaml(child / "topic.yaml")
        if not metadata:
            continue
        slug = metadata.get("slug", child.name)
        title = metadata.get("title", slug)
        collection = metadata.get("qmd_collection", slug)
        topics.append(
            {
                "folder": child.name,
                "slug": slug,
                "title": title,
                "collection": collection,
            }
        )
    return topics


def update_category_docs(vault: Path) -> None:
    category_root = vault / "yt-channels"
    topics = category_topics(category_root)
    if topics:
        table_rows = [
            f"| [{topic['folder']}/]({topic['folder']}/) | {topic['title']} | {topic['slug']} | {topic['collection']} |"
            for topic in topics
        ]
        topic_lines = [
            f"- yt-channels/{topic['folder']}/ - {topic['title']} (collection: {topic['collection']})"
            for topic in topics
        ]
    else:
        table_rows = ["| _None yet._ | Use the kb-yt-channel skill to create a channel topic. | - | - |"]
        topic_lines = ["- None yet. Use the kb-yt-channel skill to create channel topics."]
    readme = "\n".join(
        [
            "# YouTube Channels",
            "",
            "Knowledge bases generated from YouTube channel uploads.",
            "",
            "| Folder | Topic | Slug | Collection |",
            "|--------|-------|------|------------|",
            *table_rows,
            "",
        ]
    )
    (category_root / "README.md").write_text(readme, encoding="utf-8")
    for name in ("CLAUDE.md", "AGENTS.md"):
        path = category_root / name
        if not path.exists():
            continue
        text = path.read_text(encoding="utf-8")
        replacement = "## Current topics\n\n" + "\n".join(topic_lines) + "\n"
        text = re.sub(r"## Current topics\n\n.*\Z", replacement, text, flags=re.DOTALL)
        path.write_text(text, encoding="utf-8")


def patch_topic_claude(
    topic_dir: Path,
    slug: str,
    title: str,
    domain: str,
    channel_url: str,
    limit: int | None,
    transcribe: str,
    sub_langs: str,
    command_line: str,
) -> None:
    claude_path = topic_dir / "CLAUDE.md"
    text = claude_path.read_text(encoding="utf-8")
    text = text.replace("[root CLAUDE.md](../CLAUDE.md)", "[root CLAUDE.md](../../CLAUDE.md)")
    text = text.replace("../.claude/skills/karpathy-kb/SKILL.md", "../../.claude/skills/kb/SKILL.md")
    text = text.replace("the [karpathy-kb skill]", "the [kb skill]")
    scope = (
        f"**Topic scope:** Transcripts and source material from the YouTube channel {channel_url}. "
        "This topic starts as an immutable transcript corpus and can later be compiled into wiki articles about recurring themes, speakers, demos, and technical patterns."
    )
    text = re.sub(r"^\*\*Topic scope:\*\*.*$", scope, text, count=1, flags=re.MULTILINE)
    domain_line = f"**Domain:** {domain} - all notes in this topic use domain: {domain} in frontmatter."
    text = re.sub(r"^\*\*Domain:\*\*.*$", domain_line, text, count=1, flags=re.MULTILINE)
    marker_start = "<!-- kb-yt-channel:start -->"
    marker_end = "<!-- kb-yt-channel:end -->"
    selection = "all uploads" if limit is None else f"latest {limit} uploads"
    caption_languages = sub_langs.strip() or "config/default"
    section = f"""
{marker_start}
## YouTube channel ingest

- Channel URL: {channel_url}
- Topic id: yt-channels/{slug}
- QMD collection: {slug}
- Selection policy: {selection}
- Transcript policy: {transcribe}
- Caption languages: {caption_languages}
- Last ingest command: {command_line}

Raw transcripts live in raw/youtube/. Ingest reports live in outputs/reports/. This topic is functional after kb topic info, kb lint --save, and kb index --name {slug} pass for yt-channels/{slug}.
{marker_end}
""".strip()
    if marker_start in text and marker_end in text:
        pattern = re.compile(re.escape(marker_start) + r".*?" + re.escape(marker_end), re.DOTALL)
        text = pattern.sub(section, text)
    else:
        text = text.rstrip() + "\n\n" + section + "\n"
    if not text.startswith(f"# {title}"):
        text = re.sub(r"^# .+$", f"# {title}", text, count=1, flags=re.MULTILINE)
    claude_path.write_text(text, encoding="utf-8")


def table_cell(value: str) -> str:
    return value.replace("|", "\\|").replace("\n", " ").strip()


def read_frontmatter(path: Path) -> dict[str, str]:
    text = path.read_text(encoding="utf-8", errors="ignore")
    if not text.startswith("---\n"):
        return {}
    end = text.find("\n---", 4)
    if end == -1:
        return {}
    values: dict[str, str] = {}
    for line in text[4:end].splitlines():
        if ":" not in line or line.startswith(" "):
            continue
        key, value = line.split(":", 1)
        values[key.strip()] = value.strip().strip('"').strip("'")
    return values


def youtube_sources(topic_dir: Path) -> list[dict[str, str]]:
    sources: list[dict[str, str]] = []
    raw_youtube = topic_dir / "raw" / "youtube"
    if not raw_youtube.exists():
        return sources
    for path in sorted(raw_youtube.glob("*.md")):
        metadata = read_frontmatter(path)
        title = metadata.get("title") or path.stem.replace("-", " ").title()
        sources.append(
            {
                "path": path.relative_to(topic_dir).as_posix(),
                "title": title,
                "scraped": metadata.get("scraped", ""),
                "source_url": metadata.get("source_url", ""),
                "transcript_source": metadata.get("transcript_source", ""),
            }
        )
    return sources


def update_topic_inventory(topic_dir: Path, source_count: int) -> None:
    claude_path = topic_dir / "CLAUDE.md"
    text = claude_path.read_text(encoding="utf-8")
    line = f"- `raw/youtube/` - {source_count} transcript sources"
    if "`raw/youtube/`" in text:
        text = re.sub(r"^- `raw/youtube/` .*$", line, text, count=1, flags=re.MULTILINE)
    else:
        text = re.sub(r"^- `raw/github/` .*$", lambda match: match.group(0) + "\n" + line, text, count=1, flags=re.MULTILINE)
    claude_path.write_text(text, encoding="utf-8")


def write_channel_indexes(topic_dir: Path, title: str, domain: str, source_count: int) -> None:
    today = dt.date.today().isoformat()
    index_dir = topic_dir / "wiki" / "index"
    index_dir.mkdir(parents=True, exist_ok=True)
    dashboard = "\n".join(
        [
            "---",
            f"domain: {domain}",
            "title: Dashboard",
            "type: index",
            f"updated: \"{today}\"",
            "---",
            "",
            f"# {title} - Dashboard",
            "",
            f"Landing page for the {title} knowledge base.",
            "",
            "## At a glance",
            "",
            "- **Articles:** 0",
            "- **Total words:** 0",
            f"- **Raw sources:** {source_count} YouTube transcript sources",
            f"- **Last updated:** {today}",
            "",
            "## Topic scope",
            "",
            "YouTube transcript corpus for recurring themes, speakers, demos, and technical patterns from this channel.",
            "",
            "## Featured articles",
            "",
            "_No featured articles yet._",
            "",
            "## Recent additions",
            "",
            "See [[Source Index]] for the current transcript corpus.",
            "",
            "## Coverage map",
            "",
            "_Compile wiki concepts after the transcript corpus has been reviewed._",
            "",
            "## Research gaps",
            "",
            "- Identify recurring agent-engineering practices across talks",
            "- Map speakers, projects, and demos to technical themes",
            "",
            "## Related topics",
            "",
            "_No related topics linked yet._",
            "",
            "## Navigation",
            "",
            "- [[Concept Index]] - alphabetical listing of all articles",
            "- [[Source Index]] - all sources and which articles cite them",
            "",
        ]
    )
    (index_dir / "Dashboard.md").write_text(dashboard, encoding="utf-8")
    source_lines = [
        "---",
        f"domain: {domain}",
        "title: Source Index",
        "type: index",
        f"updated: \"{today}\"",
        "---",
        "",
        f"# {title} - Source Index",
        "",
        "All raw sources that inform this topic's wiki, with the articles that cite them.",
        "",
        "## YouTube transcripts",
        "",
        "| Source | Scraped | Transcript | URL | Cited by |",
        "|--------|---------|------------|-----|----------|",
    ]
    sources = youtube_sources(topic_dir)
    if sources:
        for source in sources:
            source_lines.append(
                "| "
                + f"[[{source['path']}|{table_cell(source['title'])}]]"
                + f" | {table_cell(source['scraped'] or 'unknown')}"
                + f" | {table_cell(source['transcript_source'] or 'unknown')}"
                + f" | {table_cell(source['source_url'] or 'n/a')}"
                + " | _Uncited_ |"
            )
    else:
        source_lines.append("| _None yet._ | - | - | - | - |")
    source_lines.extend(
        [
            "",
            "## Articles (raw/articles/)",
            "",
            "| Source | Scraped | Cited by |",
            "|--------|---------|----------|",
            "| _None yet._ | - | - |",
            "",
            "## GitHub (raw/github/)",
            "",
            "| Source | Scraped | Cited by |",
            "|--------|---------|----------|",
            "| _None yet._ | - | - |",
            "",
            "## Bookmark clusters (raw/bookmarks/)",
            "",
            "| Cluster | Updated | Cited by |",
            "|---------|---------|----------|",
            "| _None yet._ | - | - |",
            "",
            "## Orphan sources",
            "",
            "All YouTube transcripts are uncited until wiki articles are compiled from this corpus.",
            "",
        ]
    )
    (index_dir / "Source Index.md").write_text("\n".join(source_lines), encoding="utf-8")


def parse_json_output(text: str) -> Any:
    stripped = text.strip()
    if not stripped:
        return None
    try:
        return json.loads(stripped)
    except json.JSONDecodeError:
        return None


def write_report(topic_dir: Path, summary: dict[str, Any]) -> Path:
    reports = topic_dir / "outputs" / "reports"
    reports.mkdir(parents=True, exist_ok=True)
    now = dt.datetime.now()
    stamp = now.strftime("%Y-%m-%d-%H%M%S")
    path = reports / f"{stamp}-youtube-channel-ingest.md"
    lines = [
        "---",
        f"title: YouTube Channel Ingest Report - {summary['topic']}",
        "type: output",
        "stage: lint-report",
        f"domain: {summary.get('domain', 'youtube-channel')}",
        f"created: {now.strftime('%Y-%m-%d')}",
        f"issues_found: {len(summary['failures'])}",
        "issues_fixed: 0",
        "tags:",
        "  - youtube-channel",
        "  - ingest",
        "---",
        "",
        "# YouTube Channel Ingest Report",
        "",
        f"- Channel URL: {summary['channelUrl']}",
        f"- Normalized channel URL: {summary['normalizedChannelUrl']}",
        f"- Topic: {summary['topic']}",
        f"- Transcript policy: {summary['transcribe']}",
        f"- Caption languages: {summary.get('captionLanguages', 'config/default')}",
        f"- Selection: {summary['selection']}",
        f"- Resolved videos: {len(summary['videos'])}",
        f"- Successful ingests: {len(summary['ingested'])}",
        f"- Skipped existing videos: {len(summary['skipped'])}",
        f"- Failures: {len(summary['failures'])}",
        "",
        "## Videos",
        "",
    ]
    for video in summary["videos"]:
        lines.append(f"- {video['video_id']} - {video['title']} - {video['url']}")
    lines.extend(["", "## Ingested", ""])
    for item in summary["ingested"]:
        lines.append(f"- {item['video_id']} - {item['title']}")
    lines.extend(["", "## Skipped", ""])
    for item in summary["skipped"]:
        lines.append(f"- {item['video_id']} - {item['title']}")
    lines.extend(["", "## Failures", ""])
    if summary["failures"]:
        for item in summary["failures"]:
            lines.extend(
                [
                    f"### {item['video_id']} - {item['title']}",
                    "",
                    f"- URL: {item.get('url', '') or 'n/a'}",
                    f"- Error: {item['error']}",
                ]
            )
            for stream_name in ("stderr", "stdout"):
                stream = str(item.get(stream_name, "")).strip()
                if stream:
                    lines.extend(
                        [
                            "",
                            f"{stream_name}:",
                            "",
                            "```text",
                            stream,
                            "```",
                        ]
                    )
            lines.append("")
    else:
        lines.append("- None")
    lines.extend(["", "## Validation", ""])
    for item in summary.get("validation", []):
        lines.append(f"- {item['command']}: exit {item['exit_code']}")
    path.write_text("\n".join(lines) + "\n", encoding="utf-8")
    return path


def build_command_line(args: argparse.Namespace) -> str:
    selector = "--all" if args.all else f"--limit {args.limit}"
    embed = " --embed" if args.embed else ""
    sub_langs = f" --sub-langs {args.sub_langs}" if args.sub_langs else ""
    return (
        "python3 .agents/skills/kb-yt-channel/scripts/ingest-channel.py "
        f"--vault {args.vault} --channel-url {args.channel_url} --topic-slug {args.topic_slug} "
        f"--title {json.dumps(args.title)} --domain {args.domain} {selector} --transcribe {args.transcribe}{sub_langs}{embed}"
    )


def validation_command(label: str, args: list[str], vault: Path) -> dict[str, Any]:
    completed = run(args, cwd=vault, check=False)
    return {
        "command": label,
        "args": args,
        "exit_code": completed.returncode,
        "stdout": completed.stdout,
        "stderr": completed.stderr,
        "json": parse_json_output(completed.stdout),
    }


def run_channel_ingest(
    args: argparse.Namespace,
    normalized_channel_url: str,
    limit: int | None,
    vault: Path,
) -> dict[str, Any]:
    command = [
        args.kb_path,
        "ingest",
        "channel",
        normalized_channel_url,
        "--topic",
        f"yt-channels/{args.topic_slug}",
        "--transcribe",
        args.transcribe,
    ]
    if args.sub_langs:
        command.extend(["--sub-langs", args.sub_langs])
    if limit is not None:
        command.extend(["--limit", str(limit)])
    completed = run(command, cwd=vault, check=False)
    result = parse_json_output(completed.stdout)
    if not isinstance(result, dict):
        raise CommandError(command, completed.returncode, completed.stdout, completed.stderr)
    return result


def run_ingest(args: argparse.Namespace) -> dict[str, Any]:
    validate_inputs(args)
    vault = Path(args.vault).expanduser().resolve()
    if not vault.exists():
        raise RuntimeError(f"vault does not exist: {vault}")
    normalized_channel_url = normalize_channel_url(args.channel_url)
    limit = None if args.all else args.limit
    selection = "all uploads" if limit is None else f"latest {limit} uploads"
    summary: dict[str, Any] = {
        "channelUrl": args.channel_url,
        "normalizedChannelUrl": normalized_channel_url,
        "topic": f"yt-channels/{args.topic_slug}",
        "topicPath": str(vault / "yt-channels" / args.topic_slug),
        "domain": args.domain,
        "transcribe": args.transcribe,
        "captionLanguages": args.sub_langs or "config/default",
        "selection": selection,
        "videos": [],
        "ingested": [],
        "skipped": [],
        "failures": [],
        "validation": [],
    }
    if args.dry_run:
        videos = resolve_videos(normalized_channel_url, limit, args.yt_dlp_path, vault)
        summary["videos"] = [video.__dict__ for video in videos]
        summary["dryRun"] = True
        return summary
    topic_dir = scaffold_topic(vault, args.topic_slug, args.title, args.domain, args.kb_path)
    command_line = build_command_line(args)
    write_topic_metadata(topic_dir, args.topic_slug, args.title, args.domain)
    ensure_agents_symlink(topic_dir)
    update_category_docs(vault)
    patch_topic_claude(
        topic_dir,
        args.topic_slug,
        args.title,
        args.domain,
        normalized_channel_url,
        limit,
        args.transcribe,
        args.sub_langs,
        command_line,
    )
    # Enumeration, resume, throttling, and adaptive backoff live in the kb binary.
    channel_result = run_channel_ingest(args, normalized_channel_url, limit, vault)
    summary["videos"] = channel_result.get("videos", [])
    summary["captionLanguages"] = channel_result.get("caption_languages", summary["captionLanguages"])
    summary["ingested"] = channel_result.get("ingested", [])
    summary["skipped"] = channel_result.get("skipped", [])
    summary["failures"].extend(channel_result.get("failures", []))
    source_count = len(youtube_sources(topic_dir))
    update_topic_inventory(topic_dir, source_count)
    write_channel_indexes(topic_dir, args.title, args.domain, source_count)
    index_command = [args.kb_path, "index", "--topic", f"yt-channels/{args.topic_slug}", "--name", args.topic_slug]
    if not args.embed:
        index_command.append("--embed=false")
    validation_specs = [
        ("topic info", [args.kb_path, "topic", "info", f"yt-channels/{args.topic_slug}"]),
        ("lint", [args.kb_path, "lint", f"yt-channels/{args.topic_slug}", "--save"]),
        ("index", index_command),
        (
            "search",
            [
                args.kb_path,
                "search",
                args.topic_slug,
                "--topic",
                f"yt-channels/{args.topic_slug}",
                "--collection",
                args.topic_slug,
                "--lex",
                "--format",
                "json",
            ],
        ),
    ]
    for label, command in validation_specs:
        result = validation_command(label, command, vault)
        summary["validation"].append(result)
        if result["exit_code"] != 0:
            summary["failures"].append(
                {
                    "video_id": "validation",
                    "title": label,
                    "url": "",
                    "error": f"validation command failed: {label}",
                    "stdout": result["stdout"],
                    "stderr": result["stderr"],
                }
            )
            break
    topic_info = next((item.get("json") for item in summary["validation"] if item["command"] == "topic info"), None)
    if isinstance(topic_info, dict):
        source_count = int(topic_info.get("sourceCount", 0))
        successful = len(summary["ingested"]) + len(summary["skipped"])
        if source_count < successful:
            summary["failures"].append(
                {
                    "video_id": "validation",
                    "title": "source count",
                    "url": "",
                    "error": f"sourceCount {source_count} is lower than successful/skipped video count {successful}",
                    "stdout": "",
                    "stderr": "",
                }
            )
    report_path = write_report(topic_dir, summary)
    summary["reportPath"] = str(report_path)
    summary["dryRun"] = False
    return summary


def parse_args(argv: list[str]) -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Create or update a KB topic from a YouTube channel.")
    parser.add_argument("--vault", default=".", help="Vault root path")
    parser.add_argument("--channel-url", required=True, help="YouTube channel URL")
    parser.add_argument("--topic-slug", required=True, help="Topic slug under yt-channels/")
    parser.add_argument("--title", required=True, help="Topic title")
    parser.add_argument("--domain", required=True, help="Topic domain")
    selector = parser.add_mutually_exclusive_group(required=True)
    selector.add_argument("--limit", type=int, help="Maximum newest uploads to ingest")
    selector.add_argument("--all", action="store_true", help="Ingest all channel uploads")
    parser.add_argument("--transcribe", choices=["captions", "auto", "stt"], default="captions")
    parser.add_argument("--sub-langs", default="", help="Comma-separated caption languages; use orig for the video's original language")
    parser.add_argument("--yt-dlp-path", default=os.environ.get("YOUTUBE_YT_DLP_PATH", "yt-dlp"))
    parser.add_argument("--kb-path", default="kb")
    parser.add_argument("--embed", action="store_true", help="Run vector embedding after QMD collection sync")
    parser.add_argument("--dry-run", action="store_true", help="Resolve videos without mutating the vault")
    args = parser.parse_args(argv)
    args.sub_langs = args.sub_langs.strip()
    if args.limit is not None and args.limit < 1:
        parser.error("--limit must be greater than zero")
    return args


def main(argv: list[str]) -> int:
    try:
        args = parse_args(argv)
        summary = run_ingest(args)
        print(json.dumps(summary, indent=2, ensure_ascii=False))
        return 1 if summary.get("failures") else 0
    except (CommandError, RuntimeError, ValueError) as err:
        print(json.dumps({"error": str(err)}, indent=2), file=sys.stdout)
        return 1


if __name__ == "__main__":
    raise SystemExit(main(sys.argv[1:]))
