package models

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jmarren/hall-monitor/internal/db"
)

type Models struct {
	ctx context.Context
}

func NewModels(ctx context.Context) *Models {
	return &Models{
		ctx: ctx,
	}
}

func (m *Models) Users() *UsersModel {
	return &UsersModel{
		m.ctx,
	}
}

func Users(ctx context.Context) *UsersModel {
	return &UsersModel{
		ctx,
	}
}

type GetFunc[T any, U any] func(ctx context.Context, id U) (*T, error)
type GetListFunc[T any, U any] func(ctx context.Context, id U) ([]*T, error)
type UpdateFunc[T any, U any] func(ctx context.Context, id U, data *T) error
type DeleteFunc[U any] func(ctx context.Context, id U) error

type model[T any, U any] struct {
	ctx        context.Context
	id         U
	getFunc    GetFunc[T, U]
	updateFunc UpdateFunc[T, U]
	deleteFunc DeleteFunc[U]
}

type FetchResult[T any] struct {
	result *T
	error  error
}

// func MakeResult[T, U](getFunc GetFunc) {
//
// }

func (m *model[T, U]) Fetch() (*T, error) {
	return m.getFunc(m.ctx, m.id)
}

func (m *model[T, U]) Update(t *T) error {
	return m.updateFunc(m.ctx, m.id, t)
}

func (m *model[T, U]) Delete() error {
	return m.deleteFunc(m.ctx, m.id)
}

func (m *model[T, U]) Context() context.Context {
	return m.ctx
}

type UsernameOrInt4Id interface {
	pgtype.Int4 | string
}

type listModel[T any, U any] struct {
	id         U
	ctx        context.Context
	getFunc    GetListFunc[T, U]
	updateFunc UpdateFunc[T, U]
	deleteFunc DeleteFunc[U]
}

func (l *listModel[T, U]) Fetch() ([]*T, error) {
	return l.getFunc(l.ctx, l.id)
}

func (l *listModel[T, U]) Delete() error {
	return l.deleteFunc(l.ctx, l.id)
}

type userPostsModel struct {
	ListModel[db.Post]
	// *listModel[db.Post, int32]
}

type ListModel[T any] interface {
	Fetch() ([]*T, error)
	Delete() error
}

func newInt4(i int32) pgtype.Int4 {
	return pgtype.Int4{
		Int32: i,
		Valid: true,
	}
}

//
// func NewUserPostsModel[U UsernameOrInt4Id](ctx context.Context, ident U) ListModel[db.Post] {
//
// 	anyIdent := any(ident)
//
// 	var model ListModel[db.Post]
//
// 	switch anyIdent.(type) {
// 	case pgtype.Int4:
// 		idInt4 := anyIdent.(pgtype.Int4)
// 		model = &userPostsModel{
// 			&listModel[db.Post, pgtype.Int4]{
// 				ctx:        ctx,
// 				id:         idInt4,
// 				getFunc:    db.Query.GetPostsByUserId,
// 				deleteFunc: db.Query.DeletePostsByUserId,
// 			},
// 		}
// 	case string:
// 		username := anyIdent.(string)
//
// 		model = &userPostsModel{
// 			&listModel[db.Post, string]{
// 				ctx:     ctx,
// 				id:      username,
// 				getFunc: db.Query.GetPostsByUsername,
// 			},
// 		}
// 	}
//
// 	return model
// }

type ModelOption[T any, U any] func(m *model[T, U])

func NewModel[T any, U any](ctx context.Context, id U, opts ...ModelOption[T, U]) *model[T, U] {
	m := &model[T, U]{
		ctx: ctx,
		id:  id,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func GetFuncOpt[T any, U any](getFunc func(ctx context.Context, id U) (*T, error)) ModelOption[T, U] {
	return func(m *model[T, U]) {
		m.getFunc = getFunc
	}
}

func DeleteFuncOpt[T any, U any](deleteFunc func(ctx context.Context, id U) error) ModelOption[T, U] {
	return func(m *model[T, U]) {
		m.deleteFunc = deleteFunc
	}
}

type UserModel interface {
	Model[db.User]
	LastPost() (PostModel, error)
	Posts() UserPostsModel
}

type UserPostsModel interface {
	ListModel[db.Post]
}

type Model[T any] interface {
	Fetch() (*T, error)
	Update(data *T) error
	Delete() error
	Context() context.Context
}

type IdOrUsername interface {
	string | int32
}

type userModel[T IdOrUsername] struct {
	id T
	Model[db.User]
}

type PostModel interface {
	Model[db.Post]
	Author() (UserModel, error)
}

type postModel struct {
	id int32
	Model[db.Post]
}

// func (u *userModel[IdOrUsername]) Posts() UserPostsModel {
//
// 	var model UserPostsModel
//
// 	idAny := any(u.id)
//
// 	switch idAny.(type) {
// 	case int32:
// 		idInt32 := idAny.(int32)
// 		model = NewUserPostsModel(u.Context(), newInt4(idInt32))
// 	case string:
// 		username := idAny.(string)
// 		model = NewUserPostsModel(u.Context(), username)
// 	}
// 	return model
// }

func (p *postModel) Author() (UserModel, error) {
	authorId, err := db.Query.GetPostAuthor(p.Context(), p.id)

	if err != nil {
		return nil, err
	}

	return NewUserModel(p.Context(), authorId), nil

}

func NewPostModel(ctx context.Context, id int32) PostModel {
	baseModel := NewModel(ctx, id, GetFuncOpt(db.Query.GetPostById), DeleteFuncOpt[db.Post](db.Query.DeletePostById))
	return &postModel{
		id:    id,
		Model: baseModel,
	}
}

// func (u *userModel[IdOrUsername]) LastPost() (PostModel, error) {
//
// 	var mostRecentId int32
// 	var err error
//
// 	switch any(u.id).(type) {
// 	case int32:
// 		id := any(u.id)
// 		idInt32 := id.(int32)
// 		mostRecentId, err = db.Query.GetMostRecentUserPostById(u.Context(), newInt4(idInt32))
//
// 		if err != nil {
// 			return nil, err
// 		}
//
// 	case string:
// 		username := any(u.id).(string)
// 		mostRecentId, err = db.Query.GetMostRecentUserPostByUserName(u.Context(), username)
// 		if err != nil {
// 			return nil, err
// 		}
// 	}
//
// 	return NewPostModel(u.Context(), mostRecentId), nil
// }

func NewUserModel[U IdOrUsername](ctx context.Context, id U) UserModel {

	switch any(id).(type) {
	case int32:
		idInt32 := any(id).(int32)

		baseModel := NewModel(ctx, idInt32,
			GetFuncOpt(db.Query.GetUserById),
			DeleteFuncOpt[db.User](db.Query.DeleteUserById))

		return &UserId{
			id:    idInt32,
			Model: baseModel,
		}
	case string:
		username := any(id).(string)
		baseModel := NewModel(ctx, username,
			GetFuncOpt(db.Query.GetUserByName),
			DeleteFuncOpt[db.User](db.Query.DeleteUserByName))

		return &Username{
			id:    username,
			Model: baseModel,
		}

	default:
		panic("user identifier must be a string or an int32")
	}

}
