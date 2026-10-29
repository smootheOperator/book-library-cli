package book

import (
	"errors"
	"sort"
)

type MemoryStorage struct {
	books  map[int]Book
	nextID int
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		books:  make(map[int]Book),
		nextID: 1,
	}
}

func (s *MemoryStorage) AddBook(title string, author string, year int) (Book, error) {
	b := Book{
		ID:     s.nextID,
		Title:  title,
		Author: author,
		Year:   year,
		Status: "planned",
		Rating: 0,
	}
	if err := ValidateBook(b); err != nil {
		return Book{}, err
	}
	s.books[b.ID] = b
	s.nextID++
	return b, nil
}

func (s *MemoryStorage) GetBook(id int) (Book, bool) {
	book, ok := s.books[id]
	return book, ok
}

func (s *MemoryStorage) DeleteBook(id int) error {
	if _, exists := s.books[id]; !exists {
		return errors.New("book does not exist")
	}
	delete(s.books, id)
	return nil
}

func (s *MemoryStorage) ListBooks() []Book {
	bookList := make([]Book, 0, len(s.books))

	for _, book := range s.books {
		bookList = append(bookList, book)
	}

	sort.Slice(bookList, func(i, j int) bool {
		return bookList[i].ID < bookList[j].ID
	})

	return bookList
}

func (s *MemoryStorage) ChangeStatus(id int, status string) error {
	book, ok := s.books[id]
	if !ok {
		return errors.New("book does not exist")
	}

	if err := ValidateStatus(status); err != nil {
		return err
	}

	book.Status = status
	s.books[id] = book
	return nil
}

func (s *MemoryStorage) ChangeRating(id int, rating int) error {
	book, ok := s.books[id]
	if !ok {
		return errors.New("book does not exist")
	}

	if err := ValidateRating(rating); err != nil {
		return err
	}

	book.Rating = rating
	s.books[id] = book
	return nil
}

func (s *MemoryStorage) RenameBook(id int, title string) error {
	book, ok := s.books[id]
	if !ok {
		return errors.New("book does not exist")
	}

	if err := ValidateTitle(title); err != nil {
		return err
	}

	book.Title = title
	s.books[id] = book
	return nil
}

func (s *MemoryStorage) ChangeAuthor(id int, author string) error {
	book, ok := s.books[id]
	if !ok {
		return errors.New("book does not exist")
	}

	if err := ValidateAuthor(author); err != nil {
		return err
	}

	book.Author = author
	s.books[id] = book
	return nil
}
