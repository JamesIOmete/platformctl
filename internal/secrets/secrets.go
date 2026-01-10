package secrets

import "time"

// Secret represents a mocked secret (e.g. from Vault).
type Secret struct {
	Key       string    `json:"key"`
	Value     string    `json:"value"` // In real life, never store plaintext!
	CreatedBy string    `json:"created_by"`
	CreatedAt time.Time `json:"created_at"`
}
