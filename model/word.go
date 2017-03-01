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

func (w Word) CreateWord() (lastId, rowCnt int64) {
	stmt, err := db.Prepare("INSERT INTO word(name, phonogram, audio, explanation, example, createdAt, updatedAt) VALUES(?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
	}
	res, err := stmt.Exec(w.Name, w.Phonogram, w.Audio, w.Explanation, w.Example, w.CreatedAt.Format("2006-01-02 15:04:05"), w.UpdatedAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		log.Fatal(err)
	}
	lastId, err = res.LastInsertId()
	if err != nil {
		log.Fatal(err)
	}
	rowCnt, err = res.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	return
}

func (w Word) GetWordList() Word {
	var nt mysql.NullTime
	rows, err := db.Query("SELECT name, phonogram, createdAt FROM word WHERE id=?", 1)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&w.Name, &w.Phonogram, &nt); err != nil {
			log.Fatal(err)
		}
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	if nt.Valid {
		w.CreatedAt = nt.Time
	}
	return w
}
