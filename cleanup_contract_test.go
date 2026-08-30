package dugble

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestTopicUpdateCanClearDescription(t *testing.T) {
	client := newTestClient(t, func(req *http.Request) *http.Response {
		if got, want := req.Method, http.MethodPatch; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := req.URL.String(), "https://example.test/topics/topic%2F123"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
		body := decodeBody(t, req)
		value, exists := body["description"]
		if !exists || value != nil {
			t.Fatalf("description should be explicit null: %#v", body)
		}
		return jsonResponse(http.StatusOK, `{"success":true,"data":{"object":"topic","id":"topic_123"}}`)
	})

	var cleared *string
	_, err := client.Topics.Update(context.Background(), "topic/123", UpdateTopicParams{
		Description: &cleared,
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestBroadcastDuplicateAllowsNoParamsAndSendsObject(t *testing.T) {
	client := newTestClient(t, func(req *http.Request) *http.Response {
		if got, want := req.Method, http.MethodPost; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := req.URL.String(), "https://example.test/broadcasts/broadcast%2F123/duplicate"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
		data, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatal(err)
		}
		if strings.TrimSpace(string(data)) != "{}" {
			t.Fatalf("duplicate body = %q, want {}", data)
		}
		return broadcastResponse()
	})

	if _, err := client.Broadcasts.Duplicate(context.Background(), "broadcast/123"); err != nil {
		t.Fatal(err)
	}
}
