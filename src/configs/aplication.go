package configs

import (
	"cmp"
	"dainxor/atv/logger"
	"dainxor/atv/types"
	"os"
	"strconv"
)

const (
	DEFAULT_ROUTE_VERSION = "1"     // Default version for the API routes
	DEFAULT_API_VERSION   = "0.1.0" // Default version for the API
)

type appType struct {
	apiVersion    types.Version
	routesVersion uint32
}

var App appType

func init() {
	envInit()
}
func ReloadAppEnv() {
	envInit()
}
func envInit() {
	stringRoutesVersion := cmp.Or(os.Getenv("ATV_ROUTE_VERSION"), DEFAULT_ROUTE_VERSION)
	routeVersion, err := strconv.ParseUint(stringRoutesVersion, 10, 32)
	if err != nil {
		logger.Warningf("Invalid route version '%s': %v, falling back to default '%s'", stringRoutesVersion, err, DEFAULT_ROUTE_VERSION)
		// Parse default - this should never fail as it's a constant
		if routeVersion, err = strconv.ParseUint(DEFAULT_ROUTE_VERSION, 10, 32); err != nil {
			logger.Fatal("Invalid DEFAULT_ROUTE_VERSION constant:", err)
		}
	}
	App.routesVersion = uint32(routeVersion)

	appVersion, err := types.VersionFrom(cmp.Or(os.Getenv("ATV_API_VERSION"), DEFAULT_API_VERSION))
	if err != nil {
		logger.Warningf("Invalid API version: %v, falling back to default '%s'", err, DEFAULT_API_VERSION)
		appVersion = types.V(DEFAULT_API_VERSION)
	}
	App.apiVersion = appVersion
}

func (appType) RoutesVersion() uint32 {
	return App.routesVersion
}

func (appType) ApiVersion() types.Version {
	return App.apiVersion
}
