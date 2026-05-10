package cli

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/tae2089/agent-team/internal/output"
	"github.com/tae2089/agent-team/internal/store"
	"github.com/tae2089/agent-team/internal/version"
)

type rootOptions struct {
	params string
	body   string
}

func NewRoot() *cobra.Command {
	opts := &rootOptions{}
	cmd := &cobra.Command{
		Use:           "agent-team",
		Short:         "Daemonless agent team state coordinator",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.PersistentFlags().StringVar(&opts.params, "params", "{}", "strict JSON command arguments; supports inline JSON, @file, or -")
	cmd.PersistentFlags().StringVar(&opts.body, "json", "{}", "strict JSON payload; supports inline JSON, @file, or -")
	cmd.AddCommand(initCmd(opts))
	cmd.AddCommand(runCmd(opts))
	cmd.AddCommand(taskCmd(opts))
	cmd.AddCommand(messageCmd(opts))
	cmd.AddCommand(inboxCmd(opts))
	cmd.AddCommand(syncCmd(opts))
	cmd.AddCommand(eventCmd(opts))
	cmd.AddCommand(schemaCmd(opts))
	cmd.AddCommand(versionCmd())
	return cmd
}

func initCmd(opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use: "init",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return output.NewError("validation_error", "positional arguments are not supported; use --params JSON", nil)
			}
			if err := decodeStrict(opts.params, &struct{}{}); err != nil {
				return err
			}
			if err := decodeStrict(opts.body, &struct{}{}); err != nil {
				return err
			}
			st, version, err := store.Init(cmd.Context())
			if err != nil {
				return err
			}
			defer st.Close()
			return output.Write(cmd.OutOrStdout(), version, map[string]string{
				"state_dir": store.StateDir(),
				"db_path":   store.DBPath(),
			}, nil)
		},
	}
}

func runCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "run"}

	runCreateID := ""
	runCreateTitle := ""
	createCmd := &cobra.Command{
		Use: "create",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				ID    string `json:"id"`
				Title string `json:"title"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"id":    {},
					"title": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "id", "id", runCreateID, &params.ID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "title", "title", runCreateTitle, &params.Title, paramsData, bodyData); err != nil {
				return err
			}
			run, version, err := st.CreateRun(ctx, params.ID, params.Title)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"run": run}, nil)
		}),
	}
	createCmd.Flags().StringVar(&runCreateID, "id", "", "run identifier")
	createCmd.Flags().StringVar(&runCreateTitle, "title", "", "run title")
	cmd.AddCommand(createCmd)

	runStatusRun := ""
	statusCmd := &cobra.Command{
		Use: "status",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				RunID string `json:"run_id"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{"run_id": {}},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", runStatusRun, &params.RunID, paramsData, nil); err != nil {
				return err
			}
			data, version, err := st.RunStatus(ctx, params.RunID)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, data, nil)
		}),
	}
	statusCmd.Flags().StringVar(&runStatusRun, "run", "", "run identifier")
	cmd.AddCommand(statusCmd)

	runListStatus := ""
	runListLimit := int64(0)
	runListAfterVersion := int64(0)
	listCmd := &cobra.Command{
		Use: "list",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				Status       string `json:"status"`
				Limit        int64  `json:"limit"`
				AfterVersion int64  `json:"after_version"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"status":        {},
					"limit":         {},
					"after_version": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "status", "status", runListStatus, &params.Status, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "limit", "limit", runListLimit, &params.Limit, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "after-version", "after_version", runListAfterVersion, &params.AfterVersion, paramsData, nil); err != nil {
				return err
			}
			runs, version, err := st.ListRuns(ctx, params.Status, store.PageOptions{Limit: params.Limit, AfterVersion: params.AfterVersion})
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"runs": runs}, nil)
		}),
	}
	listCmd.Flags().StringVar(&runListStatus, "status", "", "status filter")
	listCmd.Flags().Int64Var(&runListLimit, "limit", 100, "page size")
	listCmd.Flags().Int64Var(&runListAfterVersion, "after-version", 0, "after state_version")
	cmd.AddCommand(listCmd)

	runSummaryRun := ""
	runSummaryRecentLimit := int64(0)
	summaryCmd := &cobra.Command{
		Use: "summary",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				RunID       string `json:"run_id"`
				RecentLimit int64  `json:"recent_limit"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"run_id":       {},
					"recent_limit": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", runSummaryRun, &params.RunID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "recent-limit", "recent_limit", runSummaryRecentLimit, &params.RecentLimit, paramsData, nil); err != nil {
				return err
			}
			summary, version, err := st.RunSummary(ctx, params.RunID, params.RecentLimit)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"summary": summary}, nil)
		}),
	}
	summaryCmd.Flags().StringVar(&runSummaryRun, "run", "", "run identifier")
	summaryCmd.Flags().Int64Var(&runSummaryRecentLimit, "recent-limit", 10, "recent event count")
	cmd.AddCommand(summaryCmd)

	runCloseRun := ""
	runCloseReason := ""
	closeCmd := &cobra.Command{
		Use: "close",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				RunID  string `json:"run_id"`
				Reason string `json:"reason"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"run_id": {},
					"reason": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", runCloseRun, &params.RunID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "reason", "reason", runCloseReason, &params.Reason, paramsData, bodyData); err != nil {
				return err
			}
			run, version, warnings, err := st.CloseRun(ctx, params.RunID, params.Reason)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"run": run}, warnings)
		}),
	}
	closeCmd.Flags().StringVar(&runCloseRun, "run", "", "run identifier")
	closeCmd.Flags().StringVar(&runCloseReason, "reason", "", "close reason")
	cmd.AddCommand(closeCmd)

	runCancelRun := ""
	runCancelReason := ""
	cancelCmd := &cobra.Command{
		Use: "cancel",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				RunID  string `json:"run_id"`
				Reason string `json:"reason"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"run_id": {},
					"reason": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", runCancelRun, &params.RunID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "reason", "reason", runCancelReason, &params.Reason, paramsData, bodyData); err != nil {
				return err
			}
			run, version, err := st.CancelRun(ctx, params.RunID, params.Reason)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"run": run}, nil)
		}),
	}
	cancelCmd.Flags().StringVar(&runCancelRun, "run", "", "run identifier")
	cancelCmd.Flags().StringVar(&runCancelReason, "reason", "", "cancel reason")
	cmd.AddCommand(cancelCmd)
	return cmd
}

func taskCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "task"}

	taskCreateID := ""
	taskCreateRun := ""
	taskCreateAgent := ""
	taskCreateTitle := ""
	taskCreateDepends := []string{}
	taskCreateBody := ""
	taskCreateMetadata := ""
	createCmd := &cobra.Command{
		Use: "create",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				ID        string   `json:"id"`
				RunID     string   `json:"run_id"`
				Agent     string   `json:"agent"`
				Title     string   `json:"title"`
				DependsOn []string `json:"depends_on"`
			}
			var body struct {
				Body     string          `json:"body"`
				Metadata json.RawMessage `json:"metadata"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"id":         {},
					"run_id":     {},
					"agent":      {},
					"title":      {},
					"depends_on": {},
				},
				map[string]struct{}{
					"body":     {},
					"metadata": {},
				},
				&params,
				&body,
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "id", "id", taskCreateID, &params.ID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", taskCreateRun, &params.RunID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", taskCreateAgent, &params.Agent, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "title", "title", taskCreateTitle, &params.Title, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringSliceField(cmd, "depends-on", "depends_on", taskCreateDepends, &params.DependsOn, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "body", "body", taskCreateBody, &body.Body, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeObjectField(cmd, "metadata", "metadata", taskCreateMetadata, &body.Metadata, paramsData, bodyData); err != nil {
				return err
			}
			if _, ok := bodyData["metadata"]; ok {
				if err := assertJSONObject(body.Metadata, "metadata"); err != nil {
					return err
				}
			}
			task, version, err := st.CreateTask(ctx, store.CreateTaskArgs{
				ID:        params.ID,
				RunID:     params.RunID,
				Agent:     params.Agent,
				Title:     params.Title,
				Body:      body.Body,
				Metadata:  body.Metadata,
				DependsOn: params.DependsOn,
			})
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"task": task}, nil)
		}),
	}
	createCmd.Flags().StringVar(&taskCreateID, "id", "", "task identifier")
	createCmd.Flags().StringVar(&taskCreateRun, "run", "", "run identifier")
	createCmd.Flags().StringVar(&taskCreateAgent, "agent", "", "agent")
	createCmd.Flags().StringVar(&taskCreateTitle, "title", "", "title")
	createCmd.Flags().StringSliceVar(&taskCreateDepends, "depends-on", nil, "dependent task ids")
	createCmd.Flags().StringVar(&taskCreateBody, "body", "", "task body")
	createCmd.Flags().StringVar(&taskCreateMetadata, "metadata", "", "task metadata JSON object")
	cmd.AddCommand(createCmd)

	taskListRun := ""
	taskListAgent := ""
	taskListStatus := ""
	taskListLimit := int64(0)
	taskListAfterVersion := int64(0)
	listCmd := &cobra.Command{
		Use: "list",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				RunID        string `json:"run_id"`
				Agent        string `json:"agent"`
				Status       string `json:"status"`
				Limit        int64  `json:"limit"`
				AfterVersion int64  `json:"after_version"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"run_id":        {},
					"agent":         {},
					"status":        {},
					"limit":         {},
					"after_version": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", taskListRun, &params.RunID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", taskListAgent, &params.Agent, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "status", "status", taskListStatus, &params.Status, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "limit", "limit", taskListLimit, &params.Limit, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "after-version", "after_version", taskListAfterVersion, &params.AfterVersion, paramsData, nil); err != nil {
				return err
			}
			tasks, version, err := st.ListTasks(ctx, params.RunID, params.Agent, params.Status, store.PageOptions{Limit: params.Limit, AfterVersion: params.AfterVersion})
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"tasks": tasks}, nil)
		}),
	}
	listCmd.Flags().StringVar(&taskListRun, "run", "", "run identifier")
	listCmd.Flags().StringVar(&taskListAgent, "agent", "", "agent")
	listCmd.Flags().StringVar(&taskListStatus, "status", "", "status filter")
	listCmd.Flags().Int64Var(&taskListLimit, "limit", 100, "page size")
	listCmd.Flags().Int64Var(&taskListAfterVersion, "after-version", 0, "after state_version")
	cmd.AddCommand(listCmd)

	taskShowTask := ""
	showCmd := &cobra.Command{
		Use: "show",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				TaskID string `json:"task_id"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{"task_id": {}},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", taskShowTask, &params.TaskID, paramsData, nil); err != nil {
				return err
			}
			data, version, err := st.ShowTask(ctx, params.TaskID)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, data, nil)
		}),
	}
	showCmd.Flags().StringVar(&taskShowTask, "task", "", "task identifier")
	cmd.AddCommand(showCmd)

	taskStartTask := ""
	taskStartAgent := ""
	startCmd := &cobra.Command{
		Use: "start",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				TaskID string `json:"task_id"`
				Agent  string `json:"agent"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"task_id": {},
					"agent":   {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", taskStartTask, &params.TaskID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", taskStartAgent, &params.Agent, paramsData, nil); err != nil {
				return err
			}
			task, version, err := st.StartTask(ctx, params.TaskID, params.Agent)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"task": task}, nil)
		}),
	}
	startCmd.Flags().StringVar(&taskStartTask, "task", "", "task identifier")
	startCmd.Flags().StringVar(&taskStartAgent, "agent", "", "agent")
	cmd.AddCommand(startCmd)

	taskCompleteTask := ""
	taskCompleteAgent := ""
	taskCompleteForce := false
	taskCompleteEvidence := ""
	taskCompleteArtifact := ""
	completeCmd := &cobra.Command{
		Use: "complete",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				TaskID string `json:"task_id"`
				Agent  string `json:"agent"`
				Force  bool   `json:"force"`
			}
			var body struct {
				Evidence string `json:"evidence"`
				Artifact string `json:"artifact"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"task_id": {},
					"agent":   {},
					"force":   {},
				},
				map[string]struct{}{
					"evidence": {},
					"artifact": {},
				},
				&params,
				&body,
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", taskCompleteTask, &params.TaskID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", taskCompleteAgent, &params.Agent, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeBoolField(cmd, "force", "force", taskCompleteForce, &params.Force, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "evidence", "evidence", taskCompleteEvidence, &body.Evidence, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "artifact", "artifact", taskCompleteArtifact, &body.Artifact, paramsData, bodyData); err != nil {
				return err
			}
			task, version, warnings, err := st.CompleteTask(ctx, params.TaskID, params.Agent, body.Evidence, body.Artifact, params.Force)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"task": task}, warnings)
		}),
	}
	completeCmd.Flags().StringVar(&taskCompleteTask, "task", "", "task identifier")
	completeCmd.Flags().StringVar(&taskCompleteAgent, "agent", "", "agent")
	completeCmd.Flags().BoolVar(&taskCompleteForce, "force", false, "force completion")
	completeCmd.Flags().StringVar(&taskCompleteEvidence, "evidence", "", "evidence")
	completeCmd.Flags().StringVar(&taskCompleteArtifact, "artifact", "", "artifact path")
	cmd.AddCommand(completeCmd)

	taskBlockTask := ""
	taskBlockAgent := ""
	taskBlockReason := ""
	blockCmd := &cobra.Command{
		Use: "block",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				TaskID string `json:"task_id"`
				Agent  string `json:"agent"`
			}
			var body struct {
				Reason string `json:"reason"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"task_id": {},
					"agent":   {},
				},
				map[string]struct{}{
					"reason": {},
				},
				&params,
				&body,
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", taskBlockTask, &params.TaskID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", taskBlockAgent, &params.Agent, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "reason", "reason", taskBlockReason, &body.Reason, paramsData, bodyData); err != nil {
				return err
			}
			task, version, err := st.BlockTask(ctx, params.TaskID, params.Agent, body.Reason)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"task": task}, nil)
		}),
	}
	blockCmd.Flags().StringVar(&taskBlockTask, "task", "", "task identifier")
	blockCmd.Flags().StringVar(&taskBlockAgent, "agent", "", "agent")
	blockCmd.Flags().StringVar(&taskBlockReason, "reason", "", "block reason")
	cmd.AddCommand(blockCmd)

	taskReassignTask := ""
	taskReassignAgent := ""
	taskReassignReason := ""
	reassignCmd := &cobra.Command{
		Use: "reassign",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				TaskID string `json:"task_id"`
				Agent  string `json:"agent"`
				Reason string `json:"reason"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"task_id": {},
					"agent":   {},
					"reason":  {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", taskReassignTask, &params.TaskID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", taskReassignAgent, &params.Agent, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "reason", "reason", taskReassignReason, &params.Reason, paramsData, bodyData); err != nil {
				return err
			}
			task, version, err := st.ReassignTask(ctx, params.TaskID, params.Agent, params.Reason)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"task": task}, nil)
		}),
	}
	reassignCmd.Flags().StringVar(&taskReassignTask, "task", "", "task identifier")
	reassignCmd.Flags().StringVar(&taskReassignAgent, "agent", "", "agent")
	reassignCmd.Flags().StringVar(&taskReassignReason, "reason", "", "reason")
	cmd.AddCommand(reassignCmd)

	taskRetryTask := ""
	taskRetryReason := ""
	retryCmd := &cobra.Command{
		Use: "retry",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				TaskID string `json:"task_id"`
				Reason string `json:"reason"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"task_id": {},
					"reason":  {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", taskRetryTask, &params.TaskID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "reason", "reason", taskRetryReason, &params.Reason, paramsData, bodyData); err != nil {
				return err
			}
			task, version, err := st.RetryTask(ctx, params.TaskID, params.Reason)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"task": task}, nil)
		}),
	}
	retryCmd.Flags().StringVar(&taskRetryTask, "task", "", "task identifier")
	retryCmd.Flags().StringVar(&taskRetryReason, "reason", "", "reason")
	cmd.AddCommand(retryCmd)

	taskCancelTask := ""
	taskCancelReason := ""
	cancelCmd := &cobra.Command{
		Use: "cancel",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				TaskID string `json:"task_id"`
				Reason string `json:"reason"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"task_id": {},
					"reason":  {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", taskCancelTask, &params.TaskID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "reason", "reason", taskCancelReason, &params.Reason, paramsData, bodyData); err != nil {
				return err
			}
			task, version, err := st.CancelTask(ctx, params.TaskID, params.Reason)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"task": task}, nil)
		}),
	}
	cancelCmd.Flags().StringVar(&taskCancelTask, "task", "", "task identifier")
	cancelCmd.Flags().StringVar(&taskCancelReason, "reason", "", "reason")
	cmd.AddCommand(cancelCmd)

	taskFailTask := ""
	taskFailAgent := ""
	taskFailReason := ""
	taskFailArtifact := ""
	failCmd := &cobra.Command{
		Use: "fail",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				TaskID   string `json:"task_id"`
				Agent    string `json:"agent"`
				Reason   string `json:"reason"`
				Artifact string `json:"artifact"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"task_id":  {},
					"agent":    {},
					"reason":   {},
					"artifact": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", taskFailTask, &params.TaskID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", taskFailAgent, &params.Agent, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "reason", "reason", taskFailReason, &params.Reason, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "artifact", "artifact", taskFailArtifact, &params.Artifact, paramsData, bodyData); err != nil {
				return err
			}
			task, version, err := st.FailTask(ctx, params.TaskID, params.Agent, params.Reason, params.Artifact)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"task": task}, nil)
		}),
	}
	failCmd.Flags().StringVar(&taskFailTask, "task", "", "task identifier")
	failCmd.Flags().StringVar(&taskFailAgent, "agent", "", "agent")
	failCmd.Flags().StringVar(&taskFailReason, "reason", "", "reason")
	failCmd.Flags().StringVar(&taskFailArtifact, "artifact", "", "artifact path")
	cmd.AddCommand(failCmd)

	taskStaleRun := ""
	taskStaleOlderThan := ""
	taskStaleLimit := int64(0)
	taskStaleAfterVersion := int64(0)
	staleCmd := &cobra.Command{
		Use: "stale",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				RunID        string `json:"run_id"`
				OlderThan    string `json:"older_than"`
				Limit        int64  `json:"limit"`
				AfterVersion int64  `json:"after_version"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"run_id":        {},
					"older_than":    {},
					"limit":         {},
					"after_version": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", taskStaleRun, &params.RunID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "older-than", "older_than", taskStaleOlderThan, &params.OlderThan, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "limit", "limit", taskStaleLimit, &params.Limit, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "after-version", "after_version", taskStaleAfterVersion, &params.AfterVersion, paramsData, nil); err != nil {
				return err
			}
			olderThan, err := parseDuration(params.OlderThan)
			if err != nil {
				return err
			}
			stale, version, err := st.StaleTasks(ctx, params.RunID, olderThan, store.PageOptions{Limit: params.Limit, AfterVersion: params.AfterVersion})
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"stale_tasks": stale}, nil)
		}),
	}
	staleCmd.Flags().StringVar(&taskStaleRun, "run", "", "run identifier")
	staleCmd.Flags().StringVar(&taskStaleOlderThan, "older-than", "", "duration such as 2h")
	staleCmd.Flags().Int64Var(&taskStaleLimit, "limit", 100, "page size")
	staleCmd.Flags().Int64Var(&taskStaleAfterVersion, "after-version", 0, "after state_version")
	cmd.AddCommand(staleCmd)
	return cmd
}

func messageCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "message"}
	messageSendID := ""
	messageSendRun := ""
	messageSendTask := ""
	messageSendFrom := ""
	messageSendTo := ""
	messageSendKind := ""
	messageSendBody := ""
	messageSendMetadata := ""
	sendCmd := &cobra.Command{
		Use: "send",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				ID     string `json:"id"`
				RunID  string `json:"run_id"`
				TaskID string `json:"task_id"`
				From   string `json:"from"`
				To     string `json:"to"`
				Kind   string `json:"kind"`
			}
			var body struct {
				Body     string          `json:"body"`
				Metadata json.RawMessage `json:"metadata"`
			}
			paramsData, bodyData, err := readCommandInputs(opts,
				map[string]struct{}{
					"id":      {},
					"run_id":  {},
					"task_id": {},
					"from":    {},
					"to":      {},
					"kind":    {},
				},
				map[string]struct{}{
					"body":     {},
					"metadata": {},
				},
				&params,
				&body,
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "id", "id", messageSendID, &params.ID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", messageSendRun, &params.RunID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", messageSendTask, &params.TaskID, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "from", "from", messageSendFrom, &params.From, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "to", "to", messageSendTo, &params.To, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "kind", "kind", messageSendKind, &params.Kind, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "body", "body", messageSendBody, &body.Body, paramsData, bodyData); err != nil {
				return err
			}
			if err := mergeObjectField(cmd, "metadata", "metadata", messageSendMetadata, &body.Metadata, paramsData, bodyData); err != nil {
				return err
			}
			if _, ok := bodyData["metadata"]; ok {
				if err := assertJSONObject(body.Metadata, "metadata"); err != nil {
					return err
				}
			}
			msg, version, err := st.SendMessage(ctx, store.Message{
				ID:       params.ID,
				RunID:    params.RunID,
				TaskID:   params.TaskID,
				From:     params.From,
				To:       params.To,
				Kind:     params.Kind,
				Body:     body.Body,
				Metadata: body.Metadata,
			})
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"message": msg}, nil)
		}),
	}
	sendCmd.Flags().StringVar(&messageSendID, "id", "", "message id")
	sendCmd.Flags().StringVar(&messageSendRun, "run", "", "run identifier")
	sendCmd.Flags().StringVar(&messageSendTask, "task", "", "task identifier")
	sendCmd.Flags().StringVar(&messageSendFrom, "from", "", "sender")
	sendCmd.Flags().StringVar(&messageSendTo, "to", "", "recipient")
	sendCmd.Flags().StringVar(&messageSendKind, "kind", "", "kind")
	sendCmd.Flags().StringVar(&messageSendBody, "body", "", "message body")
	sendCmd.Flags().StringVar(&messageSendMetadata, "metadata", "", "metadata JSON object")
	cmd.AddCommand(sendCmd)

	messageListRun := ""
	messageListTask := ""
	messageListFrom := ""
	messageListTo := ""
	messageListKind := ""
	messageListUnread := false
	messageListLimit := int64(0)
	messageListAfterVersion := int64(0)
	messageListCmd := &cobra.Command{
		Use: "list",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				RunID        string `json:"run_id"`
				TaskID       string `json:"task_id"`
				From         string `json:"from"`
				To           string `json:"to"`
				Kind         string `json:"kind"`
				Unread       bool   `json:"unread"`
				Limit        int64  `json:"limit"`
				AfterVersion int64  `json:"after_version"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"run_id":        {},
					"task_id":       {},
					"from":          {},
					"to":            {},
					"kind":          {},
					"unread":        {},
					"limit":         {},
					"after_version": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", messageListRun, &params.RunID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", messageListTask, &params.TaskID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "from", "from", messageListFrom, &params.From, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "to", "to", messageListTo, &params.To, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "kind", "kind", messageListKind, &params.Kind, paramsData, nil); err != nil {
				return err
			}
			if err := mergeBoolField(cmd, "unread", "unread", messageListUnread, &params.Unread, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "limit", "limit", messageListLimit, &params.Limit, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "after-version", "after_version", messageListAfterVersion, &params.AfterVersion, paramsData, nil); err != nil {
				return err
			}
			messages, version, err := st.ListMessages(ctx, params.RunID, params.TaskID, params.From, params.To, params.Kind, params.Unread, store.PageOptions{Limit: params.Limit, AfterVersion: params.AfterVersion})
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"messages": messages}, nil)
		}),
	}
	messageListCmd.Flags().StringVar(&messageListRun, "run", "", "run identifier")
	messageListCmd.Flags().StringVar(&messageListTask, "task", "", "task identifier")
	messageListCmd.Flags().StringVar(&messageListFrom, "from", "", "sender")
	messageListCmd.Flags().StringVar(&messageListTo, "to", "", "recipient")
	messageListCmd.Flags().StringVar(&messageListKind, "kind", "", "message kind")
	messageListCmd.Flags().BoolVar(&messageListUnread, "unread", false, "unread only")
	messageListCmd.Flags().Int64Var(&messageListLimit, "limit", 100, "page size")
	messageListCmd.Flags().Int64Var(&messageListAfterVersion, "after-version", 0, "after state_version")
	cmd.AddCommand(messageListCmd)
	return cmd
}

func inboxCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "inbox"}
	inboxListAgent := ""
	inboxListRun := ""
	inboxListUnread := false
	listCmd := &cobra.Command{
		Use: "list",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				Agent  string `json:"agent"`
				RunID  string `json:"run_id"`
				Unread bool   `json:"unread"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"agent":  {},
					"run_id": {},
					"unread": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", inboxListAgent, &params.Agent, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", inboxListRun, &params.RunID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeBoolField(cmd, "unread", "unread", inboxListUnread, &params.Unread, paramsData, nil); err != nil {
				return err
			}
			messages, version, err := st.InboxList(ctx, params.Agent, params.RunID, params.Unread)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"messages": messages}, nil)
		}),
	}
	listCmd.Flags().StringVar(&inboxListAgent, "agent", "", "agent")
	listCmd.Flags().StringVar(&inboxListRun, "run", "", "run identifier")
	listCmd.Flags().BoolVar(&inboxListUnread, "unread", false, "unread only")
	cmd.AddCommand(listCmd)

	inboxAckMsg := ""
	inboxAckAgent := ""
	ackCmd := &cobra.Command{
		Use: "ack",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				MsgID string `json:"msg_id"`
				Agent string `json:"agent"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"msg_id": {},
					"agent":  {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "msg", "msg_id", inboxAckMsg, &params.MsgID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", inboxAckAgent, &params.Agent, paramsData, nil); err != nil {
				return err
			}
			msg, version, err := st.InboxAck(ctx, params.MsgID, params.Agent)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"message": msg}, nil)
		}),
	}
	ackCmd.Flags().StringVar(&inboxAckMsg, "msg", "", "message identifier")
	ackCmd.Flags().StringVar(&inboxAckAgent, "agent", "", "agent")
	cmd.AddCommand(ackCmd)
	return cmd
}

func syncCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "sync"}
	syncCheckAgent := ""
	syncCheckRun := ""
	syncCheckTask := ""
	checkCmd := &cobra.Command{
		Use: "check",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				Agent  string `json:"agent"`
				RunID  string `json:"run_id"`
				TaskID string `json:"task_id"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"agent":   {},
					"run_id":  {},
					"task_id": {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "agent", "agent", syncCheckAgent, &params.Agent, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", syncCheckRun, &params.RunID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "task", "task_id", syncCheckTask, &params.TaskID, paramsData, nil); err != nil {
				return err
			}
			report, err := st.SyncCheck(ctx, params.Agent, params.RunID, params.TaskID)
			if err != nil {
				return err
			}
			version, err := st.StateVersion(ctx)
			if err != nil {
				return err
			}
			warnings := []string{}
			if report.Blocking {
				warnings = append(warnings, "blocking sync mismatch")
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"sync": report}, warnings)
		}),
	}
	checkCmd.Flags().StringVar(&syncCheckAgent, "agent", "", "agent")
	checkCmd.Flags().StringVar(&syncCheckRun, "run", "", "run identifier")
	checkCmd.Flags().StringVar(&syncCheckTask, "task", "", "task identifier")
	cmd.AddCommand(checkCmd)
	return cmd
}

func eventCmd(opts *rootOptions) *cobra.Command {
	cmd := &cobra.Command{Use: "event"}
	eventLogRun := ""
	eventLogEntityType := ""
	eventLogEntityID := ""
	eventLogEventType := ""
	eventLogAfterVersion := int64(0)
	eventLogLimit := int64(0)
	logCmd := &cobra.Command{
		Use: "log",
		RunE: withStore(opts, func(ctx context.Context, st *store.Store, cmd *cobra.Command) error {
			var params struct {
				RunID        string `json:"run_id"`
				EntityType   string `json:"entity_type"`
				EntityID     string `json:"entity_id"`
				EventType    string `json:"event_type"`
				AfterVersion int64  `json:"after_version"`
				Limit        int64  `json:"limit"`
			}
			paramsData, _, err := readCommandInputs(opts,
				map[string]struct{}{
					"run_id":        {},
					"entity_type":   {},
					"entity_id":     {},
					"event_type":    {},
					"after_version": {},
					"limit":         {},
				},
				nil,
				&params,
				&struct{}{},
			)
			if err != nil {
				return err
			}
			if err := mergeStringField(cmd, "run", "run_id", eventLogRun, &params.RunID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "entity-type", "entity_type", eventLogEntityType, &params.EntityType, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "entity", "entity_id", eventLogEntityID, &params.EntityID, paramsData, nil); err != nil {
				return err
			}
			if err := mergeStringField(cmd, "type", "event_type", eventLogEventType, &params.EventType, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "after-version", "after_version", eventLogAfterVersion, &params.AfterVersion, paramsData, nil); err != nil {
				return err
			}
			if err := mergeInt64Field(cmd, "limit", "limit", eventLogLimit, &params.Limit, paramsData, nil); err != nil {
				return err
			}
			limit := params.Limit
			if limit == 0 {
				limit = 100
			}
			if limit > 1000 {
				return output.NewError("validation_error", "limit must be 1000 or less", nil)
			}
			events, version, err := st.EventLog(ctx, params.RunID, params.EntityType, params.EntityID, params.EventType, params.AfterVersion, limit)
			if err != nil {
				return err
			}
			return output.Write(cmd.OutOrStdout(), version, map[string]any{"events": events}, nil)
		}),
	}
	logCmd.Flags().StringVar(&eventLogRun, "run", "", "run identifier")
	logCmd.Flags().StringVar(&eventLogEntityType, "entity-type", "", "entity type")
	logCmd.Flags().StringVar(&eventLogEntityID, "entity", "", "entity identifier")
	logCmd.Flags().StringVar(&eventLogEventType, "type", "", "event type")
	logCmd.Flags().Int64Var(&eventLogAfterVersion, "after-version", 0, "after state_version")
	logCmd.Flags().Int64Var(&eventLogLimit, "limit", 100, "page size")
	cmd.AddCommand(logCmd)
	return cmd
}

func schemaCmd(opts *rootOptions) *cobra.Command {
	_ = opts
	cmd := &cobra.Command{Use: "schema"}
	exportCmd := &cobra.Command{
		Use: "export",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return output.NewError("validation_error", "positional arguments are not supported; use --params JSON", nil)
			}
			return output.Write(cmd.OutOrStdout(), 0, map[string]any{"schema": cliSchema()}, nil)
		},
	}
	cmd.AddCommand(exportCmd)
	return cmd
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use: "version",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				return output.NewError("validation_error", "positional arguments are not supported", nil)
			}
			return output.Write(cmd.OutOrStdout(), 0, map[string]any{"version": version.Version}, nil)
		},
	}
}

func withStore(opts *rootOptions, fn func(context.Context, *store.Store, *cobra.Command) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if len(args) > 0 {
			return output.NewError("validation_error", "positional arguments are not supported; use --params JSON", nil)
		}
		st, err := store.Open(cmd.Context())
		if err != nil {
			return err
		}
		defer st.Close()
		return fn(cmd.Context(), st, cmd)
	}
}

func readCommandInputs(opts *rootOptions, paramsAllowed, bodyAllowed map[string]struct{}, params any, body any) (map[string]json.RawMessage, map[string]json.RawMessage, error) {
	paramsData, err := decodeStrictObject(opts.params, paramsAllowed, params)
	if err != nil {
		return nil, nil, err
	}
	bodyData, err := decodeStrictObject(opts.body, bodyAllowed, body)
	if err != nil {
		return nil, nil, err
	}
	return paramsData, bodyData, nil
}

func decodeStrict(raw string, target any) error {
	content, err := readJSONArg(raw)
	if err != nil {
		return err
	}
	decoder := json.NewDecoder(bytes.NewReader(content))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(target); err != nil {
		return output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": err.Error()})
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return output.NewError("invalid_json", "JSON input must contain a single object", nil)
	}
	return nil
}

func decodeStrictObject(raw string, allowed map[string]struct{}, target any) (map[string]json.RawMessage, error) {
	content, err := readJSONArg(raw)
	if err != nil {
		return nil, err
	}
	payload := map[string]json.RawMessage{}
	decoder := json.NewDecoder(bytes.NewReader(content))
	if err := decoder.Decode(&payload); err != nil {
		return nil, output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": err.Error()})
	}
	if decoder.Decode(&struct{}{}) != io.EOF {
		return nil, output.NewError("invalid_json", "JSON input must contain a single object", nil)
	}
	for key := range payload {
		if _, ok := allowed[key]; !ok {
			return nil, output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": fmt.Sprintf("unknown field: %s", key)})
		}
	}
	if target != nil {
		normalized, err := json.Marshal(payload)
		if err != nil {
			return nil, output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": err.Error()})
		}
		decoder = json.NewDecoder(bytes.NewReader(normalized))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(target); err != nil {
			return nil, output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": err.Error()})
		}
		if decoder.Decode(&struct{}{}) != io.EOF {
			return nil, output.NewError("invalid_json", "JSON input must contain a single object", nil)
		}
	}
	return payload, nil
}

func hasField(payload map[string]json.RawMessage, field string) bool {
	_, ok := payload[field]
	return ok
}

func mergeStringField(cmd *cobra.Command, flagName, canonical, value string, target *string, paramsRaw, bodyRaw map[string]json.RawMessage) error {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	if hasField(paramsRaw, canonical) || hasField(bodyRaw, canonical) {
		return output.NewError("input_conflict", "flag and JSON parameters provide the same field", map[string]string{"field": canonical})
	}
	*target = value
	return nil
}

func mergeBoolField(cmd *cobra.Command, flagName, canonical string, value bool, target *bool, paramsRaw, bodyRaw map[string]json.RawMessage) error {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	if hasField(paramsRaw, canonical) || hasField(bodyRaw, canonical) {
		return output.NewError("input_conflict", "flag and JSON parameters provide the same field", map[string]string{"field": canonical})
	}
	*target = value
	return nil
}

func mergeInt64Field(cmd *cobra.Command, flagName, canonical string, value int64, target *int64, paramsRaw, bodyRaw map[string]json.RawMessage) error {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	if hasField(paramsRaw, canonical) || hasField(bodyRaw, canonical) {
		return output.NewError("input_conflict", "flag and JSON parameters provide the same field", map[string]string{"field": canonical})
	}
	*target = value
	return nil
}

func mergeStringSliceField(cmd *cobra.Command, flagName, canonical string, value []string, target *[]string, paramsRaw, bodyRaw map[string]json.RawMessage) error {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	if hasField(paramsRaw, canonical) || hasField(bodyRaw, canonical) {
		return output.NewError("input_conflict", "flag and JSON parameters provide the same field", map[string]string{"field": canonical})
	}
	*target = value
	return nil
}

func mergeObjectField(cmd *cobra.Command, flagName, canonical, value string, target *json.RawMessage, paramsRaw, bodyRaw map[string]json.RawMessage) error {
	if !cmd.Flags().Changed(flagName) {
		return nil
	}
	if hasField(paramsRaw, canonical) || hasField(bodyRaw, canonical) {
		return output.NewError("input_conflict", "flag and JSON parameters provide the same field", map[string]string{"field": canonical})
	}
	var parsed map[string]any
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": fmt.Sprintf("%s cannot be empty", canonical)})
	}
	if err := json.Unmarshal([]byte(trimmed), &parsed); err != nil {
		return output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": err.Error()})
	}
	if parsed == nil {
		return output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": fmt.Sprintf("%s must be an object", canonical)})
	}
	*target = json.RawMessage(trimmed)
	return nil
}

func parseDuration(raw string) (time.Duration, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, output.NewError("validation_error", "duration is required", map[string]string{"field": "older_than"})
	}
	value, err := time.ParseDuration(raw)
	if err != nil {
		return 0, output.NewError("validation_error", "invalid duration", map[string]string{"field": "older_than", "error": err.Error()})
	}
	return value, nil
}

func assertJSONObject(raw json.RawMessage, field string) error {
	var parsed map[string]any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": err.Error(), "field": field})
	}
	if parsed == nil {
		return output.NewError("invalid_json", "invalid JSON input", map[string]string{"error": fmt.Sprintf("%s must be an object", field)})
	}
	return nil
}

func readJSONArg(raw string) ([]byte, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		raw = "{}"
	}
	switch {
	case raw == "-":
		return io.ReadAll(os.Stdin)
	case strings.HasPrefix(raw, "@"):
		path := strings.TrimPrefix(raw, "@")
		if path == "" {
			return nil, output.NewError("invalid_json_source", "file path after @ is required", nil)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("read JSON file: %w", err)
		}
		return data, nil
	default:
		return []byte(raw), nil
	}
}
