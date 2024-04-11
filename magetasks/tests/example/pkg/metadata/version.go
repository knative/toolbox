package metadata

// Version holds application version information.
var Version = "0.0.0"

// VersionPath return a path to the version variable.
func VersionPath() string {
	return importPath("Version")
}
