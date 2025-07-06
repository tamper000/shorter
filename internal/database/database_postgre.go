//go:build postgre

package database

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"urlshort/internal/models"
	"urlshort/internal/utils"
)

type Database struct {
	*pgx.Conn
}

// Setup database
func NewDatabase(postgre string) (*Database, error) {
	conn, err := pgx.Connect(context.Background(), postgre)
	if err != nil {
		return nil, err
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	database := &Database{conn}

	if err := database.createTables(); err != nil {
		conn.Close(context.Background())
		return nil, err
	}

	return database, nil
}

// Setup default tables
func (db *Database) createTables() error {
	_, err := db.Exec(context.Background(), `
    CREATE TABLE IF NOT EXISTS urlshort (
        link VARCHAR(255) NOT NULL,
        alias VARCHAR(255) NOT NULL PRIMARY KEY,
        clicks INT NOT NULL DEFAULT 0,
        username VARCHAR(255) NOT NULL
    );
	`)
	if err != nil {
		return err
	}

	_, err = db.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
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
			if !db.checkAliasExists(alias) {
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

	var pgErr *pgconn.PgError
	_, err := db.Exec(context.Background(), "INSERT INTO urlshort (link, alias, username) VALUES ($1, $2, $3)", m.Link, alias, username)
	if err != nil {
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			if customAlias {
				return "", ErrAliasExists
			}
			return db.InsertNew(models.ShortLink{Link: m.Link}, username)
		}
		return "", ErrFailedCreate
	}

	return alias, nil
}

func (db *Database) checkAliasExists(alias string) bool {
	var count int

	err := db.QueryRow(context.Background(), "SELECT COUNT(*) FROM urlshort WHERE alias = $1", alias).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}

func (db *Database) GetByAlias(alias string) (string, error) {
	var link string

	err := db.QueryRow(context.Background(), "SELECT link FROM urlshort WHERE alias = $1", alias).Scan(&link)
	return link, err
}

// Add +1 to the number of clicks on the link
func (db *Database) AddClick(alias string) error {
	_, err := db.Exec(context.Background(), "UPDATE urlshort SET clicks = clicks + 1 WHERE alias = $1", alias)
	return err
}

func (db *Database) CheckUserExists(username string) (bool, error) {
	var exists bool

	err := db.QueryRow(context.Background(), "SELECT EXISTS (SELECT 1 FROM users WHERE username = $1 LIMIT 1)", username).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *Database) AddUser(username, password string) error {
	_, err := db.Exec(context.Background(), `INSERT INTO users (username, password) VALUES ($1, $2)`, username, password)
	return err
}

func (db *Database) GetUser(username string) (string, error) {
	var password string

	err := db.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1", username).Scan(&password)
	return password, err
}

func (db *Database) GetLinksByUser(username string) ([]models.Link, error) {
	rows, err := db.Query(context.Background(), "SELECT link, alias, clicks FROM urlshort WHERE username = $1", username)
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
	result, err := db.Exec(context.Background(), "DELETE FROM urlshort WHERE alias = $1 AND username = $2;", alias, username)
	fmt.Println(err)
	affected := result.RowsAffected()
	return affected > 0, err
}

func (db *Database) CloseDatabase() {
	db.Close()
}
