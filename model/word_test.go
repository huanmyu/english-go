package model

import (
	"testing"
	"time"
)

// Get not found ID
func TestGetWordByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	i := 1000000
	w := Word{ID: int64(i)}
	err := w.GetWordByID()

	if err != nil {
		t.Error("Get word catch ", err.Error())
	}

	if w.ID != int64(i) {
		t.Error("Expected ", i, ", got ", w.ID)
	}
}

// Get max page
func TestGetWordList(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	var pageNumber, pageSize int64 = 100000, 500000
	var w Word
	words, err := w.GetWordList(pageNumber, pageSize)

	if err != nil {
		t.Error("Get word list catch exception ", err.Error())
	}

	if len(words) < 1 {
		t.Error("Expected 0, get ", len(words))
	}
}

// Create And Delete word
func TestCreateAndDeleteWord(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	word := Word{
		Name:        "peer",
		Phonogram:   "[pɪr]",
		Audio:       "https://tts.hjapi.com/en-us/9DF3BF531D7AD6B2",
		Explanation: "n. 同等的人；同辈，同事 \nv. 凝视，盯着看",
		Example:     "The WebSocket protocol defines three types of control messages: close, ping and pong. Call the connection WriteControl, WriteMessage or NextWriter methods to send a control message to the peer.",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	lastInsertId, err := word.CreateWord()
	if err != nil {
		t.Error("Create word catch exception ", err.Error())
	}

	if word.Name != "peer" {
		t.Error("Expected peer, get ", word.Name)
	}

	wordDeleted := Word{ID: lastInsertId}
	err = wordDeleted.DeleteWord()
	if err != nil {
		t.Error("Delete word catch exception ", err.Error())
	}
}
