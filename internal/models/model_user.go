package models

import (
	"fmt"

	"github.com/mitchellh/mapstructure"
)

const UserType = "user"

type User struct {
	ID         int    `json:"id,omitempty" mapstructure:"id"`
	Username   string `json:"username,omitempty" mapstructure:"username" form:"username"`
	Password   string `json:"password,omitempty" mapstructure:"password" form:"password"`
	FirstName  string `json:"first_name,omitempty" mapstructure:"first_name"`
	LastName   string `json:"last_name,omitempty" mapstructure:"last_name"`
	Email      string `json:"email,omitempty" mapstructure:"email"`
	Phone      string `json:"phone,omitempty" mapstructure:"phone"`
	UserStatus bool   `json:"user_status" mapstructure:"user_status"`
}

func (u *User) GetType() string {
	return UserType
}

func (u *User) GetList() (lst []interface{}) {
	lst = append(lst, u.Username, u.Password, u.FirstName, u.LastName, u.Email, u.Phone, u.UserStatus)
	return
}

func (u *User) Set(m map[string]interface{}) error {
	if err := mapstructure.Decode(m, &u); err != nil {
		return err
	}
	return nil
}

func (u *User) Get() map[string]interface{} {
	mUserFields := map[string]interface{}{
		"id":          u.ID,
		"username":    u.Username,
		"password":    u.Password,
		"first_name":  u.FirstName,
		"last_name":   u.LastName,
		"email":       u.Email,
		"phone":       u.Phone,
		"user_status": u.UserStatus,
	}
	return mUserFields
}

func (u *User) String() string {
	return fmt.Sprintf("{\nID: %d\nUsername: %s\nPassword: %s\nFirstName: %s\nLastName: %s\nEmail: %s\nPhone: %s\nUserStatus: %t\n}",
		u.ID, u.Username, u.Password, u.FirstName, u.LastName, u.Email, u.Phone, u.UserStatus)
}
