package application

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hash string, password string) error
}
