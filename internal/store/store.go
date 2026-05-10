package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/tae2089/agent-team/internal/output"
	_ "modernc.org/sqlite"
)

const DefaultStateDir = ".agent-team"

type Store struct {
	db *sql.DB
}

type Run struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type Task struct {
	ID             string          `json:"id"`
	RunID          string          `json:"run_id"`
	Agent          string          `json:"agent"`
	Title          string          `json:"title"`
	Body           string          `json:"body,omitempty"`
	Status         string          `json:"status"`
	Metadata       json.RawMessage `json:"metadata,omitempty"`
	Evidence       string          `json:"evidence,omitempty"`
	Artifact       string          `json:"artifact,omitempty"`
	BlockedReason  string          `json:"blocked_reason,omitempty"`
	StartedVersion int64           `json:"started_version,omitempty"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
}

type Message struct {
	ID        string          `json:"id"`
	RunID     string          `json:"run_id"`
	TaskID    string          `json:"task_id,omitempty"`
	From      string          `json:"from"`
	To        string          `json:"to"`
	Kind      string          `json:"kind"`
	Body      string          `json:"body"`
	Metadata  json.RawMessage `json:"metadata,omitempty"`
	AckedAt   string          `json:"acked_at,omitempty"`
	CreatedAt string          `json:"created_at"`
}

type Event struct {
	ID           int64           `json:"id"`
	StateVersion int64           `json:"state_version"`
	EventType    string          `json:"event_type"`
	EntityType   string          `json:"entity_type"`
	EntityID     string          `json:"entity_id"`
	RunID        string          `json:"run_id"`
	Payload      json.RawMessage `json:"payload"`
	CreatedAt    string          `json:"created_at"`
}

type SyncReport struct {
	Agent          string    `json:"agent"`
	RunID          string    `json:"run_id,omitempty"`
	TaskID         string    `json:"task_id,omitempty"`
	Blocking       bool      `json:"blocking"`
	UnreadMessages []Message `json:"unread_messages"`
	Issues         []string  `json:"issues"`
}

type PageOptions struct {
	Limit        int64
	AfterVersion int64
}

type RunSummary struct {
	Run             Run            `json:"run"`
	TaskCounts      map[string]int `json:"task_counts"`
	BlockedTasks    []Task         `json:"blocked_tasks"`
	InProgressTasks []Task         `json:"in_progress_tasks"`
	UnreadMessages  []Message      `json:"unread_messages"`
	RecentEvents    []Event        `json:"recent_events"`
	CloseReady      bool           `json:"close_ready"`
}

type StaleTask struct {
	Task            Task   `json:"task"`
	AgeSeconds      int64  `json:"age_seconds"`
	SuggestedAction string `json:"suggested_action"`
}

func StateDir() string {
	if value := os.Getenv("AGENT_TEAM_STATE_DIR"); value != "" {
		return value
	}
	return DefaultStateDir
}

func DBPath() string {
	return filepath.Join(StateDir(), "agent-team.db")
}

func Init(ctx context.Context) (*Store, int64, error) {
	if err := os.MkdirAll(StateDir(), 0o755); err != nil {
		return nil, 0, fmt.Errorf("create state dir: %w", err)
	}
	st, err := Open(ctx)
	if err != nil {
		return nil, 0, err
	}
	if err := st.migrate(ctx); err != nil {
		_ = st.Close()
		return nil, 0, err
	}
	version, err := st.StateVersion(ctx)
	return st, version, err
}

func Open(ctx context.Context) (*Store, error) {
	db, err := sql.Open("sqlite", DBPath())
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(1)
	if _, err := db.ExecContext(ctx, `PRAGMA foreign_keys = ON`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("configure sqlite: %w", err)
	}
	if _, err := db.ExecContext(ctx, `PRAGMA busy_timeout = 5000`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("configure sqlite: %w", err)
	}
	if _, err := db.ExecContext(ctx, `PRAGMA journal_mode = WAL`); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("configure sqlite: %w", err)
	}
	st := &Store{db: db}
	if err := st.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return st, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func (s *Store) migrate(ctx context.Context) error {
	statements := []string{
		`PRAGMA foreign_keys = ON`,
		`CREATE TABLE IF NOT EXISTS meta (key TEXT PRIMARY KEY, value TEXT NOT NULL)`,
		`INSERT OR IGNORE INTO meta(key, value) VALUES ('state_version', '0')`,
		`CREATE TABLE IF NOT EXISTS runs (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			status TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`,
		`CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			agent TEXT NOT NULL,
			title TEXT NOT NULL,
			body TEXT NOT NULL DEFAULT '',
			status TEXT NOT NULL,
			metadata TEXT NOT NULL DEFAULT '{}',
			evidence TEXT NOT NULL DEFAULT '',
			artifact TEXT NOT NULL DEFAULT '',
			blocked_reason TEXT NOT NULL DEFAULT '',
			started_version INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL,
			FOREIGN KEY(run_id) REFERENCES runs(id)
		)`,
		`CREATE TABLE IF NOT EXISTS task_dependencies (
			task_id TEXT NOT NULL,
			depends_on TEXT NOT NULL,
			PRIMARY KEY(task_id, depends_on),
			FOREIGN KEY(task_id) REFERENCES tasks(id),
			FOREIGN KEY(depends_on) REFERENCES tasks(id)
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			run_id TEXT NOT NULL,
			task_id TEXT NOT NULL DEFAULT '',
			from_agent TEXT NOT NULL,
			to_agent TEXT NOT NULL,
			kind TEXT NOT NULL,
			body TEXT NOT NULL,
			metadata TEXT NOT NULL DEFAULT '{}',
			acked_at TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL,
			FOREIGN KEY(run_id) REFERENCES runs(id)
		)`,
		`CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			state_version INTEGER NOT NULL,
			run_id TEXT NOT NULL DEFAULT '',
			event_type TEXT NOT NULL,
			entity_type TEXT NOT NULL,
			entity_id TEXT NOT NULL,
			payload TEXT NOT NULL,
			created_at TEXT NOT NULL
		)`,
	}
	for _, statement := range statements {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	if err := s.ensureEventsRunIDColumn(ctx); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	if err := s.backfillEventRunIDs(ctx); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	if err := s.ensureEventIndexes(ctx); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func (s *Store) StateVersion(ctx context.Context) (int64, error) {
	var raw string
	if err := s.db.QueryRowContext(ctx, `SELECT value FROM meta WHERE key = 'state_version'`).Scan(&raw); err != nil {
		return 0, err
	}
	var version int64
	_, err := fmt.Sscan(raw, &version)
	return version, err
}

func NewID(prefix string) string {
	return prefix + "_" + strings.ReplaceAll(uuid.NewString(), "-", "")
}

func (s *Store) CreateRun(ctx context.Context, id, title string) (Run, int64, error) {
	if title == "" {
		return Run{}, 0, output.NewError("validation_error", "title is required", nil)
	}
	if id == "" {
		id = NewID("run")
	}
	now := now()
	run := Run{ID: id, Title: title, Status: "open", CreatedAt: now, UpdatedAt: now}
	version, err := s.withMutation(ctx, "run_created", "run", id, run.ID, run, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO runs(id, title, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`, run.ID, run.Title, run.Status, run.CreatedAt, run.UpdatedAt)
		return err
	})
	return run, version, err
}

