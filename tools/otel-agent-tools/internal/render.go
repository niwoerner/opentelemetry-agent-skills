package otelagenttools

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"
)

type VersionEntry struct {
	Language    string
	Package     string
	Version     string
	SourceLabel string
	SourceURL   string
	DocsURL     string
	ExamplesURL string
}

type versionIndexView struct {
	Rows []versionIndexRow
}

type versionIndexRow struct {
	Language     string
	PackageCell  string
	VersionCell  string
	SourceCell   string
	DocsCell     string
	ExamplesCell string
}

var versionIndexTemplate = template.Must(template.New("version-index").Parse(`# OpenTelemetry SDK Version Index

| Language | Package/Repo | Latest | Release Source | Setup Docs | Examples |
|---|---|---|---|---|---|
{{- range .Rows }}
| {{ .Language }} | {{ .PackageCell }} | {{ .VersionCell }} | {{ .SourceCell }} | {{ .DocsCell }} | {{ .ExamplesCell }} |
{{- end }}
`))

func RenderVersionIndex(entries []VersionEntry) ([]byte, error) {
	rows := make([]versionIndexRow, 0, len(entries))
	for _, entry := range entries {
		rows = append(rows, versionIndexRow{
			Language:     entry.Language,
			PackageCell:  fmt.Sprintf("`%s`", entry.Package),
			VersionCell:  fmt.Sprintf("`%s`", entry.Version),
			SourceCell:   markdownLink(entry.SourceLabel, entry.SourceURL),
			DocsCell:     markdownLink("Docs", entry.DocsURL),
			ExamplesCell: optionalMarkdownLink("Examples", entry.ExamplesURL),
		})
	}

	slices.SortFunc(rows, func(a, b versionIndexRow) int {
		if cmp := strings.Compare(a.Language, b.Language); cmp != 0 {
			return cmp
		}
		return strings.Compare(a.PackageCell, b.PackageCell)
	})

	view := versionIndexView{Rows: rows}

	var buf bytes.Buffer
	if err := versionIndexTemplate.Execute(&buf, view); err != nil {
		return nil, fmt.Errorf("render version index: %w", err)
	}
	return buf.Bytes(), nil
}

func markdownLink(label, rawURL string) string {
	return fmt.Sprintf("[%s](%s)", label, rawURL)
}

func optionalMarkdownLink(label, rawURL string) string {
	if rawURL == "" {
		return "-"
	}
	return markdownLink(label, rawURL)
}
