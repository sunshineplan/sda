package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sunshineplan/sda-go/sda"
)

func handler(c *gin.Context) {
	start := time.Now()

	var source string
	var data, data1, data2 []string
	switch source = c.PostForm("source"); source {
	case "Data1", "Data2":
		if source == "Data1" {
			data = strings.Split(c.PostForm("data1"), "\n")
		} else {
			data = strings.Split(c.PostForm("data2"), "\n")
		}
		if len(data) == 1 && data[0] == "" {
			c.String(200, fmt.Sprintf("%s is empty.\nPlease enter something...", source))
			return
		}
	case "Data1,Data2":
		data1 = strings.Split(c.PostForm("data1"), "\n")
		data2 = strings.Split(c.PostForm("data2"), "\n")
		if len(data1) == 1 && data1[0] == "" {
			c.String(200, "Data1 is empty.\nPlease enter something...")
			return
		} else if len(data2) == 1 && data2[0] == "" {
			c.String(200, "Data2 is empty.\nPlease enter something...")
			return
		}
	}

	var result []string
	switch function := c.PostForm("func"); function {
	case "chk_duplicates":
		_, r1 := sda.CheckDuplicates(data)
		if len(r1) == 0 {
			c.String(200, fmt.Sprintf(
				"%s has no duplicate value.\n\nDuration for process: %s.", source, time.Since(start)))
			return
		}
		result = append(result,
			sda.Resultf(r1, fmt.Sprintf(
				"Duplicate values found in %s", source), "result:", time.Since(start))...)
	case "rm_duplicates":
		result = sda.RemoveDuplicates(data)
	case "chk_consecutive":
		r1, err := sda.CheckConsecutive(data)
		if err != nil {
			c.String(200, fmt.Sprintf("Error!\n%s contains non-numeric value. Please check!", source))
			return
		}
		_, tmp := sda.CheckDuplicates(data)
		if len(tmp) > 0 {
			result = append(result, fmt.Sprintf(
				"[Warning]Duplicate values found in %[1]s.\nYou can \"Check Duplicates (%[1]s)\" to check it.\n\n",
				source))
		}
		if len(r1) == 0 {
			c.String(200, fmt.Sprintf(
				"%s%s contains consecutive numbers.\n\nDuration for process: %s.",
				strings.Join(result, "\n"),
				source,
				time.Since(start)))
			return
		}
		result = append(result, sda.Resultf(
			r1, fmt.Sprintf("%s is not consecutive", source), "The following numbers are missing:", -1)...)
	case "compare":
		ignoreDuplicates := c.PostForm("ignore_duplicates")
		if ignoreDuplicates == "true" {
			data1 = sda.RemoveDuplicates(data1)
			data2 = sda.RemoveDuplicates(data2)
		}
		switch mode := c.PostForm("mode"); mode {
		case "comm":
			r1 := sda.Compare(data1, data2, "comm")
			if len(r1) == 0 {
				c.String(200, "Two data contain no common value.\n\nDuration for process: %s.", time.Since(start))
				return
			}
			result = append(result,
				sda.Resultf(r1, "Common values found between two data.", "result:", time.Since(start))...)
		case "diff":
			r1 := sda.Compare(data1, data2, "diff")
			r2 := sda.Compare(data2, data1, "diff")
			switch {
			case len(append(r1, r2...)) == 0:
				c.String(200, "Data1 is same as Data2.\n\nDuration for process: %s.", time.Since(start))
				return
			case len(r1) == 0:
				result = append(result, "Data2 completely contains Data1.\n")
				result = append(result,
					sda.Resultf(r2, "Data2 is more than Data1", "result:", time.Since(start))...)
			case len(r2) == 0:
				result = append(result, "Data1 completely contains Data2.\n")
				result = append(result,
					sda.Resultf(r1, "Data1 is more than Data2", "result:", time.Since(start))...)
			default:
				result = append(result, "Two files have inconsistent content.\n")
				result = append(result,
					sda.Resultf(r1, "Data1 is more than Data2", "result:", -1)...)
				result = append(result, "")
				result = append(result,
					sda.Resultf(r2, "Data2 is more than Data1", "result:", time.Since(start))...)
			}
		}
	case "diff":
		result = sda.Diff(data1, data2)
		result = append(result, fmt.Sprintf("\nDuration for process: %s.", time.Since(start)))
	}

	if len(result) == 0 {
		c.String(400, "")
	} else {
		c.String(200, strings.Join(result, "\n"))
	}
}