func (s *Store) RunStatus(ctx context.Context, id string) (map[string]any, int64, error) {
	if id == "" {
		return nil, 0, output.NewError("validation_error", "run_id is required", nil)
	}
	run, err := s.getRun(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	counts, err := s.taskStatusCounts(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	version, err := s.StateVersion(ctx)
	return map[string]any{"run": run, "tasks": counts}, version, err
}

func (s *Store) ListRuns(ctx context.Context, status string, page PageOptions) ([]Run, int64, error) {
	page = normalizePage(page)
	query := `SELECT id, title, status, created_at, updated_at FROM runs WHERE 1=1`
	args := []any{}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	if page.AfterVersion > 0 {
		query += ` AND id IN (SELECT entity_id FROM events WHERE entity_type = 'run' AND state_version > ?)`
		args = append(args, page.AfterVersion)
	}
	query += ` ORDER BY updated_at DESC, id LIMIT ?`
	args = append(args, page.Limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	runs, err := scanRuns(rows)
	if err != nil {
		return nil, 0, err
	}
	version, err := s.StateVersion(ctx)
	return runs, version, err
}

func (s *Store) CloseRun(ctx context.Context, runID, reason string) (Run, int64, []string, error) {
	if runID == "" {
		return Run{}, 0, nil, output.NewError("validation_error", "run_id is required", nil)
	}
	run, err := s.getRun(ctx, runID)
	if err != nil {
		return Run{}, 0, nil, err
	}
	if run.Status == "closed" {
		version, err := s.StateVersion(ctx)
		return run, version, []string{"run already closed"}, err
	}
	if run.Status != "open" {
		return Run{}, 0, nil, output.NewError("invalid_run_state", "run is not open", map[string]string{"run_id": runID, "status": run.Status})
	}
	counts, err := s.taskStatusCounts(ctx, runID)
	if err != nil {
		return Run{}, 0, nil, err
	}
	for status, count := range counts {
		if !isTerminalTaskStatus(status) && count > 0 {
			return Run{}, 0, nil, output.NewError("run_not_ready", "run has unfinished tasks", map[string]any{
				"run_id": runID,
				"tasks":  counts,
				"reason": reason,
			})
		}
	}
	run.Status = "closed"
	run.UpdatedAt = now()
	payload := map[string]any{
		"run":    run,
		"reason": reason,
	}
	version, err := s.withMutation(ctx, "run_closed", "run", run.ID, run.ID, payload, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE runs SET status = 'closed', updated_at = ? WHERE id = ?`, run.UpdatedAt, run.ID)
		return err
	})
	return run, version, nil, err
}

func (s *Store) CancelRun(ctx context.Context, runID, reason string) (Run, int64, error) {
	if runID == "" || reason == "" {
		return Run{}, 0, output.NewError("validation_error", "run_id and reason are required", nil)
	}
	run, err := s.getRun(ctx, runID)
	if err != nil {
		return Run{}, 0, err
	}
	if run.Status == "closed" {
		return Run{}, 0, output.NewError("invalid_run_state", "closed run cannot be cancelled", map[string]string{"run_id": runID, "status": run.Status})
	}
	if run.Status == "cancelled" {
		version, err := s.StateVersion(ctx)
		return run, version, err
	}
	run.Status = "cancelled"
	run.UpdatedAt = now()
	payload := map[string]any{"run": run, "reason": reason}
	version, err := s.withMutation(ctx, "run_cancelled", "run", run.ID, run.ID, payload, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE runs SET status = 'cancelled', updated_at = ? WHERE id = ?`, run.UpdatedAt, run.ID)
		return err
	})
	return run, version, err
}

func (s *Store) RunSummary(ctx context.Context, runID string, recentLimit int64) (RunSummary, int64, error) {
	if runID == "" {
		return RunSummary{}, 0, output.NewError("validation_error", "run_id is required", nil)
	}
	if recentLimit <= 0 {
		recentLimit = 10
	}
	if recentLimit > 100 {
		recentLimit = 100
	}
	run, err := s.getRun(ctx, runID)
	if err != nil {
		return RunSummary{}, 0, err
	}
	counts, err := s.taskStatusCounts(ctx, runID)
	if err != nil {
		return RunSummary{}, 0, err
	}
	blocked, _, err := s.ListTasks(ctx, runID, "", "blocked", PageOptions{Limit: 100})
	if err != nil {
		return RunSummary{}, 0, err
	}
	inProgress, _, err := s.ListTasks(ctx, runID, "", "in_progress", PageOptions{Limit: 100})
	if err != nil {
		return RunSummary{}, 0, err
	}
	unread, _, err := s.ListMessages(ctx, runID, "", "", "", "", true, PageOptions{Limit: 100})
	if err != nil {
		return RunSummary{}, 0, err
	}
	events, _, err := s.EventLog(ctx, runID, "", "", "", 0, recentLimit)
	if err != nil {
		return RunSummary{}, 0, err
	}
	closeReady := run.Status == "open"
	for status, count := range counts {
		if !isTerminalTaskStatus(status) && count > 0 {
			closeReady = false
			break
		}
	}
	version, err := s.StateVersion(ctx)
	return RunSummary{
		Run:             run,
		TaskCounts:      counts,
		BlockedTasks:    blocked,
		InProgressTasks: inProgress,
		UnreadMessages:  unread,
		RecentEvents:    events,
		CloseReady:      closeReady,
	}, version, err
}

func (s *Store) CreateTask(ctx context.Context, args CreateTaskArgs) (Task, int64, error) {
	if args.RunID == "" || args.Agent == "" || args.Title == "" {
		return Task{}, 0, output.NewError("validation_error", "run_id, agent, and title are required", nil)
	}
	if args.ID == "" {
		args.ID = NewID("task")
	}
	if len(args.Metadata) == 0 {
		args.Metadata = json.RawMessage(`{}`)
	}
	if err := validateJSON(args.Metadata); err != nil {
		return Task{}, 0, err
	}
	now := now()
	task := Task{ID: args.ID, RunID: args.RunID, Agent: args.Agent, Title: args.Title, Body: args.Body, Status: "pending", Metadata: args.Metadata, CreatedAt: now, UpdatedAt: now}
	version, err := s.withMutation(ctx, "task_created", "task", task.ID, task.RunID, task, func(tx *sql.Tx) error {
		if _, err := tx.ExecContext(ctx, `INSERT INTO tasks(id, run_id, agent, title, body, status, metadata, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, task.ID, task.RunID, task.Agent, task.Title, task.Body, task.Status, string(task.Metadata), task.CreatedAt, task.UpdatedAt); err != nil {
			return err
		}
		for _, dep := range args.DependsOn {
			if _, err := tx.ExecContext(ctx, `INSERT INTO task_dependencies(task_id, depends_on) VALUES (?, ?)`, task.ID, dep); err != nil {
				return err
			}
		}
		return nil
	})
	return task, version, err
}

type CreateTaskArgs struct {
	ID        string
	RunID     string
	Agent     string
	Title     string
	Body      string
	Metadata  json.RawMessage
	DependsOn []string
}

func (s *Store) ListTasks(ctx context.Context, runID, agent, status string, page PageOptions) ([]Task, int64, error) {
	page = normalizePage(page)
	query := `SELECT id, run_id, agent, title, body, status, metadata, evidence, artifact, blocked_reason, started_version, created_at, updated_at FROM tasks WHERE 1=1`
	args := []any{}
	if runID != "" {
		query += ` AND run_id = ?`
		args = append(args, runID)
	}
	if agent != "" {
		query += ` AND agent = ?`
		args = append(args, agent)
	}
	if status != "" {
		query += ` AND status = ?`
		args = append(args, status)
	}
	if page.AfterVersion > 0 {
		query += ` AND id IN (SELECT entity_id FROM events WHERE entity_type = 'task' AND state_version > ?)`
		args = append(args, page.AfterVersion)
	}
	query += ` ORDER BY created_at, id LIMIT ?`
	args = append(args, page.Limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	tasks, err := scanTasks(rows)
	if err != nil {
		return nil, 0, err
	}
	version, err := s.StateVersion(ctx)
	return tasks, version, err
}

func (s *Store) ShowTask(ctx context.Context, id string) (map[string]any, int64, error) {
	task, err := s.getTask(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	deps, err := s.dependencies(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	version, err := s.StateVersion(ctx)
	return map[string]any{"task": task, "depends_on": deps}, version, err
}

func (s *Store) StartTask(ctx context.Context, taskID, agent string) (Task, int64, error) {
	if taskID == "" || agent == "" {
		return Task{}, 0, output.NewError("validation_error", "task_id and agent are required", nil)
	}
	task, err := s.getTask(ctx, taskID)
	if err != nil {
		return Task{}, 0, err
	}
	if task.Agent != agent {
		return Task{}, 0, output.NewError("agent_mismatch", "task is assigned to a different agent", map[string]string{"assigned_agent": task.Agent})
	}
	if task.Status != "pending" {
		return Task{}, 0, output.NewError("invalid_task_state", "task can only be started when pending", map[string]string{"task_id": taskID, "status": task.Status})
	}
	current, err := s.StateVersion(ctx)
	if err != nil {
		return Task{}, 0, err
	}
	task.Status = "in_progress"
	task.StartedVersion = current
	task.UpdatedAt = now()
	version, err := s.withMutation(ctx, "task_started", "task", taskID, task.RunID, task, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE tasks SET status = 'in_progress', started_version = ?, updated_at = ? WHERE id = ?`, current, task.UpdatedAt, taskID)
		return err
	})
	return task, version, err
}

func (s *Store) CompleteTask(ctx context.Context, taskID, agent, evidence, artifact string, force bool) (Task, int64, []string, error) {
	if taskID == "" || agent == "" || evidence == "" || artifact == "" {
		return Task{}, 0, nil, output.NewError("validation_error", "task_id, agent, evidence, and artifact are required", nil)
	}
	task, err := s.getTask(ctx, taskID)
	if err != nil {
		return Task{}, 0, nil, err
	}
	if task.Agent != agent {
		return Task{}, 0, nil, output.NewError("agent_mismatch", "task is assigned to a different agent", map[string]string{"assigned_agent": task.Agent})
	}
	if isTerminalTaskStatus(task.Status) && task.Status != "done" {
		return Task{}, 0, nil, output.NewError("invalid_task_state", "task cannot be completed from terminal status", map[string]string{"task_id": taskID, "status": task.Status})
	}
	report, err := s.SyncCheck(ctx, agent, task.RunID, taskID)
	if err != nil {
		return Task{}, 0, nil, err
	}
	if report.Blocking && !force {
		version, _ := s.StateVersion(ctx)
		return Task{}, version, nil, output.NewError("sync_conflict", "blocking sync mismatch; rerun sync check or pass force:true", report)
	}
	task.Status = "done"
	task.Evidence = evidence
	task.Artifact = artifact
	task.UpdatedAt = now()
	version, err := s.withMutation(ctx, "task_completed", "task", taskID, task.RunID, task, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE tasks SET status = 'done', evidence = ?, artifact = ?, blocked_reason = '', updated_at = ? WHERE id = ?`, evidence, artifact, task.UpdatedAt, taskID)
		return err
	})
	warnings := []string{}
	if report.Blocking && force {
		warnings = append(warnings, "task completed with force:true despite sync mismatch")
	}
	return task, version, warnings, err
}

func (s *Store) BlockTask(ctx context.Context, taskID, agent, reason string) (Task, int64, error) {
	if taskID == "" || agent == "" || reason == "" {
		return Task{}, 0, output.NewError("validation_error", "task_id, agent, and reason are required", nil)
	}
	task, err := s.getTask(ctx, taskID)
	if err != nil {
		return Task{}, 0, err
	}
	if task.Agent != agent {
		return Task{}, 0, output.NewError("agent_mismatch", "task is assigned to a different agent", map[string]string{"assigned_agent": task.Agent})
	}
	if isTerminalTaskStatus(task.Status) {
		return Task{}, 0, output.NewError("invalid_task_state", "terminal task cannot be blocked", map[string]string{"task_id": taskID, "status": task.Status})
	}
	task.Status = "blocked"
	task.BlockedReason = reason
	task.UpdatedAt = now()
	version, err := s.withMutation(ctx, "task_blocked", "task", taskID, task.RunID, task, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE tasks SET status = 'blocked', blocked_reason = ?, updated_at = ? WHERE id = ?`, reason, task.UpdatedAt, taskID)
		return err
	})
	return task, version, err
}

