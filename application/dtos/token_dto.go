package dtos

type TokenDto struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}
