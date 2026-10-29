package book

type LibraryStats struct {
	Total         int
	ByStatus      map[string]int
	RatedCount    int
	UnratedCount  int
	AverageRating float64
}

func CalculateStats(books []Book) LibraryStats {
	var stats LibraryStats
	stats.ByStatus = make(map[string]int)

	stats.Total = len(books)
	ratedSum := 0

	for _, book := range books {
		stats.ByStatus[book.Status]++

		if book.Rating > 0 {
			stats.RatedCount++
			ratedSum += book.Rating
		} else {
			stats.UnratedCount++
		}
	}

	if stats.RatedCount > 0 {
		stats.AverageRating = float64(ratedSum) / float64(stats.RatedCount)
	}

	return stats
}
