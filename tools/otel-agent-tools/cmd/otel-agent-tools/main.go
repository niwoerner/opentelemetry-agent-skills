package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	otelagenttools "otel-agent-tools/internal"
)

const (
	defaultConfigFilename = "sdk-sources.json"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("missing command")
	}

	switch args[0] {
	case "reconcile-versions":
		return runReconcileVersions(args[1:])
	case "reconcile":
		if len(args) < 2 {
			return fmt.Errorf("missing reconcile target")
		}
		switch args[1] {
		case "versions":
			return runReconcileVersions(args[2:])
		default:
			return fmt.Errorf("unknown reconcile target: %s", args[1])
		}
	case "-h", "--help", "help":
		printUsage()
		return nil
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runReconcileVersions(args []string) error {
	fs := flag.NewFlagSet("reconcile-versions", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	configPath := fs.String("config", resolveConfigPath(), "Path to sdk-sources.json")
	outputPath := fs.String("out", "", "Path to generated markdown index")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *outputPath == "" {
		return fmt.Errorf("missing required --out path")
	}

	sources, err := otelagenttools.LoadSources(*configPath)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	entries := make([]otelagenttools.VersionEntry, 0, len(sources))
	for _, source := range sources {
		version, err := otelagenttools.FetchVersion(ctx, source.SourceKind, source.Target)
		if err != nil {
			return fmt.Errorf("%s: %w", source.Language, err)
		}
		entries = append(entries, otelagenttools.VersionEntry{
			Language:    source.Language,
			Package:     source.Package,
			Version:     version,
			SourceLabel: source.SourceLabel,
			SourceURL:   source.SourceURL,
			DocsURL:     source.DocsURL,
			ExamplesURL: source.ExamplesURL,
		})
	}

	rendered, err := otelagenttools.RenderVersionIndex(entries)
	if err != nil {
		return err
	}

	unchanged, err := writeIfChanged(*outputPath, rendered)
	if err != nil {
		return err
	}
	if unchanged {
		fmt.Printf("up to date: %s\n", *outputPath)
		return nil
	}

	fmt.Printf("wrote: %s\n", *outputPath)
	return nil
}

func writeIfChanged(path string, next []byte) (bool, error) {
	current, err := os.ReadFile(path)
	if err == nil && bytes.Equal(current, next) {
		return true, nil
	}
	if err != nil && !os.IsNotExist(err) {
		return false, fmt.Errorf("read output: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, fmt.Errorf("create output directory: %w", err)
	}
	if err := os.WriteFile(path, next, 0o644); err != nil {
		return false, fmt.Errorf("write output: %w", err)
	}
	return false, nil
}

func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  otel-agent-tools reconcile versions --out path [--config path]")
	fmt.Println("  otel-agent-tools reconcile-versions --out path [--config path]")
}

func resolveConfigPath() string {
	candidates := []string{
		defaultConfigFilename,
		filepath.Join("tools", "otel-agent-tools", defaultConfigFilename),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
	}
	return defaultConfigFilename
}
