package otelagenttools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type SDKSource struct {
	Language    string `json:"language"`
	Package     string `json:"package"`
	SourceKind  string `json:"source_kind"`
	SourceLabel string `json:"source_label"`
	SourceURL   string `json:"source_url"`
	Target      string `json:"target"`
	DocsURL     string `json:"docs_url"`
	ExamplesURL string `json:"examples_url"`
}

func LoadSources(path string) ([]SDKSource, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read sources: %w", err)
	}

	var sources []SDKSource
	if err := json.Unmarshal(data, &sources); err != nil {
		return nil, fmt.Errorf("parse sources json: %w", err)
	}

	if len(sources) == 0 {
		return nil, fmt.Errorf("no sdk sources configured")
	}

	seenKeys := map[string]struct{}{}
	for i := range sources {
		sources[i].Language = strings.TrimSpace(sources[i].Language)
		sources[i].Package = strings.TrimSpace(sources[i].Package)
		sources[i].SourceKind = strings.TrimSpace(sources[i].SourceKind)
		sources[i].SourceLabel = strings.TrimSpace(sources[i].SourceLabel)
		sources[i].SourceURL = strings.TrimSpace(sources[i].SourceURL)
		sources[i].Target = strings.TrimSpace(sources[i].Target)
		sources[i].DocsURL = strings.TrimSpace(sources[i].DocsURL)
		sources[i].ExamplesURL = strings.TrimSpace(sources[i].ExamplesURL)

		if sources[i].Language == "" || sources[i].Package == "" || sources[i].SourceKind == "" || sources[i].SourceLabel == "" || sources[i].SourceURL == "" || sources[i].Target == "" {
			return nil, fmt.Errorf("source %d has empty required fields", i)
		}
		key := strings.Join([]string{
			sources[i].Language,
			sources[i].Package,
			sources[i].SourceKind,
			sources[i].Target,
		}, "\x00")
		if _, exists := seenKeys[key]; exists {
			return nil, fmt.Errorf("duplicate source entry for language=%s package=%s target=%s", sources[i].Language, sources[i].Package, sources[i].Target)
		}
		seenKeys[key] = struct{}{}
	}

	return sources, nil
}
