package main

import (
	"fmt"
	"github.com/JMTyler/homestuck-watcher/internal/db"
	"github.com/JMTyler/homestuck-watcher/internal/fcm"
	"github.com/kataras/iris/v12"
	"os"
)

func main() {
	fmt.Println()
	defer db.CloseDatabase()

	app := iris.Default()

	app.Use(func(ctx iris.Context) {
		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header("Access-Control-Allow-Headers", "*")

		//reqBytes, _ := ctx.GetBody()
		//ctx.GetViewData()
		//var req map[string]interface{}
		//_ = json.Unmarshal(reqBytes, &req)

		ctx.Next()
	})

	app.Options("*", func(ctx iris.Context) {})

	app.Post("/v1/subscribe", func(ctx iris.Context) {
		var req map[string]interface{}
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusUnprocessableEntity)
			//ctx.StopWithError(iris.StatusUnprocessableEntity, err)
			return
		}

		if req["token"] == nil || req["token"] == false || req["token"] == "" {
			ctx.StatusCode(iris.StatusUnprocessableEntity)
			ctx.WriteString("Required field `token` was empty")
			return
		}

		token := req["token"].(string)
		err := fcm.Subscribe("v1", token)
		if err != nil {
			// TODO: Gotta start using log.Fatal() and its ilk.
			fmt.Println(err)
			ctx.StatusCode(iris.StatusInternalServerError)
			// TODO: Problem?
			return
		}

		ctx.JSON(map[string]interface{}{
			"token": token,
		})
		//ctx.StopWithJSON(iris.StatusOK, map[string]interface{}{
		//	"token": token,
		//})
	})

	app.Post("/v1/unsubscribe", func(ctx iris.Context) {
		var req map[string]interface{}
		if err := ctx.ReadJSON(&req); err != nil {
			ctx.StatusCode(iris.StatusUnprocessableEntity)
			//ctx.StopWithError(iris.StatusUnprocessableEntity, err)
			return
		}

		if req["token"] == nil || req["token"] == false || req["token"] == "" {
			ctx.StatusCode(iris.StatusUnprocessableEntity)
			ctx.WriteString("Required field `token` was empty")
			return
		}

		token := req["token"].(string)
		err := fcm.Unsubscribe("v1", token)
		if err != nil {
			// TODO: Gotta start using log.Fatal() and its ilk.
			fmt.Println(err)
			ctx.StatusCode(iris.StatusInternalServerError)
			// TODO: Problem?
			return
		}

		// TODO: manual 200OK?
	})

	app.Get("/v1/stories", func(ctx iris.Context) {
		stories := new(db.Story).FindAll("v1")
		scrubbed := make([]map[string]interface{}, len(stories))
		for i, model := range stories {
			scrubbed[i] = model.Scrub("v1")
		}
		ctx.JSON(scrubbed)
		//ctx.StopWithJSON(iris.StatusOK, scrubbed)
	})

	app.Get("/error", func(ctx iris.Context) {
		ctx.Problem(iris.NewProblem().Status(iris.StatusNotFound).Detail("Bleep Bloop").Key("message", "Bleep Bloop"))
	})

	port, exists := os.LookupEnv("PORT")
	if !exists {
		port = "80"
	}

	app.Listen(":" + port)
}