func (s *Store) ReassignTask(ctx context.Context, taskID, agent, reason string) (Task, int64, error) {
	if taskID == "" || agent == "" || reason == "" {
		return Task{}, 0, output.NewError("validation_error", "task_id, agent, and reason are required", nil)
	}
	task, err := s.getTask(ctx, taskID)
	if err != nil {
		return Task{}, 0, err
	}
	if task.Status != "pending" && task.Status != "blocked" {
		return Task{}, 0, output.NewError("invalid_task_state", "task can only be reassigned when pending or blocked", map[string]string{
			"task_id": taskID,
			"status":  task.Status,
		})
	}
	previousStatus := task.Status
	previousAgent := task.Agent
	task.Agent = agent
	task.UpdatedAt = now()
	if task.Status == "blocked" {
		task.Status = "pending"
		task.BlockedReason = ""
		task.StartedVersion = 0
	}
	payload := map[string]any{
		"task":            task,
		"from_agent":      previousAgent,
		"to_agent":        agent,
		"reason":          reason,
		"previous_status": previousStatus,
	}
	version, err := s.withMutation(ctx, "task_reassigned", "task", taskID, task.RunID, payload, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE tasks SET agent = ?, status = ?, blocked_reason = '', started_version = ?, updated_at = ? WHERE id = ?`,
			task.Agent, task.Status, task.StartedVersion, task.UpdatedAt, taskID)
		return err
	})
	return task, version, err
}

func (s *Store) RetryTask(ctx context.Context, taskID, reason string) (Task, int64, error) {
	if taskID == "" || reason == "" {
		return Task{}, 0, output.NewError("validation_error", "task_id and reason are required", nil)
	}
	task, err := s.getTask(ctx, taskID)
	if err != nil {
		return Task{}, 0, err
	}
	if task.Status != "blocked" && task.Status != "in_progress" && task.Status != "failed" {
		return Task{}, 0, output.NewError("invalid_task_state", "task can only be retried when blocked, in_progress, or failed", map[string]string{
			"task_id": taskID,
			"status":  task.Status,
		})
	}
	previousStatus := task.Status
	previousAgent := task.Agent
	task.Status = "pending"
	task.Evidence = ""
	task.Artifact = ""
	task.BlockedReason = ""
	task.StartedVersion = 0
	task.UpdatedAt = now()
	payload := map[string]any{
		"task":            task,
		"reason":          reason,
		"previous_status": previousStatus,
		"previous_agent":  previousAgent,
	}
	version, err := s.withMutation(ctx, "task_retried", "task", taskID, task.RunID, payload, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE tasks SET status = 'pending', evidence = '', artifact = '', blocked_reason = '', started_version = 0, updated_at = ? WHERE id = ?`, task.UpdatedAt, taskID)
		return err
	})
	return task, version, err
}

