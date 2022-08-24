package objrepo

import (
	"context"
	"fmt"

	"github.com/opentracing/opentracing-go"
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
	Update(ctx context.Context, obj T) (T, error)
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

//func UsersNew(s Storage[*models.User], l *zap.Logger, t opentracing.Tracer) *Users {
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
		//span.LogFields(log.Error(err))
		return nil, fmt.Errorf("cannot read user: %w", err)
	}
	return user, nil
}

type Links struct {
	store  Storage[*models.Link]
	logger *zap.Logger
	tracer opentracing.Tracer
}

func LinksNew(s Storage[*models.Link], l *zap.Logger, t opentracing.Tracer) *Links {
	return &Links{
		store:  s,
		logger: l,
		tracer: t,
	}
}

//type Object[T Modelable] struct {
//	t string
//}

//type Create interface {
//	Create(ctx context.Context, obj Modelable) (int, error)
//}
//
//type Read interface {
//	Read(ctx context.Context, id int) (*Modelable, error)
//}
//
//type Update interface {
//	Update(ctx context.Context, obj Modelable) (*Modelable, error)
//}
//
//type Delete interface {
//	Delete(ctx context.Context, id int) error
//}
//
//type Info interface {
//	Info(ctx context.Context, obj Modelable) error
//}
//
////type UserSearch interface {
////	Search()
////}
//
//type Storage interface {
//	//Create
//	Read
//	ReadUser(ctx context.Context, key int) (interface{}, interface{})
//	//Update
//	//Delete
//	//UserSearch
//	//Info
//}
//
//type Objects[T Modelable] struct {
//	store Storage
//}
//
//func ObjectsNew[T Modelable](s Storage) *Objects[T] {
//	//func ObjectsNew[T Modelable]() *Objects[T] {
//	return &Objects[T]{
//		store: s,
//	}
//}
//
//func (o *Objects[T]) ReadUser(ctx context.Context, obj T, key int) (t T) {
//	panic("aaaaa!")
//}

//
//func (o *Objects[T]) Create(ctx context.Context) (Object[T], error) {
//	id, err := o.store.Create(ctx, T)
//	if err != nil {
//		u.logger.Error(fmt.Sprintf(`cannot read user: %s`, err))
//		span.LogFields(log.Error(err))
//		return nil, fmt.Errorf("cannot create user: %w", err)
//	}
//	user.Id = id
//	span.LogFields(
//		log.String("User result", user.String()),
//	)
//	return &user, nil
//}
//
//func (u Users) Read(ctx context.Context, id int) (*models.User, error) {
//	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, u.tracer,
//		"UserRepo.Read")
//	defer span.Finish()
//	span.LogFields(
//		log.String("id", strconv.Itoa(id)),
//	)
//	user, err := u.us.Read(ctx, id)
//	if err != nil {
//		u.logger.Error(fmt.Sprintf(`cannot read user: %s`, err))
//		span.LogFields(log.Error(err))
//		return nil, fmt.Errorf("cannot read user: %w", err)
//	}
//	span.LogFields(
//		log.String("User result", user.String()),
//	)
//	return user, nil
//}
//
//func (u Users) Update(ctx context.Context, updateUser models.User) (*models.User, error) {
//	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, u.tracer,
//		"UserRepo.Update")
//	defer span.Finish()
//	span.LogFields(
//		log.String("id", strconv.Itoa(updateUser.Id)),
//		log.String("updateUser", updateUser.String()),
//	)
//	user, err := u.us.Update(ctx, updateUser)
//	if err != nil {
//		u.logger.Error(fmt.Sprintf(`cannot update user: %s`, err))
//		span.LogFields(log.Error(err))
//		return nil, fmt.Errorf("cannot update user: %w", err)
//	}
//	return user, nil
//}
//
//func (u Users) Delete(ctx context.Context, id int) (*models.User, error) {
//	span, ctx := opentracing.StartSpanFromContextWithTracer(ctx, u.tracer,
//		"UserRepo.Delete")
//	defer span.Finish()
//	span.LogFields(
//		log.String("id", strconv.Itoa(id)),
//	)
//	user, err := u.us.Read(ctx, id)
//	if err != nil {
//		u.logger.Error(fmt.Sprintf(`cannot read user: %s`, err))
//		span.LogFields(log.Error(err))
//		return nil, fmt.Errorf("cannot read user: %w", err)
//	}
//	span.LogFields(
//		log.String("User delete", user.String()),
//	)
//	return user, u.us.Delete(ctx, id)
//}
