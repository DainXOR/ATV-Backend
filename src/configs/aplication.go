package configs

import (
	"cmp"
	"dainxor/atv/logger"
	"dainxor/atv/utils"
	"os"
	"strconv"
)

const (
	DEFAULT_ROUTE_VERSION = "1"     // Default version for the API routes
	DEFAULT_API_VERSION   = "0.1.2" // Default version for the API
)

type appType struct {
	routesVersion uint64

	apiVersion      string
	apiMajorVersion uint64
	apiMinorVersion uint64
	apiPatchVersion uint64
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
	App.routesVersion, _ = strconv.ParseUint(stringRoutesVersion, 10, 64)

	App.apiVersion = cmp.Or(os.Getenv("ATV_API_VERSION"), DEFAULT_API_VERSION)
	App.apiMajorVersion = versionMajor(App.apiVersion)
	App.apiMinorVersion = versionMinor(App.apiVersion)
	App.apiPatchVersion = versionPatch(App.apiVersion)

	logger.Info("Application initialized with API version:", App.apiVersion)
	logger.Info("Application initialized with Routes version:", App.routesVersion)
}

func versionMajor(version string) uint64 {
	extractedStr := utils.Extract("", version, ".")
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	return num
}
func versionMinor(version string) uint64 {
	extractedStr := utils.Extract(".", version, ".")
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	return num
}
func versionPatch(version string) uint64 {
	extractedStr := utils.Extract(".", version, "")
	num, _ := strconv.ParseUint(extractedStr, 10, 64)
	return num
}

func (appType) RoutesVersion() uint64 {
	return App.routesVersion
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
