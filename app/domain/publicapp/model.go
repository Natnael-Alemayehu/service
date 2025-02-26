package publicapp

import (
	"encoding/json"
	"fmt"
	"net/mail"
	"time"

	"github.com/ardanlabs/service/app/sdk/errs"
	"github.com/ardanlabs/service/business/domain/publicuesrbus"
	"github.com/ardanlabs/service/business/types/name"
	"github.com/ardanlabs/service/business/types/role"
)

// User represents information about an individual user.
type User struct {
	ID           string   `json:"id"`
	Name         string   `json:"name"`
	Email        string   `json:"email"`
	Roles        []string `json:"roles"`
	PasswordHash []byte   `json:"-"`
	Department   string   `json:"department"`
	Enabled      bool     `json:"enabled"`
	DateCreated  string   `json:"dateCreated"`
	DateUpdated  string   `json:"dateUpdated"`
}

// Encode implements the encoder interface.
func (app User) Encode() ([]byte, string, error) {
	data, err := json.Marshal(app)
	return data, "application/json", err
}

func toAppUser(bus publicuesrbus.PublicUser) User {
	return User{
		ID:           bus.ID.String(),
		Name:         bus.Name.String(),
		Email:        bus.Email.Address,
		Roles:        role.ParseToString(bus.Roles),
		PasswordHash: bus.PasswordHash,
		Department:   bus.Department.String(),
		Enabled:      bus.Enabled,
		DateCreated:  bus.DateCreated.Format(time.RFC3339),
		DateUpdated:  bus.DateUpdated.Format(time.RFC3339),
	}
}

// =============================================================================

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

func toBusPublicNewUser(req RegisterRequest) (publicuesrbus.NewPublicUser, error) {
	roles, err := role.ParseMany(req.Roles)
	if err != nil {
		return publicuesrbus.NewPublicUser{}, fmt.Errorf("parse: %w", err)
	}

	addr, err := mail.ParseAddress(req.Email)
	if err != nil {
		return publicuesrbus.NewPublicUser{}, fmt.Errorf("parse: %w", err)
	}

	nme, err := name.Parse(req.Name)
	if err != nil {
		return publicuesrbus.NewPublicUser{}, fmt.Errorf("parse: %w", err)
	}

	department, err := name.ParseNull(req.Department)
	if err != nil {
		return publicuesrbus.NewPublicUser{}, fmt.Errorf("parse: %w", err)
	}

	bus := publicuesrbus.NewPublicUser{
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
