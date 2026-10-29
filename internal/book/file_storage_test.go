package book

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestFileStoragePersistsBooks(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "test.json")
	firstStore, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	orig, err := firstStore.AddBook("title", "author", 2000)
	if err != nil {
		t.Fatal(err)
	}

	secontStore, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	loaded, ok := secontStore.GetBook(orig.ID)
	if !ok {
		t.Fatal("expected ok, got unexist")
	}

	if loaded != orig {
		t.Fatalf("%+v != %+v", loaded, orig)
	}
}

func TestFileStoragePersistsChanges(t *testing.T) {
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

	err = firstStore.ChangeRating(orig.ID, 5)
	if err != nil {
		t.Fatal(err)
	}

	secondStore, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	loaded, ok := secondStore.GetBook(orig.ID)
	if !ok {
		t.Fatal("expected true, got false")
	}

	if loaded.Status != "done" || loaded.Rating != 5 {
		t.Fatalf("want done and 5, got status = %v, rating = %d", loaded.Status, loaded.Rating)
	}
}

func TestFileStoragePersistsDelete(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "test.json")
	firstStore, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	orig, err := firstStore.AddBook("title", "author", 2000)
	if err != nil {
		t.Fatal(err)
	}

	err = firstStore.DeleteBook(orig.ID)
	if err != nil {
		t.Fatal(err)
	}

	secondStore, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	_, ok := secondStore.GetBook(orig.ID)
	if ok {
		t.Fatal("book didn't delete (got true), we expected false")
	}
}

func TestFileStorageMissingFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "test.json")
	store, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}
	list := store.ListBooks()

	if len(list) != 0 {
		t.Fatalf("want 0, got %d", len(list))
	}
}

func TestFileStorageInvalidJSON(t *testing.T) {
	text := []byte(`[
		{
			"id": 1,
			"title": "",
			"author": "Frank Herbert",
			"year": 1965,
			"status": "planned",
			"rating": 0
		}
	]`)

	filePath := filepath.Join(t.TempDir(), "test.json")
	err := os.WriteFile(filePath, text, 0644)
	if err != nil {
		t.Fatal(err)
	}

	_, err = NewFileStorage(filePath)
	if err == nil {
		t.Fatal("want error, got nil")
	}
}

func TestFileStorageInvalidUpdateDoesNotOverwriteFile(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "books.json")

	store, err := NewFileStorage(filePath)
	if err != nil {
		t.Fatal(err)
	}

	orig, err := store.AddBook("title", "author", 2000)
	if err != nil {
		t.Fatal(err)
	}

	before, err := os.ReadFile(filePath)

	err = store.ChangeRating(orig.ID, 100)
	if err == nil {
		t.Fatal("want error, got nil")
	}

	after, err := os.ReadFile(filePath)

	ok := bytes.Equal(before, after)
	if !ok {
		t.Fatalf("Book changed: \nbefore: %v\nafter: %v\n", before, after)
	}
}
