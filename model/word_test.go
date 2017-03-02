package model

import "testing"

// Get not found ID
func TestGetWordByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	i := 10000
	var w Word
	err := w.GetWordByID(int64(i))

	if err != nil {
		t.Error("Get word catch ", err.Error())
	}

	if w.ID == i {
		t.Error("Expected ", i, ", got ", w.ID)
	}
}

// Get max page
func TestGetWordList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var pageNumber, pageSize int64 = 10000, 50000
	var w Word
	words, err := w.GetWordList(pageNumber, pageSize)

	if err != nil {
		t.Error("Get word list catch ", err.Error())
	}

	if len(words) < 1 {
		t.Error("Expected 0, get ", len(words))
	}
}

// Get max page
func TestGetWordList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var pageNumber, pageSize int64 = 10000, 50000
	var w Word
	words, err := w.GetWordList(pageNumber, pageSize)

	if err != nil {
		t.Error("Get word list catch ", err.Error())
	}

	if len(words) < 1 {
		t.Error("Expected 0, get ", len(words))
	}
}
