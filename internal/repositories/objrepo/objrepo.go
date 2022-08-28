package objrepo

import (
	"context"
	"fmt"

	"github.com/speps/go-hashids/v2"
	"go.uber.org/zap"

	"github.com/ptsypyshev/shortlink/internal/models"
)

const (
	HashSalt          = "SaltForTheProject2022"
	HashSmallAlphabet = "abcdefghijklmnopqrstuvwxyz1234567890"
	HashMinLength     = 6
)

type Modelable interface {
	*models.User | *models.Link | *models.ShortLink
	GetType() string
	GetList() []interface{}
	Set(m map[string]interface{}) error
	Get() map[string]interface{}
}

type Create[T Modelable] interface {
	Create(ctx context.Context, obj T) (int, error)
}

type Read[T Modelable] interface {
	Read(ctx context.Context, id int, obj T) (T, error)
}

type Search[T Modelable] interface {
	Search(ctx context.Context, field any, value any, obj T) ([]T, error)
}

type Update[T Modelable] interface {
	Update(ctx context.Context, obj T, newObj T) error
}

type Delete[T Modelable] interface {
	Delete(ctx context.Context, id int) error
}

type Check[T Modelable] interface {
	Check(ctx context.Context, obj T) (T, bool)
}

type Storage[T Modelable] interface {
	Create[T]
	Read[T]
	Search[T]
	Update[T]
	Delete[T]
	Check[T]
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

func (u Users) Search(ctx context.Context, field any, value any) ([]*models.User, error) {
	links, err := u.store.Search(ctx, field, value, &models.User{})
	if err != nil {
		u.logger.Error(fmt.Sprintf(`cannot search links: %s`, err))
		return nil, fmt.Errorf("cannot search links: %w", err)
	}
	return links, nil
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

func (u Users) Check(ctx context.Context, checkUser *models.User) (*models.User, bool) {
	user, ok := u.store.Check(ctx, checkUser)
	if !ok {
		u.logger.Error(fmt.Sprintf(`check failed for user: %s`, checkUser.Username))
		return nil, false
	}
	return user, true
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
		l.logger.Error(fmt.Sprintf(`cannot create link: %s`, err))
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

func (l Links) Search(ctx context.Context, field any, value any) ([]*models.Link, error) {
	links, err := l.store.Search(ctx, field, value, &models.Link{})
	if err != nil {
		l.logger.Error(fmt.Sprintf(`cannot search links: %s`, err))
		return nil, fmt.Errorf("cannot search links: %w", err)
	}
	return links, nil
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

type ShortLinks struct {
	store  Storage[*models.ShortLink]
	logger *zap.Logger
	//tracer opentracing.Tracer
}

func ShortLinksNew(s Storage[*models.ShortLink], l *zap.Logger) *ShortLinks {
	return &ShortLinks{
		store:  s,
		logger: l,
		//tracer: t,
	}
}

func (s ShortLinks) Create(ctx context.Context, longLinkID int) (*models.ShortLink, error) {
	token, err := GenerateShortLinkToken(longLinkID)
	if err != nil {
		return nil, err
	}
	shortlink := &models.ShortLink{
		Token:      token,
		LongLinkID: longLinkID,
	}
	id, err := s.store.Create(ctx, shortlink)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`cannot read shortlink: %s`, err))
		return nil, fmt.Errorf("cannot create shortlink: %w", err)
	}
	shortlink.ID = id
	return shortlink, nil
}

func (s ShortLinks) Read(ctx context.Context, id int) (*models.ShortLink, error) {
	shortlink, err := s.store.Read(ctx, id, &models.ShortLink{})
	if err != nil {
		s.logger.Error(fmt.Sprintf(`cannot read shortlink: %s`, err))
		return nil, fmt.Errorf("cannot read shortlink: %w", err)
	}
	return shortlink, nil
}

func (s ShortLinks) Search(ctx context.Context, field any, value any) ([]*models.ShortLink, error) {
	shortLinks, err := s.store.Search(ctx, field, value, &models.ShortLink{})
	if err != nil {
		s.logger.Error(fmt.Sprintf(`cannot search shortLinks: %s`, err))
		return nil, fmt.Errorf("cannot search shortLinks: %w", err)
	}
	return shortLinks, nil
}

func (s ShortLinks) Update(ctx context.Context, id int, updateShortLink *models.ShortLink) (*models.ShortLink, error) {
	shortlink, err := s.store.Read(ctx, id, &models.ShortLink{})
	if err != nil {
		s.logger.Error(fmt.Sprintf(`cannot find shortlink with id %d: %s`, id, err))
		return nil, fmt.Errorf("cannot find shortlink with id %d: %w", id, err)
	}
	err = s.store.Update(ctx, shortlink, updateShortLink)
	if err != nil {
		s.logger.Error(fmt.Sprintf(`cannot update shortlink: %s`, err))
		return nil, fmt.Errorf("cannot update shortlink: %w", err)
	}
	return s.store.Read(ctx, id, &models.ShortLink{})
}

func (s ShortLinks) Delete(ctx context.Context, id int) (*models.ShortLink, error) {
	shortlink, err := s.store.Read(ctx, id, &models.ShortLink{})
	if err != nil {
		s.logger.Error(fmt.Sprintf(`search shortlink error: %s`, err))
		return nil, fmt.Errorf("search shortlink error: %w", err)
	}
	return shortlink, s.store.Delete(ctx, id)
}

func GenerateShortLinkToken(id int) (string, error) {
	hd := hashids.NewData()
	hd.Alphabet = HashSmallAlphabet
	hd.Salt = HashSalt
	hd.MinLength = HashMinLength
	h, err := hashids.NewWithData(hd)
	if err != nil {
		return "", err
	}
	return h.Encode([]int{id})
}
