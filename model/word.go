package model

import (
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Word struct {
	ID                                           int
	Name, Phonogram, Audio, Explanation, Example string
	CreatedAt, UpdatedAt                         time.Time
}

func (w Word) GetWordById(id int) Word {
	var createdAt, updatedAt mysql.NullTime
	rows, err := db.Query("SELECT name, phonogram, createdAt, updatedAt FROM word WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&w.Name, &w.Phonogram, &createdAt, &updatedAt); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if createdAt.Valid {
		w.CreatedAt = createdAt.Time
	}
	return w
}

func (w Word) GetWordList() Word {
	var createdAt, updatedAt mysql.NullTime
	rows, err := db.Query("SELECT * FROM word")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&w.ID, &w.Name, &w.Phonogram, &w.Audio, &w.Explanation, &w.Example, &createdAt, &updatedAt); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if createdAt.Valid {
		w.CreatedAt = createdAt.Time
	}

	if updatedAt.Valid {
		w.UpdatedAt = updatedAt.Time
	}
	return w
}

func (w Word) CreateWord() (lastID, rowCnt int64) {
	stmt, err := db.Prepare("INSERT INTO word(name, phonogram, audio, explanation, example, createdAt, updatedAt) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(w.Name, w.Phonogram, w.Audio, w.Explanation, w.Example, w.CreatedAt.Format(timeLayout), w.UpdatedAt.Format(timeLayout))
	if err != nil {
		log.Fatal(err)
	}
	lastID, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (w Word) UpdateWord() (rowCnt int64) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}
	defer tx.Rollback()
	stmt, err := tx.Prepare("UPDATE word SET name = ?, phonogram = ?, explanation = ?, example = ?, updatedAt = ? WHERE id = ?")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(w.Name, w.Phonogram, w.Explanation, w.Example, time.Now().Format(timeLayout), w.ID)
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close() // danger!
	rowCnt, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
	}
	return
}
