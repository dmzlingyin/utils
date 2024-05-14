package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func DefaultCors() gin.HandlerFunc {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AddAllowHeaders("Authorization")
	config.AddAllowHeaders("X-App-Name")
	config.AddAllowHeaders("X-App-Version")
	return cors.New(config)
}
