package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./emails.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `
	CREATE TABLE IF NOT EXISTS emails (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		message_id TEXT,
		from_addr TEXT,
		to_addr TEXT,
		subject TEXT,
		body TEXT,
		status TEXT,
		retry_count INTEGER DEFAULT 0,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err = db.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
}

func saveEmail(from, to, subject, body string) int64 {
	res, err := db.Exec(`
		INSERT INTO emails(message_id, from_addr, to_addr, subject, body, status)
		VALUES (?, ?, ?, ?, ?, ?)`,
		time.Now().String(), from, to, subject, body, "PENDING",
	)

	if err != nil {
		log.Println("DB insert error:", err)
		return -1
	}

	id, _ := res.LastInsertId()
	return id
}

func updateStatus(id int64, status string, retry int) {
	_, err := db.Exec(`
		UPDATE emails SET status=?, retry_count=? WHERE id=?`,
		status, retry, id,
	)

	if err != nil {
		log.Println("DB update error:", err)
	}
}
