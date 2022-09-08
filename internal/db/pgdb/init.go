package pgdb

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/zapadapter"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/ptsypyshev/shortlink/internal/models"
	"github.com/ptsypyshev/shortlink/internal/repositories/objrepo"
)

const (
	CreateQuery = iota
	ReadQuery
	UpdateQuery
	DeleteQuery
	SearchQuery
	CheckQuery

	EnvVarUserDB     = "DB_USER"
	EnvVarPasswordDB = "DB_PASS"
	EnvVarHostPortDB = "DB_HOST_PORT"
	EnvVarNameDB     = "DB_NAME"

	DefaultUserDB     = "usr"
	DefaultPasswordDB = "pwd"
	DefaultHostPortDB = "localhost:5432"
	DefaultNameDB     = "shortlink"

	UserTable      = "users"
	UserSelectByID = `SELECT * FROM users WHERE id = $1;`
	UserDeleteByID = `DELETE FROM users WHERE id = $1;`
	UserCreate     = `
INSERT INTO users(username, password, first_name, last_name, email, phone, user_status)
VALUES
    ($1, crypt($2, gen_salt('bf', 8)), $3, $4, $5, $6, $7)
RETURNING id;
`
	//UserSelectByField = `SELECT * FROM users WHERE $1 = $2;`
	UserSelectByField = `SELECT * FROM users WHERE username = $1;`
	UserSelectAll     = `SELECT * FROM users ORDER BY id;`
	UserCheckPW       = `SELECT * FROM users WHERE username = $1 AND password = crypt($2, password);`

	LinkTable      = "links"
	LinkSelectByID = `SELECT * FROM links WHERE id = $1;`
	LinkDeleteByID = `DELETE FROM links WHERE id = $1;`
	LinkCreate     = `
INSERT INTO links(long_link, click_counter, owner_id, is_active)
VALUES
    ($1, $2, $3, $4)
RETURNING id;
`

	//LinkSelectByField = `SELECT * FROM links WHERE $1 = $2;`
	LinkSelectByField = `
SELECT long_link, token, click_counter, is_active FROM links, shortlinks 
WHERE links.id = shortlinks.long_link_id AND owner_id = $1 ORDER BY links.id DESC;`

	ShortLinkTable      = "shortlinks"
	ShortLinkSelectByID = `SELECT * FROM shortlinks WHERE id = $1;`
	ShortLinkDeleteByID = `DELETE FROM shortlinks WHERE id = $1;`
	ShortLinkCreate     = `
INSERT INTO shortlinks(token, long_link_id)
VALUES
    ($1, $2)
RETURNING id;
`
	//ShortLinkSelectByField = `SELECT * FROM shortlinks WHERE $1 = $2;`
	ShortLinkSelectByField = `SELECT * FROM shortlinks WHERE token = $1;`

	DatabaseDefaultURL = "postgres://usr:pwd@localhost:5432/shortlink?sslmode=disable"
	InitDBQuery        = `
-- Create some required DB settings (Only first time)
-- Set timezone to Yekaterinburg (GMT+05)
set timezone = 'Asia/Yekaterinburg';
-- Create extension to use cryptography functions in queries
CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- Drop All Tables and Extensions
DROP TABLE IF EXISTS users CASCADE;
DROP TABLE IF EXISTS links CASCADE;
DROP TABLE IF EXISTS shortlinks;
 
-- Create New Tables
CREATE TABLE IF NOT EXISTS users
(
	id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	username VARCHAR(100) NOT NULL UNIQUE,
	password VARCHAR(100) NOT NULL,
	first_name VARCHAR(100),
	last_name VARCHAR(100),
	email VARCHAR(100),
	phone VARCHAR(100),
	user_status BOOL
);

CREATE TABLE IF NOT EXISTS links
(
	id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	long_link TEXT NOT NULL,
	click_counter INT DEFAULT 0 NOT NULL,
	owner_id INT REFERENCES users (id) ON DELETE SET NULL ON UPDATE CASCADE,
	is_active BOOL
);

CREATE TABLE IF NOT EXISTS shortlinks
(
	id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	token VARCHAR(255) NOT NULL UNIQUE,
	long_link_id INT NOT NULL UNIQUE REFERENCES links (id) ON DELETE CASCADE ON UPDATE CASCADE
	
);

-- Insert Administrator user
INSERT INTO users(username, password, first_name, last_name, email, phone, user_status)
VALUES
	('admin', crypt('admin', gen_salt('bf', 8)), 'Administrator', 'TaskSystem', 'admin@example.loc', '111', 'true');
`
	InitDemoQuery = `
-- Insert Users
INSERT INTO users(username, password, first_name, last_name, email, phone, user_status)
VALUES
	('test', crypt('test', gen_salt('bf', 8)), 'Pavel', 'Tsypyshev', 'ptsypyshev@example.loc', '222', 'true'),
	('user', crypt('pass', gen_salt('bf', 8)), 'Vasiliy', 'Pupkin', 'vpupkin@example.loc', '333', 'false'),
	('iivanov', crypt('ivantest', gen_salt('bf', 8)), 'Ivan', 'Ivanov', 'iivanov@example.loc', '444', 'true'),
	('ppetrov', crypt('petrtest', gen_salt('bf', 8)), 'Petr', 'Petrov', 'ppetrov@example.loc', '555', 'true'),
	('ssidorov', crypt('sidrtest', gen_salt('bf', 8)), 'Sidor', 'Sidorov', 'ssidorov@example.loc', '666', 'true');

-- Insert Links
INSERT INTO links(long_link, click_counter, owner_id, is_active)
VALUES
	('https://ya.ru', 100, 2, true),
	('https://mail.ru', 33, 3, true),
	('https://gb.ru', 1, 4, true),
	('https://google.com', 60, 5, true),
	('https://oracle.com', 5, 6, false),
	('https://aws.com', 18, 6, true),
	('https://reg.ru', 7, 5, true),
	('https://timeweb.ru', 23, 4, true),
	('https://ozon.ru', 44, 3, true),
	('https://stackoverflow.com', 58, 2, true);

-- Insert Links
INSERT INTO shortlinks(token, long_link_id)
VALUES
	('p2z68d', 1),
	('08ky2q', 2),
	('429785', 3),
	('z86w2k', 4),
	('l8wxrd', 5),
	('y2ld8p', 6),
	('wrvdr9', 7),
	('q85w8m', 8),
	('6rn32e', 9),
	('yr7grg', 10);
`
)

