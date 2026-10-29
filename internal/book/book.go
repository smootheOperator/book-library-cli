package book

import (
	"fmt"
	"strconv"
)

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Status string `json:"status"`
	Rating int    `json:"rating"`
}

func (b Book) Display() string {
	var rating string
	if b.Rating == 0 {
		rating = "-"
	} else {
		rating = strconv.Itoa(b.Rating)
	}

	return fmt.Sprintf("[%d] %s — %s (%d) | %s | rating: %s", b.ID, b.Title, b.Author, b.Year, b.Status, rating)
}

func (b *Book) ChangeStatus(status string) error {
	if err := ValidateStatus(status); err != nil {
		return fmt.Errorf("failed to change status to %q: %w", status, err)
	}

	b.Status = status
	return nil
}

func (b *Book) ChangeRating(rating int) error {
	if err := ValidateRating(rating); err != nil {
		return fmt.Errorf("failed to change rating to %d: %w", rating, err)
	}

	b.Rating = rating
	return nil
}

func (b *Book) Rename(title string) error {
	if err := ValidateTitle(title); err != nil {
		return fmt.Errorf("failed to change title to %q: %w", title, err)
	}

	b.Title = title
	return nil
}

func (b *Book) ChangeAuthor(author string) error {
	if err := ValidateAuthor(author); err != nil {
		return fmt.Errorf("failed to change author to %q: %w", author, err)
	}

	b.Author = author
	return nil
}

func (b Book) IsReading() bool {
	if b.Status == "reading" {
		return true
	}

	return false
}

func (b Book) IsDone() bool {
	if b.Status == "done" {
		return true
	}

	return false
}
