package middleware

import (
	"net/http"

	"github.com/torngkab/grit-api-gateway-service/app/accounts/signin"
	"github.com/torngkab/grit-api-gateway-service/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func BasicAuth(config config.Config) echo.MiddlewareFunc {
	return middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		signInAdapter := signin.NewSignInAdapter(config)
		signInService := signin.NewSignInService(signInAdapter)
		signInResponse, err := signInService.SignIn(c.Request().Context(), signin.SignInRequest{
			Username: username,
			Password: password,
		})
		if err != nil {
			return false, echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized")
		}

		c.Set("userId", signInResponse)

		return true, nil
	})
}
