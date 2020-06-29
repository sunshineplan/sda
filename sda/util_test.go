package sda

import (
	"reflect"
	"strings"
	"testing"
)

func TestPrecheck(t *testing.T) {
	data := `
3
1
2
a`
	got := precheck(strings.Split(data, "\n"))
	if !reflect.DeepEqual(got, []string{"1", "2", "3", "a"}) {
		t.Error("precheck wrong")
	}
}

func TestContains(t *testing.T) {
	if !contains([]string{"1", "2", "3"}, "1") {
		t.Error("contains wrong")
	}

	if contains([]string{"1", "2", "3"}, "a") {
		t.Error("contains wrong")
	}
}

func TestMakeRange(t *testing.T) {
	got := makeRange(2, 5)
	if !reflect.DeepEqual(got, []string{"2", "3", "4", "5"}) {
		t.Error("makeRange wrong")
	}
}
