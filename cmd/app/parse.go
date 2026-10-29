package main

import (
	"fmt"
	"strconv"
)

func parseID(raw string) (int, error) {
	id, err := strconv.Atoi(raw)
	if err != nil {
		return 0, fmt.Errorf("didn't convert to id")
	}

	if id <= 0 {
		return 0, fmt.Errorf("invalid id")
	}

	return id, nil
}
