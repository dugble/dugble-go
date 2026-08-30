# dugble-go

Official Go SDK for the Dugble API.

## Requirements

- Go 1.26.6 or newer
- A Dugble team API key

## Install

```bash
go get github.com/dugble/dugble-go
```

## Quickstart

```go
package main

import (
	"context"
	"log"

	dugble "github.com/dugble/dugble-go"
)

func main() {
	ctx := context.Background()
	client := dugble.New("dgb_team_...")

	topics, err := client.Topics.List(ctx, dugble.ListTopicsParams{Limit: 20})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("topics: %d", len(topics.Data))
}
```

The client defaults to `https://api.dugble.com`, sends bearer authentication automatically, unwraps Dugble API envelopes, and returns structured `*dugble.APIError` values for API failures.

## Client configuration

Use `WithBaseURL` for self-hosted or test environments and `WithHTTPClient` when you need a custom transport, proxy, timeout, or tracing setup.

```go
httpClient := &http.Client{Timeout: 10 * time.Second}

client := dugble.New(
	"dgb_team_...",
	dugble.WithBaseURL("https://api.example.com"),
	dugble.WithHTTPClient(httpClient),
)
```

## Error handling

```go
message, err := client.Emails.Get(ctx, "email_123")
if err != nil {
	var apiErr *dugble.APIError
	if errors.As(err, &apiErr) {
		log.Printf("status=%d code=%s request_id=%s", apiErr.StatusCode, apiErr.Code, apiErr.RequestID)
	}
	return
}

log.Println(message.ID)
```

## Email

Send an email. The idempotency key is optional; pass an empty string when you do not need one.

```go
sent, err := client.Emails.Send(ctx, dugble.SendEmailParams{
	From:    "Dugble <hello@example.com>",
	To:      "person@example.com",
	Subject: "Welcome",
	HTML:    "<h1>Welcome</h1>",
}, "checkout-email-123")
if err != nil {
	log.Fatal(err)
}

log.Println(sent.ID)
```

Addresses may be supplied in the API-supported string or object forms. Attachments use Base64 content plus a filename; attachment paths are not supported by the public API.

```go
sent, err := client.Emails.Send(ctx, dugble.SendEmailParams{
	To:      dugble.EmailAddress{Email: "person@example.com", Name: "Ada"},
	Subject: "Your receipt",
	Text:    "Receipt attached",
	Attachments: []dugble.EmailAttachment{
		{
			Filename:    "receipt.txt",
			Content:     base64.StdEncoding.EncodeToString([]byte("receipt")),
			ContentType: "text/plain",
		},
	},
}, "receipt-123")
```

Scheduling, events, analytics, cancellation, and batch sending are also available:

```go
_, err = client.Emails.Send(ctx, dugble.SendEmailParams{
	To:          "person@example.com",
	Subject:     "Reminder",
	Text:        "Your appointment is tomorrow",
	ScheduledAt: "2026-09-01T10:00:00Z",
}, "reminder-123")

_, err = client.Emails.Events(ctx, "email_123", dugble.ListEmailEventsParams{
	Limit:  50,
	Offset: 0,
})

analytics, err := client.Emails.Analytics(ctx)
```

`scheduled_at` on send accepts the API's RFC3339/RFC3339Nano and supported relative schedule forms. Rescheduling with `Emails.Update` uses an RFC3339 timestamp.

## SMS

```go
sent, err := client.SMS.Send(ctx, dugble.SendSmsParams{
	To:   "+233200000000",
	From: "Dugble",
	Body: "Your code is 123456",
}, "sms-123")
if err != nil {
	log.Fatal(err)
}

log.Println(sent.ID)
```

List messages with the current server-side filters:

```go
messages, err := client.SMS.List(ctx, dugble.ListSmsParams{
	Limit:     50,
	Status:    dugble.SmsStatusDelivered,
	Sender:    "Dugble",
	StartDate: "2026-08-01",
	EndDate:   "2026-08-31",
	Search:    "233",
})
```

SMS also supports batch sending, scheduling/rescheduling, cancellation, events, analytics, and status synchronization.

