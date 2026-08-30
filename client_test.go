package dugble

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (fn roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return fn(req)
}

func newTestClient(t *testing.T, handler func(*http.Request) *http.Response) *Client {
	t.Helper()
	httpClient := &http.Client{Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
		return handler(req), nil
	})}
	return New("dgb_team_test", WithBaseURL("https://example.test"), WithHTTPClient(httpClient))
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     make(http.Header),
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

func TestTopicsListUsesOffsetPagination(t *testing.T) {
	client := newTestClient(t, func(req *http.Request) *http.Response {
		if got, want := req.URL.String(), "https://example.test/topics?limit=20&offset=40"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
		if got := req.Header.Get("Authorization"); got != "Bearer dgb_team_test" {
			t.Fatalf("authorization = %q", got)
		}
		return jsonResponse(http.StatusOK, `{"success":true,"data":{"object":"list","has_more":false,"data":[]}}`)
	})

	result, err := client.Topics.List(context.Background(), ListTopicsParams{Limit: 20, Offset: 40})
	if err != nil {
		t.Fatal(err)
	}
	if result.Object != "list" || result.HasMore {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestSegmentsAudienceSizeEncodesID(t *testing.T) {
	client := newTestClient(t, func(req *http.Request) *http.Response {
		if got, want := req.URL.String(), "https://example.test/segments/segment%2F123/audience-size"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
		return jsonResponse(http.StatusOK, `{"success":true,"data":{"segment_id":"segment_123","count":42}}`)
	})

	result, err := client.Segments.AudienceSize(context.Background(), "segment/123")
	if err != nil {
		t.Fatal(err)
	}
	if result.Count != 42 {
		t.Fatalf("count = %d, want 42", result.Count)
	}
}

func TestAPIErrorIncludesStatusAndRequestID(t *testing.T) {
	client := newTestClient(t, func(req *http.Request) *http.Response {
		resp := jsonResponse(http.StatusBadRequest, `{"success":false,"error":{"code":"BAD_REQUEST","message":"invalid request"}}`)
		resp.Header.Set("x-request-id", "req_123")
		return resp
	})

	_, err := client.Topics.Get(context.Background(), "topic_123")
	if err == nil {
		t.Fatal("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("error type = %T", err)
	}
	if apiErr.StatusCode != http.StatusBadRequest || apiErr.RequestID != "req_123" {
		t.Fatalf("unexpected api error: %#v", apiErr)
	}
}
