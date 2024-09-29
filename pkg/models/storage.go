package models

type RegisterRequest struct {
	Email     string `json:"email" db:"email" default:"your email"`
	Phone     string `json:"phone" db:"phone" default:"+123456789123456"`
	FirstName string `json:"first_name" db:"first_name" default:"Tom"`
	LastName  string `json:"last_name" db:"last_name" default:"Joe"`
	Username  string `json:"username" db:"username" default:"tom0011"`
	Country   string `json:"country" db:"country" default:"Uzbekistan"`
	Password  string `json:"password" db:"password" default:"123456"`
	Bio       string `json:"bio" db:"bio" default:"holasela berish shartmas"`
}
type RegisterRequest1 struct {
	FirstName string `json:"first_name" db:"first_name"`
	LastName  string `json:"last_name" db:"last_name"`
	Email     string `json:"email" db:"email"`
	Password  string `json:"password" db:"password"`
	Phone     string `json:"phone" db:"phone"`
	Username  string `json:"username" db:"username"`
	Country   string `json:"country" db:"country"`
	Bio       string `json:"bio" db:"bio"`
	Code      string `json:"code" binding:"required"`
}

type RegisterResponse struct {
	Id           string `json:"id" db:"id"`
	Email        string `json:"email" db:"email"`
	Flag         string `json:"flag" db:"flag"`
	AccessToken  string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

type LoginEmailRequest struct {
	Email    string `json:"email" db:"email" default:"registerdagi email ni kiritng"`
	Password string `json:"password" db:"password" default:"123456"`
}

type LoginResponse struct {
	Id       string `json:"id" db:"id"`
	Email    string `json:"email" db:"email"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	Role     string `json:"role" db:"role"`
	Country  string `json:"country" db:"country"`
}

type LoginUsernameRequest struct {
	Username string `json:"username" db:"username" default:"admin"`
	Password string `json:"password" db:"password" default:"123321"`
}

type Tokens struct {
	AccessToken  string `json:"access_token" db:"access_token"`
	RefreshToken string `json:"refresh_token" db:"refresh_token"`
}

type Error struct {
	Error string `json:"error" db:"error"`
}
type AcceptCode struct {
	Email string `json:"email" default:"code cogan email ni kiriting"`
	Code  string `json:"code" default:"12369"`
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
	Password string `json:"new_password" default:"123369"`
	Code     string `json:"code"`
}
type UpdatePasswordReq struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}
