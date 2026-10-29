package book

import (
	"fmt"
	"path/filepath"
	"testing"
)

func checkStorageAddBook(t *testing.T, store Storage) {
	t.Helper()

	created, err := store.AddBook("Dune", "Frank Herbert", 1965)
	if err != nil {
		t.Fatalf("AddBook() error: %v", err)
	}

	if created.ID != 1 {
		t.Errorf("ID = %d, want 1", created.ID)
	}

	if created.Status != "planned" {
		t.Errorf("Status = %q, want planned", created.Status)
	}

	if created.Rating != 0 {
		t.Errorf("Rating = %d, want 0", created.Rating)
	}
}

func TestStorageAddBookContract(t *testing.T) {
	t.Run("MemoryStorage", func(t *testing.T) {
		store := NewMemoryStorage()

		checkStorageAddBook(t, store)
	})

	t.Run("FileStorage", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "books.json")

		store, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}

		checkStorageAddBook(t, store)
	})
}

func TestStorageInvalidAddDoesNotConsumeID(t *testing.T) {
	t.Run("MemoryStorage", func(t *testing.T) {
		store := NewMemoryStorage()

		_, err := store.AddBook("", "author", 2000)
		if err == nil {
			t.Fatal("want error, got nil")
		}

		book, err := store.AddBook("title", "author", 2000)
		if err != nil {
			t.Fatal(err)
		}

		if book.ID != 1 {
			t.Fatalf("want 1, got %d", book.ID)
		}

		books := store.ListBooks()
		if len(books) != 1 {
			t.Fatalf("want 1, got %d", len(books))
		}
	})

	t.Run("FileStorage", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "test.json")
		store, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}

		_, err = store.AddBook("", "author", 2000)
		if err == nil {
			t.Fatal("want error, got nil")
		}

		book, err := store.AddBook("title", "author", 2000)
		if err != nil {
			t.Fatal(err)
		}

		if book.ID != 1 {
			t.Fatalf("want 1, got %d", book.ID)
		}

		books := store.ListBooks()
		if len(books) != 1 {
			t.Fatalf("want 1, got %d", len(books))
		}
	})
}

func TestStorageInvalidUpdatePreservesBook(t *testing.T) {
	t.Run("MemoryStorage", func(t *testing.T) {
		store := NewMemoryStorage()

		orig, err := store.AddBook("title", "author", 2000)
		if err != nil {
			t.Fatal(err)
		}

		orig, ok := store.GetBook(orig.ID)
		if !ok {
			t.Fatal("expected True, got False")
		}

		err = store.ChangeStatus(orig.ID, "banana")
		if err == nil {
			t.Fatal("want error, got nil")
		}

		loaded, ok := store.GetBook(orig.ID)
		if !ok {
			t.Fatal("expected True, got False")
		}

		if loaded != orig {
			t.Fatalf("loaded-%v != orig-%v", loaded, orig)
		}
	})

	t.Run("FileStorage", func(t *testing.T) {
		filePath := filepath.Join(t.TempDir(), "test.json")
		store, err := NewFileStorage(filePath)
		if err != nil {
			t.Fatal(err)
		}

		orig, err := store.AddBook("title", "author", 2000)
		if err != nil {
			t.Fatal(err)
		}

		orig, ok := store.GetBook(orig.ID)
		if !ok {
			t.Fatal("expected True, got False")
		}

		err = store.ChangeStatus(orig.ID, "banana")
		if err == nil {
			t.Fatal("want error, got nil")
		}

		loaded, ok := store.GetBook(orig.ID)
		if !ok {
			t.Fatal("expected True, got False")
		}

		if loaded != orig {
			t.Fatalf("loaded-%v != orig-%v", loaded, orig)
		}
	})
}

func TestValidateRatingInvalid(t *testing.T) {
	tests := []struct {
		name   string
		rating int
	}{
		{"negative one", -1},
		{"negative ten", -10},
		{"above maximum", 6},
		{"large rating", 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateRating(tt.rating)

			if err == nil {
				t.Errorf("ValidateRating(%d): expected error", tt.rating)
			}
		})
	}
}

func TestStorageChangeRatingInvalid(t *testing.T) {
	tests := []struct {
		name   string
		rating int
	}{
		{"Negative", -1},
		{"TooLarge", 6},
		{"veryLarge", 100},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprint("MemoryStorage:", tt.name), func(t *testing.T) {
			store := NewMemoryStorage()

			orig, err := store.AddBook("title", "author", 2000)
			if err != nil {
				t.Fatal(err)
			}

			err = store.ChangeRating(orig.ID, 4)
			if err != nil {
				t.Fatal(err)
			}

			err = store.ChangeRating(orig.ID, tt.rating)
			if err == nil {
				t.Fatal("want error, got nil")
			}

			loaded, ok := store.GetBook(orig.ID)
			if !ok {
				t.Fatal("for some reason book doesn't exist")
			}

			if loaded.Rating != 4 {
				t.Fatalf("expected 4, got %d", loaded.Rating)
			}
		})

		t.Run(fmt.Sprint("FileStorage:", tt.name), func(t *testing.T) {
			filePath := filepath.Join(t.TempDir(), "test.json")
			store, err := NewFileStorage(filePath)
			if err != nil {
				t.Fatal(err)
			}

			orig, err := store.AddBook("title", "author", 2000)
			if err != nil {
				t.Fatal(err)
			}

			err = store.ChangeRating(orig.ID, 4)
			if err != nil {
				t.Fatal(err)
			}

			err = store.ChangeRating(orig.ID, tt.rating)
			if err == nil {
				t.Fatal("want error, got nil")
			}

			loaded, ok := store.GetBook(orig.ID)
			if !ok {
				t.Fatal("for some reason book doesn't exist")
			}

			if loaded.Rating != 4 {
				t.Fatalf("expected 4, got %d", loaded.Rating)
			}
		})
	}
}

func TestStorageDeleteDoesNotReuseID(t *testing.T) {
	t.Run("MemoryStorage", func(t *testing.T) {
		store := NewMemoryStorage()

		_, err := store.AddBook("title1", "author", 2000)
		if err != nil {
			t.Fatal(err)
		}

		orig, err := store.AddBook("title2", "author", 2001)
		if err != nil {
			t.Fatal(err)
		}

		_, err = store.AddBook("title3", "author", 2002)
		if err != nil {
			t.Fatal(err)
		}

		err = store.DeleteBook(orig.ID)
		if err != nil {
			t.Fatal(err)
		}

		if got, ok := store.GetBook(orig.ID); ok {
			t.Fatalf("book was deleted, but got %v", got)
		}

		fourth, err := store.AddBook("title4", "author", 2003)
		if err != nil {
			t.Fatal(err)
		}

		if fourth.ID != 4 {
			t.Fatalf("want 4, got %d", fourth.ID)
		}

		if store.books[1].ID != 1 || store.books[3].ID != 3 || store.books[4].ID != 4 {
			t.Fatalf("order is not correct: %v - %v - %v", store.books[1].ID, store.books[3].ID, store.books[4].ID)
		}
	})
}
