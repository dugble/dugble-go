package dugble

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestContactTopicsListUsesCursorPagination(t *testing.T) {
	client := newTestClient(t, func(req *http.Request) *http.Response {
		if got, want := req.URL.String(), "https://example.test/contacts/contact_1/topics?after=topic_a&limit=25"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
		return jsonResponse(http.StatusOK, `{"success":true,"data":{"object":"list","has_more":false,"data":[]}}`)
	})
	_, err := client.Contacts.Topics.List(context.Background(), "contact_1", ListContactTopicsParams{Limit: 25, After: "topic_a"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestContactPropertyUpdateCanSendNullFallback(t *testing.T) {
	client := newTestClient(t, func(req *http.Request) *http.Response {
		body, _ := io.ReadAll(req.Body)
		if got, want := string(body), `{"fallback_value":null}`; got != want {
			t.Fatalf("body = %s, want %s", got, want)
		}
		return jsonResponse(http.StatusOK, `{"success":true,"data":{"object":"contact_property","id":"prop_1"}}`)
	})
	_, err := client.ContactProperties.Update(context.Background(), "prop_1", UpdateContactPropertyParams{FallbackValue: nil})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSuppressionsListSerializesCursorAndOrigin(t *testing.T) {
	client := newTestClient(t, func(req *http.Request) *http.Response {
		if got, want := req.URL.String(), "https://example.test/suppressions?after=sup_1&limit=20&origin=complaint"; got != want {
			t.Fatalf("url = %q, want %q", got, want)
		}
		return jsonResponse(http.StatusOK, `{"success":true,"data":{"object":"list","has_more":false,"data":[]}}`)
	})
	_, err := client.Suppressions.List(context.Background(), ListSuppressionsParams{Limit: 20, After: "sup_1", Origin: SuppressionOriginComplaint})
	if err != nil {
		t.Fatal(err)
	}
}

func TestSenderIDCreateSerializesCountryCode(t *testing.T) {
	client := newTestClient(t, func(req *http.Request) *http.Response {
		body, _ := io.ReadAll(req.Body)
		got := string(body)
		for _, want := range []string{`"name":"DUGBLE"`, `"country_code":"GH"`, `"purpose":"transactional"`} {
			if !strings.Contains(got, want) {
				t.Fatalf("body %s missing %s", got, want)
			}
		}
		return jsonResponse(http.StatusOK, `{"success":true,"data":{"id":"sid_1","team_id":"team_1","name":"DUGBLE","country_code":"GH","purpose":"transactional","status":"pending","created_at":"x","updated_at":"x"}}`)
	})
	_, err := client.SenderIDs.Create(context.Background(), CreateSenderIDParams{Name: "DUGBLE", CountryCode: "GH", Purpose: "transactional"})
	if err != nil {
		t.Fatal(err)
	}
}
