package main

import (
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/database"
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/services"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	
	clientRouter := router.Group("/api")
	clientCORS := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return database.IsAllowedOrigin(origin)
		},

		AllowCredentials: true,

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Content-Type",
			"Authorization",
			"X-CSRF-Token",
		},
	}
	clientRouter.Use(cors.New(clientCORS))
	// route client api

	// router.POST("/admin/login", services.AdminLogin)
	adminRouter := router.Group("/admin")
	adminCORS := cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:5173" || origin == "http://127.0.0.1:5173"
		},

		AllowCredentials: true,

		AllowMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},

		AllowHeaders: []string{
			"Content-Type",
			"X-CSRF-Token",
		},
	}
	adminRouter.Use(cors.New(adminCORS))
	// respond to preflight requests for any admin route
	adminRouter.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(200)
	})
	services.RouteUsersAPI(adminRouter)
	services.RouteGroupsAPI(adminRouter)
	services.RouteAppsAPI(adminRouter)

	router.Run()
}