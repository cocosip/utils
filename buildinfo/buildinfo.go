package buildinfo

import "strconv"

var (
	devMode = "true"
	version = "DEVELOPMENT Version"
	build   = "DEVELOPMENT Build"
)

func Version() string {
	return version
}

func Build() string {
	return build
}

func DevMode() bool {
	b, _ := strconv.ParseBool(devMode)
	return b
}
