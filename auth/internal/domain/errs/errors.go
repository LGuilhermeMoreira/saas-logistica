package errs

import "errors"

var (
	ErrInvalidUserName       = errors.New("client must have a name")
	ErrInvalidEmail          = errors.New("email is invalid")
	ErrWeakPassword          = errors.New("password does not meet the security requirements")
	ErrInvalidUserID         = errors.New("user id is invalid")
	ErrInvalidRoleID         = errors.New("role id is invalid")
	ErrInvalidRoleName       = errors.New("name is invalid")
	ErrRoleWithoutPermission = errors.New("role must have one permission")

	ErrInvalidCredentials = errors.New("invalid credentials")

	ErrUserNotFound = errors.New("user not found")
	ErrRoleNotFound = errors.New("role not found")

	ErrUserAlreadyExists = errors.New("user already exists")
	ErrRoleAlreadyExists = errors.New("role already exists")
	ErrRoleInUse         = errors.New("role is assigned to one or more users")

	ErrInternal = errors.New("internal server error")
)
