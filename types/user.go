package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	cost               int = 12
	minFirstNameLength int = 2
	minLastNameLength  int = 2
	minPwdLength       int = 8
)

type CreateUserDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Pwd       string `json:"password"`
}

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string             `bson:"firstname" json:"firstname"`
	LastName  string             `bson:"lastname" json:"lastname"`
	Email     string             `bson:"email" json:"email"`
	EncPwd    string             `bson:"encpwd" json:"-"`
}

func (dto CreateUserDTO) Validate() []string {
	errors := []string{}
	if len(dto.FirstName) < minFirstNameLength {
		errors = append(errors, fmt.Sprintf("first name should be at least %d characters", minFirstNameLength))
	}
	if len(dto.LastName) < minLastNameLength {
		errors = append(errors, fmt.Sprintf("last name should be at least %d characters", minLastNameLength))
	}
	if len(dto.Pwd) < minPwdLength {
		errors = append(errors, fmt.Sprintf("password should be at least %d characters", minPwdLength))
	}
	if !isEmailValid(dto.Email) {
		errors = append(errors, fmt.Sprintf("email isn't valid"))
	}
	return errors
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func NewUserFromDTO(dto CreateUserDTO) (*User, error) {
	encpwd, err := bcrypt.GenerateFromPassword([]byte(dto.Pwd), cost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName: dto.FirstName,
		LastName:  dto.LastName,
		Email:     dto.Email,
		EncPwd:    string(encpwd),
	}, nil
}
