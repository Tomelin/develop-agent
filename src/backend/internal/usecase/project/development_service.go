package project

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"

	domain "github.com/develop-agent/backend/internal/domain/project"
)

type DevelopmentService struct {
	projects    domain.ProjectRepository
	tasks       domain.TaskRepository
	files       domain.CodeFileRepository
	accumulator *CodeContextAccumulator
}

func NewDevelopmentService(projects domain.ProjectRepository, tasks domain.TaskRepository, files domain.CodeFileRepository) *DevelopmentService {
	return &DevelopmentService{projects: projects, tasks: tasks, files: files, accumulator: NewCodeContextAccumulator()}
}

func (s *DevelopmentService) SetExecutionMode(ctx context.Context, projectID, ownerID string, mode domain.ExecutionMode) (domain.ExecutionMode, error) {
	if !mode.IsValid() {
		return "", errors.New("invalid execution mode")
	}
	p, err := s.projects.FindByID(ctx, projectID)
	if err != nil {
		return "", err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return "", errors.New("project not found")
	}
	p.Phase5Mode = mode
	if err := s.projects.Update(ctx, p); err != nil {
		return "", err
	}
	return mode, nil
}

func (s *DevelopmentService) ExecuteTask(ctx context.Context, projectID, ownerID, taskID string) error {
	p, err := s.projects.FindByID(ctx, projectID)
	if err != nil {
		return err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return errors.New("project not found")
	}
	tasks, err := s.tasks.ListByProject(ctx, domain.TaskListFilter{ProjectID: projectID})
	if err != nil {
		return err
	}
	var target *domain.Task
	for _, t := range tasks {
		if t.ID.Hex() == taskID {
			target = t
			break
		}
	}
	if target == nil {
		return errors.New("task not found")
	}
	if target.Status == domain.TaskDone {
		return nil
	}
	if err := s.tasks.UpdateStatus(ctx, projectID, taskID, domain.TaskInProgress); err != nil {
		return err
	}

	generated := s.generateTaskArtifact(*target)
	if err := validateCompilationLike(generated.Path, generated.Content); err != nil {
		if err := s.tasks.UpdateStatus(ctx, projectID, taskID, domain.TaskBlocked); err != nil {
			return err
		}
		return fmt.Errorf("task blocked due compilation validation: %w", err)
	}
	projectOID, _ := bson.ObjectIDFromHex(projectID)
	if err := s.files.Upsert(ctx, &domain.CodeFile{
		ProjectID:   projectOID,
		Path:        generated.Path,
		Content:     generated.Content,
		TaskID:      target.ID.Hex(),
		Language:    generated.Language,
		PhaseNumber: 5,
	}); err != nil {
		return err
	}
	return s.tasks.UpdateStatus(ctx, projectID, taskID, domain.TaskDone)
}

func (s *DevelopmentService) ExecuteAllPending(ctx context.Context, projectID, ownerID string) (int, error) {
	items, err := s.tasks.ListByProject(ctx, domain.TaskListFilter{ProjectID: projectID, Status: domain.TaskTodo})
	if err != nil {
		return 0, err
	}
	executed := 0
	for _, task := range items {
		if err := s.ExecuteTask(ctx, projectID, ownerID, task.ID.Hex()); err != nil {
			return executed, err
		}
		executed++
	}
	return executed, nil
}

func (s *DevelopmentService) Summary(ctx context.Context, projectID, ownerID string) (*domain.Phase5Summary, error) {
	p, err := s.projects.FindByID(ctx, projectID)
	if err != nil {
		return nil, err
	}
	if p.OwnerUserID.Hex() != ownerID {
		return nil, errors.New("project not found")
	}
	tasks, err := s.tasks.ListByProject(ctx, domain.TaskListFilter{ProjectID: projectID})
	if err != nil {
		return nil, err
	}
	files, err := s.files.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	autoRejections, err := s.files.CountAutoRejections(ctx, projectID)
	if err != nil {
		return nil, err
	}

	summary := &domain.Phase5Summary{ExecutionMode: string(p.Phase5Mode), TotalPhaseTokens: p.TotalTokensUsed, AutoRejections: autoRejections}
	if summary.ExecutionMode == "" {
		summary.ExecutionMode = string(domain.ExecutionModeManual)
	}
	var lastTime time.Time
	for _, t := range tasks {
		summary.TotalTasks++
		switch t.Status {
		case domain.TaskDone:
			summary.DoneTasks++
		case domain.TaskInProgress:
			summary.InProgressTasks++
		case domain.TaskBlocked:
			summary.BlockedTasks++
		case domain.TaskTodo:
			summary.TodoTasks++
		}
		if t.UpdatedAt.After(lastTime) {
			lastTime = t.UpdatedAt
		}
	}
	if summary.TotalTasks > 0 {
		summary.CompletionPercent = float64(summary.DoneTasks) * 100.0 / float64(summary.TotalTasks)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Path, ".go") {
			summary.BackendFiles++
		} else {
			summary.FrontendFiles++
		}
		summary.GeneratedLinesOfCode += int64(strings.Count(f.Content, "\n") + 1)
		if f.UpdatedAt.After(lastTime) {
			lastTime = f.UpdatedAt
		}
	}
	if !lastTime.IsZero() {
		summary.LastExecutionUnixTime = lastTime.Unix()
	}
	if summary.DoneTasks > 0 {
		durationMin := time.Since(p.CreatedAt).Minutes()
		summary.AverageTaskMinutes = durationMin / float64(summary.DoneTasks)
	}
	return summary, nil
}

