
type Username userModel[string]
type UserId userModel[int32]

func (u *Username) Posts() UserPostsModel {
	return &userPostsModel{
		&listModel[db.Post, string]{
			ctx:     u.Context(),
			id:      u.id,
			getFunc: db.Query.GetPostsByUsername,
		},
	}
}

func (u *Username) LastPost() (PostModel, error) {
	mostRecentId, err := db.Query.GetMostRecentUserPostByUserName(u.Context(), u.id)
	if err != nil {
		return nil, err
	}

	return NewPostModel(u.Context(), mostRecentId), nil
}

func (u *UserId) Posts() UserPostsModel {
	return &userPostsModel{
		&listModel[db.Post, pgtype.Int4]{
			ctx:        u.Context(),
			id:         newInt4(u.id),
			getFunc:    db.Query.GetPostsByUserId,
			deleteFunc: db.Query.DeletePostsByUserId,
		},
	}
}

func (u *UserId) LastPost() (PostModel, error) {
	mostRecentId, err := db.Query.GetMostRecentUserPostById(u.Context(), newInt4(u.id))

	if err != nil {
		return nil, err
	}

	return NewPostModel(u.Context(), mostRecentId), nil

}

type UsersModel struct {
	context.Context
}

func (u *UsersModel) ByUsername(username string) UserModel {
	baseModel := NewModel(u, username,
		GetFuncOpt(db.Query.GetUserByName),
		DeleteFuncOpt[db.User](db.Query.DeleteUserByName))

	return &Username{
		id:    username,
		Model: baseModel,
	}
}

func (u *UsersModel) ById(ident int32) UserModel {
	baseModel := NewModel(u, ident,
		GetFuncOpt(db.Query.GetUserById),
		DeleteFuncOpt[db.User](db.Query.DeleteUserById))

	return &UserId{
		id:    ident,
		Model: baseModel,
	}
}

