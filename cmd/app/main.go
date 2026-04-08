package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/user/go-devstack/internal/config"
	"github.com/user/go-devstack/internal/logger"
	"github.com/user/go-devstack/internal/version"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if err := execute(args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}

func execute(args []string) error {
	command, rest := splitCommand(args)
	if command == "version" {
		fmt.Println(version.String())
		return nil
	}
	if command != "run" {
		return fmt.Errorf("unknown command %q", command)
	}
	return runCommand(rest)
}

func splitCommand(args []string) (string, []string) {
	if len(args) == 0 {
		return "run", nil
	}
	return args[0], args[1:]
}

func runCommand(args []string) error {
	if err := config.LoadDotEnvIfPresent(config.DefaultDotEnvPath); err != nil {
		return err
	}

	configPath, err := parseRunFlags(args)
	if err != nil {
		return err
	}

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	log, err := logger.New(cfg.Log.Level)
	if err != nil {
		return err
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	log.Info("starting", "app", cfg.App.Name, "env", cfg.App.Env)

	// Wire your application services below using the launch pattern.
	// Each service should accept a context.Context and return an error.
	var wg sync.WaitGroup
	errCh := make(chan error, 1)
	launch := func(fn func() error) {
		wg.Add(1)
		go func() {
			defer wg.Done()
			errCh <- fn()
		}()
	}
	_ = launch

	// TODO: launch your services here. Example:
	// launch(func() error { return myService.Run(ctx) })

	<-ctx.Done()
	log.Info("shutdown signal received")
	cancel()

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
		close(errCh)
	}()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		log.Warn("shutdown deadline exceeded")
	}

	var runErr error
	for err := range errCh {
		if err == nil || errors.Is(err, context.Canceled) {
			continue
		}
		if runErr == nil {
			runErr = err
		}
	}
	return runErr
}

func parseRunFlags(args []string) (string, error) {
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	configPath := fs.String("config", os.Getenv(config.EnvConfigPath), "path to config file")
	if err := fs.Parse(args); err != nil {
		return "", err
	}
	return *configPath, nil
}