var (
	_                objrepo.Storage[*models.User] = &DB[*models.User]{}
	_                objrepo.Storage[*models.Link] = &DB[*models.Link]{}
	ErrNotFound                                    = errors.New("not found")
	ErrMultipleFound                               = errors.New("multiple found")
)

func InitDB(ctx context.Context, connectionString string, logger *zap.Logger) (*pgxpool.Pool, error) {
	if connectionString == "" {
		connectionString = DatabaseDefaultURL
	}
	config, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conn string (%s): %w", connectionString, err)
	}
	config.ConnConfig.LogLevel = pgx.LogLevelDebug
	config.ConnConfig.Logger = zapadapter.NewLogger(logger) // логгер запросов в БД
	pool, err := pgxpool.ConnectConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}
	return pool, nil
}

func InitSchema(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, InitDBQuery)
	return err
}

func AddDemoData(ctx context.Context, pool *pgxpool.Pool) error {
	_, err := pool.Exec(ctx, InitDemoQuery)
	return err
}

func MakeConnectionStringFromEnv() string {
	dbUser := getEnv(EnvVarUserDB, DefaultUserDB)
	dbPass := getEnv(EnvVarPasswordDB, DefaultPasswordDB)
	dbHost := getEnv(EnvVarHostPortDB, DefaultHostPortDB)
	dbName := getEnv(EnvVarNameDB, DefaultNameDB)
	return fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", dbUser, dbPass, dbHost, dbName)
}

func getEnv(envVarName, defaultValue string) string {
	result := os.Getenv(envVarName)
	if result == "" {
		result = defaultValue
	}
	return result
}
