package configs

import (
	"cmp"
	"dainxor/atv/utils"
	"os"
	"strconv"
)

const (
	DEFAULT_ROUTE_VERSION = 1       // Default version for the API routes
	DEFAULT_API_VERSION   = "0.1.1" // Default version for the API
	DEFAULT_APP_VERSION   = 2       // Default version for the application
)

type appType struct {
	routesVersion uint64
	appVersion    uint64 // This represents the number of significant versions of the API (major + minor)

	apiVersion      string
	apiMajorVersion uint64
	apiMinorVersion uint64
	apiPatchVersion uint64
}

var App appType

func (appType) EnvInit() error {
	App.apiVersion = os.Getenv("ATV_API_VERSION")
	routesVersion, _ := strconv.ParseInt(os.Getenv("ATV_ROUTE_VERSION"), 10, 64)

	App.routesVersion = cmp.Or(uint64(routesVersion), DEFAULT_ROUTE_VERSION)

	extractedStr := os.Getenv("ATV_APP_VERSION")
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	App.appVersion = cmp.Or(num, DEFAULT_APP_VERSION)

	extractedStr = utils.Extract("", App.apiVersion, ".")
	num, _ = strconv.ParseUint(extractedStr, 10, 64)
	App.apiMajorVersion = num

	extractedStr = utils.Extract(".", App.apiVersion, ".")
	num, _ = strconv.ParseUint(extractedStr, 10, 64)
	App.apiMinorVersion = num

	extractedStr = utils.Extract(".", App.apiVersion, ".")
	num, _ = strconv.ParseUint(extractedStr, 10, 64)
	App.apiPatchVersion = num

	return nil
}

func (appType) RoutesVersion() uint64 {
	return App.routesVersion
}
func (appType) Version() uint64 {
	return App.appVersion
}

func (appType) ApiVersion() string {
	return App.apiVersion
}
func (appType) ApiMajor() uint64 {
	return App.apiMajorVersion
}
func (appType) ApiMinor() uint64 {
	return App.apiMinorVersion
}
func (appType) ApiPatch() uint64 {
	return App.apiPatchVersion
}