func (s *DevelopmentService) BuildCodeContext(ctx context.Context, projectID string) (domain.CodeContextManifest, error) {
	files, err := s.files.ListByProject(ctx, projectID)
	if err != nil {
		return domain.CodeContextManifest{}, err
	}
	return s.accumulator.Build(files), nil
}

func (s *DevelopmentService) ListFiles(ctx context.Context, projectID string) ([]*domain.CodeFile, error) {
	return s.files.ListByProject(ctx, projectID)
}

func (s *DevelopmentService) GetFile(ctx context.Context, projectID, fileID string) (*domain.CodeFile, error) {
	return s.files.FindByID(ctx, projectID, fileID)
}

func (s *DevelopmentService) TriggerAutoRejection(ctx context.Context, projectID, ownerID string, report domain.CatastrophicFailureReport) error {
	if !isCatastrophic(report) {
		return errors.New("report does not match catastrophic criteria")
	}
	if err := s.files.IncrementAutoRejections(ctx, projectID); err != nil {
		return err
	}
	noteTaskID := bson.NewObjectID().Hex()
	title := "Auto-fix triggered from quality gate"
	if len(report.FailureDescriptions) > 0 {
		title = fmt.Sprintf("Auto-fix: %s", report.FailureDescriptions[0])
	}
	projectOID, _ := bson.ObjectIDFromHex(projectID)
	return s.files.Upsert(ctx, &domain.CodeFile{
		ProjectID:   projectOID,
		Path:        "reports/auto_rejection_phase5.md",
		TaskID:      noteTaskID,
		Language:    "markdown",
		PhaseNumber: 5,
		Content:     fmt.Sprintf("# Auto-Rejection Trigger\n\nSource phase: %d\n\nFailures:\n- %s\n\nUser feedback counter preserved.\n", report.SourcePhase, strings.Join(report.FailureDescriptions, "\n- ")) + "\n" + title,
	})
}

func (s *DevelopmentService) DownloadZIP(ctx context.Context, projectID string) ([]byte, error) {
	files, err := s.files.ListByProject(ctx, projectID)
	if err != nil {
		return nil, err
	}
	buf := &bytes.Buffer{}
	zw := zip.NewWriter(buf)
	for _, f := range files {
		entry, err := zw.Create(strings.TrimPrefix(f.Path, "/"))
		if err != nil {
			return nil, err
		}
		if _, err := entry.Write([]byte(f.Content)); err != nil {
			return nil, err
		}
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func isCatastrophic(r domain.CatastrophicFailureReport) bool {
	if r.CoveragePercent < 50 {
		return true
	}
	if r.MaxCVSS >= 9.0 {
		return true
	}
	return r.CredentialsExposed || r.CompilationFailed
}

type generatedArtifact struct {
	Path     string
	Language string
	Content  string
}

func (s *DevelopmentService) generateTaskArtifact(task domain.Task) generatedArtifact {
	safeName := strings.ToLower(strings.ReplaceAll(strings.ReplaceAll(task.Title, " ", "_"), "-", "_"))
	if task.Type == domain.TaskTypeFrontend {
		path := filepath.ToSlash(fmt.Sprintf("src/frontend/components/%s.tsx", safeName))
		content := fmt.Sprintf("import React from 'react';\n\nexport function %sComponent(): JSX.Element {\n  return <section aria-label=\"%s\">%s</section>;\n}\n", toPascal(safeName), task.Title, task.Title)
		return generatedArtifact{Path: path, Language: "typescript", Content: content}
	}
	path := filepath.ToSlash(fmt.Sprintf("src/backend/internal/generated/%s.go", safeName))
	content := fmt.Sprintf("package generated\n\n// %sTask executa a task %s.\nfunc %sTask() string {\n\treturn %q\n}\n", toPascal(safeName), task.ID.Hex(), toPascal(safeName), task.Title)
	return generatedArtifact{Path: path, Language: "golang", Content: content}
}

func toPascal(in string) string {
	parts := strings.FieldsFunc(in, func(r rune) bool { return r == '_' || r == '-' || r == ' ' })
	for i := range parts {
		if parts[i] == "" {
			continue
		}
		parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
	}
	return strings.Join(parts, "")
}

func validateCompilationLike(path, content string) error {
	if strings.HasSuffix(path, ".go") {
		_, err := parser.ParseFile(token.NewFileSet(), path, content, parser.AllErrors)
		return err
	}
	if strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".tsx") {
		if !strings.Contains(content, "export ") && !strings.Contains(content, "import ") {
			return errors.New("typescript source must include import or export")
		}
	}
	return nil
}
