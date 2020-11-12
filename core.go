package sda

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/pmezard/go-difflib/difflib"
)

// CheckDuplicates check slice duplicate value and count
func CheckDuplicates(s []string) (map[string]int, []string) {
	var output []string
	result := make(map[string]int)
	data := precheck(s)
	for _, v := range data {
		result[v] = result[v] + 1
	}

	type dvc struct {
		Value string
		Count int
	}
	var ss []dvc
	for v, c := range result {
		if c > 1 {
			ss = append(ss, dvc{v, c})
		}
	}
	sort.SliceStable(ss, func(i, j int) bool {
		if ss[i].Count == ss[j].Count {
			ni, err1 := strconv.ParseFloat(ss[i].Value, 64)
			nj, err2 := strconv.ParseFloat(ss[j].Value, 64)
			if err1 == nil && err2 == nil {
				return ni < nj
			}
			return ss[i].Value < ss[j].Value
		}
		return ss[i].Count > ss[j].Count
	})

	for _, i := range ss {
		output = append(output, fmt.Sprintf("%s\t\tappears %d times", i.Value, i.Count))
	}
	return result, output
}

// RemoveDuplicates remove duplicate value in slice
func RemoveDuplicates(s []string) []string {
	var output []string
	keys := make(map[string]bool)
	for _, i := range s {
		if _, value := keys[i]; !value {
			keys[i] = true
			output = append(output, i)
		}
	}
	return precheck(output)
}

// Compare two slices common or diff value
func Compare(data1, data2 []string, mode string) (result []string) {
	if mode == "diff" {
		for _, i := range data1 {
			if !contains(data2, i) {
				result = append(result, i)
			}
		}
	} else {
		data1 = RemoveDuplicates(data1)
		for _, i := range data1 {
			if contains(data2, i) {
				result = append(result, i)
			}
		}
	}
	Sort(result)
	return
}

// CheckConsecutive check number slice is consecutive or not
func CheckConsecutive(data []string) (result []string, err error) {
	data = precheck(data)
	var intSlice []int
	var num int
	for _, i := range data {
		num, err = strconv.Atoi(i)
		if err != nil {
			return
		}
		intSlice = append(intSlice, num)
	}
	sort.Ints(intSlice)
	full := makeRange(intSlice[0], intSlice[len(data)-1])
	result = Compare(full, data, "diff")
	return
}

// Diff get unified diff from two slices
func Diff(data1, data2 []string) []string {
	diff := difflib.UnifiedDiff{
		A:        difflib.SplitLines(strings.Join(data1, "\n")),
		B:        difflib.SplitLines(strings.Join(data2, "\n")),
		FromFile: "Data1",
		ToFile:   "Data2",
		Context:  3,
	}
	result, _ := difflib.GetUnifiedDiffString(diff)
	return []string{result}
}
