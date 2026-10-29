package book

import (
	"errors"
	"sort"
	"strings"
)

func (s *MemoryStorage) SearchByTitle(query string) []Book {
	foundList := []Book{}

	for _, book := range s.books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(query)) {
			foundList = append(foundList, book)
		}
	}
	sort.Slice(foundList, func(i, j int) bool {
		return foundList[i].ID < foundList[j].ID
	})

	return foundList
}

func (s *MemoryStorage) FilterByAuthor(author string) []Book {
	foundList := []Book{}

	for _, book := range s.books {
		if strings.ToLower(book.Author) == strings.ToLower(author) {
			foundList = append(foundList, book)
		}
	}
	sort.Slice(foundList, func(i, j int) bool {
		return foundList[i].ID < foundList[j].ID
	})

	return foundList
}

func (s *MemoryStorage) FilterByStatus(status string) ([]Book, error) {
	if err := ValidateStatus(status); err != nil {
		return nil, err
	}

	foundList := []Book{}

	for _, book := range s.books {
		if status == book.Status {
			foundList = append(foundList, book)
		}
	}

	sort.Slice(foundList, func(i, j int) bool {
		return foundList[i].ID < foundList[j].ID
	})

	return foundList, nil
}

func (s *MemoryStorage) FilterByMinRating(minRating int) ([]Book, error) {
	if err := ValidateRating(minRating); err != nil {
		return nil, err
	}

	foundList := []Book{}

	for _, book := range s.books {
		if book.Rating >= minRating {
			foundList = append(foundList, book)
		}
	}

	sort.Slice(foundList, func(i, j int) bool {
		return foundList[i].ID < foundList[j].ID
	})

	return foundList, nil
}

func SortBooks(books []Book, by string) ([]Book, error) {
	// id     → ascending
	// title  → A → Z
	// year   → oldest → newest
	// rating → lowest → highest
	result := make([]Book, len(books))
	copy(result, books)

	switch by {
	case "id":
		sort.Slice(result, func(i, j int) bool {
			return result[i].ID < result[j].ID
		})
		return result, nil

	case "title":
		sort.Slice(result, func(i, j int) bool {
			return strings.ToLower(result[i].Title) < strings.ToLower(result[j].Title)
		})
		return result, nil

	case "year":
		sort.Slice(result, func(i, j int) bool {
			return result[i].Year < result[j].Year
		})
		return result, nil

	case "rating":
		sort.Slice(result, func(i, j int) bool {
			return result[i].Rating < result[j].Rating
		})
		return result, nil

	default:
		return nil, errors.New("such category doesn't exist")
	}
}

func (s *MemoryStorage) FilterByYearRange(from, to int) ([]Book, error) {
	if from > to {
		return nil, errors.New("FilterByYearRange: from > to")
	}
	if from <= 0 {
		return nil, errors.New("FilterByYearRange: from <= 0")
	}
	if to <= 0 {
		return nil, errors.New("FilterByYearRange: to <= 0")
	}
	foundList := []Book{}

	for _, book := range s.ListBooks() {
		if from <= book.Year && book.Year <= to {
			foundList = append(foundList, book)
		}
	}

	return foundList, nil
}
