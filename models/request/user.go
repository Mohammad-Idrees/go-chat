package request

type CreateUserRequest struct {
	Username string  `json:"username" binding:"required"`
	Email    string  `json:"email" binding:"required,isEmail"`
	Password string  `json:"password" binding:"required,min=6,max=10"`
	Phone    *string `json:"phone"`
}

type GetUserByEmailRequest struct {
	Email string `uri:"email" binding:"required"`
}

type GetUserByIdRequest struct {
	Id int64 `uri:"id" binding:"required"`
}

type LoginUserRequest struct {
	Email     string `json:"email" binding:"required"`
	Password  string `json:"password" binding:"required"`
	UserAgent string
	ClientIp  string
}

type LogoutUserRequest struct {
	RefreshToken string `json:"refreshToken" bindind:"required"`
}
