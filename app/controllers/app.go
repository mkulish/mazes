package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/revel/revel"
	jwt "github.com/dgrijalva/jwt-go"
	gorpController "github.com/revel/modules/orm/gorp/app/controllers"

	"github.com/mkulish/mazes/app/models"
)

// App base controller
type App struct {
	gorpController.Controller
}

// Auth interceptor ensures authentication and stores user in session (if any)
func (c App) Auth() revel.Result {
	authData := strings.Split(c.Request.Header.Get("Authorization"), " ")
	if len(authData) != 2 || authData[0] != "Bearer" {
		return c.unauthorizedError()
	}

	claims, err := decodeToken(authData[1])
	username, found := claims["username"]
	if err != nil || ! found {
		return c.unauthorizedError()
	}

	user, err := c.getUser(username.(string))
	if user == nil || err != nil {
		return c.unauthorizedError()
	}

	c.Session.Set("user", user)
	return nil
}

// ValidationError returns JSON error response with validation errors and HTTP 400
func (c App) validationError(errors []*revel.ValidationError) revel.Result {
	c.Response.Status = http.StatusBadRequest
	return c.RenderJSON(models.ValidationError{ Errors: errors })
}

// unauthorizedError returns JSON error response with HTTP 401
func (c App) unauthorizedError() revel.Result {
	c.Response.Status = http.StatusUnauthorized
	return c.RenderJSON(models.UnauthorizedError{ Error: "Unauthorized" })
}

// InternalError returns JSON error response with HTTP 500
func (c App) internalError() revel.Result {
	c.Response.Status = http.StatusInternalServerError
	return c.RenderJSON(models.InternalError{ Error: "Internal error" })
}

// encodeToken returns JWT auth token
func encodeToken(user *models.User) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id": user.ID,
		"username": user.Username,
		"exp":   time.Now().Add(time.Duration(24) * time.Hour).Unix(),
	})

	tokenString, _ := token.SignedString(hmacSecret)
	return tokenString
}

// decodeToken parses and verifies JWT auth token
func decodeToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return hmacSecret, nil
	})

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		return claims, nil
	}

	return jwt.MapClaims{}, err
}
