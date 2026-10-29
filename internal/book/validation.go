package book

import (
	"errors"
	"strings"
)

func ValidateTitle(title string) error {
	if strings.TrimSpace(title) == "" {
		return errors.New("title is empty")
	}
	return nil
}

func ValidateAuthor(author string) error {
	if strings.TrimSpace(author) == "" {
		return errors.New("author is empty")
	}
	return nil
}

func ValidateYear(year int) error {
	if year < 0 {
		return errors.New("year is too low")
	}
	return nil
}

func ValidateStatus(status string) error {
	switch status {
	case "planned", "done", "reading", "dropped":
		return nil
	default:
		return errors.New("wrong status")
	}
}

func ValidateRating(rating int) error {
	if !(-1 < rating && rating < 6) {
		return errors.New("wrong rating")
	}
	return nil
}

func ValidateBook(book Book) error {
	if err := ValidateTitle(book.Title); err != nil {
		return err
	}

	if err := ValidateAuthor(book.Author); err != nil {
		return err
	}

	if err := ValidateYear(book.Year); err != nil {
		return err
	}

	if err := ValidateStatus(book.Status); err != nil {
		return err
	}

	if err := ValidateRating(book.Rating); err != nil {
		return err
	}

	return nil
}
