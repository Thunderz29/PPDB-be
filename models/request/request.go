package request

type RegisterRequest struct {
	Username       string `json:"username" binding:"required"`
	Fullname       string `json:"fullname" binding:"required"`
	Password       string `json:"password" binding:"required"`
	RetypePassword string `json:"retypePassword" binding:"required"`
}

type ToggleFavoriteRequest struct {
	UserId       uint `json:"userId" binding:"required"`
}
