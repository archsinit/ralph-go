# Implement ntfy client

Goal: POST notifications to an ntfy server topic with optional auth.

Context:
- Default server https://ntfy.sh.

Do:
1. In internal/ntfy/ntfy.go define Client { Server, Topic, Token string; HTTP *http.Client }.
2. func New(cfg config.NtfyConfig) *Client.
3. func (c *Client) Publish(title, body string, priority int, tags []string) error: POST to Server/Topic with headers Title, Priority, Tags; Authorization Bearer if Token set.
4. Short timeout; return error but caller treats as non-fatal.

Done when:
- go build ./... succeeds.
- Request hits Server/Topic with correct headers.
- Token adds Authorization header.

Files: internal/ntfy/ntfy.go
