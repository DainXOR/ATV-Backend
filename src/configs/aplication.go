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
	DEFAULT_ADDRESS       = ":8080" // Default server address
)

type appEnvNS struct{}
type appModeNS struct{}
type appNS struct {
	environment   string
	apiVersion    types.Version
	routesVersion uint32
	address       string
}

var App appNS
var _appEnvVars appEnvNS
var _appModeVars appModeNS

func (appNS) Env() appEnvNS {
	return _appEnvVars
}
func (appNS) Mode() appModeNS {
	return _appModeVars
}

func (appEnvNS) Mode() string {
	return "APP_ENV"
}
func (appEnvNS) APIVersion() string {
	return "APP_API_VERSION"
}
func (appEnvNS) RouteVersion() string {
	return "APP_ROUTE_VERSION"
}
func (appEnvNS) Address() string {
	return "SERVER_ADDRESS"
}

func (appModeNS) Production() string {
	return "prod"
}
func (appModeNS) Development() string {
	return "dev"
}
func (appModeNS) Debug() string {
	return "debug"
}
func (appModeNS) Test() string {
	return "test"
}
func (appModeNS) Default() string {
	return App.Mode().Debug()
}

func init() {
	envInit()
}
func ReloadAppEnv() {
	envInit()
}

func envInit() {
	envMode := cmp.Or(os.Getenv(App.Env().Mode()), "none")

	switch envMode {
	case "prod":
		App.environment = App.Mode().Production()
	case "dev":
		App.environment = App.Mode().Development()
	case "debug":
		App.environment = App.Mode().Debug()
	case "test":
		App.environment = App.Mode().Test()
	default:
		logger.Warningf("Unknown %s '%s', defaulting to '%s'", App.Env().Mode(), envMode, App.Mode().Default())
		App.environment = App.Mode().Default()
	}

	logger.Info("Application environment set to:", App.Environment())

	stringRoutesVersion := cmp.Or(os.Getenv(App.Env().RouteVersion()), DEFAULT_ROUTE_VERSION)
	routeVersion, err := strconv.ParseUint(stringRoutesVersion, 10, 32)
	if err != nil {
		logger.Warningf("Invalid route version '%s': %v, falling back to default '%s'", stringRoutesVersion, err, DEFAULT_ROUTE_VERSION)
		// Parse default - this should never fail as it's a constant
		if routeVersion, err = strconv.ParseUint(DEFAULT_ROUTE_VERSION, 10, 32); err != nil {
			logger.Fatal("Invalid DEFAULT_ROUTE_VERSION constant:", err)
		}
	}
	App.routesVersion = uint32(routeVersion)

	appVersion, err := types.VersionFrom(cmp.Or(os.Getenv(App.Env().APIVersion()), DEFAULT_API_VERSION))
	if err != nil {
		logger.Warningf("Invalid API version: %v, falling back to default '%s'", err, DEFAULT_API_VERSION)
		appVersion = types.V(DEFAULT_API_VERSION)
	}
	App.apiVersion = appVersion

	App.address = cmp.Or(os.Getenv(App.Env().Address()), DEFAULT_ADDRESS)
}

func (appNS) Environment() string {
	return App.environment
}
func (appNS) RoutesVersion() uint32 {
	return App.routesVersion
}
func (appNS) ApiVersion() types.Version {
	return App.apiVersion
}
func (appNS) Address() string {
	return App.address
}
