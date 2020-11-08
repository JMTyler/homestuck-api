package homestuck_watcher

import (
	"encoding/json"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/router"
)

func AttachRoutes(party router.Party) {
	party.Use(prepare)
	party.Done(respond)

	party.Post("/subscribe", func(ctx iris.Context) {
		service := ctx.Values().Get("service").(*Service)
		token := service.Body["token"]
		// TODO: All these checks probably aren't necessary.
		if token == nil || token == false || token == "" {
			ctx.StatusCode(iris.StatusUnprocessableEntity)
			ctx.Values().Set("res", map[string]string{
				"message": "Required field `token` was empty",
			})
			ctx.Next()
			return
		}

		if err := service.Subscribe(token.(string)); err != nil {
			// TODO: Problem?
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("res", map[string]interface{}{
			"token": token,
		})
		ctx.Next()
	})

	party.Post("/unsubscribe", func(ctx iris.Context) {
		service := ctx.Values().Get("service").(*Service)
		token := service.Body["token"]
		// TODO: All these checks probably aren't necessary.
		if token == nil || token == false || token == "" {
			ctx.StatusCode(iris.StatusUnprocessableEntity)
			ctx.Values().Set("res", map[string]string{
				"message": "Required field `token` was empty",
			})
			ctx.Next()
			return
		}

		if err := service.Unsubscribe(token.(string)); err != nil {
			// TODO: Problem?
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		ctx.StatusCode(iris.StatusOK)
		ctx.Next()
	})

	party.Get("/stories", func(ctx iris.Context) {
		service := ctx.Values().Get("service").(*Service)

		stories, err := service.GetStories()
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			return
		}

		ctx.StatusCode(iris.StatusOK)
		ctx.Values().Set("res", stories)
		ctx.Next()
	})
}

func prepare(ctx iris.Context) {
	// TODO: Replace with version-specific request struct.
	var req map[string]interface{}

	if body, err := ctx.GetBody(); err == nil && len(body) > 0 {
		if err := json.Unmarshal(body, &req); err != nil {
			ctx.StatusCode(iris.StatusUnprocessableEntity)
			return
		}
	}

	ctx.Values().Set("service", &Service{ Body: req })
	ctx.Next()
}

func respond(ctx iris.Context) {
	body := ctx.Values().Get("res")
	if body == nil {
		return
	}

	res, err := json.Marshal(body)
	if err != nil {
		// TODO: Problem?
		ctx.StatusCode(iris.StatusInternalServerError)
		return
	}

	ctx.Write(res)
}
