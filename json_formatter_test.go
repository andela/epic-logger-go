package epiclogger

import (
	"encoding/json"
	"errors"
	"testing"

	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

var epicLogger = WithFields(map[string]interface{}{"service": "test-service", "version": "123"})

func TestErrorNotLost(t *testing.T) {
	formatter := &EpicFormatter{}

	b, err := formatter.Format(epicLogger.WithField("error", errors.New("wild walrus")).Entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}
	if entry["error"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestErrorNotLostOnFieldNotNamedError(t *testing.T) {
	formatter := &EpicFormatter{}

	b, err := formatter.Format(epicLogger.WithField("omg", errors.New("wild walrus")).Entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["omg"] != "wild walrus" {
		t.Fatal("Error field not set")
	}
}

func TestContextNotLost(t *testing.T) {
	formatter := &EpicFormatter{}

	ctx := metadata.NewIncomingContext(
		context.Background(),
		metadata.Pairs(
			"author_id", "this_author_id",
			"author_name", "this_author_name",
		),
	)
	b, err := formatter.Format(epicLogger.WithCtx(ctx).Entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}
	if entry["userId"] != "this_author_id" {
		t.Fatal("userId field not set")
	}
	context := entry["context"].(map[string]interface{})
	if context["user"] != "this_author_name" {
		t.Fatal("context.user field not set")
	}
}

func TestFieldClashWithTime(t *testing.T) {
	formatter := &EpicFormatter{}

	b, err := formatter.Format(epicLogger.WithField("time", "right now!").Entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.time"] != "right now!" {
		t.Fatal("fields.time not set to original time field")
	}

	if entry["time"] != "0001-01-01T00:00:00Z" {
		t.Fatal("time field not set to current time, was: ", entry["time"])
	}
}

func TestFieldClashWithMessage(t *testing.T) {
	formatter := &EpicFormatter{}

	b, err := formatter.Format(epicLogger.WithField("message", "something").Entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.message"] != "something" {
		t.Fatal("fields.message not set to original msg field")
	}
}

func TestFieldClashWithSeverity(t *testing.T) {
	formatter := &EpicFormatter{}

	b, err := formatter.Format(epicLogger.WithField("severity", "something").Entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	entry := make(map[string]interface{})
	err = json.Unmarshal(b, &entry)
	if err != nil {
		t.Fatal("Unable to unmarshal formatted entry: ", err)
	}

	if entry["fields.severity"] != "something" {
		t.Fatal("fields.severity not set to original level field")
	}
}

func TestJSONEntryEndsWithNewline(t *testing.T) {
	formatter := &EpicFormatter{}

	b, err := formatter.Format(epicLogger.WithField("level", "something").Entry)
	if err != nil {
		t.Fatal("Unable to format entry: ", err)
	}

	if b[len(b)-1] != '\n' {
		t.Fatal("Expected JSON log entry to end with a newline")
	}
}
