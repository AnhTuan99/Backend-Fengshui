// Package reporting binds the reporting domain set of routes into the specified app.
package reporting

import (
	"fengshui.com/back-fengshui/app/domain/checkapp"
	"fengshui.com/back-fengshui/app/sdk/mux"
	"fengshui.com/back-fengshui/foundation/web"
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

	checkapp.Routes(app, checkapp.Config{
		Build: cfg.Build,
		Log:   cfg.Log,
		DB:    cfg.DB,
	})
}
