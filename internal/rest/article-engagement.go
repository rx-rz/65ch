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
	info, err := api.models.Articles.Like(ctx, req.UserID, req.ArticleID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": info}, "Article liked successfully")

}

func (api *API) unlikeArticleHandler(w http.ResponseWriter, r *http.Request) {
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
	info, err := api.models.Articles.Unlike(ctx, req.UserID, req.ArticleID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": info}, "Article unliked successfully")

}

func (api *API) saveArticleHandler(w http.ResponseWriter, r *http.Request) {
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

	info, err := api.models.Articles.Save(ctx, req.UserID, req.ArticleID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": info}, "Article saved successfully")

}

func (api *API) unsaveArticleHandler(w http.ResponseWriter, r *http.Request) {
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

	info, err := api.models.Articles.Unsave(ctx, req.UserID, req.ArticleID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	api.writeSuccessResponse(w, http.StatusOK, envelope{"data": info}, "Article unsaved successfully")
}

func (api *API) listArticleLikesHandler() {}

func (api *API) viewArticleStatisticsHandler() {}

func (api *API) increaseArticleViewsHandler() {}
