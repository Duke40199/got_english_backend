package models

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Account model struct
type Account struct {
	//Login
	ID       uuid.UUID `gorm:"size:255;column:id;not null;unique; primaryKey;" json:"id"`
	Username *string   `gorm:"size:255;unique" json:"username"`
	Fullname *string   `gorm:"size:255;" json:"fullname"`
	Email    *string   `gorm:"size:100;not null;unique;" json:"email"`
	Password *string   `gorm:"-" json:"password,omitempty"`
	RoleName string    `gorm:"size:100;not null;" json:"role_name"`
	//Info
	AvatarURL   *string    `gorm:"size:255" json:"avatar_url"`
	Address     *string    `gorm:"size:255;" json:"address"`
	PhoneNumber *string    `gorm:"column:phone_number;autoCreateTime" json:"phone_number"`
	Birthday    *string    `gorm:"column:birthday;type:date" json:"birthday" sql:"date"`
	IsSuspended bool       `gorm:"column:is_suspended;default:false;" json:"is_suspended"`
	SuspendedAt *time.Time `gorm:"column:suspended_at" json:"suspended_at"`
	//Role
	Learner   *Learner   `gorm:"foreignKey:AccountID" json:"learner_details,omitempty"`
	Expert    *Expert    `gorm:"foreignKey:AccountID" json:"expert_details,omitempty"`
	Moderator *Moderator `gorm:"foreignKey:AccountID" json:"moderator_details,omitempty"`
	Admin     *Admin     `gorm:"foreignKey:AccountID" json:"admin_details,omitempty"`
	//default timestamps
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime" json:"updated_at"`
}

type PermissionStruct struct {
	//learner info
	AvailableCoinCount uint `gorm:"column:account_id" json:"available_coin_count"`
	//expert perms
	CanChat                   bool `json:"can_chat,omitempty"`
	CanJoinTranslationSession bool `json:"can_join_translation_session,omitempty"`
	CanJoinLiveCallSession    bool `json:"can_join_live_call_session,omitempty"`
	//admin perms
	CanManageExpert    bool `json:"can_manage_expert,omitempty"`
	CanManageLearner   bool `json:"can_manage_learner,omitempty"`
	CanManageAdmin     bool `json:"can_manage_admin,omitempty"`
	CanManageModerator bool `json:"can_manage_moderator,omitempty"`
	//moderator perms
	CanManageCoinBundle      bool `json:"can_manage_coin_bundle,omitempty"`
	CanManagePricing         bool `json:"can_manage_pricing,omitempty"`
	CanManageApplicationForm bool `json:"can_manage_application_form,omitempty"`
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
