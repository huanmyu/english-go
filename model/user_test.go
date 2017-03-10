package model

import "testing"

// Get not found ID
func TestGetUserByID(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	i := 1000000
	u := User{ID: int64(i)}
	err := u.GetUserByID()

	if err != nil {
		t.Error("Get user catch ", err.Error())
	}

	if u.ID != int64(i) {
		t.Error("Expected ", i, ", got ", u.ID)
	}
}
