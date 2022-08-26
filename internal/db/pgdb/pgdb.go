package pgdb

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgconn"
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

	LinkTable      = "links"
	LinkSelectByID = `SELECT * FROM links WHERE id = $1;`
	LinkDeleteByID = `DELETE FROM links WHERE id = $1;`
	LinkCreate     = `
INSERT INTO links(short_link, long_link, click_counter, owner_id, is_active)
VALUES
    ($1, $2, $3, $4, $5)
RETURNING id;
`

	//LinkSelectByField = `SELECT * FROM links WHERE $1 = $2;`
	LinkSelectByField = `SELECT * FROM links WHERE owner_id = $1;`

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

	//DatabaseURL = "postgres://usr:pwd@postgres:5432/simpleblog?sslmode=disable"
	DatabaseURL = "postgres://usr:pwd@localhost:5432/shortlink?sslmode=disable"
	InitDBQuery = `
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
	short_link VARCHAR(255) NOT NULL UNIQUE,
	long_link TEXT NOT NULL,
	click_counter INT,
	owner_id INT REFERENCES users (id) ON DELETE SET NULL ON UPDATE CASCADE,
	is_active BOOL
);

CREATE TABLE IF NOT EXISTS shortlinks
(
	id INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	token VARCHAR(255) NOT NULL UNIQUE,
	long_link_id INT NOT NULL UNIQUE REFERENCES links (id) ON DELETE CASCADE ON UPDATE CASCADE
	
);
`
	InitDemoQuery = `
-- Insert Users
INSERT INTO users(username, password, first_name, last_name, email, phone, user_status)
VALUES
	('admin', crypt('password', gen_salt('bf', 8)), 'Administrator', 'TaskSystem', 'admin@example.loc', '111', 'true'),
	('ptsypyshev', crypt('testpass', gen_salt('bf', 8)), 'Pavel', 'Tsypyshev', 'ptsypyshev@example.loc', '111', 'true'),
	('vpupkin', crypt('puptest', gen_salt('bf', 8)), 'Vasiliy', 'Pupkin', 'vpupkin@example.loc', '111', 'false'),
	('iivanov', crypt('ivantest', gen_salt('bf', 8)), 'Ivan', 'Ivanov', 'iivanov@example.loc', '111', 'true'),
	('ppetrov', crypt('petrtest', gen_salt('bf', 8)), 'Petr', 'Petrov', 'ppetrov@example.loc', '111', 'true'),
	('ssidorov', crypt('sidrtest', gen_salt('bf', 8)), 'Sidor', 'Sidorov', 'ssidorov@example.loc', '111', 'true');

-- Insert Links
INSERT INTO links(short_link, long_link, click_counter, owner_id, is_active)
VALUES
	('5tng0asdf', 'https://ya.ru', 100, 2, true),
	('gt3sdahu7', 'https://mail.ru', 33, 3, true),
	('9mmsgtt20', 'https://gb.ru', 1, 4, true),
	('n3opkvjn9', 'https://google.com', 60, 5, true),
	('q34ssng91', 'https://oracle.com', 5, 6, false),
	('f82npaq2v', 'https://aws.com', 18, 6, true),
	('x5wst0rq2', 'https://reg.ru', 7, 5, true),
	('tr17bwgh9', 'https://timeweb.ru', 23, 4, true),
	('s4yhrr9qz', 'https://ozon.ru', 44, 3, true),
	('y6tm3rpas', 'https://stackoverflow.com', 58, 2, true);

-- Insert Links
INSERT INTO shortlinks(token, long_link_id)
VALUES
	('5tng0asdf', 1),
	('gt3sdahu7', 2),
	('9mmsgtt20', 3),
	('n3opkvjn9', 4),
	('q34ssng91', 5),
	('f82npaq2v', 6),
	('x5wst0rq2', 7),
	('tr17bwgh9', 8),
	('s4yhrr9qz', 9),
	('y6tm3rpas', 10);
`
)

var (
	_                objrepo.Storage[*models.User] = &DB[*models.User]{}
	_                objrepo.Storage[*models.Link] = &DB[*models.Link]{}
	ErrNotFound                                    = errors.New("not found")
	ErrMultipleFound                               = errors.New("multiple found")
)

func InitDB(ctx context.Context, logger *zap.Logger) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse conn string (%s): %w", DatabaseURL, err)
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

type DB[T objrepo.Modelable] struct {
	pool *pgxpool.Pool
}

