package pbusrbus

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"github.com/ardanlabs/service/business/sdk/delegate"
	"github.com/ardanlabs/service/business/sdk/sqldb"
	"github.com/ardanlabs/service/foundation/logger"
	"github.com/ardanlabs/service/foundation/otel"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Set of error variables for CRUD operations.
var (
	ErrNotFound              = errors.New("user not found")
	ErrUniqueEmail           = errors.New("email is not unique")
	ErrAuthenticationFailure = errors.New("authentication failed")
)

// Storer interface declares the behavior this package needs to persist and
// retrieve data.
type Storer interface {
	NewWithTx(tx sqldb.CommitRollbacker) (Storer, error)
	Create(ctx context.Context, usr NewPublicUser) error
	Update(ctx context.Context, usr NewPublicUser) error
	Delete(ctx context.Context, usr NewPublicUser) error
	QueryByEmail(ctx context.Context, email mail.Address) (NewPublicUser, error)
}

// Business manages the set of APIs for user access.
type Business struct {
	log      *logger.Logger
	storer   Storer
	delegate *delegate.Delegate
}

// NewBusiness constructs a user business API for use.
func NewBusiness(log *logger.Logger, delegate *delegate.Delegate, storer Storer) *Business {
	return &Business{
		log:      log,
		delegate: delegate,
		storer:   storer,
	}
}

// Create adds a new user to the system.
func (b *Business) Create(ctx context.Context, nu NewPublicUser) (PublicUser, error) {
	ctx, span := otel.AddSpan(ctx, "business.publicnewuser.create")
	defer span.End()

	hash, err := bcrypt.GenerateFromPassword([]byte(nu.Password), bcrypt.DefaultCost)
	if err != nil {
		return PublicUser{}, fmt.Errorf("generatefrompassword: %w", err)
	}

	now := time.Now()

	usr := PublicUser{
		ID:           uuid.New(),
		Name:         nu.Name,
		Email:        nu.Email,
		PasswordHash: hash,
		Roles:        nu.Roles,
		Department:   nu.Department,
		Enabled:      true,
		DateCreated:  now,
		DateUpdated:  now,
	}

	if err := b.storer.Create(ctx, usr); err != nil {
		return PublicUser{}, fmt.Errorf("create: %w", err)
	}

	return usr, nil
}

// QueryByEmail finds the user by a specified user email.
func (b *Business) QueryByEmail(ctx context.Context, email mail.Address) (User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.querybyemail")
	defer span.End()

	user, err := b.storer.QueryByEmail(ctx, email)
	if err != nil {
		return PublicUser{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	return user, nil
}

// Authenticate finds a user by their email and verifies their password. On
// success it returns a Claims User representing this user. The claims can be
// used to generate a token for future authentication.
func (b *Business) Authenticate(ctx context.Context, email mail.Address, password string) (User, error) {
	ctx, span := otel.AddSpan(ctx, "business.userbus.authenticate")
	defer span.End()

	usr, err := b.QueryByEmail(ctx, email)
	if err != nil {
		return User{}, fmt.Errorf("query: email[%s]: %w", email, err)
	}

	if err := bcrypt.CompareHashAndPassword(usr.PasswordHash, []byte(password)); err != nil {
		return User{}, fmt.Errorf("comparehashandpassword: %w", ErrAuthenticationFailure)
	}

	return usr, nil
}
