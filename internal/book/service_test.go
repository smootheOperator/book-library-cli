package book

import (
	"path/filepath"
	"testing"
)

func TestServiceSortedBooks(t *testing.T) {
	store := NewMemoryStorage()
	service := NewService(store)

	_, err := service.AddBook("1", "1", 2020)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.AddBook("2", "2", 1999)
	if err != nil {
		t.Fatal(err)
	}

	_, err = service.AddBook("3", "3", 2015)
	if err != nil {
		t.Fatal(err)
	}

	books, err := service.SortedBooks("year")
	if err != nil {
		t.Fatal(err)
	}

	if len(books) != 3 {
		t.Fatalf("expected 3 books, got %d", len(books))
	}

	if books[0].Year != 1999 ||
		books[1].Year != 2015 ||
		books[2].Year != 2020 {
		t.Fatalf("Year's order is not correct")
	}
}

func TestServiceStatistics(t *testing.T) {
	store := MemoryStorage{
		books:  make(map[int]Book),
		nextID: 1,
	}

	b := []Book{
		{
			ID:     1,
			Title:  "Clean Code",
			Author: "Author1",
			Year:   2001,
			Status: "done",
			Rating: 5,
		},
		{
			ID:     2,
			Title:  "Dune",
			Author: "Author2",
			Year:   2002,
			Status: "reading",
			Rating: 0,
		},
		{
			ID:     3,
			Title:  "Learning Go",
			Author: "Author3",
			Year:   2003,
			Status: "done",
			Rating: 4,
		},
	}

	err := store.ImportBooks(b)
	if err != nil {
		t.Fatal(err)
	}

	service := NewService(&store)

	stats := service.Statistics()

	if stats.Total != 3 || stats.RatedCount != 2 || stats.UnratedCount != 1 || stats.AverageRating != 4.5 {
		t.Fatal("results doesn't correct")
	}
}

func TestServicePersistsChanges(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "test.json")
	firstStore, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	orig, err := firstStore.AddBook("title", "author", 2000)
	if err != nil {
		t.Fatal(err)
	}

	err = firstStore.ChangeStatus(orig.ID, "done")
	if err != nil {
		t.Fatal(err)
	}

	secondStore, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	if secondStore.memory.books[orig.ID].Status != "done" {
		t.Fatalf("expected Done, got %v", secondStore.memory.books[orig.ID].Status)
	}
}

func TestServiceWithMemoryStorage(t *testing.T) {
	memory := NewMemoryStorage()
	service := NewService(memory)

	orig, err := service.AddBook("title", "author", 2004)
	if err != nil {
		t.Fatal(err)
	}

	err = service.ChangeStatus(orig.ID, "done")
	if err != nil {
		t.Fatal(err)
	}

	loaded, ok := service.GetBook(orig.ID)
	if !ok {
		t.Fatal("expected true got false")
	}

	if loaded.Status != "done" {
		t.Fatalf("expected Done, got %v", loaded.Status)
	}
}

func TestServiceWithFileStorage(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "books.json")

	fileStore, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	service := NewService(fileStore)

	orig, err := service.AddBook("title", "author", 2004)
	if err != nil {
		t.Fatal(err)
	}

	err = service.ChangeStatus(orig.ID, "done")
	if err != nil {
		t.Fatal(err)
	}

	secondFileStore, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	loaded, ok := secondFileStore.GetBook(orig.ID)
	if !ok {
		t.Fatal("expected true got false")
	}

	if loaded.Status != "done" {
		t.Fatalf("expected Done, got %v", loaded.Status)
	}
}