func (s *Store) CancelTask(ctx context.Context, taskID, reason string) (Task, int64, error) {
	if taskID == "" || reason == "" {
		return Task{}, 0, output.NewError("validation_error", "task_id and reason are required", nil)
	}
	task, err := s.getTask(ctx, taskID)
	if err != nil {
		return Task{}, 0, err
	}
	if isTerminalTaskStatus(task.Status) {
		return Task{}, 0, output.NewError("invalid_task_state", "terminal task cannot be cancelled", map[string]string{"task_id": taskID, "status": task.Status})
	}
	task.Status = "cancelled"
	task.BlockedReason = reason
	task.UpdatedAt = now()
	payload := map[string]any{"task": task, "reason": reason}
	version, err := s.withMutation(ctx, "task_cancelled", "task", taskID, task.RunID, payload, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE tasks SET status = 'cancelled', blocked_reason = ?, updated_at = ? WHERE id = ?`, reason, task.UpdatedAt, taskID)
		return err
	})
	return task, version, err
}

func (s *Store) FailTask(ctx context.Context, taskID, agent, reason, artifact string) (Task, int64, error) {
	if taskID == "" || agent == "" || reason == "" {
		return Task{}, 0, output.NewError("validation_error", "task_id, agent, and reason are required", nil)
	}
	task, err := s.getTask(ctx, taskID)
	if err != nil {
		return Task{}, 0, err
	}
	if task.Agent != agent {
		return Task{}, 0, output.NewError("agent_mismatch", "task is assigned to a different agent", map[string]string{"assigned_agent": task.Agent})
	}
	if isTerminalTaskStatus(task.Status) {
		return Task{}, 0, output.NewError("invalid_task_state", "terminal task cannot be failed", map[string]string{"task_id": taskID, "status": task.Status})
	}
	task.Status = "failed"
	task.BlockedReason = reason
	task.Artifact = artifact
	task.UpdatedAt = now()
	payload := map[string]any{"task": task, "reason": reason}
	version, err := s.withMutation(ctx, "task_failed", "task", taskID, task.RunID, payload, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE tasks SET status = 'failed', blocked_reason = ?, artifact = ?, updated_at = ? WHERE id = ?`, reason, artifact, task.UpdatedAt, taskID)
		return err
	})
	return task, version, err
}

