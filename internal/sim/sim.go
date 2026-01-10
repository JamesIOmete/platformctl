package sim

// Job represents a cloud simulation job.
type Job struct {
	ID        string `json:"id"`
	Scenario  string `json:"scenario"`
	Status    string `json:"status"` // Pending, Running, Success, Failed
	CreatedAt string `json:"created_at"`
}
