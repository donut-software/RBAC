package service

import "rbac/internal/tokenmaker"

func (a *RBAC) CreateToken(username string) (string, error) {
	return a.token.CreateToken(username)
}
func (a *RBAC) VerifyToken(token string) (*tokenmaker.Payload, error) {
	return a.token.VerifyToken(token)
}
