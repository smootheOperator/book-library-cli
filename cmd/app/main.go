package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/smootheOperator/book-library-cli/internal/book"
)

const dataFile = "books.json"

func main() {
	fileStore, err := book.NewFileStorage(dataFile)
	if err != nil {
		fmt.Println("Storage error:", err)
		return
	}

	service := book.NewService(fileStore)

	if len(os.Args) < 2 {
		printUsage()
		return
	}

	switch os.Args[1] {
	case "help", "--help":
		printUsage()
	case "list":
		handleList(service, os.Args)
	case "add":
		handleAdd(service, os.Args)
	case "get":
		handleGet(service, os.Args)
	case "delete":
		handleDelete(service, os.Args)
	case "status":
		handleStatus(service, os.Args)
	case "rating":
		handleRating(service, os.Args)
	case "rename":
		handleRename(service, os.Args)
	case "author":
		handleAuthor(service, os.Args)
	case "search":
		handleSearch(service, os.Args)
	case "filter-author":
		handleFilterAuthor(service, os.Args)
	case "filter-status":
		handleFilterStatus(service, os.Args)
	case "min-rating":
		handleMinRating(service, os.Args)
	case "sort":
		handleSort(service, os.Args)
	case "stats":
		handleStats(service, os.Args)
	default:
		printUsage()
	}
}

func handleList(service *book.Service, args []string) {
	if len(args) != 2 {
		printUsage()
		return
	}

	books := service.ListBooks()
	printBooks(books)
}

func handleAdd(service *book.Service, args []string) {
	if len(args) != 5 {
		printUsage()
		return
	}
	title, author := args[2], args[3]
	year, err := strconv.Atoi(args[4])
	if err != nil {
		printError("string didn't convert to year", err)
		return
	}

	result, err := service.AddBook(title, author, year)
	if err != nil {
		printError("new book didn't add", err)
		return
	}

	fmt.Println("Book added:\n", result.Display())
}

func handleGet(service *book.Service, args []string) {
	if len(args) != 3 {
		printUsage()
		return
	}

	id, err := parseID(args[2])
	if err != nil {
		printError("didn't parse", err)
		return
	}

	book, ok := service.GetBook(id)
	if !ok {
		fmt.Printf("Book with ID %d not found\n", id)
		return
	}

	fmt.Println("FOUND:", book.Display())
}

func handleDelete(service *book.Service, args []string) {
	if len(args) != 3 {
		printUsage()
		return
	}

	id, err := parseID(args[2])
	if err != nil {
		printError("didn't parse", err)
		return
	}

	err = service.DeleteBook(id)
	if err != nil {
		printError("didn't delete", err)
		return
	}

	fmt.Println("COMPLETED")
}

func handleStatus(service *book.Service, args []string) {
	if len(args) != 4 {
		printUsage()
		return
	}

	id, err := parseID(args[2])
	if err != nil {
		printError("string didn't convert to year", err)
		return
	}

	err = service.ChangeStatus(id, args[3])
	if err != nil {
		printError("status didn't change", err)
		return
	}

	fmt.Println("COMPLETED")
}

func handleRating(service *book.Service, args []string) {
	if len(args) != 4 {
		printUsage()
		return
	}

	id, err := parseID(args[2])
	if err != nil {
		printError("string didn't convert to id", err)
		return
	}

	rating, err := strconv.Atoi(args[3])
	if err != nil {
		printError("string didn't convert to rating", err)
		return
	}

	err = service.ChangeRating(id, rating)
	if err != nil {
		printError("rating didn't change", err)
		return
	}

	fmt.Println("COMPLETED")
}

func handleRename(service *book.Service, args []string) {
	if len(args) != 4 {
		printUsage()
		return
	}

	id, err := parseID(args[2])
	if err != nil {
		printError("didn't parse", err)
		return
	}

	err = service.RenameBook(id, args[3])
	if err != nil {
		printError("title didn't change", err)
		return
	}

	fmt.Println("COMPLETED")
}

