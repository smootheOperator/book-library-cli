package book

import (
	"path/filepath"
	"testing"
)

func TestCalculateStats(t *testing.T) {
	books := []Book{
		{
			Status: "planned",
			Rating: 0,
		},
		{
			Status: "reading",
			Rating: 0,
		},
		{
			Status: "done",
			Rating: 5,
		},
		{
			Status: "done",
			Rating: 4,
		},
		{
			Status: "done",
			Rating: 0,
		},
	}

	library := CalculateStats(books)

	//Total Books
	if library.Total != 5 {
		t.Errorf("want 5, got Total: %d", library.Total)
	}

	// Status
	if library.ByStatus["planned"] != 1 {
		t.Errorf("want 1, got Status-planned: %d", library.ByStatus["planned"])
	}
	if library.ByStatus["reading"] != 1 {
		t.Errorf("want 1, got Status-reading: %d", library.ByStatus["reading"])
	}
	if library.ByStatus["done"] != 3 {
		t.Errorf("want 3, got Status-done: %d", library.ByStatus["done"])
	}
	if library.ByStatus["dropped"] != 0 {
		t.Errorf("want 0, got Status-dropped: %d", library.ByStatus["dropped"])
	}

	// Rating
	if library.RatedCount != 2 {
		t.Errorf("want 2, got RatedCount: %d", library.RatedCount)
	}
	if library.UnratedCount != 3 {
		t.Errorf("want 3, got UnratedCount: %d", library.UnratedCount)
	}
	if library.AverageRating != 4.5 {
		t.Errorf("want 4.5, got AverageRating: %.2f", library.AverageRating)
	}
}

func TestCalculateStatsEmptyLibrary(t *testing.T) {
	books := []Book{}

	library := CalculateStats(books)

	//Total Books
	if library.Total != 0 {
		t.Errorf("want 0, got Total: %d", library.Total)
	}

	// Status
	if library.ByStatus["planned"] != 0 {
		t.Errorf("want 0, got Status-planned: %d", library.ByStatus["planned"])
	}
	if library.ByStatus["reading"] != 0 {
		t.Errorf("want 0, got Status-reading: %d", library.ByStatus["reading"])
	}
	if library.ByStatus["done"] != 0 {
		t.Errorf("want 0, got Status-done: %d", library.ByStatus["done"])
	}
	if library.ByStatus["dropped"] != 0 {
		t.Errorf("want 0, got Status-dropped: %d", library.ByStatus["dropped"])
	}

	// Rating
	if library.RatedCount != 0 {
		t.Errorf("want 0, got RatedCount: %d", library.RatedCount)
	}
	if library.UnratedCount != 0 {
		t.Errorf("want 0, got UnratedCount: %d", library.UnratedCount)
	}
	if library.AverageRating != 0 {
		t.Errorf("want 0, got AverageRating: %.2f", library.AverageRating)
	}
}

func TestCalculateStatsWithoutRatings(t *testing.T) {
	books := []Book{
		{
			Status: "planned",
			Rating: 0,
		},
		{
			Status: "reading",
			Rating: 0,
		},
		{
			Status: "done",
			Rating: 0,
		},
	}

	library := CalculateStats(books)

	//Total Books
	if library.Total != 3 {
		t.Errorf("want 3, got Total: %d", library.Total)
	}

	// Rating
	if library.RatedCount != 0 {
		t.Errorf("want 0, got RatedCount: %d", library.RatedCount)
	}
	if library.UnratedCount != 3 {
		t.Errorf("want 3, got UnratedCount: %d", library.UnratedCount)
	}
	if library.AverageRating != 0 {
		t.Errorf("want 0, got AverageRating: %.2f", library.AverageRating)
	}
}

func TestServicePersistsChange(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "books.json")

	store, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}
	service := NewService(store)

	origbook, err := service.AddBook("Book 1", "Author 1", 1999)
	if err != nil {
		t.Fatal(err)
	}

	err = service.ChangeStatus(origbook.ID, "reading")
	if err != nil {
		t.Fatal(err)
	}

	// Создаём независимое хранилище
	store2 := NewMemoryStorage()

	err = store2.LoadFromFile(filePath)
	if err != nil {
		t.Fatal(err)
	}

	// Получаем книгу по её настоящему ID
	loaded, ok := store2.GetBook(origbook.ID)
	if !ok {
		t.Fatal("book not found after loading")
	}

	// Проверяем, что изменение сохранилось
	if loaded.Status != "reading" {
		t.Errorf("expected status reading, got %q", loaded.Status)
	}

	// Проверяем, что остальные данные сохранились
	if loaded.Title != origbook.Title {
		t.Errorf("expected title %q, got %q", origbook.Title, loaded.Title)
	}
}
