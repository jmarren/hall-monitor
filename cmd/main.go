package main

import (
	"context"
	"fmt"

	"github.com/a-h/templ"
	"github.com/jmarren/hall-monitor/internal/db"
	"github.com/jmarren/hall-monitor/internal/db/models"
	"github.com/jmarren/hall-monitor/internal/views"
	"github.com/jmarren/hypergo"
	"github.com/joho/godotenv"
)

func LogRequest(h hypergo.HandlerFunc) hypergo.HandlerFunc {
	return func(rw *hypergo.RW) {
		fmt.Printf("%s %s\n", rw.Method, rw.URL.Path)
		h(rw)
	}
}

func LoggerOne(h hypergo.HandlerFunc) hypergo.HandlerFunc {
	return func(rw *hypergo.RW) {
		fmt.Println("Logger one")
		h(rw)
	}
}

func LoggerTwo(h hypergo.HandlerFunc) hypergo.HandlerFunc {
	return func(rw *hypergo.RW) {
		fmt.Println("Logger two ")
		h(rw)
	}
}

func User(rw *hypergo.RW) templ.Component {

	ctx := rw.Context()

	models.Users(ctx).ById(1).Posts().Delete()

	// models := models.NewModels(context.Background())
	//
	// models.UserId(1).Fetch()
	//
	// models.Username("john").Posts().Delete()
	// user, err := models.Username("john").Fetch()
	//
	// if err != nil {
	//
	// }
	//
	// fmt.Printf("err: %v\n", err)
	//
	// if user != nil {
	//
	// }
	//
	// userModel := models.NewUserModel(rw.Context(), int32(1))
	// user, err := userModel.Fetch()
	// userModel.Delete()
	//
	// postModel, err := userModel.LastPost()
	// if err != nil {
	// 	panic(err)
	// }
	//
	// authorModel, err := postModel.Author()
	//
	// author, err := authorModel.Fetch()
	//
	// authorPosts := authorModel.Posts()
	//
	// authorPosts.Delete()
	//
	// fmt.Printf("author. = %s\n", author.Name)
	//
	// post, err := postModel.Fetch()
	//
	// fmt.Printf("post.Content = %s ", post.Content)
	//
	// post.Get()
	// user, err := db.Query.GetUserById(rw.Context(), 1)

	// if errors.Is(err, pgx.ErrNoRows) {
	// 	fmt.Println("no rows")
	// } else {
	// 	panic(err)
	// }

	// fmt.Printf("user = %v\n", user)
	//
	// if err != nil {
	// 	panic(err)
	// }

	return views.About()
}

func main() {

	godotenv.Load()

	ctx := context.Background()

	db.InitDb(ctx)

	app := hypergo.New("#content")

	app.Use(LogRequest)
	app.Use(LoggerOne)
	app.Use(LoggerTwo)
	app.HxWrap(hypergo.SimpleWrapper(views.Base))

	app.GetComponent("", hypergo.SimpleComponent(views.Home))
	app.GetComponent("about", hypergo.SimpleComponent(views.About))
	app.GetComponent("user", User)

	// spew.Dump(app)
	app.Listen(":5040")
}
