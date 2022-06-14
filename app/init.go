
//go:generate swagger generate spec -m -o ../swagger.json

// Maze API
//
// Backend #2 demo project
//
//     Schemes: https
//     Host: mazes.demo.pics
//     BasePath: /
//     Version: 0.0.1
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: Max<max-kulishev@gmail.com>
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     SecurityDefinitions:
//     oauth2:
//         type: oauth2
//         in: header
//         scopes:
//           read: read access
//           write: write access
//         flow: accessCode
//
// swagger:meta
package app

import (
	"github.com/revel/revel"

	rgorp "github.com/revel/modules/orm/gorp/app"

	"github.com/mkulish/mazes/app/controllers"
	"github.com/mkulish/mazes/app/models"
)

var (
	// AppVersion revel app version (ldflags)
	AppVersion string

	// BuildTime revel app build-time (ldflags)
	BuildTime string
)

func init() {
	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		HeaderFilter,                  // Add some security based headers
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.CompressFilter,          // Compress the result.
		revel.BeforeAfterFilter,       // Call the before and after filter functions
		revel.ActionInvoker,           // Invoke the action.
	}

	revel.InterceptMethod(controllers.Maze.Auth, revel.BEFORE)

	revel.OnAppStart(InitSQLite)
}

// HeaderFilter adds common security headers
// There is a full implementation of a CSRF filter in
// https://github.com/revel/modules/tree/master/csrf
var HeaderFilter = func(c *revel.Controller, fc []revel.Filter) {
	c.Response.Out.Header().Add("X-Frame-Options", "SAMEORIGIN")
	c.Response.Out.Header().Add("X-XSS-Protection", "1; mode=block")
	c.Response.Out.Header().Add("X-Content-Type-Options", "nosniff")
	c.Response.Out.Header().Add("Referrer-Policy", "strict-origin-when-cross-origin")

	fc[0](c, fc[1:]) // Execute the next filter stage.
}

// InitSQLite initializes sqlite schema
func InitSQLite() {
	Dbm := rgorp.Db.Map

	t := Dbm.AddTable(models.User{}).SetKeys(true, "ID")
	t.AddIndex("UsernameIndex", "Btree", []string{"Username"}).SetUnique(true)
	t.ColMap("Password").Transient = true

	t = Dbm.AddTable(models.Maze{}).SetKeys(true, "ID")
	t.AddIndex("OwnerIDIndex", "Btree", []string{"OwnerID"})
	t.ColMap("Walls").Transient = true

	rgorp.Db.TraceOn(revel.AppLog)
	Dbm.CreateTables()
}
