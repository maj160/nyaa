package form

import (
	"regexp"

	"github.com/ewhal/nyaa/util/log"
)

const EMAIL_REGEX = `(\w[-._\w]*\w@\w[-._\w]*\w\.\w{2,3})`
const USERNAME_REGEX = `(\W)`

func EmailValidation(email string, err map[string][]string) (bool, map[string][]string) {
	exp, errorRegex := regexp.Compile(EMAIL_REGEX)
	if regexpCompiled := log.CheckError(errorRegex); regexpCompiled {
		if exp.MatchString(email) {
			return true, err
		}
	} 
	err["email"] = append(err["email"], "Email Address is not valid")
	return false, err
}
func ValidateUsername(username string, err map[string][]string) (bool, map[string][]string)  {
	exp, errorRegex := regexp.Compile(USERNAME_REGEX)
	if regexpCompiled := log.CheckError(errorRegex); regexpCompiled {
		if exp.MatchString(username) {
			err["username"] = append(err["username"], "Username contains illegal characters")
			return false, err
		}
	} else {
		return false, err
	}
	return true, err
}

func NewErrors() map[string][]string {
	err := make(map[string][]string)
	return err
}
func NewInfos() map[string][]string {
	infos := make(map[string][]string)
	return infos
}
func IsAgreed(t_and_c string) bool {
	if t_and_c == "1" {
		return true
	}
	return false
}

// RegistrationForm is used when creating a user.
type RegistrationForm struct {
	Username  string `form:"username" needed:"true" len_min:"3" len_max:"20"`
	Email     string `form:"email" needed:"true"`
	Password  string `form:"password" needed:"true" len_min:"6" len_max:"25" equalInput:"Confirm_Password"`
	Confirm_Password string `form:"password_confirmation" omit:"true" needed:"true"`
	CaptchaID string `form:"captchaID" omit:"true" needed:"true"`
	T_and_C   bool   `form:"t_and_c" omit:"true" needed:"true" equal:"true" hum_name:"Terms and Conditions"`
}

// LoginForm is used when a user logs in.
type LoginForm struct {
	Username    string `form:"username" needed="true"`
	Password string `form:"password" needed="true"`
}

// UserForm is used when updating a user.
type UserForm struct {
	Username  string `form:"username" needed:"true" len_min:"3" len_max:"20"`
	Email     string `form:"email" needed:"true"`
	Language  string `form:"language" default:"en-us"`
	CurrentPassword  string `form:"password" len_min:"6" len_max:"25" omit:"true"`
	Password  string `form:"password" len_min:"6" len_max:"25" equalInput:"Confirm_Password"`
	Confirm_Password string `form:"password_confirmation" omit:"true"`
	Status 	  int `form:"language" default:"0"`
}

// PasswordForm is used when updating a user password.
type PasswordForm struct {
	CurrentPassword string `form:"currentPassword"`
	Password        string `form:"newPassword"`
}

// SendPasswordResetForm is used when sending a password reset token.
type SendPasswordResetForm struct {
	Email string `form:"email"`
}

// PasswordResetForm is used when reseting a password.
type PasswordResetForm struct {
	PasswordResetToken string `form:"token"`
	Password           string `form:"newPassword"`
}
