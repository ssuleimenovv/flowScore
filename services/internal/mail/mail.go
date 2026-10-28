package mail

import (
	"bytes"
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// Message is one email, as plain text and as HTML: a mail app shows the one
// it can
type Message struct {
	To, Subject, Text, HTML string
}

// Sender sends an email: Brevo in production, Log locally.
type Sender interface {
	Send(ctx context.Context, m Message) error
}

// Brevo sends through Brevo's HTTP API. Its free plan sends 300 emails a day
// from a sender address confirmed in its dashboard; no domain of our own is
// needed. HTTP and not SMTP: free hosts often close the SMTP ports
type Brevo struct {
	Key  string // the API key
	From string // the confirmed sender address
	URL  string // empty for Brevo's own; a test server in tests
}

var client = &http.Client{Timeout: 10 * time.Second}

func (b Brevo) Send(ctx context.Context, m Message) error {
	body, err := json.Marshal(map[string]any{
		"sender":      map[string]string{"name": "FlowScore", "email": b.From},
		"to":          []map[string]string{{"email": m.To}},
		"subject":     m.Subject,
		"textContent": m.Text,
		"htmlContent": m.HTML,
	})
	if err != nil {
		return err
	}
	url := cmp.Or(b.URL, "https://api.brevo.com/v3/smtp/email")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("api-key", b.Key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode >= 300 {
		why, _ := io.ReadAll(io.LimitReader(res.Body, 512))
		return fmt.Errorf("brevo: %s: %s", res.Status, why)
	}
	return nil
}

// Log writes the email to the server's log instead of sending it: locally
// the code is read from the terminal
type Log struct{}

func (Log) Send(_ context.Context, m Message) error {
	log.Printf("mail to %s: %s\n%s", m.To, m.Subject, m.Text)
	return nil
}

// CodeMail is the email with a code, to confirm the address or to set a new
// password. The code is in the subject too, so it shows in the notification
func CodeMail(to, code string, reset bool) Message {
	what := "подтвердить почту"
	if reset {
		what = "задать новый пароль"
	}
	text := fmt.Sprintf("Твой код, чтобы %s: %s\n\n"+
		"Он действует 10 минут. Если это был не ты, просто не обращай внимания на письмо.\n\nFlowScore", what, code)
	html := fmt.Sprintf(`<div style="font-family:system-ui,sans-serif;font-size:15px;line-height:1.5;color:#161213">
<p>Твой код, чтобы %s:</p>
<p style="margin:16px 0;font-family:ui-monospace,monospace;font-size:32px;font-weight:700;letter-spacing:6px">%s</p>
<p style="color:#5e5557">Он действует 10 минут. Если это был не ты, просто не обращай внимания на письмо.</p>
<p style="color:#5e5557">FlowScore</p>
</div>`, what, code)
	return Message{To: to, Subject: "Код FlowScore: " + code, Text: text, HTML: html}
}
