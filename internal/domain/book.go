package domain

type Book struct {
	ID     int
	Title  string
	Author string
	Year   int
	Pages  int

	Readers []ReadingInfo
}
