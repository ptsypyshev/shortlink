package pgdb

import (
	"context"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/ptsypyshev/shortlink/internal/models"
)

const AllObjects = "all"

type NGDB struct {
	pool *pgxpool.Pool
}

func NGDBNew(p *pgxpool.Pool) *NGDB {
	return &NGDB{
		pool: p,
	}
}

func (n *NGDB) SearchUsers(ctx context.Context, field any, value any) ([]*models.User, error) {
	var (
		rows pgx.Rows
		err  error
	)

	if field == AllObjects {
		query := UserSelectAll
		rows, err = n.pool.Query(ctx, query)
	} else {
		query := UserSelectByField
		rows, err = n.pool.Query(ctx, query, value)
	}

	if err != nil {
		return nil, err
	}

	sliceUsers := make([]*models.User, 0)
	for rows.Next() {
		newlink, err := setUserFieldsNG(rows)
		if err != nil {
			return nil, err
		}
		sliceUsers = append(sliceUsers, &newlink)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sliceUsers, nil
}

func (n *NGDB) SearchLinks(ctx context.Context, field any, value any) ([]*models.Link, error) {
	query := LinkSelectByField

	rows, err := n.pool.Query(ctx, query, value)
	if err != nil {
		return nil, err
	}

	sliceLinks := make([]*models.Link, 0)
	for rows.Next() {
		newlink, err := setLinkFieldsNG(rows)
		if err != nil {
			return nil, err
		}
		sliceLinks = append(sliceLinks, &newlink)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sliceLinks, nil
}

func setUserFieldsNG(rows pgx.Rows) (models.User, error) {
	var (
		id                                                    int
		username, password, firstName, lastName, email, phone string
		userstatus                                            bool
		userStruct                                            models.User
	)
	if err := rows.Scan(&id, &username, &password, &firstName, &lastName, &email, &phone, &userstatus); err != nil {
		return userStruct, err
	}
	mUserFields := map[string]interface{}{
		"id":          id,
		"username":    username,
		"password":    password,
		"first_name":  firstName,
		"last_name":   lastName,
		"email":       email,
		"phone":       phone,
		"user_status": userstatus,
	}

	err := userStruct.Set(mUserFields)
	return userStruct, err
}

func setLinkFieldsNG(rows pgx.Rows) (models.Link, error) {
	var (
		clickCounter, ownerID    int
		longLink, shortLinkToken string
		isActive                 bool
		linkStruct               models.Link
	)
	if err := rows.Scan(&longLink, &shortLinkToken, &clickCounter, &isActive); err != nil {
		return linkStruct, err
	}
	mLinkFields := map[string]interface{}{
		"long_link":     longLink,
		"click_counter": clickCounter,
		"owner_id":      ownerID,
		"is_active":     isActive,
		"short_link":    shortLinkToken,
	}
	err := linkStruct.Set(mLinkFields)
	return linkStruct, err
}
