package rest

import (
	"github.com/go-playground/validator/v10"
	"github.com/rx-rz/65ch/internal/data"
	"net/http"
	"strconv"
)

func (api *API) initializeArticleRoutes() {
	api.router.HandlerFunc(http.MethodPost, "/v1/articles", api.publishArticleHandler)
	api.router.HandlerFunc(http.MethodPost, "/v1/articles/draft", api.createDraftHandler)
}

type CreateArticleRequest struct {
	ID         *string `json:"id"`
	AuthorID   string  `json:"author_id" validate:"required"`
	Title      string  `json:"title" validate:"required"`
	Content    string  `json:"content" validate:"required"`
	TagIDs     []int   `json:"tag_ids"`
	CategoryID int     `json:"category_id" validate:"required"`
}

func (api *API) publishArticleHandler(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := api.CreateContext()
	defer cancel()

	var req CreateArticleRequest
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
	_, err = api.models.Users.GetByID(ctx, req.AuthorID)
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	_, err = api.models.Categories.GetByID(strconv.Itoa(req.CategoryID))
	if err != nil {
		api.handleDBError(w, r, err)
		return
	}
	if req.ID != nil {
		_, err = api.models.Articles.Update(ctx, &data.Article{
			ID:         *req.ID,
			Title:      req.Title,
			Content:    req.Content,
			CategoryID: req.CategoryID,
			Status:     "published",
			TagIDs:     req.TagIDs,
		})
	} else {
		_, err = api.models.Articles.Create(ctx, &data.Article{
			AuthorID:   req.AuthorID,
			Title:      req.Title,
			Content:    req.Content,
			CategoryID: req.CategoryID,
			Status:     "published",
			TagIDs:     req.TagIDs,
		})
		if err != nil {
			api.handleDBError(w, r, err)
			return
		}
	}
	api.writeSuccessResponse(w, http.StatusCreated, nil, "Article successfully published")
}

type CreateDraftRequest struct {
	ID         *string `json:"id"`
	AuthorID   string  `json:"author_id" validate:"required"`
	Title      *string `json:"title"`
	Content    *string `json:"content"`
	TagIDs     []int   `json:"tag_ids"`
	CategoryID *int    `json:"category_id"`
}

func (api *API) createDraftHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := api.CreateContext()
	defer cancel()

	var req CreateDraftRequest
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
	article := &data.Article{
		AuthorID: req.AuthorID,
		TagIDs:   req.TagIDs,
	}

	if req.ID != nil {
		article.ID = *req.ID
	}

	if req.Title != nil {
		article.Title = *req.Title
	}
	if req.Content != nil {
		article.Content = *req.Content
	}
	if req.CategoryID != nil {
		article.CategoryID = *req.CategoryID
	}
	article.Status = "draft"
	if req.ID != nil {
		updateInfo, err := api.models.Articles.Update(ctx, article)
		if err != nil {
			api.handleDBError(w, r, err)
			return
		}
		api.writeSuccessResponse(w, http.StatusOK, envelope{"data": updateInfo}, "Draft successfully updated")
	} else {
		_, err := api.models.Articles.Create(ctx, article)
		if err != nil {
			api.handleDBError(w, r, err)
			return
		}
		api.writeSuccessResponse(w, http.StatusOK, nil, "Draft created successfully")
	}
}

func (api *API) deleteArticleHandler(w http.ResponseWriter, r *http.Request) {
	//ctx, cancel := api.CreateContext()
	//defer cancel()

}

func (api *API) getArticleDetailsHandler(w http.ResponseWriter, r *http.Request) {
	//ctx, cancel := api.CreateContext()
	//defer cancel()
}