```go
message, err := client.SMS.SyncStatus(ctx, "sms_123")
```

## Domains

```go
domain, err := client.Domains.Create(ctx, dugble.CreateDomainParams{
	Name:   "mail.example.com",
	Region: dugble.DomainRegionUSEast1,
})
if err != nil {
	log.Fatal(err)
}

if domain.Provisioning != nil {
	log.Printf("provisioning; retry after %d seconds", domain.Provisioning.RetryAfterSeconds)
}
if domain.Domain != nil {
	log.Println(domain.Domain.ID)
}
```

Verify a domain or enforce TLS:

```go
_, err = client.Domains.Verify(ctx, "domain_123")

tls := dugble.DomainTLSEnforced
_, err = client.Domains.Update(ctx, "domain_123", dugble.UpdateDomainParams{TLS: &tls})
```

## Contacts

```go
unsubscribed := false
contact, err := client.Contacts.Create(ctx, dugble.CreateContactParams{
	Email:            "person@example.com",
	Phone:            "+233200000000",
	SMSConsentStatus: dugble.SMSConsentOptedIn,
	SMSConsentSource: dugble.SMSConsentSourceAPI,
	FirstName:        "Ada",
	Unsubscribed:     &unsubscribed,
	Properties: map[string]any{
		"plan": "pro",
	},
})
```

Contacts expose topic subscriptions and explicit segment membership as nested services.

```go
topics, err := client.Contacts.Topics.List(ctx, contact.ID, dugble.ListContactTopicsParams{
	Limit: 50,
})

_, err = client.Contacts.Topics.Update(ctx, contact.ID, []dugble.UpdateContactTopic{
	{ID: "topic_123", Subscription: dugble.ContactTopicOptIn},
})

_, err = client.Contacts.Segments.Add(ctx, contact.ID, "segment_123")
err = client.Contacts.Segments.Remove(ctx, contact.ID, "segment_123")
```

## Contact properties

```go
property, err := client.ContactProperties.Create(ctx, dugble.CreateContactPropertyParams{
	Key:           "lifetime_value",
	Type:          dugble.ContactPropertyNumber,
	FallbackValue: 0,
})
```

Clear a fallback value by sending `nil` explicitly:

```go
_, err = client.ContactProperties.Update(ctx, property.ID, dugble.UpdateContactPropertyParams{
	FallbackValue: nil,
})
```

## Topics

```go
created, err := client.Topics.Create(ctx, dugble.CreateTopicParams{
	Name:                "Product updates",
	DefaultSubscription: dugble.TopicSubscriptionOptOut,
	Visibility:          dugble.TopicVisibilityPublic,
})
```

Topic descriptions are nullable on update. Pointer-to-pointer fields distinguish omitted from explicit `null`:

```go
var cleared *string
_, err = client.Topics.Update(ctx, created.ID, dugble.UpdateTopicParams{
	Description: &cleared, // sends "description": null
})
```

## Segments

```go
segment, err := client.Segments.Create(ctx, dugble.CreateSegmentParams{Name: "Customers"})
if err != nil {
	log.Fatal(err)
}

size, err := client.Segments.AudienceSize(ctx, segment.ID)
log.Printf("explicit members: %d", size.Count)
```

`AudienceSize` reports explicit segment membership for the current public contract.

## Templates

```go
template, err := client.Templates.Create(ctx, dugble.CreateTemplateParams{
	Name:     "Welcome",
	Category: dugble.TemplateCategoryWelcome,
	Subject:  "Welcome, {{name}}",
	HTML:     "<h1>Welcome, {{name}}</h1>",
	Variables: []dugble.TemplateVariableParam{
		{Key: "name", Type: dugble.TemplateVariableString},
	},
})
```

Publish, preview, test-send, and inspect versions:

