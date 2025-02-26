package publicapp

import (
	"encoding/json"
	"fmt"
	"net/mail"

	"github.com/ardanlabs/service/app/sdk/errs"
	"github.com/ardanlabs/service/business/domain/pbusrbus"
	"github.com/ardanlabs/service/business/types/name"
	"github.com/ardanlabs/service/business/types/role"
)

// RegisterRequest contains information needed for creating a new user.
type RegisterRequest struct {
	Name            string   `json:"name" validate:"required"`
	Email           string   `json:"email" validate:"required,email"`
	Department      string   `json:"department"`
	Password        string   `json:"password" validate:"required,min=8"`
	PasswordConfirm string   `json:"passwordConfirm" validate:"required,eqfield=Password"`
	Roles           []string `json:"roles"` // Will be overridden to USER role only
}

// Decode implements the decoder interface.
func (reg *RegisterRequest) Decode(data []byte) error {
	return json.Unmarshal(data, reg)
}

// Validate checks the data in the model is considered clean.
func (req RegisterRequest) Validate() error {
	if err := errs.Check(req); err != nil {
		return fmt.Errorf("validate: %w", err)
	}

	return nil
}

func toBusPublicNewUser(req RegisterRequest) (pbusrbus.NewUser, error) {
	roles, err := role.ParseMany(req.Roles)
	if err != nil {
		return pbusrbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	addr, err := mail.ParseAddress(req.Email)
	if err != nil {
		return pbusrbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	nme, err := name.Parse(req.Name)
	if err != nil {
		return pbusrbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	department, err := name.ParseNull(req.Department)
	if err != nil {
		return pbusrbus.NewUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := pbusrbus.NewUser{
		Name:       nme,
		Email:      *addr,
		Roles:      roles,
		Department: department,
		Password:   req.Password,
	}

	return bus, nil
}

// LoginRequest contains information needed for logging in.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// Decode implements the decoder interface.
func (log *LoginRequest) Decode(data []byte) error {
	return json.Unmarshal(data, log)
}

// LoginResponse contains the response data for login.
type LoginResponse struct {
	Token string `json:"token"`
}

// Encode implements the encoder interface.
func (res LoginResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(res)
	return data, "application/json", err
}
