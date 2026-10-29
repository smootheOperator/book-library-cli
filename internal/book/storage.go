package book

type Storage interface {
	AddBook(title, author string, year int) (Book, error)
	GetBook(id int) (Book, bool)
	DeleteBook(id int) error
	ListBooks() []Book
	ChangeStatus(id int, status string) error
	ChangeRating(id, rating int) error
	RenameBook(id int, title string) error
	ChangeAuthor(id int, author string) error
	SearchByTitle(query string) []Book
	FilterByAuthor(author string) []Book
	FilterByStatus(status string) ([]Book, error)
	FilterByMinRating(rating int) ([]Book, error)
	FilterByYearRange(from, to int) ([]Book, error)
}
