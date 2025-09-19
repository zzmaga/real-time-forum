package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB(filepath string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filepath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		return nil, err
	}
	// keep global for legacy parts if any
	DB = db
	return db, nil
}

func migrate(db *sql.DB) error {
	stmts := []string{
		// users
		`CREATE TABLE IF NOT EXISTS users(
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			nickname NVARCHAR(32) UNIQUE NOT NULL CHECK(LENGTH(nickname) <= 32),
			email NVARCHAR(320) UNIQUE NOT NULL CHECK(LENGTH(email) <= 320),
			password TEXT NOT NULL,
			first_name NVARCHAR(50) NOT NULL,
			last_name NVARCHAR(50) NOT NULL,
			age INTEGER NOT NULL,
			gender NVARCHAR(10) NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		// posts
		`CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			title NVARCHAR(100) NOT NULL CHECK(LENGTH(title) <= 100),
			content TEXT NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
		`CREATE TABLE IF NOT EXISTS sessions (
			uuid TEXT PRIMARY KEY NOT NULL,
			user_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			expired_at DATETIME NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		// categories
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			name NVARCHAR(50) UNIQUE NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);`,
		// post_categories junction table
		`CREATE TABLE IF NOT EXISTS post_categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			post_id INTEGER NOT NULL,
			category_id INTEGER NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY(category_id) REFERENCES categories(id) ON DELETE CASCADE,
			UNIQUE(post_id, category_id)
		);`,
		// post_comments
		`CREATE TABLE IF NOT EXISTS post_comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		// private_messages
		`CREATE TABLE IF NOT EXISTS private_messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			sender_id INTEGER NOT NULL,
			recipient_id INTEGER NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE CASCADE,
			FOREIGN KEY(recipient_id) REFERENCES users(id) ON DELETE CASCADE
		);`,
		// post_votes
		`CREATE TABLE IF NOT EXISTS post_votes (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			vote INTEGER NOT NULL CHECK(vote IN (-1, 1)),
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(post_id, user_id)
		);`,
		// post_comment_votes
		`CREATE TABLE IF NOT EXISTS post_comment_votes (
			id INTEGER PRIMARY KEY AUTOINCREMENT NOT NULL,
			comment_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			vote INTEGER NOT NULL CHECK(vote IN (-1, 1)),
			created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(comment_id) REFERENCES post_comments(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE,
			UNIQUE(comment_id, user_id)
		);`,
	}

	for _, q := range stmts {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	// compatibility adjustments for legacy schemas
	// TODO add updated_at to posts
	ensureColumn(db, "posts", "updated_at", "TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP")
	ensureColumn(db, "users", "first_name", "NVARCHAR(50) DEFAULT ''")
	ensureColumn(db, "users", "last_name", "NVARCHAR(50) DEFAULT ''")
	ensureColumn(db, "users", "age", "INTEGER DEFAULT 0")
	ensureColumn(db, "users", "gender", "NVARCHAR(10) DEFAULT 'other'")
	_, _ = db.Exec(`UPDATE users SET nickname = email WHERE (nickname IS NULL OR nickname = '') AND email IS NOT NULL AND email != ''`)
	return nil
}

func ensureColumn(db *sql.DB, table, column, decl string) {
	rows, err := db.Query("PRAGMA table_info(" + table + ")")
	if err != nil {
		return
	}
	defer rows.Close()
	var (
		cid     int
		name    string
		ctype   string
		notnull int
		dflt    interface{}
		pk      int
	)
	exists := false
	for rows.Next() {
		_ = rows.Scan(&cid, &name, &ctype, &notnull, &dflt, &pk)
		if name == column {
			exists = true
			break
		}
	}
	if !exists {
		_, _ = db.Exec("ALTER TABLE " + table + " ADD COLUMN " + column + " " + decl)
	}
}
