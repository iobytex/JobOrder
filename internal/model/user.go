package model

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"strings"
)

type User struct {
	gorm.Model
	Name string
	Code string
	PhoneNumber string `gorm:"column:phone_no;unique" `
	Passcode string
	Role string
	Order []Order `gorm:"constraint:OnDelete:RESTRICT;"`
}

type UserRequest struct{
	Name string `form:"name" binding:"required"`
	PhoneNumber string `form:"phone_no" binding:"required"`
	Passcode string `form:"passcode" binding:"required"`
	Role string `form:"role" binding:"required"`
}

type Login struct {
	PhoneNumber string `form:"phone_no" binding:"required"`
	Passcode string `form:"passcode" binding:"required"`
}

func (User) TableName() string {
	return "users"
}


func (request *UserRequest) HashPassword() error {
	password, err := bcrypt.GenerateFromPassword([]byte(request.Passcode),bcrypt.DefaultCost)
	if err != nil {
		return errors.Wrap(err,"Error hashing passcode")
	}
	request.Passcode = string(password)
	return  nil
}

func (user *User) CompareHashPassCode(hashedPasscode []byte)  error {
	err := bcrypt.CompareHashAndPassword(hashedPasscode, []byte(user.Passcode))
	if err != nil {
		return errors.Wrap(err,"Error comparing hashed passcode")
	}

	return nil
}

// Sanitize user password
func (user *User) SanitizePassword() {
	user.Passcode  = ""
}

// Prepare user for register
func (request *UserRequest) PrepareCreate() error {
	request.Passcode = strings.TrimSpace(request.Passcode )

	if err := request.HashPassword(); err != nil {
		return err
	}
	return nil
}

