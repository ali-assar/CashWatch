package types

import (
	"fmt"
	"regexp"

	"golang.org/x/crypto/bcrypt"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const BcryptCost = 12

const minFirstNameLen = 2
const minLastNameLen = 2
const minPasswordLen = 7

type User struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ID                int32  `protobuf:"varint,1,opt,name=iD,proto3" json:"iD,omitempty"`
	Email             string `protobuf:"bytes,2,opt,name=email,proto3" json:"email,omitempty"`
	FirstName         string `protobuf:"bytes,3,opt,name=first_name,json=firstName,proto3" json:"first_name,omitempty"`
	LastName          string `protobuf:"bytes,4,opt,name=last_name,json=lastName,proto3" json:"last_name,omitempty"`
	EncryptedPassword string `protobuf:"bytes,5,opt,name=encrypted_password,json=encryptedPassword,proto3" json:"encrypted_password,omitempty"`
}
type UpdateUserParams struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}

func (params UpdateUserParams) ValidateUpdate() map[string]string {
	errors := make(map[string]string, 2)
	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("first name length should at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("last name length should at least %d characters", minLastNameLen)
	}
	return errors
}

type CreateUserParams struct {
	ID        string `json:"id,omitempty"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

func (params CreateUserParams) Validate() map[string]string {
	errors := make(map[string]string, 4)

	if len(params.FirstName) < minFirstNameLen {
		errors["firstName"] = fmt.Sprintf("first name length should at least %d characters", minFirstNameLen)
	}
	if len(params.LastName) < minLastNameLen {
		errors["lastName"] = fmt.Sprintf("last name length should at least %d characters", minLastNameLen)
	}
	if len(params.Password) < minPasswordLen {
		errors["password"] = fmt.Sprintf("password length should at least %d characters", minPasswordLen)
	}
	if !isEmailValid(params.Email) {
		errors["email"] = fmt.Sprintf("email %s is invalid", params.Email)
	}
	return errors
}

func IsPasswordValid(encpw, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(encpw), []byte(pw)) == nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}

func NewUserFormParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), BcryptCost)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
