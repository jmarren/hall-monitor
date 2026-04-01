package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/a-h/templ"
	"github.com/jackc/pgx/v5"
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

	userModel := models.NewUserModel(int32(1))
	user, err := userModel.Get(rw.Context())
	userModel.Delete(rw.Context())
	// user, err := db.Query.GetUserById(rw.Context(), 1)

	if errors.Is(err, pgx.ErrNoRows) {
		fmt.Println("no rows")
	} else {
		panic(err)
	}

	fmt.Printf("user = %v\n", user)
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
