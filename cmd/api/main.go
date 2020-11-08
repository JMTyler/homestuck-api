package main

import (
	v1 "github.com/JMTyler/homestuck-watcher"
	"github.com/JMTyler/homestuck-watcher/internal/db"
	"github.com/kataras/iris/v12"
	"os"
)

func main() {
	defer db.CloseDatabase()

	app := iris.Default()

	allowCORS(app)

	v1.AttachRoutes(app.Party("/v1"))

	//app.Get("/error", func(ctx iris.Context) {
	//	ctx.Problem(iris.NewProblem().Status(iris.StatusNotFound).Detail("Bleep Bloop").Key("message", "Bleep Bloop"))
	//})

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}

	app.Listen(":" + port)
}

func allowCORS(app *iris.Application) {
	app.Use(func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")
		ctx.Next()
	})

	app.Options("*", func(ctx iris.Context) {})
}
