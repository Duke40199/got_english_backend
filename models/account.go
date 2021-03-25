package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Account model struct
type Account struct {
	gorm.Model `json:"-"`
	//Login
	ID       uuid.UUID `gorm:"size:255;column:id;not null;unique; primaryKey;" json:"id"`
	Username string    `gorm:"size:255;not null;unique" json:"username"`
	Fullname string    `gorm:"size:255;" json:"fullname"`
	Email    string    `gorm:"size:100;not null;unique" json:"email"`
	Password string    `gorm:"size:100;" json:"password"`
	RoleName string    `gorm:"size:100;not null;" json:"role_name"`
	//Info
	AvatarURL   string     `gorm:"size:255" json:"avatar_url"`
	Address     string     `gorm:"size:255;" json:"address"`
	PhoneNumber string     `gorm:"column:phone_number;autoCreateTime" json:"phone_number"`
	Birthday    time.Time  `gorm:"column:birthday;type:date" json:"birthday" sql:"date"`
	IsSuspended bool       `gorm:"column:isSuspended" json:"is_suspended"`
	SuspendedAt *time.Time `gorm:"column:SuspendedAt" json:"suspended_at"`
	//default timestamps
	CreatedAt time.Time  `gorm:"column:CreatedAt;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:UpdatedAt;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:DeletedAt" json:"deleted_at";sql:"index"`
}

//Hash password
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

//VerifyPassword .
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

//BeforeSave checks Hash
func (u *Account) BeforeSave(*gorm.DB) error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *Account) BeforeUpdate(*gorm.DB) error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}
