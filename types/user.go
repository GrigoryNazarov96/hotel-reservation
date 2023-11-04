package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson"
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

type UpdateUserDTO struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

type LoginDTO struct {
	Email string `json:"email"`
	Pwd   string `json:"password"`
}

type LoginResponse struct {
	User  *User  `json:"user"`
	Token string `json:"token"`
}

func (d UpdateUserDTO) ToBSONM() bson.M {
	m := bson.M{}
	if len(d.FirstName) > 0 {
		m["firstname"] = d.FirstName
	}
	if len(d.LastName) > 0 {
		m["lastname"] = d.LastName
	}
	return m
}

type User struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	FirstName string             `bson:"firstname" json:"firstname"`
	LastName  string             `bson:"lastname" json:"lastname"`
	Email     string             `bson:"email" json:"email"`
	EncPwd    string             `bson:"encpwd" json:"-"`
}

func (dto CreateUserDTO) Validate() map[string]string {
	errors := map[string]string{}
	if len(dto.FirstName) < minFirstNameLength {
		errors["firstname"] = fmt.Sprintf("first name should be at least %d characters", minFirstNameLength)
	}
	if len(dto.LastName) < minLastNameLength {
		errors["lastname"] = fmt.Sprintf("last name should be at least %d characters", minLastNameLength)
	}
	if len(dto.Pwd) < minPwdLength {
		errors["pwd"] = fmt.Sprintf("password should be at least %d characters", minPwdLength)
	}
	if !isValidEmail(dto.Email) {
		errors["email"] = "email isn't valid"
	}
	return errors
}

func isValidEmail(e string) bool {
	emailRegex := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	return emailRegex.MatchString(e)
}

func IsValidPassword(pwd, candidate string) bool {
	return bcrypt.CompareHashAndPassword([]byte(pwd), []byte(candidate)) == nil
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
