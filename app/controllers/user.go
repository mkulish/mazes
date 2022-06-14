package controllers

import (
	"database/sql"

	"github.com/revel/revel"
	"golang.org/x/crypto/bcrypt"

	"github.com/mkulish/mazes/app/models"
)

var hmacSecret = []byte("demo")

// User controller
type User struct {
	App
}

// Register creates a new user
// swagger:route POST /user user registerUser
//
// Registers a new user
//
//     Parameters:
//     + name: user
//       in: body
//       description: User signup data
//       required: true
//       type: User
//
//     Responses:
//       200: LoginResponse
//       400: ValidationError
//       500: InternalError
func (c App) Register(user models.User) revel.Result {
	user.Validate(c.Validation)

	if ! c.Validation.HasErrors() {
		// ensure username not taken
		existing, err := c.getUser(user.Username)
		if err != nil {
			return c.internalError()
		}
		if existing != nil {
			c.Validation.Error("Already taken").Key("username")
		}
	}

	if c.Validation.HasErrors() {
		return c.validationError(c.Validation.Errors)
	}

	user.HashedPassword, _ = bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	err := c.Txn.Insert(&user)
	if err != nil {
		c.Log.Errorf("user '%s' insert: %v", user.Username, err)
		return c.internalError()
	}

	return c.RenderJSON(models.LoginResponse{OK: true, Token: encodeToken(&user)})
}

// Login performs login and returns JWT token
// swagger:route POST /login user login
//
// Performs login
//
//     Parameters:
//     + name: user
//       in: body
//       description: login data
//       required: true
//       type: User
//
//     Responses:
//       200: LoginResponse
//       400: ValidationError
//       500: InternalError
func (c App) Login(loginData models.User) revel.Result {
	var user *models.User
	var err error

	loginData.Validate(c.Validation)
	if ! c.Validation.HasErrors() {
		// user lookup and password check
		if user, err = c.getUser(loginData.Username); err != nil {
			return c.internalError()
		}
		if user == nil {
			c.Validation.Error("Not found").Key("username")
		} else if bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(loginData.Password)) != nil {
			c.Validation.Error("Incorrect username or password").Key("password")
		}
	}

	if c.Validation.HasErrors() {
		return c.validationError(c.Validation.Errors)
	}

	return c.RenderJSON(models.LoginResponse{OK: true, Token: encodeToken(user)})
}

// getUser performs user lookup by username
func (c App) getUser(username string) (*models.User, error) {
	user := &models.User{}
	err := c.Txn.SelectOne(user, c.Db.SqlStatementBuilder.Select("*").From("User").Where("Username=?", username))

	if err == sql.ErrNoRows {
		// not found
		return nil, nil
	}
	if err != nil {
		c.Log.Error("Failed to find user", "user", username, "error", err)
	}
	return user, err
}
