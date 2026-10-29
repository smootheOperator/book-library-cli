package book

type Service struct {
	storage Storage
}

func NewService(storage Storage) *Service {
	return &Service{
		storage: storage,
	}
}

func (s *Service) ListBooks() []Book {
	return s.storage.ListBooks()
}

func (s *Service) GetBook(id int) (Book, bool) {
	return s.storage.GetBook(id)
}

func (s *Service) SearchByTitle(query string) []Book {
	return s.storage.SearchByTitle(query)
}

func (s *Service) FilterByAuthor(author string) []Book {
	return s.storage.FilterByAuthor(author)
}

func (s *Service) FilterByStatus(status string) ([]Book, error) {
	return s.storage.FilterByStatus(status)
}

func (s *Service) FilterByMinRating(rating int) ([]Book, error) {
	return s.storage.FilterByMinRating(rating)
}

func (s *Service) SortedBooks(by string) ([]Book, error) {
	return SortBooks(s.storage.ListBooks(), by)

}

func (s *Service) Statistics() LibraryStats {
	return CalculateStats(s.storage.ListBooks())
}

func (s *Service) ChangeStatus(id int, status string) error {
	return s.storage.ChangeStatus(id, status)
}

func (s *Service) DeleteBook(id int) error {
	return s.storage.DeleteBook(id)
}

func (s *Service) ChangeRating(id int, rating int) error {
	return s.storage.ChangeRating(id, rating)
}

func (s *Service) RenameBook(id int, title string) error {
	return s.storage.RenameBook(id, title)
}

func (s *Service) ChangeAuthor(id int, author string) error {
	return s.storage.ChangeAuthor(id, author)
}

func (s *Service) AddBook(title string, author string, year int) (Book, error) {
	return s.storage.AddBook(title, author, year)
}
