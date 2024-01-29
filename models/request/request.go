package request

import "book-recipe-be-go/models/response"

type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Fullname       string `json:"fullname" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RetypePassword string `json:"retypePassword" binding:"required"`
}

type ToggleFavoriteRequest struct {
	UserId string `json:"userId" binding:"required"`
}

type SignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type CreateRecipeRequest struct {
	UserId     int                   `json:"userId" binding:"required"`
	RecipeName string                `json:"recipeName" binding:"required"`
	Categories response.CategoryInfo `json:"categories" binding:"required"`
	Levels     response.LevelInfo    `json:"levels" binding:"required"`
	TimeCook   int                   `json:"timeCook" binding:"required"`
	Ingridient string                `json:"ingridient" binding:"required"`
	HowToCook  string                `json:"howToCook" binding:"required"`
}

type RecipeFilter struct {
	PageNumber string `form:"pageNumber"`
	PageSize   string `form:"pageSize"`
	RecipeName string `form:"recipeName"`
	LevelID    string `form:"levelId"`
	CategoryID string `form:"categoryId"`
	Time       string `form:"time"`
	SortBy     string `form:"sortBy"`
	UserID     string `form:"userId"`
}
