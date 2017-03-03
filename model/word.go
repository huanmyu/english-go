package model

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

// Word word model
type Word struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Phonogram   string    `json:"phonogram"`
	Audio       string    `json:"audio"`
	Explanation string    `json:"explanation"`
	Example     string    `json:"example"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetWordList find page word
func (w *Word) GetWordList(pageNumber, pageSize int64) (words []Word, err error) {
	var total, offset int64
	var word Word
	err = db.QueryRow("SELECT count(*) as total FROM word").Scan(&total)
	if err != nil {
		return
	}

	pages := total / pageSize
	if total > pages*pageSize {
		pages++
	}

	if pageNumber > pages {
		pageNumber = pages
	}

	if pageNumber <= 1 {
		offset = 0
	} else {
		offset = pageNumber * pageSize
	}

	var createdAt, updatedAt mysql.NullTime
	rows, err := db.Query("SELECT * FROM word limit ?,?", offset, pageSize)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&word.ID, &word.Name, &word.Phonogram, &word.Audio, &word.Explanation, &word.Example, &createdAt, &updatedAt); err != nil {
			return
		}

		if createdAt.Valid {
			word.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			word.UpdatedAt = updatedAt.Time
		}

		words = append(words, word)
	}

	err = rows.Err()
	return
}

// GetWordByID find word by ID
func (w *Word) GetWordByID() error {
	var createdAt, updatedAt mysql.NullTime
	rows, err := db.Query("SELECT * FROM word WHERE id=?", w.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&w.ID, &w.Name, &w.Phonogram, &w.Audio, &w.Explanation, &w.Example, &createdAt, &updatedAt); err != nil {
			return err
		}

		if createdAt.Valid {
			w.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			w.UpdatedAt = updatedAt.Time
		}
	}

	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// CreateWord create word
func (w *Word) CreateWord() (lastInsertID int64, err error) {
	stmt, err := db.Prepare("INSERT INTO word(name, phonogram, audio, explanation, example, created_at, updated_at) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		return
	}

	res, err := stmt.Exec(w.Name, w.Phonogram, w.Audio, w.Explanation, w.Example, w.CreatedAt.Format(timeLayout), w.UpdatedAt.Format(timeLayout))
	if err != nil {
		return
	}
	defer stmt.Close() // danger!

	lastInsertID, err = res.LastInsertId()
	if err != nil {
		return
	}

	_, err = res.RowsAffected()
	return
}

// UpdateWord update word
func (w *Word) UpdateWord() (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE word SET name = ?, phonogram = ?, explanation = ?, example = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return
	}

	res, err := stmt.Exec(w.Name, w.Phonogram, w.Explanation, w.Example, w.UpdatedAt.Format(timeLayout), w.ID)
	if err != nil {
		return
	}
	defer stmt.Close() // danger!

	_, err = res.RowsAffected()
	if err != nil {
		return
	}

	err = tx.Commit()
	return
}

// DeleteWord delete word
func (w *Word) DeleteWord() (err error) {
	stmt, err := db.Prepare("DELETE FROM word where id = ?")
	if err != nil {
		return
	}

	res, err := stmt.Exec(w.ID)
	if err != nil {
		return
	}

	_, err = res.RowsAffected()
	return
}
