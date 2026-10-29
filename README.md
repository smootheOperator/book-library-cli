## Usage
==================================================
|       Book management:                         |
|               add <title> <author> <year>      |
|               list                             |
|               get <id>                         |
|               delete <id>                      |
|               status <id> <status>             |
|               rating <id> <rating>             |
|               rename <id> <title>              |
|               author <id> <author>             |
|       Queries:                                 |
|               search <title>                   |
|               filter-author <author>           |
|               filter-status <status>           |
|               min-rating <rating>              |
|               sort <id|title|year|rating>      |
|       Statistics:                              |
|               stats                            |
|       Other:                                   |
|               help                             |
|       Valid statuses:                          |
|               planned, reading, done, dropped  |
|       Rating:                                  |
|               0 = unrated                      |
|               1-5 = book rating                |
==================================================

## Run Application
go run ./cmd/app help

## Statistics
    go run ./cmd/app stats

## Search and filtering
    go run ./cmd/app search "go"
    go run ./cmd/app filter-status done
    go run ./cmd/app min-rating 4
    go run ./cmd/app sort year

## Field
id     - Unique book identifier
title  - Non-empty book title
author - Non-empty author name
year   - Positive publication year
status - planned, reading, done, dropped
rating - Integer from 0 to 5; 0 means unrated

## Features
Add, view, update, and delete books.
Search books by title.
Filter books by author, reading status, or minimum rating.
Sort books by ID, title, publication year, or rating.
View library statistics, including reading statuses and average ratings.
Persist library data in a local JSON file.
Use interchangeable storage implementations through a Go interface.
Run unit tests for domain logic, storage, and application services.

## Book Library CLI
A command-line application for managing a personal book library, written in Go.

## Requirements
Go (1.22.2+)
Git (for cloning the repository)

## Installation
git clone <repository-url>
cd book-library-cli

## Architecture
CLI
 |
 v
Service
 |      \
 v       \
Storage   Interface
 |           |
 v           v
Memory      File
Storage     Storage
                |
                v
             Memory
             Storage

## Explanation
Book              - Represents a book and contains validation and state-changing methods.
MemoryStorage     - Stores books in memory using a map.
FileStorage       - Uses MemoryStorage internally and persists changes to a JSON file.
Storage interface - Defines the operations required by Service.
Service           - Provides application operations and coordinates sorting and statistics.
CLI               - Parses command-line arguments and displays results.

## Known Limitations
Data is stored locally in a JSON file.
Book IDs are reconstructed from existing books on startup, so deleted IDs may be reused after a restart.
Concurrent access by multiple application processes is not supported.