package public

import "fmt"

// Notice represents a Bithumb announcement
type Notice struct {
	Categories  []string `json:"categories"`
	Title       string   `json:"title"`
	PCURL       string   `json:"pc_url"`
	PublishedAt string   `json:"published_at"`
	ModifiedAt  string   `json:"modified_at"`
}

// GetNoticesRequest represents request parameters for GetNotices
type GetNoticesRequest struct {
	Count int // Number of notices to retrieve (1-20, default: 5)
}

// Validate validates the request parameters
func (r *GetNoticesRequest) Validate() error {
	if r.Count < 1 || r.Count > 20 {
		return fmt.Errorf("count must be between 1 and 20, got %d", r.Count)
	}
	return nil
}
