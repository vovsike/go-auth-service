package jwtInternal

type Service interface {
	GenerateToken(userId int) ([]byte, error)
	ValidateToken(t []byte) (bool, error)
	ExtractUserID(token []byte) (int, error)
}
