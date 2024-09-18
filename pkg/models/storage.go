package models

type RegisterRequest struct {
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	Phone       string `json:"phone" db:"phone"`
	Username    string `json:"username" db:"username"`
	Nationality string `json:"nationality" db:"nationality"`
	Bio         string `json:"bio" db:"bio"`
}
type RegisterRequest1 struct {
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password"`
	Phone       string `json:"phone" db:"phone"`
	Username    string `json:"username" db:"username"`
	Nationality string `json:"nationality" db:"nationality"`
	Bio         string `json:"bio" db:"bio"`
	Code        string `json:"code" binding:"required"`
}

type RegisterResponse struct {
	Id          string `json:"id" db:"id"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	Email       string `json:"email" db:"email"`
	Phone       string `json:"phone" db:"phone"`
	Username    string `json:"username" db:"username"`
	Nationality string `json:"nationality" db:"nationality"`
	Bio         string `json:"bio" db:"bio"`
}

type LoginEmailRequest struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type LoginResponse struct {
	Id       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
}

type LoginUsernameRequest struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

type Tokens struct {
	AccessToken  string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

type Error struct {
	Error string `json:"error" db:"error"`
}
type AcceptCode struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}
type GetProfileResponse struct {
	Id        string `json:"id"`
	CreatedAt string `json:"created_at"`
}
type Message struct {
	Message string `json:"message"`
}
type ResetPassReq struct {
	Email    string `json:"email"`
	Password string `json:"new_password"`
	Code     string `json:"code"`
}
type UpdatePasswordReq struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
