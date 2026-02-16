package domain

type User struct {
	ID        int64  `json:"id"`
	LibraryID int64  `json:"library_id,omitempty"`
	Name      string `json:"full_name"`
}
