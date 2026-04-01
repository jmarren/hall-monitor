package models

import (
	"context"

	"github.com/jmarren/hall-monitor/internal/db"
)

type GetFunc[T any, U any] func(ctx context.Context, id U) (*T, error)
type UpdateFunc[T any, U any] func(ctx context.Context, id U, data *T) error
type DeleteFunc[U any] func(ctx context.Context, id U) error

type Model[T any, U any] struct {
	ctx        context.Context
	id         U
	getFunc    GetFunc[T, U]
	updateFunc UpdateFunc[T, U]
	deleteFunc DeleteFunc[U]
}

func (m *Model[T, U]) Get() (*T, error) {
	return m.getFunc(m.ctx, m.id)
}

func (m *Model[T, U]) Update(t *T) error {
	return m.updateFunc(m.ctx, m.id, t)
}

func (m *Model[T, U]) Delete() error {
	return m.deleteFunc(m.ctx, m.id)
}

// func CreateModel[T any, U any](id U, opts ...ModelOption[T, U]) *Model[T, U] {
// 	// do creation
// 	// id := "123"
// 	return NewModel[T, U](id, opts...)
// }

type ModelOption[T any, U any] func(m *Model[T, U])

func NewModel[T any, U any](ctx context.Context, id U, opts ...ModelOption[T, U]) *Model[T, U] {
	m := &Model[T, U]{
		ctx: ctx,
		id:  id,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func GetFuncOpt[T any, U any](getFunc func(ctx context.Context, id U) (*T, error)) ModelOption[T, U] {
	return func(m *Model[T, U]) {
		m.getFunc = getFunc
	}
}

func DeleteFuncOpt[T any, U any](deleteFunc func(ctx context.Context, id U) error) ModelOption[T, U] {
	return func(m *Model[T, U]) {
		m.deleteFunc = deleteFunc
	}
}

type UserModel interface {
	Get() (*db.User, error)
	Update(user *db.User) error
	Delete() error
}

type IdOrUsername interface {
	string | int32
}

func NewUserModel[U string | int32](ctx context.Context, id U) UserModel {

	switch any(id).(type) {
	case int32:
		idInt32 := any(id).(int32)
		return NewModel(ctx, idInt32,
			GetFuncOpt(db.Query.GetUserById),
			DeleteFuncOpt[db.User](db.Query.DeleteUserById))
	case string:
		username := any(id).(string)
		return NewModel(ctx, username,
			GetFuncOpt(db.Query.GetUserByName),
			DeleteFuncOpt[db.User](db.Query.DeleteUserByName))
	default:
		panic("user identifier must be a string or an int32")
	}

}
