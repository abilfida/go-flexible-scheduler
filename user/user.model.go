package user

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"-"` // Jangan pernah tampilkan password di JSON
}

// HashPassword mengenkripsi password sebelum disimpan
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	user.Password = string(bytes)
	return nil
}

// CheckPassword memvalidasi password yang diberikan
func (user *User) CheckPassword(providedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(providedPassword))
}
