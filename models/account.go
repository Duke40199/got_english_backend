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
	Username *string   `gorm:"size:255;unique" json:"username"`
	Fullname *string   `gorm:"size:255;" json:"fullname"`
	Email    *string   `gorm:"size:100;not null;unique" json:"email"`
	Password *string   `gorm:"size:100;" json:"password"`
	RoleName string    `gorm:"size:100;not null;" json:"role_name"`
	//credentials
	//Access Token will be created when use logs in with Google
	AccessToken *string `gorm:"column:access_token;size:255;" json:"access_token"`
	//Info
	AvatarURL   *string    `gorm:"size:255" json:"avatar_url"`
	Address     *string    `gorm:"size:255;" json:"address"`
	PhoneNumber *string    `gorm:"column:phone_number;autoCreateTime" json:"phone_number"`
	Birthday    *time.Time `gorm:"column:birthday;type:date" json:"birthday" sql:"date"`
	IsSuspended *bool      `gorm:"column:isSuspended" json:"is_suspended"`
	SuspendedAt *time.Time `gorm:"column:SuspendedAt" json:"suspended_at"`
	//default timestamps
	CreatedAt time.Time  `gorm:"column:CreatedAt;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time  `gorm:"column:UpdatedAt;autoCreateTime" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:DeletedAt" json:"deleted_at"`
}

type AccountFullInfo struct {
	//Login
	ID          uuid.UUID `gorm:"size:255;column:id;not null;unique; primaryKey;" json:"id"`
	Username    string    `gorm:"size:255;not null;unique" json:"username"`
	Fullname    string    `gorm:"size:255;not null;unique" json:"fullname"`
	Email       string    `gorm:"size:100;not null;unique" json:"email"`
	RoleName    string    `gorm:"size:100;not null;" json:"role_name"`
	LearnerID   uint      `json:"learner_id,omitempty"`
	ExpertID    uint      `json:"expert_id,omitempty"`
	ModeratorID uint      `json:"moderator_id,omitempty"`
	AdminID     uint      `json:"admin_id,omitempty"`
	//Info
	AvatarURL   string     `json:"avatar_url"`
	Address     string     `json:"address"`
	PhoneNumber string     `json:"phone_number"`
	Birthday    string     `json:"birthday" sql:"date"`
	IsSuspended bool       `json:"is_suspended"`
	SuspendedAt *time.Time `json:"suspended_at"`
	//default timestamps
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	//expert perms
	CanChat                   *bool `json:"can_chat,omitempty"`
	CanJoinTranslationSession *bool `json:"can_join_translation_session,omitempty"`
	CanJoinPrivateCallSession *bool `json:"can_private_call_session,omitempty"`
	//admin perms
	CanManageExpert  *bool `json:"can_manage_expert,omitempty"`
	CanManageLearner *bool `json:"can_manage_learner,omitempty"`
	CanManageAdmin   *bool `json:"can_manage_admin,omitempty"`
	//moderator perms
	CanManageCoinBundle      *bool `json:"can_manage_coin_bundle,omitempty"`
	CanManagePricing         *bool `json:"can_manage_pricing,omitempty"`
	CanManageApplicationForm *bool `json:"can_manage_application_form,omitempty"`
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
	if u.Password != nil {
		hashedPassword, err := Hash(*u.Password)
		if err != nil {
			return err
		}
		*u.Password = string(hashedPassword)
	}
	return nil
}

func (u *Account) BeforeUpdate(*gorm.DB) error {
	if u.Password != nil {
		hashedPassword, err := Hash(*u.Password)
		if err != nil {
			return err
		}
		*u.Password = string(hashedPassword)
	}
	return nil
}
