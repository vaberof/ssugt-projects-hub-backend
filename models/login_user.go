package models

type LoginUserRequestParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginUserResponse struct {
	AccessToken string `json:"accessToken"`
}
