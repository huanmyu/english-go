package model

import (
	"time"

	"github.com/go-sql-driver/mysql"
)

// User user model
type User struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Password   string    `json:"password"`
	IsRemember int       `json:"is_remember"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// GetUserList find page user
func (u *User) GetUserList(pageNumber, pageSize int64) (users []User, err error) {
	var total, offset int64
	var user User
	err = db.QueryRow("SELECT count(*) as total FROM user").Scan(&total)
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
	rows, err := db.Query("SELECT id, name, password, created_at, updated_at FROM user limit ?,?", offset, pageSize)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&user.ID, &user.Name, &user.Password, &createdAt, &updatedAt); err != nil {
			return
		}

		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}

		users = append(users, user)
	}

	err = rows.Err()
	return
}

// GetUserByID find user by ID
func (u *User) GetUserByID() (err error) {
	var createdAt, updatedAt mysql.NullTime
	rows, err := db.Query("SELECT id, name, password, is_remember, created_at, updated_at FROM user WHERE id=?", u.ID)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.Name, &u.Password, &u.IsRemember, &createdAt, &updatedAt); err != nil {
			return
		}

		if createdAt.Valid {
			u.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			u.UpdatedAt = updatedAt.Time
		}
	}

	err = rows.Err()
	return
}

// GetUserByNameAndPassword find user by Name and Password
func (u *User) GetUserByNameAndPassword() (err error) {
	rows, err := db.Query("SELECT id, is_remember FROM user WHERE name=? AND password=? ", u.Name, u.Password)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&u.ID, &u.IsRemember); err != nil {
			return
		}
	}

	err = rows.Err()
	return
}

// CreateUser create user
func (u *User) CreateUser() (lastInsertID int64, err error) {
	stmt, err := db.Prepare("INSERT INTO user(name, password, is_remember, created_at, updated_at) VALUES(?,?,?,?,?)")
	if err != nil {
		return
	}

	res, err := stmt.Exec(u.Name, u.Password, u.IsRemember, u.CreatedAt.Format(timeLayout), u.UpdatedAt.Format(timeLayout))
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

// UpdateUser update user
func (u *User) UpdateUser() (err error) {
	stmt, err := db.Prepare("UPDATE user SET name = ?, password = ?, is_remember = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return
	}

	res, err := stmt.Exec(u.Name, u.Password, u.IsRemember, u.UpdatedAt.Format(timeLayout), u.ID)
	if err != nil {
		return
	}
	defer stmt.Close() // danger!

	_, err = res.RowsAffected()
	return
}

// DeleteUser delete user
func (u *User) DeleteUser() (err error) {
	stmt, err := db.Prepare("DELETE FROM user where id = ?")
	if err != nil {
		return
	}

	res, err := stmt.Exec(u.ID)
	if err != nil {
		return
	}

	_, err = res.RowsAffected()
	return
}
