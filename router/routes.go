package router

import (
	"gin-rest-api/config"
	"gin-rest-api/controllers"
	"gin-rest-api/middleware"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func GetRoute(r *gin.Engine, cfg *config.Config, db *gorm.DB, rds *redis.Client) {
	middle := middleware.NewMiddleware(cfg, rds)
	user := controllers.NewUserAPI(cfg, db, rds)
	category := controllers.NewCategoryAPI(db)
	post := controllers.NewPostAPI(db)
	comment := controllers.NewCommentAPI(db)

	// User routes
	r.POST("/api/register", user.Register)
	r.POST("/api/login", user.Login)

	r.Use(middle.CheckAuth)
	r.GET("/api/me", user.Me)
	r.POST("/api/logout", user.Logout)
	userRouter := r.Group("/api/users")
	{
		userRouter.GET("/", user.Gets)
		userRouter.GET("/:id/show", user.Get)
		userRouter.PUT("/:id/update", user.Update)
		userRouter.DELETE("/:id/delete", user.Delete)
		userRouter.GET("/all-trash", user.Trashed)
		userRouter.DELETE("/:id/delete-trash", user.EmptyTrash)
	}

	// Category routes
	catRouter := r.Group("/api/categories")
	{
		catRouter.GET("/", category.Gets)
		catRouter.POST("/create", category.Create)
		catRouter.GET("/:id/show", category.Get)
		catRouter.PUT("/:id/update", category.Update)
		catRouter.DELETE("/:id/delete", category.Delete)
		catRouter.GET("/all-trash", category.Trashed)
		catRouter.DELETE("/:id/delete-trash", category.EmptyTrash)
	}

	// Post routes
	postRouter := r.Group("/api/posts")
	{
		postRouter.GET("/", post.Gets)
		postRouter.POST("/create", post.Create)
		postRouter.GET("/:id/show", post.Get)
		postRouter.PUT("/:id/update", post.Update)
		postRouter.DELETE("/:id/delete", post.Delete)
		postRouter.GET("/all-trash", post.Trashed)
		postRouter.DELETE("/:id/delete-trash", post.EmptyTrash)
	}

	// Comment routes
	commentRouter := r.Group("/api/posts/:id/comment")
	{
		commentRouter.POST("/create", comment.Create)
		commentRouter.GET("/:comment_id/show", comment.Get)
		commentRouter.PUT("/:comment_id/update", comment.Update)
		commentRouter.DELETE("/:comment_id/delete", comment.Delete)
	}
}
