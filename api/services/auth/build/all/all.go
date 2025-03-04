// Package all binds all the routes into the specified app.
package all

import (
	"time"

	"github.com/ardanlabs/service/app/domain/authapp"
	"github.com/ardanlabs/service/app/domain/checkapp"
	"github.com/ardanlabs/service/app/domain/publicapp"
	"github.com/ardanlabs/service/app/sdk/mux"
	"github.com/ardanlabs/service/business/domain/publicuesrbus"
	publicuserdb "github.com/ardanlabs/service/business/domain/publicuesrbus/stores/publicuserdb"
	pubicusercache "github.com/ardanlabs/service/business/domain/publicuesrbus/stores/usercache"
	"github.com/ardanlabs/service/business/domain/userbus"
	"github.com/ardanlabs/service/business/domain/userbus/stores/usercache"
	"github.com/ardanlabs/service/business/domain/userbus/stores/userdb"
	"github.com/ardanlabs/service/business/sdk/delegate"
	"github.com/ardanlabs/service/foundation/web"
)

// Routes constructs the add value which provides the implementation of
// of RouteAdder for specifying what routes to bind to this instance.
func Routes() add {
	return add{}
}

type add struct{}

// Add implements the RouterAdder interface.
func (add) Add(app *web.App, cfg mux.Config) {

	// Construct the business domain packages we need here so we are using the
	// sames instances for the different set of domain apis.
	delegate := delegate.New(cfg.Log)
	userBus := userbus.NewBusiness(cfg.Log, delegate, usercache.NewStore(cfg.Log, userdb.NewStore(cfg.Log, cfg.DB), time.Minute))
	publicUserBus := publicuesrbus.NewBusiness(cfg.Log, delegate, pubicusercache.NewStore(cfg.Log, publicuserdb.NewStore(cfg.Log, cfg.DB), time.Minute))

	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})

	authapp.Routes(app, authapp.Config{
		UserBus: userBus,
		Auth:    cfg.Auth,
	})

	publicapp.Routes(app, publicapp.Config{
		Auth:             cfg.Auth,
		PublicNewUserBus: publicUserBus,
		Kid:              "54bb2165-71e1-41a6-af3e-7da4a0e1e2c1",
	})
}
