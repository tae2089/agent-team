package output

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestWriteErrorIncludesRecoveryGuide(t *testing.T) {
	var out bytes.Buffer
	WriteError(&out, 12, NewError("sync_conflict", "blocking sync mismatch", map[string]string{"task_id": "task_a"}))

	var got map[string]any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	body := got["error"].(map[string]any)
	if body["code"] != "sync_conflict" {
		t.Fatalf("unexpected code: %#v", body["code"])
	}
	if _, ok := body["details"].(map[string]any); !ok {
		t.Fatalf("missing details: %#v", body)
	}
	recovery := body["recovery"].(map[string]any)
	if recovery["summary"] == "" {
		t.Fatalf("missing recovery summary: %#v", recovery)
	}
	if len(recovery["actions"].([]any)) == 0 {
		t.Fatalf("missing recovery actions: %#v", recovery)
	}
	if len(recovery["skills"].([]any)) == 0 {
		t.Fatalf("missing recovery skills: %#v", recovery)
	}
}

func TestWriteErrorFallsBackToInternalRecovery(t *testing.T) {
	var out bytes.Buffer
	WriteError(&out, 0, assertErr("sqlite unavailable"))

	var got map[string]any
	if err := json.Unmarshal(out.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, out.String())
	}
	body := got["error"].(map[string]any)
	if body["code"] != "internal_error" {
		t.Fatalf("unexpected code: %#v", body["code"])
	}
	recovery := body["recovery"].(map[string]any)
	docs := recovery["docs"].([]any)
	if len(docs) == 0 {
		t.Fatalf("missing recovery docs: %#v", recovery)
	}
}

type assertErr string

func (e assertErr) Error() string {
	return string(e)
}
