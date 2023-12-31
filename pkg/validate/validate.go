package validate

import "regexp"

type (
	Value           string
	Error           string
	ValidationError map[Value]Error
)

const (
	EmailValue    Value = "email"
	UsernameValue Value = "username"
	PasswordValue Value = "password"
	TitleValue    Value = "title"
	LinkValue     Value = "link"
)

var (
	EmailExistsError    Error = "Sorry, this email is already in use"
	EmailInvalidError   Error = "Please enter a valid email address"
	UsernameExistsError Error = "Sorry, this username is already in use"
	UsernameLengthError Error = "Username must be between 4 and 20 characters"
	UsernameCharError   Error = "Username must only contain letters, numbers, and underscores"
	PasswordLengthError Error = "Password must be between 8 and 50 characters"
	PasswordCharError   Error = "Password must contain at least one special character"
	PasswordMatchError  Error = "Passwords do not match"
	TitleLengthError    Error = "Title must be between 1 and 100 characters"
	LinkInvalidError    Error = "Please enter a valid URL"
)

var (
	emailRegex       = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	usernameRegex    = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	specialCharRegex = regexp.MustCompile(`[^a-zA-Z0-9]`)
	linkRegex        = regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}(/.*)?$`)
)

func Email(email string) ValidationError {
	if !emailRegex.MatchString(email) {
		return ValidationError{EmailValue: EmailInvalidError}
	}
	return nil
}

func Username(username string) ValidationError {
	if len(username) < 4 || len(username) > 25 {
		return ValidationError{UsernameValue: UsernameLengthError}
	}
	if !usernameRegex.MatchString(username) {
		return ValidationError{UsernameValue: UsernameCharError}
	}
	return nil
}

func Password(password string) ValidationError {
	if len(password) < 8 || len(password) > 50 {
		return ValidationError{PasswordValue: PasswordLengthError}
	}
	if !specialCharRegex.MatchString(password) {
		return ValidationError{PasswordValue: PasswordCharError}
	}
	return nil
}

func Title(title string) ValidationError {
	if len(title) < 1 || len(title) > 100 {
		return ValidationError{TitleValue: TitleLengthError}
	}
	return nil
}

func Link(link string) ValidationError {
	if !linkRegex.MatchString(link) {
		return ValidationError{LinkValue: LinkInvalidError}
	}
	return nil
}

func (e Error) String() string {
	return string(e)
}
