package router

import (
	"gin-rest-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Docs routes
func GetDocs(r *gin.Engine) {
	docs.SwaggerInfo.BasePath = "/api"

	r.GET("swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