func handleAuthor(service *book.Service, args []string) {
	if len(args) != 4 {
		printUsage()
		return
	}

	id, err := parseID(args[2])
	if err != nil {
		printError("didn't parse", err)
		return
	}

	err = service.ChangeAuthor(id, args[3])
	if err != nil {
		printError("author didn't change", err)
		return
	}

	fmt.Println("COMPLETED")
}

func printUsage() {
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("|	Book management:                         |")
	fmt.Println("|		add <title> <author> <year>      |")
	fmt.Println("|		list                             |")
	fmt.Println("|		get <id>                         |")
	fmt.Println("|		delete <id>                      |")
	fmt.Println("|		status <id> <status>             |")
	fmt.Println("|		rating <id> <rating>             |")
	fmt.Println("|		rename <id> <title>              |")
	fmt.Println("|		author <id> <author>             |")
	fmt.Println("|	Queries:                                 |")
	fmt.Println("|		search <title>                   |")
	fmt.Println("|		filter-author <author>           |")
	fmt.Println("|		filter-status <status>           |")
	fmt.Println("|		min-rating <rating>              |")
	fmt.Println("|		sort <id|title|year|rating>      |")
	fmt.Println("|	Statistics:                              |")
	fmt.Println("|		stats                            |")
	fmt.Println("|	Other:                                   |")
	fmt.Println("|		help                             |")
	fmt.Println("|	Valid statuses:                          |")
	fmt.Println("|		planned, reading, done, dropped  |")
	fmt.Println("|	Rating:                                  |")
	fmt.Println("|		0 = unrated                      |")
	fmt.Println("|		1-5 = book rating                |")
	fmt.Println(strings.Repeat("=", 50))
}

func printBooks(books []book.Book) {
	if len(books) == 0 {
		fmt.Println("No books found")
		return
	}

	for _, book := range books {
		fmt.Println(book.Display())
	}
}

func handleSearch(service *book.Service, args []string) {
	if len(args) != 3 {
		printUsage()
		return
	}
	if strings.TrimSpace(args[2]) == "" {
		fmt.Println("Error: search query cannot be empty")
		return
	}
	result := service.SearchByTitle(args[2])

	printBooks(result)
}

func handleFilterAuthor(service *book.Service, args []string) {
	if len(args) != 3 {
		printUsage()
		return
	}

	result := service.FilterByAuthor(args[2])

	printBooks(result)
}

func handleFilterStatus(service *book.Service, args []string) {
	if len(args) != 3 {
		printUsage()
		return
	}

	result, err := service.FilterByStatus(args[2])
	if err != nil {
		fmt.Println(err)
	}

	printBooks(result)
}

func handleMinRating(service *book.Service, args []string) {
	if len(args) != 3 {
		printUsage()
		return
	}

	query, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := service.FilterByMinRating(query)
	if err != nil {
		fmt.Println(err)
		return
	}

	printBooks(result)
}

func handleSort(service *book.Service, args []string) {
	if len(args) != 3 {
		printUsage()
		return
	}

	sorted, err := service.SortedBooks(args[2])
	if err != nil {
		fmt.Println(err)
		return
	}

	printBooks(sorted)
}

func handleStats(service *book.Service, args []string) {
	if len(args) != 2 {
		printUsage()
		return
	}

	stats := service.Statistics()

	fmt.Println("<< Library statistics >>")
	fmt.Printf("\nTotal books: %d\n\n", stats.Total)
	fmt.Println("By status:")
	fmt.Printf("    Planned: %d\n", stats.ByStatus["planned"])
	fmt.Printf("    Reading: %d\n", stats.ByStatus["reading"])
	fmt.Printf("    Done: %d\n", stats.ByStatus["done"])
	fmt.Printf("    Dropped: %d\n\n", stats.ByStatus["dropped"])

	fmt.Printf("    Rated books: %d\n", stats.RatedCount)
	fmt.Printf("    Unrated books: %d\n\n", stats.UnratedCount)
	if stats.RatedCount > 0 {
		fmt.Printf("    Average rating: %.2f\n", stats.AverageRating)
	} else {
		fmt.Println("    Average rating: -")
	}
}

func printError(context string, err error) {
	fmt.Fprintf(os.Stderr, "Error: %s: %v\n", context, err)
}
