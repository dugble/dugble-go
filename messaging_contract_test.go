package dugble

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func decodeBody(t *testing.T, req *http.Request) map[string]any {
	t.Helper()
	if req.Body == nil {
		return nil
	}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("decode body: %v (%s)", err, data)
	}
	return body
}

func TestDomainsUpdateContract(t *testing.T) {
	tls := DomainTLSEnforced
	client := newTestClient(t, func(req *http.Request) *http.Response {
		if got, want := req.Method, http.MethodPatch; got != want {
			t.Fatalf("method = %q, want %q", got, want)
		}
		if got, want := req.URL.String(), "https://example.test/domains/domain%2F123"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
		body := decodeBody(t, req)
		if got := body["tls"]; got != "enforced" {
			t.Fatalf("tls = %#v", got)
		}
		return jsonResponse(http.StatusOK, `{"success":true,"data":{"id":"domain_123","team_id":"team_123","name":"example.com","region":"us-east-1","status":"verified","records":[],"tls":"enforced","health_status":"healthy","consecutive_health_failures":0,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}}`)
	})

	if _, err := client.Domains.Update(context.Background(), "domain/123", UpdateDomainParams{TLS: &tls}); err != nil {
		t.Fatal(err)
	}
}

func TestEmailSendAndEventsContract(t *testing.T) {
	calls := 0
	client := newTestClient(t, func(req *http.Request) *http.Response {
		calls++
		switch calls {
		case 1:
			if got, want := req.URL.Path, "/emails"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			if got := req.Header.Get("Idempotency-Key"); got != "idem_email" {
				t.Fatalf("idempotency = %q", got)
			}
			body := decodeBody(t, req)
			if body["subject"] != "Hello" || body["scheduled_at"] != "2026-09-01T10:00:00Z" {
				t.Fatalf("unexpected email body: %#v", body)
			}
			return jsonResponse(http.StatusOK, `{"success":true,"data":{"object":"email","id":"email_123"}}`)
		case 2:
			if got, want := req.URL.String(), "https://example.test/emails/email%2F123/events?limit=20&offset=40"; got != want {
				t.Fatalf("url = %q, want %q", got, want)
			}
			return jsonResponse(http.StatusOK, `{"success":true,"data":{"object":"list","data":[]}}`)
		default:
			t.Fatalf("unexpected request %d", calls)
		}
		return nil
	})

	_, err := client.Emails.Send(context.Background(), SendEmailParams{
		To:          "person@example.com",
		Subject:     "Hello",
		ScheduledAt: "2026-09-01T10:00:00Z",
	}, "idem_email")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Emails.Events(context.Background(), "email/123", ListEmailEventsParams{Limit: 20, Offset: 40}); err != nil {
		t.Fatal(err)
	}
}