func (s *Store) SendMessage(ctx context.Context, msg Message) (Message, int64, error) {
	if msg.RunID == "" || msg.From == "" || msg.To == "" || msg.Kind == "" || msg.Body == "" {
		return Message{}, 0, output.NewError("validation_error", "run_id, from, to, kind, and body are required", nil)
	}
	if msg.ID == "" {
		msg.ID = NewID("msg")
	}
	if len(msg.Metadata) == 0 {
		msg.Metadata = json.RawMessage(`{}`)
	}
	if err := validateJSON(msg.Metadata); err != nil {
		return Message{}, 0, err
	}
	msg.CreatedAt = now()
	version, err := s.withMutation(ctx, "message_sent", "message", msg.ID, msg.RunID, msg, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `INSERT INTO messages(id, run_id, task_id, from_agent, to_agent, kind, body, metadata, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`, msg.ID, msg.RunID, msg.TaskID, msg.From, msg.To, msg.Kind, msg.Body, string(msg.Metadata), msg.CreatedAt)
		return err
	})
	return msg, version, err
}

func (s *Store) InboxList(ctx context.Context, agent, runID string, unread bool) ([]Message, int64, error) {
	if agent == "" {
		return nil, 0, output.NewError("validation_error", "agent is required", nil)
	}
	query := `SELECT id, run_id, task_id, from_agent, to_agent, kind, body, metadata, acked_at, created_at FROM messages WHERE to_agent = ?`
	args := []any{agent}
	if runID != "" {
		query += ` AND run_id = ?`
		args = append(args, runID)
	}
	if unread {
		query += ` AND acked_at = ''`
	}
	query += ` ORDER BY created_at, id`
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()
	messages, err := scanMessages(rows)
	if err != nil {
		return nil, 0, err
	}
	version, err := s.StateVersion(ctx)
	return messages, version, err
}

