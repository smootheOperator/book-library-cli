package book

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBooksJSONRoundTrip(t *testing.T) {
	b := []Book{
		{
			ID:     1,
			Title:  "Title1",
			Author: "Author1",
			Year:   2001,
			Status: "planned",
			Rating: 0,
		},
		{
			ID:     2,
			Title:  "Title2",
			Author: "Author2",
			Year:   2002,
			Status: "planned",
			Rating: 2,
		},
		{
			ID:     3,
			Title:  "Title3",
			Author: "Author3",
			Year:   2003,
			Status: "planned",
			Rating: 3,
		},
	}

	text, err := BooksToJSON(b)
	if err != nil {
		t.Fatal(err)
	}

	books, err := BooksFromJSON(text)
	if err != nil {
		t.Fatal(err)
	}

	if len(books) != len(b) {
		t.Fatalf("expected %d books, got %d", len(b), len(books))
	}

	for i := range b {
		if books[i] != b[i] {
			t.Errorf("Different books: %+v\n%+v", books[i], b[i])
		}
	}
}

func TestImportBooksRestoresNextID(t *testing.T) {
	store := MemoryStorage{
		books:  make(map[int]Book),
		nextID: 1,
	}

	b := []Book{
		{
			ID:     1,
			Title:  "Title1",
			Author: "Author1",
			Year:   2001,
			Status: "planned",
			Rating: 0,
		},
		{
			ID:     4,
			Title:  "Title2",
			Author: "Author2",
			Year:   2002,
			Status: "planned",
			Rating: 2,
		},
		{
			ID:     9,
			Title:  "Title3",
			Author: "Author3",
			Year:   2003,
			Status: "planned",
			Rating: 3,
		},
	}

	err := store.ImportBooks(b)
	if err != nil {
		t.Fatal(err)
	}

	_, err = store.AddBook("Title4", "Author4", 2004)
	if err != nil {
		t.Fatal(err)
	}

}

func TestImportBooksRejectsDuplicateIDs(t *testing.T) {
	store := MemoryStorage{
		books:  make(map[int]Book),
		nextID: 1,
	}

	b := []Book{
		{
			ID:     1,
			Title:  "Title1",
			Author: "Author1",
			Year:   2001,
			Status: "planned",
			Rating: 0,
		},
		{
			ID:     2,
			Title:  "Title2",
			Author: "Author2",
			Year:   2002,
			Status: "planned",
			Rating: 2,
		},
		{
			ID:     1,
			Title:  "Title3",
			Author: "Author3",
			Year:   2003,
			Status: "planned",
			Rating: 3,
		},
	}

	err := store.ImportBooks(b)
	if err == nil {
		t.Fatal("expected nil error, got nil")
	}
}

func TestImportBooksDoesNotCorruptExistingStorage(t *testing.T) {
	// 1. Инициализируем новое чистое хранилище
	storage := NewMemoryStorage()

	// 2. Добавляем валидную книгу. Она должна получить ID = 1
	originalBook, err := storage.AddBook("Clean Code", "Robert C. Martin", 2008)
	if err != nil {
		t.Fatalf("Не удалось добавить исходную книгу: %v", err)
	}

	// 3. Подготавливаем невалидные данные для импорта:
	// Первая книга корректна (ID 5), а вторая — с ошибкой (ID 6, невалидный статус "banana")
	invalidImportData := []Book{
		{
			ID:     5,
			Title:  "Dune",
			Author: "Frank Herbert",
			Year:   1965,
			Status: "planned",
			Rating: 0,
		},
		{
			ID:     6,
			Title:  "Invalid Book",
			Author: "Unknown",
			Year:   2026,
			Status: "banana", // Спровоцирует ошибку валидации
			Rating: 0,
		},
	}

	// Запускаем импорт
	err = storage.ImportBooks(invalidImportData)

	// 4. Проверяем, что метод вернул ошибку
	if err == nil {
		t.Error("Ожидалась ошибка импорта из-за статуса 'banana', но ошибка не была возвращена")
	}

	// 5. Проверяем, что исходная книга (ID 1) осталась нетронутой в хранилище
	book1, exists1 := storage.GetBook(1)
	if !exists1 {
		t.Error("Исходная книга с ID 1 пропала из хранилища после неудачного импорта")
	}
	if book1.Title != originalBook.Title {
		t.Errorf("Данные книги ID 1 повредились. Ожидалось '%s', получили '%s'", originalBook.Title, book1.Title)
	}

	// 6. Проверяем, что книга с ID 5 НЕ попала в хранилище (так как весь импорт должен был отмениться)
	_, exists5 := storage.GetBook(5)
	if exists5 {
		t.Error("Книга с ID 5 была сохранена в хранилище, несмотря на сбой импорта (нарушена транзакционность)")
	}

	// 7. Проверяем работу счетчика ID после сбоя.
	// Новая книга должна получить ID = 2 (потому что ID 1 занят, а импорт не должен был сдвинуть nextID)
	newBook, err := storage.AddBook("Refactoring", "Martin Fowler", 1999)
	if err != nil {
		t.Fatalf("Не удалось добавить книгу после сбоя импорта: %v", err)
	}

	if newBook.ID != 2 {
		t.Errorf("Неверное назначение ID после сбоя импорта. Ожидался ID = 2, получили ID = %d", newBook.ID)
	}
}

func TestSaveAndLoadBooks(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "books.json")

	s1, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	book1, err := s1.AddBook("Title1", "Author1", 2001)
	if err != nil {
		t.Fatal(err)
	}

	// ====================
	s2, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	book2, ok := s2.GetBook(book1.ID)
	if !ok {
		t.Fatal("book doesnt exist")
	}

	if book1 != book2 {
		t.Fatal("book1 != book2")
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "books.json")

	s := NewMemoryStorage()
	err := os.WriteFile(filePath, []byte("not valid json"), 0644)
	if err != nil {
		t.Fatal(err)
	}

	err = s.LoadFromFile(filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestLoadMissingFile(t *testing.T) {
	td := t.TempDir()
	filePath := filepath.Join(td, "pens.json")

	s := NewMemoryStorage()

	err := s.LoadFromFile(filePath)
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !os.IsNotExist(err) {
		t.Errorf("expected nil, got err: %v", err)
	}
}
