package sda

import (
	"reflect"
	"testing"
)

func TestCheckDuplicates(t *testing.T) {
	got, _ := CheckDuplicates([]string{"1", "2", "3", "1"})
	if !reflect.DeepEqual(got, map[string]int{"1": 2, "2": 1, "3": 1}) {
		t.Error("Check Duplicates wrong")
	}
}
func TestRemoveDuplicates(t *testing.T) {
	got := RemoveDuplicates([]string{"1", "2", "3", "1"})
	if !reflect.DeepEqual(got, []string{"1", "2", "3"}) {
		t.Error("Remove Duplicates wrong")
	}
}

func TestCompare(t *testing.T) {
	data1 := []string{"1", "2", "3"}
	data2 := []string{"1", "2", "3", "4"}

	got := Compare(data1, data2, "diff")
	if len(got) != 0 {
		t.Error("Compare wrong")
	}

	got = Compare(data2, data1, "diff")
	if !reflect.DeepEqual(got, []string{"4"}) {
		t.Error("Compare wrong")
	}

	got = Compare(data1, data2, "comm")
	if !reflect.DeepEqual(got, []string{"1", "2", "3"}) {
		t.Error("Compare wrong")
	}

	got = Compare(data2, data1, "comm")
	if !reflect.DeepEqual(got, []string{"1", "2", "3"}) {
		t.Error("Compare wrong")
	}
}

func TestCheckConsecutive(t *testing.T) {
	got, _ := CheckConsecutive([]string{"5", "1", "3"})
	if !reflect.DeepEqual(got, []string{"2", "4"}) {
		t.Error("Check Consecutive wrong")
	}

	got, _ = CheckConsecutive([]string{"3", "2", "1"})
	if len(got) != 0 {
		t.Error("Check Consecutive wrong")
	}

	_, err := CheckConsecutive([]string{"1", "2", "a"})
	if err == nil {
		t.Error("Check Consecutive wrong")
	}
}
