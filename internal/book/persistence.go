package book

import (
	"encoding/json"
	"errors"
	"os"
)

func BooksToJSON(books []Book) ([]byte, error) {
	jsonData, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func BooksFromJSON(data []byte) ([]Book, error) {
	var books []Book
	err := json.Unmarshal(data, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (s *MemoryStorage) ImportBooks(inputBooks []Book) error {
	// Шаг 4: Создаем новую временную карту
	newBooks := make(map[int]Book)
	maxID := 0

	// Перебираем импортируемые книги
	for _, book := range inputBooks {
		// Шаг 1: Проверяем валидность полей (Title, Author и т.д.)
		if err := ValidateBook(book); err != nil {
			return err
		}

		// Шаг 2: Проверяем, что ID строго больше 0
		if book.ID <= 0 {
			return errors.New("invalid book ID: must be greater than 0")
		}

		// Шаг 3: Проверка на дубликаты!
		// Мы проверяем, не занимал ли КТО-ТО ДРУГОЙ этот ID ранее в этом же импорте
		if _, exists := newBooks[book.ID]; exists {
			return errors.New("duplicate book ID found")
		}

		// Если проверки прошли, спокойно добавляем книгу во временную карту
		newBooks[book.ID] = book

		// Шаг 6: Параллельно ищем самый большой ID, чтобы вычислить nextID
		if book.ID > maxID {
			maxID = book.ID
		}
	}

	// Шаг 7: Только если ВЕСЬ цикл завершился без ошибок,
	// мы одним движением перезаписываем состояние старого хранилища
	s.books = newBooks
	s.nextID = maxID + 1

	return nil
}

func (s *MemoryStorage) SaveToFile(filename string) error {
	list := s.ListBooks()

	text, err := BooksToJSON(list)
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, text, 0644)
	if err != nil {
		return err
	}

	return nil
}

func (s *MemoryStorage) LoadFromFile(filename string) error {
	text, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	books, err := BooksFromJSON(text)
	if err != nil {
		return err
	}

	err = s.ImportBooks(books)
	if err != nil {
		return err
	}

	return nil
}
