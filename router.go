package main

import (
	"fund/handlers"
	"fund/middleware"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc func(c *gin.Context)
}

type Routers []Router

var apiRouters = Routers{
	Router{
		Name:        "reptile",
		Method:      "POST",
		Pattern:     "/reptile",
		HandlerFunc: handlers.Reptile,
	},
	Router{
		Name:        "reptileData",
		Method:      "GET",
		Pattern:     "/reptile",
		HandlerFunc: handlers.ReptileData,
	},
	Router{
		Name:        "AddFundData",
		Method:      "POST",
		Pattern:     "/fund",
		HandlerFunc: handlers.FundAdd,
	},
	Router{
		Name:        "DelFundData",
		Method:      "POST",
		Pattern:     "/fund/del",
		HandlerFunc: handlers.FundDel,
	},
}

var webRouters = Routers{
	Router{
		Name:        "index",
		Pattern:     "/",
		Method:      "GET",
		HandlerFunc: handlers.Index,
	},
	Router{
		Name:        "reptile",
		Pattern:     "/reptile",
		Method:      "GET",
		HandlerFunc: handlers.ReptileTable,
	},
}

func LoadRouter() (r *gin.Engine) {
	r = gin.Default()
	gin.ForceConsoleColor()
	r.LoadHTMLGlob("view/html/*")
	r.Static("/static/js", "view/js")
	r.Static("/static/css", "view/css")
	r.Static("/image", "view/static")
	r.StaticFile("/favicon.ico", "view/static/favicon.ico")

	for _, router := range webRouters {
		r.Handle(router.Method, router.Pattern, router.HandlerFunc)
	}

	api := r.Group("/api")
	authMiddleware := middleware.AuthMiddleware()
	api.POST("/login", authMiddleware.LoginHandler)
	api.POST("/logout", authMiddleware.LogoutHandler)
	api.POST("/register", func(c *gin.Context) {
		middleware.Register(c)
	})

	// 需要验证token
	auth := api.Use(authMiddleware.MiddlewareFunc())
	{
		for _, router := range apiRouters {
			auth.Handle(router.Method, router.Pattern, router.HandlerFunc)
		}
	}

	return
}
