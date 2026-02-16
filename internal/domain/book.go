package domain

type Book struct {
	ID        int64  `json:"id"`
	LibraryID int64  `json:"library_id"`
	Title     string `json:"title"`
	Author    string `json:"author"`
	Year      int    `json:"year"`
	Pages     int    `json:"pages"`

	Readers []ReadingInfo `json:"readers,omitempty"`
}