func DBNew[T objrepo.Modelable](p *pgxpool.Pool) *DB[T] {
	return &DB[T]{
		pool: p,
	}
}

func (db *DB[T]) Create(ctx context.Context, obj T) (id int, err error) {
	fields := obj.GetList()
	query := switchQuery(obj, CreateQuery)
	res := db.pool.QueryRow(
		ctx, query, fields...,
	)
	err = res.Scan(&id)
	return
}

func (db *DB[T]) Read(ctx context.Context, id int, obj T) (T, error) {
	query := switchQuery(obj, ReadQuery)
	rows, err := db.pool.Query(ctx, query, id)
	if err != nil {
		return nil, err
	}
	return getObjectFromRows(rows, obj, id)
}

func (db *DB[T]) Search(ctx context.Context, field any, value any, obj T) ([]T, error) {
	//var obj T // obj всегда будет nil
	//var obj = &T{} // obj нельзя инициализировать таким образом, будет ошибка
	//Нужно передавать заранее созданный объект через параметры функции

	query := switchQuery(obj, SearchQuery)

	fmt.Printf("Field: %v, value: %v\n", field, value)
	//rows, err := db.pool.Query(ctx, query, field, value)
	rows, err := db.pool.Query(ctx, query, value)

	if err != nil {
		return nil, err
	}
	return getObjectsFromRows(rows, obj)
}

func (db *DB[T]) Update(ctx context.Context, obj T, newObj T) error {
	dbTable := switchQuery(obj, UpdateQuery)
	UpdateQuery, err := updateQueryCompilation(dbTable, obj, newObj)
	if err != nil {
		err = fmt.Errorf("cannot compile query: %w", err)
		return err
	}
	res, err := db.pool.Exec(ctx, UpdateQuery)
	if err != nil {
		return err
	}
	return checkRowsAffected(res, "update", obj)
}

func (db *DB[T]) Delete(ctx context.Context, id int) error {
	var (
		obj T
		err error
	)
	query := switchQuery(obj, DeleteQuery)

	res, err := db.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	return checkRowsAffected(res, "delete", obj)
}

func switchQuery[T objrepo.Modelable](obj T, queryType int) (query string) {
	switch obj.GetType() {
	case models.UserType:
		query = switchUserQuery(queryType)
	case models.LinkType:
		query = switchLinkQuery(queryType)
	case models.ShortLinkType:
		query = switchShortLinkQuery(queryType)
	default:
		panic("Non Modelable type received")
	}
	return
}

func switchUserQuery(queryType int) (query string) {
	switch queryType {
	case CreateQuery:
		query = UserCreate
	case ReadQuery:
		query = UserSelectByID
	case UpdateQuery:
		query = UserTable
	case DeleteQuery:
		query = UserDeleteByID
	case SearchQuery:
		query = UserSelectByField
	default:
		panic("Unknown query type")
	}
	return
}

func switchLinkQuery(queryType int) (query string) {
	switch queryType {
	case CreateQuery:
		query = LinkCreate
	case ReadQuery:
		query = LinkSelectByID
	case UpdateQuery:
		query = LinkTable
	case DeleteQuery:
		query = LinkDeleteByID
	case SearchQuery:
		query = LinkSelectByField
	default:
		panic("Unknown query type")
	}
	return
}

func switchShortLinkQuery(queryType int) (query string) {
	switch queryType {
	case CreateQuery:
		query = ShortLinkCreate
	case ReadQuery:
		query = ShortLinkSelectByID
	case UpdateQuery:
		query = ShortLinkTable
	case DeleteQuery:
		query = ShortLinkDeleteByID
	case SearchQuery:
		query = ShortLinkSelectByField
	default:
		panic("Unknown query type")
	}
	return
}

