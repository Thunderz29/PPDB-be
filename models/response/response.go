package response

type MessageResponse struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Status     string `json:"status"`
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
