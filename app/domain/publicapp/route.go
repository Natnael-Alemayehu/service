package publicapp

import (
	"net/http"

	"github.com/ardanlabs/service/app/sdk/auth"
	"github.com/ardanlabs/service/business/domain/pbusrbus"
	"github.com/ardanlabs/service/foundation/web"
)

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Auth             *auth.Auth
	PublicNewUserBus *pbusrbus.Business
	Kid              string
}

// Routes adds specific routes for this group.
func Routes(app *web.App, cfg Config) {
	const version = "v1"

	api := newApp(cfg.Auth, cfg.PublicNewUserBus, cfg.Kid)

	app.HandlerFunc(http.MethodPost, version, "/register", api.register)
	app.HandlerFunc(http.MethodPost, version, "/login", api.login)
}
