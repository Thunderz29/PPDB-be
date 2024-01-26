package response

type MessageResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
}

type MyRecipeResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Details     string `json:"details"`
}

type DataResponse struct {
	Total int64         `json:"total"`
	Data  []RecipeEntry `json:"data"`
	Message     string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
}

type RecipeDetailsResponse struct {
	Total int64         `json:"total"`
	Data  interface{} `json:"data"`
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

type RecipeDetailsEntry struct {
	RecipeId   int    `json:"recipeId"`
	Categories CategoryInfo `json:"categories"`
	Levels     LevelInfo    `json:"levels"`
	RecipeName string       `json:"recipeName"`
	ImageFilename   string       `json:"imageFilename"`
	TimeCook       *int          `json:"timeCook"`
	Ingridient   string       `json:"ingridient"`
	HowToCook   string       `json:"howToCook"`
	IsFavorite bool         `json:"isFavorite"`
}

type CategoryListResponse struct {
	Data       []CategoryInfo `json:"data"`
	Message    string        `json:"message"`
	StatusCode int           `json:"statusCode"`
	Status     string        `json:"status"`
}

type CategoryInfo struct {
	CategoryId   int    `json:"categoryId"`
	CategoryName string `json:"categoryName"`
}

type LevelListResponse struct {
	Data       []LevelInfo `json:"data"`
	Message    string        `json:"message"`
	StatusCode int           `json:"statusCode"`
	Status     string        `json:"status"`
}

type LevelInfo struct {
	LevelId   int    `json:"levelId"`
	LevelName string `json:"levelName"`
}

type UserData struct {
	ID       uint   `json:"id"`
	Token    string `json:"token"`
	Type     string `json:"type"`
	Username string `json:"username"`
	Role     string `json:"role"`
}

type SignInResponse struct {
	Data       UserData `json:"data"`
	Message    string   `json:"message"`
	StatusCode int      `json:"statusCode"`
	Status     string   `json:"status"`
}
