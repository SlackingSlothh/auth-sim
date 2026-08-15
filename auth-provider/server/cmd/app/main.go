package main

import (
	"github.com/SlackingSlothh/auth-sim/auth-provider/server/internal/services"
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	services.RouteUsersAPI(router)
	services.RouteGroupsAPI(router)
	services.RouteAppsAPI(router)
}