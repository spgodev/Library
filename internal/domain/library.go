package domain

type Library struct {
	Name  string
	Books []Book
}

func (l *Library) HasTitle(title string) bool {
	for _, b := range l.Books {
		if b.Title == title {
			return true
		}
	}
	return false
}

func (l *Library) FindByTitle(title string) []Book {
	out := make([]Book, 0)
	for _, b := range l.Books {
		if b.Title == title {
			out = append(out, b)
		}
	}
	return out
}

func (l *Library) FilterByYear(year int) []Book {
	out := make([]Book, 0)
	for _, b := range l.Books {
		if b.Year == year {
			out = append(out, b)
		}
	}
	return out
}

func (l *Library) FilterByAuthor(author string) []Book {
	out := make([]Book, 0)
	for _, b := range l.Books {
		if b.Author == author {
			out = append(out, b)
		}
	}
	return out
}

func (l *Library) SortedByYear(asc bool) []Book {
	cp := make([]Book, len(l.Books))
	copy(cp, l.Books)

	less := func(i, j int) bool {
		if cp[i].Year != cp[j].Year {
			if asc {
				return cp[i].Year < cp[j].Year
			}
			return cp[i].Year > cp[j].Year
		}
		return cp[i].ID < cp[j].ID
	}

	for i := 0; i < len(cp); i++ {
		for j := i + 1; j < len(cp); j++ {
			if less(j, i) {
				cp[i], cp[j] = cp[j], cp[i]
			}
		}
	}
	return cp
}
