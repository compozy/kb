package vault

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/user/go-devstack/internal/models"
)

func RenderBaseFiles(metrics models.MetricsResult) []models.BaseFile {
	_ = metrics

	return []models.BaseFile{
		createBaseFile(GetBaseFilePath("symbol-explorer"), models.BaseDefinition{
			Filters: andFilter(exprFilter(`file.hasTag("symbol")`)),
			Views: []models.BaseView{
				{
					GroupBy: &models.BaseGroupBy{Direction: "ASC", Property: "symbol_kind"},
					Name:    "All Symbols",
					Order:   []string{"file.name", "symbol_kind", "source_path", "language", "exported"},
					Type:    models.ViewTable,
				},
				{
					Filters: &models.BaseFilter{
						Or: []models.BaseFilter{
							exprFilter(`symbol_kind == "function"`),
							exprFilter(`symbol_kind == "method"`),
						},
					},
					Name:  "Functions",
					Order: []string{"file.name", "symbol_kind", "source_path", "language", "exported"},
					Type:  models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("file-explorer"), models.BaseDefinition{
			Filters: andFilter(exprFilter(`file.hasTag("file")`)),
			Views: []models.BaseView{
				{
					GroupBy: &models.BaseGroupBy{Direction: "ASC", Property: "file.folder"},
					Name:    "Files",
					Order: []string{
						"file.name",
						"language",
						"symbol_count",
						"incoming_relation_count",
						"outgoing_relation_count",
					},
					Type: models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("most-connected"), models.BaseDefinition{
			Filters: andFilter(
				exprFilter(`file.hasTag("symbol")`),
				exprFilter("outgoing_relation_count > 5"),
			),
			Formulas: map[string]string{
				"total_relations": "incoming_relation_count + outgoing_relation_count",
			},
			Properties: map[string]models.BaseProperty{
				"formula.total_relations": {DisplayName: "Total Relations"},
			},
			Views: []models.BaseView{
				{
					Name: "Most Connected",
					Order: []string{
						"file.name",
						"symbol_kind",
						"formula.total_relations",
						"incoming_relation_count",
						"outgoing_relation_count",
						"source_path",
					},
					Type: models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("entry-points"), models.BaseDefinition{
			Filters: andFilter(
				exprFilter(`file.hasTag("symbol")`),
				exprFilter("exported == true"),
			),
			Views: []models.BaseView{
				{
					GroupBy: &models.BaseGroupBy{Direction: "ASC", Property: "source_path"},
					Name:    "Entry Points",
					Order:   []string{"file.name", "symbol_kind", "source_path", "incoming_relation_count"},
					Type:    models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("isolated-symbols"), models.BaseDefinition{
			Filters: andFilter(
				exprFilter(`file.hasTag("symbol")`),
				exprFilter("exported == true"),
				exprFilter("external_reference_count <= 0"),
			),
			Views: []models.BaseView{
				{
					Name:  "Isolated Symbols",
					Order: []string{"file.name", "symbol_kind", "source_path", "external_reference_count"},
					Type:  models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("complexity-hotspots"), models.BaseDefinition{
			Filters: andFilter(
				exprFilter(`file.hasTag("function")`),
				exprFilter("cyclomatic_complexity > 5"),
			),
			Views: []models.BaseView{
				{
					Name:  "Complexity Hotspots",
					Order: []string{"file.name", "cyclomatic_complexity", "loc", "source_path"},
					Summaries: map[string]string{
						"cyclomatic_complexity": "Average",
					},
					Type: models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("danger-zone"), models.BaseDefinition{
			Filters: andFilter(
				exprFilter(`file.hasTag("symbol")`),
				exprFilter("blast_radius > 10"),
			),
			Views: []models.BaseView{
				{
					Name:  "Danger Zone",
					Order: []string{"file.name", "blast_radius", "cyclomatic_complexity", "source_path", "smells"},
					Type:  models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("module-health"), models.BaseDefinition{
			Filters: andFilter(exprFilter(`file.hasTag("file")`)),
			Views: []models.BaseView{
				{
					Name:  "Module Health",
					Order: []string{"file.name", "afferent_coupling", "efferent_coupling", "instability"},
					Summaries: map[string]string{
						"instability": "Average",
					},
					Type: models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("code-smells"), models.BaseDefinition{
			Filters: andFilter(exprFilter("has_smells == true")),
			Views: []models.BaseView{
				{
					Name:  "Code Smells",
					Order: []string{"file.name", "smells", "source_path", "symbol_kind"},
					Type:  models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("circular-deps"), models.BaseDefinition{
			Filters: andFilter(
				exprFilter(`file.hasTag("file")`),
				exprFilter("has_circular_dependency == true"),
			),
			Views: []models.BaseView{
				{
					Name:  "Circular Dependencies",
					Order: []string{"file.name", "source_path", "afferent_coupling", "efferent_coupling"},
					Type:  models.ViewTable,
				},
			},
		}),
		createBaseFile(GetBaseFilePath("dead-code"), models.BaseDefinition{
			Filters: &models.BaseFilter{
				Or: []models.BaseFilter{
					{And: []models.BaseFilter{
						exprFilter(`file.hasTag("symbol")`),
						exprFilter("is_dead_export == true"),
					}},
					{And: []models.BaseFilter{
						exprFilter(`file.hasTag("file")`),
						exprFilter("is_orphan_file == true"),
					}},
				},
			},
			Views: []models.BaseView{
				{
					Name:  "Dead Code",
					Order: []string{"file.name", "symbol_kind", "source_path", "is_dead_export", "is_orphan_file"},
					Type:  models.ViewTable,
				},
			},
		}),
	}
}

func RenderBaseDefinition(definition models.BaseDefinition) string {
	lines := renderYAMLValue(baseDefinitionValue(definition), 0)
	return strings.Join(append(lines, ""), "\n")
}

func createBaseFile(relativePath string, definition models.BaseDefinition) models.BaseFile {
	return models.BaseFile{
		Definition:   definition,
		RelativePath: relativePath,
	}
}

func exprFilter(expression string) models.BaseFilter {
	return models.BaseFilter{Expression: expression}
}

func andFilter(conditions ...models.BaseFilter) *models.BaseFilter {
	return &models.BaseFilter{And: conditions}
}

func baseDefinitionValue(definition models.BaseDefinition) map[string]interface{} {
	value := map[string]interface{}{
		"views": baseViewsValue(definition.Views),
	}

	if definition.Filters != nil {
		value["filters"] = baseFilterValue(*definition.Filters)
	}
	if len(definition.Formulas) > 0 {
		value["formulas"] = stringMapToAny(definition.Formulas)
	}
	if len(definition.Properties) > 0 {
		properties := make(map[string]interface{}, len(definition.Properties))
		for key, property := range definition.Properties {
			properties[key] = map[string]interface{}{
				"displayName": property.DisplayName,
			}
		}
		value["properties"] = properties
	}

	return value
}

func baseViewsValue(views []models.BaseView) []interface{} {
	items := make([]interface{}, 0, len(views))
	for _, view := range views {
		value := map[string]interface{}{
			"name":  view.Name,
			"order": stringSliceToAny(view.Order),
			"type":  string(view.Type),
		}

		if view.Filters != nil {
			value["filters"] = baseFilterValue(*view.Filters)
		}
		if view.GroupBy != nil {
			value["groupBy"] = map[string]interface{}{
				"direction": view.GroupBy.Direction,
				"property":  view.GroupBy.Property,
			}
		}
		if len(view.Summaries) > 0 {
			value["summaries"] = stringMapToAny(view.Summaries)
		}

		items = append(items, value)
	}
	return items
}

func baseFilterValue(filter models.BaseFilter) interface{} {
	if filter.Expression != "" {
		return filter.Expression
	}

	value := map[string]interface{}{}
	if len(filter.And) > 0 {
		items := make([]interface{}, 0, len(filter.And))
		for _, nested := range filter.And {
			items = append(items, baseFilterValue(nested))
		}
		value["and"] = items
	}
	if len(filter.Or) > 0 {
		items := make([]interface{}, 0, len(filter.Or))
		for _, nested := range filter.Or {
			items = append(items, baseFilterValue(nested))
		}
		value["or"] = items
	}
	if filter.Not != nil {
		value["not"] = baseFilterValue(*filter.Not)
	}

	return value
}

func renderYAMLValue(value interface{}, indent int) []string {
	padding := strings.Repeat(" ", indent)

	switch typed := value.(type) {
	case []interface{}:
		if len(typed) == 0 {
			return []string{padding + "[]"}
		}

		lines := make([]string, 0)
		for _, item := range typed {
			switch nested := item.(type) {
			case string:
				lines = append(lines, padding+"- "+renderYAMLScalar(nested))
			case bool:
				lines = append(lines, padding+"- "+renderYAMLScalar(nested))
			case int:
				lines = append(lines, padding+"- "+renderYAMLScalar(nested))
			case float64:
				lines = append(lines, padding+"- "+renderYAMLScalar(nested))
			default:
				nestedLines := renderYAMLValue(item, indent+2)
				if len(nestedLines) == 0 {
					lines = append(lines, padding+"- {}")
					continue
				}

				lines = append(lines, padding+"- "+strings.TrimLeft(nestedLines[0], " "))
				lines = append(lines, nestedLines[1:]...)
			}
		}
		return lines
	case map[string]interface{}:
		keys := sortedMapKeys(typed)
		if len(keys) == 0 {
			return []string{padding + "{}"}
		}

		lines := make([]string, 0)
		for _, key := range keys {
			entryValue := typed[key]
			switch nested := entryValue.(type) {
			case string:
				lines = append(lines, padding+key+": "+renderYAMLScalar(nested))
			case bool:
				lines = append(lines, padding+key+": "+renderYAMLScalar(nested))
			case int:
				lines = append(lines, padding+key+": "+renderYAMLScalar(nested))
			case float64:
				lines = append(lines, padding+key+": "+renderYAMLScalar(nested))
			default:
				lines = append(lines, padding+key+":")
				lines = append(lines, renderYAMLValue(entryValue, indent+2)...)
			}
		}
		return lines
	case string:
		return []string{padding + renderYAMLScalar(typed)}
	case bool:
		return []string{padding + renderYAMLScalar(typed)}
	case int:
		return []string{padding + renderYAMLScalar(typed)}
	case float64:
		return []string{padding + renderYAMLScalar(typed)}
	default:
		return []string{padding + renderYAMLScalar(typed)}
	}
}

func renderYAMLScalar(value interface{}) string {
	switch typed := value.(type) {
	case string:
		return strconv.Quote(typed)
	case bool:
		return strconv.FormatBool(typed)
	case int:
		return strconv.Itoa(typed)
	case float64:
		return strconv.FormatFloat(typed, 'f', -1, 64)
	default:
		return strconv.Quote(fmt.Sprint(value))
	}
}

func stringMapToAny(values map[string]string) map[string]interface{} {
	converted := make(map[string]interface{}, len(values))
	for key, value := range values {
		converted[key] = value
	}
	return converted
}

func stringSliceToAny(values []string) []interface{} {
	items := make([]interface{}, 0, len(values))
	for _, value := range values {
		items = append(items, value)
	}
	return items
}
