package echor

import (
	"net/http"

	"github.com/labstack/echo"
)

// Response func
func Response(c echo.Context, data interface{}) error {
	outputFormat := c.QueryParam("respfmt")
	if outputFormat == "xml" {
		return c.XML(http.StatusOK, data)
	}
	return c.JSON(http.StatusOK, data)
}

// ResponseOK func
func ResponseOK(c echo.Context) error {
	outputFormat := c.QueryParam("respfmt")
	data := struct {
		OK int `json:"ok" xml:"ok"`
	}{
		OK: 1,
	}
	if outputFormat == "xml" {
		return c.XML(http.StatusOK, data)
	}
	return c.JSON(http.StatusOK, data)
}

// GetAuthToken func
func GetAuthToken(c echo.Context) string {
	tk := c.QueryParam("authtoken")
	if len(tk) == 0 {
		tk = c.Request().Header.Get("Authorization")
	}
	return tk
}
