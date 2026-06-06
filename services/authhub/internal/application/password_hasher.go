package application

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	ComparePassword(hash, password string) error
}
