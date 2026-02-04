package domain

type Book struct {
	ID     int64
	Title  string
	Author string
	Year   int
	Pages  int

	Readers []ReadingInfo
}
