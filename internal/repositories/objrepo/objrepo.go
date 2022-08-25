package objrepo

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/ptsypyshev/shortlink/internal/models"
)

type Modelable interface {
	*models.User | *models.Link
	GetType() string
	GetList() []interface{}
	Set(m map[string]interface{}) error
}

type Create[T Modelable] interface {
	Create(ctx context.Context, obj T) (int, error)
}

type Read[T Modelable] interface {
	Read(ctx context.Context, id int, obj T) (T, error)
}

type Update[T Modelable] interface {
	Update(ctx context.Context, obj T, newObj T) error
}

type Delete[T Modelable] interface {
	Delete(ctx context.Context, id int) error
}

//type UserSearch interface {
//	Search()
//}

type Storage[T Modelable] interface {
	Create[T]
	Read[T]
	Update[T]
	Delete[T]
	//UserSearch
}

type Users struct {
	store  Storage[*models.User]
	logger *zap.Logger
	//tracer opentracing.Tracer
}

func UsersNew(s Storage[*models.User], l *zap.Logger) *Users {
	return &Users{
		store:  s,
		logger: l,
		//tracer: t,
	}
}

func (u Users) Create(ctx context.Context, user *models.User) (*models.User, error) {
	id, err := u.store.Create(ctx, user)
	if err != nil {
		u.logger.Error(fmt.Sprintf(`cannot read user: %s`, err))
		return nil, fmt.Errorf("cannot create user: %w", err)
	}
	user.ID = id
	return user, nil
}

func (u Users) Read(ctx context.Context, id int) (*models.User, error) {
	user, err := u.store.Read(ctx, id, &models.User{})
	if err != nil {
		u.logger.Error(fmt.Sprintf(`cannot read user: %s`, err))
		return nil, fmt.Errorf("cannot read user: %w", err)
	}
	return user, nil
}

func (u Users) Update(ctx context.Context, id int, updateUser *models.User) (*models.User, error) {
	user, err := u.store.Read(ctx, id, &models.User{})
	if err != nil {
		u.logger.Error(fmt.Sprintf(`cannot find user with id %d: %s`, id, err))
		return nil, fmt.Errorf("cannot find user with id %d: %w", id, err)
	}
	err = u.store.Update(ctx, user, updateUser)
	if err != nil {
		u.logger.Error(fmt.Sprintf(`cannot update user: %s`, err))
		return nil, fmt.Errorf("cannot update user: %w", err)
	}
	return u.store.Read(ctx, id, &models.User{})
}

func (u Users) Delete(ctx context.Context, id int) (*models.User, error) {
	user, err := u.store.Read(ctx, id, &models.User{})
	if err != nil {
		u.logger.Error(fmt.Sprintf(`search user error: %s`, err))
		return nil, fmt.Errorf("search user error: %w", err)
	}
	return user, u.store.Delete(ctx, id)
}

type Links struct {
	store  Storage[*models.Link]
	logger *zap.Logger
	//tracer opentracing.Tracer
}

func LinksNew(s Storage[*models.Link], l *zap.Logger) *Links {
	return &Links{
		store:  s,
		logger: l,
		//tracer: t,
	}
}

func (l Links) Create(ctx context.Context, link *models.Link) (*models.Link, error) {
	id, err := l.store.Create(ctx, link)
	if err != nil {
		l.logger.Error(fmt.Sprintf(`cannot read link: %s`, err))
		return nil, fmt.Errorf("cannot create link: %w", err)
	}
	link.ID = id
	return link, nil
}

func (l Links) Read(ctx context.Context, id int) (*models.Link, error) {
	link, err := l.store.Read(ctx, id, &models.Link{})
	if err != nil {
		l.logger.Error(fmt.Sprintf(`cannot read link: %s`, err))
		return nil, fmt.Errorf("cannot read link: %w", err)
	}
	return link, nil
}

func (l Links) Update(ctx context.Context, id int, updateLink *models.Link) (*models.Link, error) {
	link, err := l.store.Read(ctx, id, &models.Link{})
	if err != nil {
		l.logger.Error(fmt.Sprintf(`cannot find link with id %d: %s`, id, err))
		return nil, fmt.Errorf("cannot find link with id %d: %w", id, err)
	}
	err = l.store.Update(ctx, link, updateLink)
	if err != nil {
		l.logger.Error(fmt.Sprintf(`cannot update link: %s`, err))
		return nil, fmt.Errorf("cannot update link: %w", err)
	}
	return l.store.Read(ctx, id, &models.Link{})
}

func (l Links) Delete(ctx context.Context, id int) (*models.Link, error) {
	link, err := l.store.Read(ctx, id, &models.Link{})
	if err != nil {
		l.logger.Error(fmt.Sprintf(`search link error: %s`, err))
		return nil, fmt.Errorf("search link error: %w", err)
	}
	return link, l.store.Delete(ctx, id)
}
