package project

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	domain "github.com/develop-agent/backend/internal/domain/project"
)

const maxCodeContextTokens = 3000

type CodeContextAccumulator struct{}

func NewCodeContextAccumulator() *CodeContextAccumulator { return &CodeContextAccumulator{} }

func (a *CodeContextAccumulator) Build(files []*domain.CodeFile) domain.CodeContextManifest {
	manifest := domain.CodeContextManifest{
		Files:            make([]domain.CodeContextFile, 0, len(files)),
		Symbols:          make([]domain.CodeSymbol, 0),
		Dependencies:     make([]string, 0),
		EnvironmentHints: make([]string, 0),
	}

	depSet := map[string]struct{}{}
	envSet := map[string]struct{}{}
	for _, f := range files {
		manifest.Files = append(manifest.Files, domain.CodeContextFile{
			Path:     f.Path,
			Language: normalizeLanguage(f.Path, f.Language),
			Purpose:  inferPurpose(f.Path),
		})
		manifest.Symbols = append(manifest.Symbols, extractSymbols(f.Path, f.Content)...)
		for _, dep := range extractDependencies(f.Content) {
			depSet[dep] = struct{}{}
		}
		for _, env := range extractEnvironmentHints(f.Content) {
			envSet[env] = struct{}{}
		}
	}

	for dep := range depSet {
		manifest.Dependencies = append(manifest.Dependencies, dep)
	}
	for env := range envSet {
		manifest.EnvironmentHints = append(manifest.EnvironmentHints, env)
	}
	sort.Strings(manifest.Dependencies)
	sort.Strings(manifest.EnvironmentHints)

	manifest.ApproxTokens = approximateTokens(manifest)
	if manifest.ApproxTokens > maxCodeContextTokens {
		manifest = trimManifest(manifest, maxCodeContextTokens)
	}
	return manifest
}

func approximateTokens(m domain.CodeContextManifest) int {
	size := 0
	for _, f := range m.Files {
		size += len(f.Path) + len(f.Purpose) + len(f.Language)
	}
	for _, s := range m.Symbols {
		size += len(s.Name) + len(s.Kind) + len(s.Source)
	}
	for _, d := range m.Dependencies {
		size += len(d)
	}
	for _, e := range m.EnvironmentHints {
		size += len(e)
	}
	if size == 0 {
		return 0
	}
	return size / 4
}

func trimManifest(m domain.CodeContextManifest, maxTokens int) domain.CodeContextManifest {
	out := m
	for approximateTokens(out) > maxTokens {
		trimmed := false
		switch {
		case len(out.Symbols) > 0:
			out.Symbols = out.Symbols[:len(out.Symbols)-1]
			trimmed = true
		case len(out.Files) > 0:
			out.Files = out.Files[:len(out.Files)-1]
			trimmed = true
		case len(out.Dependencies) > 0:
			out.Dependencies = out.Dependencies[:len(out.Dependencies)-1]
			trimmed = true
		case len(out.EnvironmentHints) > 0:
			out.EnvironmentHints = out.EnvironmentHints[:len(out.EnvironmentHints)-1]
			trimmed = true
		default:
			return out
		}
		if !trimmed {
			return out
		}
	}
	out.ApproxTokens = approximateTokens(out)
	return out
}

func normalizeLanguage(path, fallback string) string {
	if strings.TrimSpace(fallback) != "" {
		return strings.ToLower(strings.TrimSpace(fallback))
	}
	switch filepath.Ext(path) {
	case ".go":
		return "golang"
	case ".ts", ".tsx":
		return "typescript"
	case ".js", ".jsx":
		return "javascript"
	default:
		return "text"
	}
}

func inferPurpose(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) == 0 {
		return "generated file"
	}
	return fmt.Sprintf("artifact for %s", parts[len(parts)-1])
}

var (
	goStructFnPattern = regexp.MustCompile(`(?m)^type\s+([A-Z]\w*)\s+struct|^func\s+([A-Z]\w*)\s*\(`)
	tsExportPattern   = regexp.MustCompile(`(?m)^export\s+(interface|type|const|function|class)\s+([A-Za-z_]\w*)`)
	dependencyPattern = regexp.MustCompile(`(?m)^(?:import|require)\s*[\("']\s*([^\)"']+)`)
	envPattern        = regexp.MustCompile(`(?i)([A-Z][A-Z0-9_]{2,})`)
)

func extractSymbols(path, content string) []domain.CodeSymbol {
	symbols := make([]domain.CodeSymbol, 0)
	ext := filepath.Ext(path)
	switch ext {
	case ".go":
		matches := goStructFnPattern.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			name := m[1]
			if name == "" {
				name = m[2]
			}
			if name == "" {
				continue
			}
			symbols = append(symbols, domain.CodeSymbol{Name: name, Kind: "export", Source: path, Backend: true})
		}
	case ".ts", ".tsx":
		matches := tsExportPattern.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			symbols = append(symbols, domain.CodeSymbol{Name: m[2], Kind: m[1], Source: path, Backend: false})
		}
	}
	return symbols
}

func extractDependencies(content string) []string {
	matches := dependencyPattern.FindAllStringSubmatch(content, -1)
	out := make([]string, 0, len(matches))
	for _, m := range matches {
		dep := strings.TrimSpace(m[1])
		if dep != "" {
			out = append(out, dep)
		}
	}
	return out
}

func extractEnvironmentHints(content string) []string {
	matches := envPattern.FindAllStringSubmatch(content, -1)
	seen := map[string]struct{}{}
	out := make([]string, 0)
	for _, m := range matches {
		k := m[1]
		if strings.HasPrefix(k, "TODO") || strings.HasPrefix(k, "HTTP") {
			continue
		}
		if _, ok := seen[k]; ok {
			continue
		}
		seen[k] = struct{}{}
		out = append(out, k)
	}
	return out
}
