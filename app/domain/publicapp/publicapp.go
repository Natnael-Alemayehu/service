package publicapp

import (
	"context"
	"errors"
	"net/http"
	"net/mail"
	"time"

	"github.com/ardanlabs/service/app/sdk/auth"
	"github.com/ardanlabs/service/app/sdk/errs"
	"github.com/ardanlabs/service/business/domain/pbusrbus"
	"github.com/ardanlabs/service/business/domain/userbus"
	"github.com/ardanlabs/service/business/types/role"
	"github.com/ardanlabs/service/foundation/web"
	"github.com/golang-jwt/jwt/v4"
)

type app struct {
	auth     *auth.Auth
	pbusrbus *pbusrbus.Business
	kid      string
}

func newApp(auth *auth.Auth, pbusrbus *pbusrbus.Business, kid string) *app {
	return &app{
		auth:     auth,
		pbusrbus: pbusrbus,
		kid:      kid,
	}
}

// register handles creating a new user in the system.
func (a *app) register(ctx context.Context, r *http.Request) web.Encoder {
	var req RegisterRequest
	if err := web.Decode(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	nu, err := toBusPublicNewUser(req)
	if err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	usr, err := a.pbusrbus.Create(ctx, nu)
	if err != nil {
		if errors.Is(err, userbus.ErrUniqueEmail) {
			return errs.NewFieldErrors("email", errors.New("email already exists"))
		}
		return errs.New(errs.Internal, err)
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: role.ParseToString(usr.Roles),
	}

	token, err := a.auth.GenerateToken(a.kid, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	_ = toAppUser(usr)

	return LoginResponse{Token: token}
}

// login authenticates a user and provides a JWT token.
func (a *app) login(ctx context.Context, r *http.Request) web.Encoder {
	var req LoginRequest
	if err := web.Decode(r, &req); err != nil {
		return errs.New(errs.InvalidArgument, err)
	}

	addr, err := mail.ParseAddress(req.Email)
	if err != nil {
		return errs.NewFieldErrors("email", errors.New("invalid email format"))
	}

	// Authenticate the user
	usr, err := a.pbusrbus.Authenticate(ctx, *addr, req.Password)
	if err != nil {
		return errs.New(errs.Unauthenticated, errors.New("invalid email or password"))
	}

	claims := auth.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   usr.ID.String(),
			Issuer:    a.auth.Issuer(),
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(8760 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		},
		Roles: role.ParseToString(usr.Roles),
	}

	token, err := a.auth.GenerateToken(a.kid, claims)
	if err != nil {
		return errs.New(errs.Internal, err)
	}

	return LoginResponse{Token: token}
}
