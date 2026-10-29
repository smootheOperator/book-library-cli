package book

import "os"

type FileStorage struct {
	memory   *MemoryStorage
	filename string
}

func NewFileStorage(filename string) (*FileStorage, error) {
	store := NewMemoryStorage()
	filestore := FileStorage{
		memory:   store,
		filename: filename,
	}

	err := filestore.memory.LoadFromFile(filename)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	return &filestore, nil
}

func (s *FileStorage) save() error {
	books := s.memory.ListBooks()
	byteData, err := BooksToJSON(books)
	if err != nil {
		return err
	}

	err = os.WriteFile(s.filename, byteData, 0644)
	if err != nil {
		return nil
	}
	return nil
}

func (s *FileStorage) load() error {
	byteData, err := os.ReadFile(s.filename)
	if err != nil {
		return err
	}

	books, err := BooksFromJSON(byteData)
	if err != nil {
		return err
	}

	err = s.memory.ImportBooks(books)
	if err != nil {
		return err
	}

	return nil
}

func (s *FileStorage) GetBook(id int) (Book, bool) {
	return s.memory.GetBook(id)
}

func (s *FileStorage) ListBooks() []Book {
	return s.memory.ListBooks()
}

func (s *FileStorage) SearchByTitle(query string) []Book {
	return s.memory.SearchByTitle(query)
}

func (s *FileStorage) FilterByAuthor(author string) []Book {
	return s.memory.FilterByAuthor(author)
}

func (s *FileStorage) FilterByStatus(status string) ([]Book, error) {
	return s.memory.FilterByStatus(status)
}

func (s *FileStorage) FilterByMinRating(rating int) ([]Book, error) {
	return s.memory.FilterByMinRating(rating)
}

func (s *FileStorage) FilterByYearRange(from, to int) ([]Book, error) {
	return s.memory.FilterByYearRange(from, to)
}

func (s *FileStorage) ChangeStatus(id int, status string) error {
	err := s.memory.ChangeStatus(id, status)
	if err != nil {
		return err
	}

	return s.save()
}

func (s *FileStorage) DeleteBook(id int) error {
	err := s.memory.DeleteBook(id)
	if err != nil {
		return err
	}

	return s.save()
}

func (s *FileStorage) ChangeRating(id int, rating int) error {
	err := s.memory.ChangeRating(id, rating)
	if err != nil {
		return err
	}

	return s.save()
}

func (s *FileStorage) RenameBook(id int, title string) error {
	err := s.memory.RenameBook(id, title)
	if err != nil {
		return err
	}

	return s.save()
}

func (s *FileStorage) ChangeAuthor(id int, author string) error {
	err := s.memory.ChangeAuthor(id, author)
	if err != nil {
		return err
	}

	return s.save()
}

func (s *FileStorage) AddBook(title string, author string, year int) (Book, error) {
	book, err := s.memory.AddBook(title, author, year)
	if err != nil {
		return Book{}, err
	}

	if err = s.save(); err != nil {
		return Book{}, err
	}

	return book, nil
}
