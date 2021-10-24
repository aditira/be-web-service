package model

type Register struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type ResUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Roles    string `json:"roles"`
	Type     string `json:"tokenType"`
	Token    string `json:"accessToken"`
}
