package response

type MessageResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
}

type DataResponse struct {
	Total int64         `json:"total"`
	Data  []RecipeEntry `json:"data"`
	Message     string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
}

type RecipeEntry struct {
	RecipeId   int    `json:"recipeId"`
	Categories CategoryInfo `json:"categories"`
	Levels     LevelInfo    `json:"levels"`
	RecipeName string       `json:"recipeName"`
	ImageUrl   string       `json:"imageUrl"`
	Time       *int          `json:"time"`
	IsFavorite bool         `json:"isFavorite"`
}

type CategoryInfo struct {
	CategoryId   int    `json:"categoryId"`
	CategoryName string `json:"categoryName"`
}

type LevelInfo struct {
	LevelId   int    `json:"levelId"`
	LevelName string `json:"levelName"`
}