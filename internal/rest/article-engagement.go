package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/rx-rz/65ch/internal/data"
	"net/http"
)

type CommentOnArticleRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	ArticleID string `json:"article_id" validate:"required"`
	Content   string `json:"content" validate:"required"`
}

func (api *API) commentOnArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req CommentOnArticleRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}

	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}

	_, err = api.models.Comments.Create(ctx, &data.Comment{UserID: req.UserID, ArticleID: req.ArticleID, Content: req.Content})
}

type EngageArticleRequest struct {
	UserID    string `json:"user_id" validate:"required"`
	ArticleID string `json:"article_id" validate:"required"`
}

func (api *API) likeArticleHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req EngageArticleRequest
	err := api.readJSON(w, r, &req)
	if err != nil {
		api.badRequestResponse(w, err, err.Error())
		return
	}
	v := validator.New()
	if validationError := v.Struct(req); validationError != nil {
		api.failedValidationResponse(w, validationError)
		return
	}
	err = api.models.Articles.Like(ctx)

}

func unlikeArticleHandler() {}

func saveArticleHandler() {}

func unsaveArticleHandler() {}

func listArticleLikesHandler() {}