func getObjectFromRows[T objrepo.Modelable](rows pgx.Rows, obj T, id int) (T, error) {
	var (
		found bool
		err   error
	)
	for rows.Next() {
		if found {
			err = fmt.Errorf("%w: user id %d", ErrMultipleFound, id)
			return nil, err
		}
		obj, err = setObjFields(rows, obj)
		if err != nil {
			return nil, err
		}
		found = true
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if !found {
		err := fmt.Errorf("%w: user id %d", ErrNotFound, id)
		return nil, err
	}
	return obj, nil
}

func getObjectsFromRows[T objrepo.Modelable](rows pgx.Rows, obj T) ([]T, error) {
	var err error
	sliceObjectsT := make([]T, 0)
	for rows.Next() {
		obj, err = setObjFields(rows, obj)
		if err != nil {
			return nil, err
		}
		sliceObjectsT = append(sliceObjectsT, obj)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sliceObjectsT, nil
}

func setObjFields[T objrepo.Modelable](rows pgx.Rows, obj T) (T, error) {
	switch obj.GetType() {
	case models.UserType:
		return setUserFields(rows, obj)
	case models.LinkType:
		return setLinkFields(rows, obj)
	case models.ShortLinkType:
		return setShortLinkFields(rows, obj)
	default:
		panic("Non Modelable type received")
	}
}

func setUserFields[T objrepo.Modelable](rows pgx.Rows, obj T) (T, error) {
	var (
		id                                                    int
		username, password, firstName, lastName, email, phone string
		userstatus                                            bool
	)
	if err := rows.Scan(&id, &username, &password, &firstName, &lastName, &email, &phone, &userstatus); err != nil {
		return nil, err
	}
	mObjFields := map[string]interface{}{
		"id":          id,
		"username":    username,
		"password":    password,
		"first_name":  firstName,
		"last_name":   lastName,
		"email":       email,
		"phone":       phone,
		"user_status": userstatus,
	}
	err := obj.Set(mObjFields)
	return obj, err
}

func setLinkFields[T objrepo.Modelable](rows pgx.Rows, obj T) (T, error) {
	var (
		id, clickCounter, ownerID int
		shortLink, longLink       string
		isActive                  bool
	)
	if err := rows.Scan(&id, &shortLink, &longLink, &clickCounter, &ownerID, &isActive); err != nil {
		return nil, err
	}
	mObjFields := map[string]interface{}{
		"id":            id,
		"short_link":    shortLink,
		"long_link":     longLink,
		"click_counter": clickCounter,
		"owner_id":      ownerID,
		"is_active":     isActive,
	}
	err := obj.Set(mObjFields)
	return obj, err
}

func setShortLinkFields[T objrepo.Modelable](rows pgx.Rows, obj T) (T, error) {
	var (
		id, longLinkID int
		token          string
	)
	if err := rows.Scan(&id, &token, &longLinkID); err != nil {
		return nil, err
	}
	mObjFields := map[string]interface{}{
		"id":           id,
		"token":        token,
		"long_link_id": longLinkID,
	}
	err := obj.Set(mObjFields)
	return obj, err
}

func updateQueryCompilation(dbTable string, obj interface{}, newObj interface{}) (string, error) {
	objMap, err1 := structToMap(obj)
	newObjMap, err2 := structToMap(newObj)
	if err1 != nil || err2 != nil {
		return "", fmt.Errorf("convert object error")
	}
	id, ok := objMap["id"]
	if !ok {
		return "", fmt.Errorf("no id specified for %v", obj)
	}

	var fmtStr string
	fields, values := getChangedFieldsAndValues(objMap, newObjMap)
	if len(values) < 2 {
		fmtStr = "UPDATE %s SET %s = (%s) WHERE id = %.0f;"
	} else {
		fmtStr = "UPDATE %s SET (%s) = (%s) WHERE id = %.0f;"
	}

	query := fmt.Sprintf(
		fmtStr,
		dbTable,
		strings.Join(fields, ","),
		strings.Join(values, ","),
		id,
	)
	return query, nil
}

func structToMap(s interface{}) (m map[string]interface{}, err error) {
	j, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("cannot parse json: %w", err)
	}
	err = json.Unmarshal(j, &m)
	if err != nil {
		return nil, fmt.Errorf("cannot unmarshal to map: %w", err)
	}
	return m, nil
}

func getChangedFieldsAndValues(objMap, newObjMap map[string]interface{}) (fields, values []string) {
	for k, v := range newObjMap {
		if k == "id" {
			continue
		}
		if v != objMap[k] {
			fields = append(fields, k)
			var (
				vStr           string
				strconvBitSize = 64
			)
			switch v := v.(type) {
			case bool:
				vStr = strconv.FormatBool(v)
			case float64:
				vStr = strconv.FormatFloat(v, 'f', 0, strconvBitSize)
			case string:
				vStr = v
			}
			values = append(values, fmt.Sprintf("'%s'", vStr))
		}
	}
	return
}

func checkRowsAffected[T objrepo.Modelable](res pgconn.CommandTag, operation string, obj T) error {
	if rowsAffected := res.RowsAffected(); rowsAffected != 1 {
		err := fmt.Errorf("%s %s error: %d rows affected", operation, obj.GetType(), rowsAffected)
		return err
	}
	return nil
}