func TestSmsListAndSyncContract(t *testing.T) {
	calls := 0
	client := newTestClient(t, func(req *http.Request) *http.Response {
		calls++
		switch calls {
		case 1:
			want := "https://example.test/sms?end_date=2026-08-31&limit=25&offset=50&search=233&sender=Dugble&start_date=2026-08-01&status=delivered"
			if got := req.URL.String(); got != want {
				t.Fatalf("url = %q, want %q", got, want)
			}
			return jsonResponse(http.StatusOK, `{"success":true,"data":[]}`)
		case 2:
			if got, want := req.URL.String(), "https://example.test/sms/sms%2F123/sync-status"; got != want {
				t.Fatalf("url = %q, want %q", got, want)
			}
			return jsonResponse(http.StatusOK, `{"success":true,"data":{"object":"sms","id":"sms_123","message_id":null,"to":"+233200000000","from":"Dugble","body":"hello","last_event":"delivered","destination":{"country":"GH"},"segments":1,"metadata":null,"scheduled_at":null,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}}`)
		default:
			t.Fatalf("unexpected request %d", calls)
		}
		return nil
	})

	_, err := client.SMS.List(context.Background(), ListSmsParams{
		Limit: 25, Offset: 50, Status: SmsStatusDelivered, Sender: "Dugble",
		StartDate: "2026-08-01", EndDate: "2026-08-31", Search: "233",
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.SMS.SyncStatus(context.Background(), "sms/123"); err != nil {
		t.Fatal(err)
	}
}

func TestTemplateCreateAndVersionsContract(t *testing.T) {
	calls := 0
	client := newTestClient(t, func(req *http.Request) *http.Response {
		calls++
		switch calls {
		case 1:
			if got, want := req.URL.Path, "/templates"; got != want {
				t.Fatalf("path = %q, want %q", got, want)
			}
			body := decodeBody(t, req)
			if body["category"] != "welcome" || body["name"] != "Welcome" {
				t.Fatalf("unexpected template body: %#v", body)
			}
			return jsonResponse(http.StatusOK, `{"success":true,"data":{"object":"template","id":"template_123"}}`)
		case 2:
			if got, want := req.URL.String(), "https://example.test/templates/welcome%2Femail/versions?limit=20&offset=40"; got != want {
				t.Fatalf("url = %q, want %q", got, want)
			}
			return jsonResponse(http.StatusOK, `{"success":true,"data":[]}`)
		default:
			t.Fatalf("unexpected request %d", calls)
		}
		return nil
	})

	_, err := client.Templates.Create(context.Background(), CreateTemplateParams{
		Name: "Welcome", HTML: "<h1>Hello</h1>", Category: TemplateCategoryWelcome,
	})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := client.Templates.Versions.List(context.Background(), "welcome/email", ListTemplateVersionsParams{Limit: 20, Offset: 40}); err != nil {
		t.Fatal(err)
	}
}

func TestBroadcastOwnedContentAndNullableUpdate(t *testing.T) {
	calls := 0
	client := newTestClient(t, func(req *http.Request) *http.Response {
		calls++
		switch calls {
		case 1:
			body := decodeBody(t, req)
			for key, want := range map[string]any{
				"segment_id": "segment_123",
				"subject":    "Launch",
				"html":       "<p>Launch</p>",
			} {
				if got := body[key]; got != want {
					t.Fatalf("%s = %#v, want %#v", key, got, want)
				}
			}
			if _, exists := body["template_id"]; exists {
				t.Fatalf("broadcast request unexpectedly contains template_id: %#v", body)
			}
			return broadcastResponse()
		case 2:
			body := decodeBody(t, req)
			if body["revision"] != float64(2) {
				t.Fatalf("revision = %#v", body["revision"])
			}
			if value, exists := body["topic_id"]; !exists || value != nil {
				t.Fatalf("topic_id should be explicit null: %#v", body)
			}
			return broadcastResponse()
		case 3:
			data, err := io.ReadAll(req.Body)
			if err != nil {
				t.Fatal(err)
			}
			if strings.TrimSpace(string(data)) != "{}" {
				t.Fatalf("duplicate body = %q, want {}", data)
			}
			return broadcastResponse()
		default:
			t.Fatalf("unexpected request %d", calls)
		}
		return nil
	})

	if _, err := client.Broadcasts.Create(context.Background(), CreateBroadcastParams{
		SegmentID: "segment_123", Subject: "Launch", HTML: "<p>Launch</p>",
	}); err != nil {
		t.Fatal(err)
	}
	var nilString *string
	if _, err := client.Broadcasts.Update(context.Background(), "broadcast_123", UpdateBroadcastParams{
		Revision: 2, TopicID: &nilString,
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := client.Broadcasts.Duplicate(context.Background(), "broadcast_123", DuplicateBroadcastParams{}); err != nil {
		t.Fatal(err)
	}
}

func broadcastResponse() *http.Response {
	return jsonResponse(http.StatusOK, `{"success":true,"data":{"id":"broadcast_123","team_id":"team_123","name":"Launch","status":"draft","segment_id":"segment_123","from_email":"","subject":"Launch","html":"<p>Launch</p>","variable_bindings":{},"audience_count":0,"eligible_count":0,"suppressed_count":0,"queued_count":0,"failed_count":0,"revision":1,"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-01T00:00:00Z"}}`)
}
