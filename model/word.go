package model

import (
	"fmt"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
)

type Word struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Phonogram   string    `json:"phonogram"`
	Audio       string    `json:"audio"`
	Explanation string    `json:"explanation"`
	Example     string    `json:"example"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (w Word) GetWordById(id int64) Word {
	var createdAt, updatedAt mysql.NullTime
	rows, err := db.Query("SELECT * FROM word WHERE id=?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&w.ID, &w.Name, &w.Phonogram, &w.Audio, &w.Explanation, &w.Example, &createdAt, &updatedAt); err != nil {
			log.Fatal(err)
		}

		if createdAt.Valid {
			w.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			w.UpdatedAt = updatedAt.Time
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

func (w Word) GetWordList(pageNumber, pageSize int64) (words []Word) {
	var total, offset int64
	err := db.QueryRow("SELECT count(*) as total FROM word").Scan(&total)
	if err != nil {
		log.Fatal(err)
	}

	pages := total / pageSize
	if total > pages*pageSize {
		pages += 1
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
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err := rows.Scan(&w.ID, &w.Name, &w.Phonogram, &w.Audio, &w.Explanation, &w.Example, &createdAt, &updatedAt); err != nil {
			log.Fatal(err)
		}

		if createdAt.Valid {
			w.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			w.UpdatedAt = updatedAt.Time
		}
		fmt.Println(w.Name)
		words = append(words, w)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return
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
