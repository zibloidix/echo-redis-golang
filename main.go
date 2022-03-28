package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/", home)
	e.GET("/q/:id", queue)

	e.Logger.Fatal(e.Start(":3000"))
}

func home(c echo.Context) error {
	ip1, ip2, ip3 := getIPs(c)
	agent := getAgent(c)
	hash := getHash(ip1, ip2, ip3, agent)

	return c.String(http.StatusOK, "IPs and User agent:\n"+ip1+"\n"+ip2+"\n"+ip3+"\n"+agent+"\n"+"\nHash: "+hash)
}

func queue(c echo.Context) error {
	id := c.Param("id")
	return c.String(http.StatusOK, id)
}

func getIPs(c echo.Context) (string, string, string) {
	ip1 := echo.ExtractIPDirect()(c.Request())
	ip2 := echo.ExtractIPFromXFFHeader()(c.Request())
	ip3 := echo.ExtractIPFromRealIPHeader()(c.Request())
	return ip1, ip2, ip3
}

func getAgent(c echo.Context) string {
	return c.Request().UserAgent()
}

func getHash(ip1, ip2, ip3, agent string) string {
	hash := sha1.New()
	hash.Write([]byte([]byte(ip1 + ip2 + ip3 + agent)))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