```go
_, err = client.Templates.Publish(ctx, template.ID)

preview, err := client.Templates.Preview(ctx, template.ID, dugble.PreviewTemplateParams{
	Variables: map[string]any{"name": "Ada"},
})

_, err = client.Templates.TestSend(ctx, template.ID, dugble.TestSendTemplateParams{
	To:        "person@example.com",
	Variables: map[string]any{"name": "Ada"},
})

versions, err := client.Templates.Versions.List(ctx, template.ID, dugble.ListTemplateVersionsParams{
	Limit: 20,
})
```

## Broadcasts

Broadcasts own the exact message content they send. There is no `template_id` coupling in the public broadcast contract.

```go
broadcast, err := client.Broadcasts.Create(ctx, dugble.CreateBroadcastParams{
	Name:      "September launch",
	SegmentID: "segment_123",
	TopicID:   "topic_123",
	FromEmail: "hello@example.com",
	Subject:   "We shipped something new",
	HTML:      "<h1>Launch</h1>",
})
```

Queue immediately or schedule for later:

```go
_, err = client.Broadcasts.Send(ctx, broadcast.ID, dugble.SendBroadcastParams{})

_, err = client.Broadcasts.Send(ctx, broadcast.ID, dugble.SendBroadcastParams{
	ScheduledAt: "2026-09-01T10:00:00Z",
})
```

Updates are revision-based. Nullable fields use pointer-to-pointer values so the SDK can distinguish “leave unchanged” from “clear this value.”

```go
var clearedTopic *string
updated, err := client.Broadcasts.Update(ctx, broadcast.ID, dugble.UpdateBroadcastParams{
	Revision: broadcast.Revision,
	TopicID:  &clearedTopic,
})
```

Duplicate with the default generated copy name or provide one:

```go
copy, err := client.Broadcasts.Duplicate(ctx, broadcast.ID)

namedCopy, err := client.Broadcasts.Duplicate(ctx, broadcast.ID, dugble.DuplicateBroadcastParams{
	Name: "September launch - resend",
})
```

## Suppressions

```go
_, err := client.Suppressions.Create(ctx, dugble.CreateSuppressionParams{
	Email: "person@example.com",
})

list, err := client.Suppressions.List(ctx, dugble.ListSuppressionsParams{
	Limit:  50,
	Origin: dugble.SuppressionOriginBounce,
})
```

Batch operations support up to the API's current batch limit. Batch removal must identify records by either emails or IDs, not both.

```go
_, err = client.Suppressions.BatchAdd(ctx, dugble.BatchAddSuppressionsParams{
	Emails: []string{"a@example.com", "b@example.com"},
})

_, err = client.Suppressions.BatchRemove(ctx, dugble.BatchRemoveSuppressionsParams{
	IDs: []string{"suppression_123", "suppression_456"},
})
```

## Sender IDs

```go
senderID, err := client.SenderIDs.Create(ctx, dugble.CreateSenderIDParams{
	Name:        "Dugble",
	CountryCode: "GH",
	Purpose:     "Transactional notifications",
})

senderIDs, err := client.SenderIDs.List(ctx)
```

Sender ID lifecycle statuses currently include `pending`, `approved`, `rejected`, `suspended`, and `inactive`.

## Pagination

Dugble's public resources currently use two pagination styles. The SDK mirrors each API exactly rather than hiding the distinction.

| Resource | Pagination |
| --- | --- |
| Topics | `limit` / `offset` |
| Contacts | `limit` / `offset` |
| Segments | `limit` / `offset` |
| Domains | `limit` / `offset` |
| Emails | `limit` / `offset` |
| SMS | `limit` / `offset` |
| Templates and versions | `limit` / `offset` |
| Broadcasts and recipients | `limit` / `offset` |
| Contact topics | `limit` / `after` / `before` |
| Contact properties | `limit` / `after` / `before` |
| Suppressions | `limit` / `after` / `before` |
| Sender IDs | no pagination currently |

## Current service surface

`Client` currently exposes:

- `Topics`
- `Segments`
- `Domains`
- `Emails`
- `SMS`
- `Templates`
- `Broadcasts`
- `Contacts`
- `ContactProperties`
- `Suppressions`
- `SenderIDs`

The SDK intentionally targets Dugble's public team-token API surface. Dashboard/session-only endpoints are not exposed as SDK services.
