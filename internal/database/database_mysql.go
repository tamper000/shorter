//go:build !postgre

package database

import (
	"database/sql"

	"github.com/go-sql-driver/mysql"

	"urlshort/internal/models"
	"urlshort/internal/utils"

	_ "github.com/go-sql-driver/mysql"
)

type Database struct {
	*sql.DB
}

// Setup database
func NewDatabase(MySql string) (*Database, error) {
	db, err := sql.Open("mysql", MySql)
	if err != nil {
		return nil, err
	}

	// check db
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	database := &Database{db}

	if err := database.createTables(); err != nil {
		database.Close()
		return nil, err
	}

	return database, nil
}

// Setup default tables
func (db *Database) createTables() error {
	_, err := db.Exec(`
    CREATE TABLE IF NOT EXISTS urlshort (
		link VARCHAR(255) NOT NULL,
		alias VARCHAR(255) NOT NULL PRIMARY KEY,
		clicks INT NOT NULL DEFAULT 0,
		user VARCHAR(255) NOT NULL
	);
    `)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		username VARCHAR(50) NOT NULL UNIQUE,
		password VARCHAR(255) NOT NULL
	);`)
	return err
}

// Insert new link
func (db *Database) InsertNew(m models.ShortLink, username string) (string, error) {
	var alias string
	var customAlias bool

	if m.Alias == "" {
		const maxAttempts = 10
		for i := 0; i < maxAttempts; i++ {
			alias = utils.GenerateRandomString()
			if !db.сheckAliasExists(alias) {
				break
			}
		}
		if alias == "" {
			return "", ErrFailedGenerate
		}
	} else {
		alias = m.Alias
		customAlias = true
	}

	_, err := db.Exec("INSERT INTO urlshort (link, alias, user) VALUES (?, ?, ?)", m.Link, alias, username)
	if err != nil {
		if mysqlErr, ok := err.(*mysql.MySQLError); ok && mysqlErr.Number == 1062 {
			if customAlias {
				return "", ErrAliasExists
			}
			return db.InsertNew(models.ShortLink{Link: m.Link}, username)
		}
		return "", ErrFailedCreate
	}

	return alias, nil
}

func (db *Database) сheckAliasExists(alias string) bool {
	var count int

	err := db.QueryRow("SELECT COUNT(*) FROM urlshort WHERE alias = ?", alias).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (db *Database) GetByAlias(alias string) (string, error) {
	var link string

	err := db.QueryRow("SELECT link FROM urlshort WHERE alias = ?", alias).Scan(&link)
	return link, err
}

// Add +1 to the number of clicks on the link
func (db *Database) AddClick(alias string) error {
	_, err := db.Exec("UPDATE urlshort SET clicks = clicks + 1 WHERE alias = ?", alias)
	return err
}

func (db *Database) CheckUserExists(username string) (bool, error) {
	var exists bool

	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ? LIMIT 1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *Database) AddUser(username, password string) error {
	_, err := db.Exec(`INSERT INTO users (username, password) VALUES (?, ?)`, username, password)
	return err
}

func (db *Database) GetUser(username string) (string, error) {
	var password string

	err := db.QueryRow("SELECT password FROM users WHERE username = ?", username).Scan(&password)
	return password, err
}

func (db *Database) GetLinksByUser(username string) ([]models.Link, error) {
	rows, err := db.Query("SELECT link, alias, clicks FROM urlshort WHERE user = ?", username)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.Link

	for rows.Next() {
		var record models.Link
		err := rows.Scan(&record.Original, &record.Alias, &record.Clicks)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (db *Database) DeleteLink(alias, username string) (bool, error) {
	result, err := db.Exec("DELETE FROM urlshort WHERE alias = ? AND user = ?;", alias, username)
	affected, _ := result.RowsAffected()
	return affected > 0, err
}

func (db *Database) CloseDatabase() {
	db.Close()
}
