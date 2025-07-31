package controllers

import (
	"gin-rest-api/models"
	"gin-rest-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CommentAPI struct {
	db *gorm.DB
}

func NewCommentAPI(db *gorm.DB) *CommentAPI {
	return &CommentAPI{db}
}

func (a *CommentAPI) Create(c *gin.Context) {
	var req models.CommentAddRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.StatusUnprocessable(c, utils.FormatErrors(errs))
			return
		}

		utils.StatusBadRequest(c, err.Error())
		return
	}

	if !utils.IsExist(a.db, "posts", "id", req.PostId) {
		utils.StatusUnprocessable(c, map[string]any{
			"PostId": "the post does not exist",
		})
		return
	}

	comment := models.Comment{
		UserID: utils.GetUserID(c),
		PostID: req.PostId,
		Body:   req.Body,
	}

	result := a.db.Create(&comment)
	if result.Error != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, comment)
}

func (a *CommentAPI) Get(c *gin.Context) {
	var comment models.Comment

	result := a.db.First(&comment, c.Param("comment_id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	utils.StatusOK(c, comment)
}

func (a *CommentAPI) Update(c *gin.Context) {
	var req models.CommentEditRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			utils.StatusUnprocessable(c, utils.FormatErrors(errs))
			return
		}

		utils.StatusBadRequest(c, err.Error())
		return
	}

	var comment models.Comment

	result := a.db.First(&comment, c.Param("comment_id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	comment.Body = req.Body

	result = a.db.Save(&comment)
	if result.Error != nil {
		utils.StatusServerError(c)
		return
	}

	utils.StatusOK(c, comment)
}

func (a *CommentAPI) Delete(c *gin.Context) {
	var comment models.Comment

	result := a.db.First(&comment, c.Param("comment_id"))
	if err := result.Error; err != nil {
		utils.StatusNotFound(c, err)
		return
	}

	a.db.Delete(&comment)

	utils.StatusOK(c, nil, "the comment has been deleted successfully")
}
