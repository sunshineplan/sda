package sda

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

func precheck(data []string) (result []string) {
	for _, i := range data {
		if v := strings.TrimSpace(i); v != "" {
			result = append(result, v)
		}
	}
	sort.Strings(result)
	return
}

// Resultf formart result
func Resultf(result []string, title1, title2 string, elapsed time.Duration) (output []string) {
	output = append(output, fmt.Sprintf("%s\n\nTotal %d record(s)\n", title1, len(result)))
	output = append(output, title2)
	output = append(output, result...)
	if elapsed != -1 {
		output = append(output, fmt.Sprintf("\nDuration for process: %s.", elapsed))
	}
	return
}

func contains(s []string, v string) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}

func makeRange(min, max int) []string {
	s := make([]string, max-min+1)
	for i := range s {
		s[i] = strconv.Itoa(min + i)
	}
	return s
}