func (s *Store) InboxAck(ctx context.Context, msgID, agent string) (Message, int64, error) {
	if msgID == "" || agent == "" {
		return Message{}, 0, output.NewError("validation_error", "msg_id and agent are required", nil)
	}
	msg, err := s.getMessage(ctx, msgID)
	if err != nil {
		return Message{}, 0, err
	}
	if msg.To != agent {
		return Message{}, 0, output.NewError("agent_mismatch", "message belongs to a different recipient", map[string]string{"recipient": msg.To})
	}
	msg.AckedAt = now()
	version, err := s.withMutation(ctx, "message_acked", "message", msgID, msg.RunID, msg, func(tx *sql.Tx) error {
		_, err := tx.ExecContext(ctx, `UPDATE messages SET acked_at = ? WHERE id = ?`, msg.AckedAt, msgID)
		return err
	})
	return msg, version, err
}

func (s *Store) ListMessages(ctx context.Context, runID, taskID, from, to, kind string, unread bool, page PageOptions) ([]Message, int64, error) {
	if runID == "" {
		return nil, 0, output.NewError("validation_error", "run_id is required", nil)
	}
	page = normalizePage(page)
	query := `SELECT id, run_id, task_id, from_agent, to_agent, kind, body, metadata, acked_at, created_at FROM messages WHERE run_id = ?`
	args := []any{runID}
	if taskID != "" {
		query += ` AND task_id = ?`
		args = append(args, taskID)
	}
	if from != "" {
		query += ` AND from_agent = ?`
		args = append(args, from)
	}
	if to != "" {
		query += ` AND to_agent = ?`
		args = append(args, to)
	}
	if kind != "" {
		query += ` AND kind = ?`
		args = append(args, kind)
	}
	if unread {
		query += ` AND acked_at = ''`
	}
	if page.AfterVersion > 0 {
		query += ` AND id IN (SELECT entity_id FROM events WHERE entity_type = 'message' AND state_version > ?)`
		args = append(args, page.AfterVersion)
	}
	query += ` ORDER BY created_at DESC, id LIMIT ?`
	args = append(args, page.Limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	messages, err := scanMessages(rows)
	if err != nil {
		return nil, 0, err
	}
	version, err := s.StateVersion(ctx)
	return messages, version, err
}

func (s *Store) StaleTasks(ctx context.Context, runID string, olderThan time.Duration, page PageOptions) ([]StaleTask, int64, error) {
	if runID == "" {
		return nil, 0, output.NewError("validation_error", "run_id is required", nil)
	}
	if olderThan <= 0 {
		return nil, 0, output.NewError("validation_error", "older_than must be greater than 0", nil)
	}
	page = normalizePage(page)
	cutoff := time.Now().UTC().Add(-olderThan).Format(time.RFC3339Nano)
	query := `SELECT id, run_id, agent, title, body, status, metadata, evidence, artifact, blocked_reason, started_version, created_at, updated_at FROM tasks WHERE run_id = ? AND status IN ('in_progress', 'blocked') AND updated_at < ?`
	args := []any{runID, cutoff}
	if page.AfterVersion > 0 {
		query += ` AND id IN (SELECT entity_id FROM events WHERE entity_type = 'task' AND state_version > ?)`
		args = append(args, page.AfterVersion)
	}
	query += ` ORDER BY updated_at ASC, id LIMIT ?`
	args = append(args, page.Limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	tasks, err := scanTasks(rows)
	if err != nil {
		return nil, 0, err
	}
	stale := make([]StaleTask, 0, len(tasks))
	now := time.Now().UTC()
	for _, task := range tasks {
		updatedAt, err := time.Parse(time.RFC3339Nano, task.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		action := "sync_check_or_reassign"
		if task.Status == "blocked" {
			action = "retry_reassign_or_cancel"
		}
		stale = append(stale, StaleTask{
			Task:            task,
			AgeSeconds:      int64(now.Sub(updatedAt).Seconds()),
			SuggestedAction: action,
		})
	}
	version, err := s.StateVersion(ctx)
	return stale, version, err
}

func (s *Store) EventLog(ctx context.Context, runID, entityType, entityID, eventType string, afterVersion, limit int64) ([]Event, int64, error) {
	if limit <= 0 {
		return nil, 0, output.NewError("validation_error", "limit must be greater than 0", nil)
	}
	if limit > 1000 {
		return nil, 0, output.NewError("validation_error", "limit must be <= 1000", nil)
	}
	query := `SELECT id, state_version, event_type, entity_type, entity_id, run_id, payload, created_at FROM events WHERE 1=1`
	args := []any{}
	if runID != "" {
		query += ` AND run_id = ?`
		args = append(args, runID)
	}
	if entityType != "" {
		query += ` AND entity_type = ?`
		args = append(args, entityType)
	}
	if entityID != "" {
		query += ` AND entity_id = ?`
		args = append(args, entityID)
	}
	if eventType != "" {
		query += ` AND event_type = ?`
		args = append(args, eventType)
	}
	if afterVersion > 0 {
		query += ` AND state_version > ?`
		args = append(args, afterVersion)
	}
	query += ` ORDER BY state_version ASC, id ASC LIMIT ?`
	args = append(args, limit)
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	events, err := scanEvents(rows)
	if err != nil {
		return nil, 0, err
	}
	version, err := s.StateVersion(ctx)
	return events, version, err
}

func (s *Store) SyncCheck(ctx context.Context, agent, runID, taskID string) (SyncReport, error) {
	if agent == "" {
		return SyncReport{}, output.NewError("validation_error", "agent is required", nil)
	}
	report := SyncReport{Agent: agent, RunID: runID, TaskID: taskID, Issues: []string{}}
	query := `SELECT id, run_id, task_id, from_agent, to_agent, kind, body, metadata, acked_at, created_at FROM messages WHERE to_agent = ? AND acked_at = ''`
	args := []any{agent}
	if runID != "" {
		query += ` AND run_id = ?`
		args = append(args, runID)
	}
	if taskID != "" {
		query += ` AND (task_id = '' OR task_id = ?)`
		args = append(args, taskID)
	}
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return report, err
	}
	messages, err := scanMessages(rows)
	if err != nil {
		return report, err
	}
	report.UnreadMessages = messages
	if len(messages) > 0 {
		report.Blocking = true
		report.Issues = append(report.Issues, "unread relevant messages")
	}
	if taskID != "" {
		blockedDeps, err := s.incompleteDependencies(ctx, taskID)
		if err != nil {
			return report, err
		}
		if len(blockedDeps) > 0 {
			report.Blocking = true
			report.Issues = append(report.Issues, "incomplete dependencies: "+strings.Join(blockedDeps, ","))
		}
	}
	return report, nil
}

func (s *Store) ensureEventsRunIDColumn(ctx context.Context) error {
	hasColumn, err := s.hasColumn(ctx, "events", "run_id")
	if err != nil {
		return err
	}
	if hasColumn {
		return nil
	}
	_, err = s.db.ExecContext(ctx, `ALTER TABLE events ADD COLUMN run_id TEXT NOT NULL DEFAULT ''`)
	return err
}

func (s *Store) backfillEventRunIDs(ctx context.Context) error {
	migrations := []string{
		`UPDATE events SET run_id = entity_id WHERE run_id = '' AND entity_type = 'run'`,
		`UPDATE events SET run_id = (SELECT run_id FROM tasks WHERE tasks.id = events.entity_id) WHERE run_id = '' AND entity_type = 'task'`,
		`UPDATE events SET run_id = (SELECT run_id FROM messages WHERE messages.id = events.entity_id) WHERE run_id = '' AND entity_type = 'message'`,
	}
	for _, statement := range migrations {
		if _, err := s.db.ExecContext(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) ensureEventIndexes(ctx context.Context) error {
	indexes := []string{
		`CREATE INDEX IF NOT EXISTS idx_events_state_version ON events(state_version)`,
		`CREATE INDEX IF NOT EXISTS idx_events_run_state_version ON events(run_id, state_version)`,
		`CREATE INDEX IF NOT EXISTS idx_events_entity_version ON events(entity_type, entity_id, state_version)`,
	}
	for _, idx := range indexes {
		if _, err := s.db.ExecContext(ctx, idx); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) hasColumn(ctx context.Context, tableName, columnName string) (bool, error) {
	rows, err := s.db.QueryContext(ctx, `PRAGMA table_info(`+tableName+`)`)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			cid       int
			name      string
			colType   string
			notNull   int
			dfltValue sql.NullString
			pk        int
		)
		if err := rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk); err != nil {
			return false, err
		}
		if name == columnName {
			return true, nil
		}
	}
	return false, rows.Err()
}

func (s *Store) taskStatusCounts(ctx context.Context, runID string) (map[string]int, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT status, COUNT(*) FROM tasks WHERE run_id = ? GROUP BY status`, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	counts := map[string]int{}
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, rows.Err()
}

func (s *Store) withMutation(ctx context.Context, eventType, entityType, entityID, eventRunID string, payload any, fn func(*sql.Tx) error) (int64, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	defer func() {
		_ = tx.Rollback()
	}()
	if err := fn(tx); err != nil {
		return 0, err
	}
	version, err := nextVersion(ctx, tx)
	if err != nil {
		return 0, err
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return 0, err
	}
	if _, err := tx.ExecContext(ctx, `INSERT INTO events(state_version, run_id, event_type, entity_type, entity_id, payload, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`, version, eventRunID, eventType, entityType, entityID, string(raw), now()); err != nil {
		return 0, err
	}
	if err := tx.Commit(); err != nil {
		return 0, err
	}
	return version, nil
}

func nextVersion(ctx context.Context, tx *sql.Tx) (int64, error) {
	var raw string
	if err := tx.QueryRowContext(ctx, `SELECT value FROM meta WHERE key = 'state_version'`).Scan(&raw); err != nil {
		return 0, err
	}
	var version int64
	if _, err := fmt.Sscan(raw, &version); err != nil {
		return 0, err
	}
	version++
	if _, err := tx.ExecContext(ctx, `UPDATE meta SET value = ? WHERE key = 'state_version'`, fmt.Sprint(version)); err != nil {
		return 0, err
	}
	return version, nil
}

func (s *Store) getRun(ctx context.Context, id string) (Run, error) {
	var run Run
	err := s.db.QueryRowContext(ctx, `SELECT id, title, status, created_at, updated_at FROM runs WHERE id = ?`, id).Scan(&run.ID, &run.Title, &run.Status, &run.CreatedAt, &run.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Run{}, output.NewError("not_found", "run not found", map[string]string{"run_id": id})
	}
	return run, err
}

func (s *Store) getTask(ctx context.Context, id string) (Task, error) {
	var task Task
	var metadata string
	err := s.db.QueryRowContext(ctx, `SELECT id, run_id, agent, title, body, status, metadata, evidence, artifact, blocked_reason, started_version, created_at, updated_at FROM tasks WHERE id = ?`, id).Scan(&task.ID, &task.RunID, &task.Agent, &task.Title, &task.Body, &task.Status, &metadata, &task.Evidence, &task.Artifact, &task.BlockedReason, &task.StartedVersion, &task.CreatedAt, &task.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Task{}, output.NewError("not_found", "task not found", map[string]string{"task_id": id})
	}
	task.Metadata = json.RawMessage(metadata)
	return task, err
}

func (s *Store) getMessage(ctx context.Context, id string) (Message, error) {
	var msg Message
	var metadata string
	err := s.db.QueryRowContext(ctx, `SELECT id, run_id, task_id, from_agent, to_agent, kind, body, metadata, acked_at, created_at FROM messages WHERE id = ?`, id).Scan(&msg.ID, &msg.RunID, &msg.TaskID, &msg.From, &msg.To, &msg.Kind, &msg.Body, &metadata, &msg.AckedAt, &msg.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return Message{}, output.NewError("not_found", "message not found", map[string]string{"msg_id": id})
	}
	msg.Metadata = json.RawMessage(metadata)
	return msg, err
}

func (s *Store) dependencies(ctx context.Context, id string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT depends_on FROM task_dependencies WHERE task_id = ? ORDER BY depends_on`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var deps []string
	for rows.Next() {
		var dep string
		if err := rows.Scan(&dep); err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}
	return deps, rows.Err()
}

func (s *Store) incompleteDependencies(ctx context.Context, id string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `SELECT d.depends_on FROM task_dependencies d JOIN tasks t ON t.id = d.depends_on WHERE d.task_id = ? AND t.status NOT IN ('done', 'cancelled', 'failed') ORDER BY d.depends_on`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var deps []string
	for rows.Next() {
		var dep string
		if err := rows.Scan(&dep); err != nil {
			return nil, err
		}
		deps = append(deps, dep)
	}
	return deps, rows.Err()
}

func scanTasks(rows *sql.Rows) ([]Task, error) {
	tasks := []Task{}
	for rows.Next() {
		var task Task
		var metadata string
		if err := rows.Scan(&task.ID, &task.RunID, &task.Agent, &task.Title, &task.Body, &task.Status, &metadata, &task.Evidence, &task.Artifact, &task.BlockedReason, &task.StartedVersion, &task.CreatedAt, &task.UpdatedAt); err != nil {
			return nil, err
		}
		task.Metadata = json.RawMessage(metadata)
		tasks = append(tasks, task)
	}
	return tasks, rows.Err()
}

func scanRuns(rows *sql.Rows) ([]Run, error) {
	defer rows.Close()
	runs := []Run{}
	for rows.Next() {
		var run Run
		if err := rows.Scan(&run.ID, &run.Title, &run.Status, &run.CreatedAt, &run.UpdatedAt); err != nil {
			return nil, err
		}
		runs = append(runs, run)
	}
	return runs, rows.Err()
}

func scanMessages(rows *sql.Rows) ([]Message, error) {
	defer rows.Close()
	messages := []Message{}
	for rows.Next() {
		var msg Message
		var metadata string
		if err := rows.Scan(&msg.ID, &msg.RunID, &msg.TaskID, &msg.From, &msg.To, &msg.Kind, &msg.Body, &metadata, &msg.AckedAt, &msg.CreatedAt); err != nil {
			return nil, err
		}
		msg.Metadata = json.RawMessage(metadata)
		messages = append(messages, msg)
	}
	return messages, rows.Err()
}

func scanEvents(rows *sql.Rows) ([]Event, error) {
	defer rows.Close()
	events := []Event{}
	for rows.Next() {
		var event Event
		var payload string
		if err := rows.Scan(&event.ID, &event.StateVersion, &event.EventType, &event.EntityType, &event.EntityID, &event.RunID, &payload, &event.CreatedAt); err != nil {
			return nil, err
		}
		event.Payload = json.RawMessage(payload)
		events = append(events, event)
	}
	return events, rows.Err()
}

func validateJSON(raw json.RawMessage) error {
	var value any
	if err := json.Unmarshal(raw, &value); err != nil {
		return output.NewError("invalid_json", "payload must be valid JSON", map[string]string{"error": err.Error()})
	}
	return nil
}

func normalizePage(page PageOptions) PageOptions {
	if page.Limit <= 0 {
		page.Limit = 100
	}
	if page.Limit > 1000 {
		page.Limit = 1000
	}
	return page
}

func isTerminalTaskStatus(status string) bool {
	return status == "done" || status == "cancelled" || status == "failed"
}

func now() string {
	return time.Now().UTC().Format(time.RFC3339Nano)
}
