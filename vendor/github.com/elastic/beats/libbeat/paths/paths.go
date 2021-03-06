// Package libbeat.paths provides a common way to handle paths
// configuration for all Beats.
//
// Currently the following paths are defined:
//
// path.home - It’s the default folder for everything that doesn't fit in
// the categories below
//
// path.data - Contains things that are expected to change often during normal
// operations (“registry” files, UUID file, etc.)
//
// path.config - Configuration files and Elasticsearch template default location
//
// These settings can be set via the configuration file or via command line flags.
// The CLI flags overwrite the configuration file options.
//
// Use the Resolve function to resolve files to their absolute paths. For example,
// to look for a file in the config path:
//
// cfgfilePath := paths.Resolve(paths.Config, "beat.yml"
package paths

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

var (
	homePath   = flag.String("path.home", "", "Home path")
	configPath = flag.String("path.config", "", "Configuration path")
	dataPath   = flag.String("path.data", "", "Data path")
	logsPath = flag.String("path.logs", "", "Logs path")
)

type Path struct {
	Home   string
	Config string
	Data   string
	Logs   string
}

// FileType is an enumeration type representing the file types.
// Currently existing file types are: Home, Config, Data
type FileType string

const (
	Home   FileType = "home"
	Config FileType = "config"
	Data   FileType = "data"
	Logs FileType = "logs"
)

// Paths is the Path singleton on which the top level functions from this
// package operate.
var Paths = New()

// New creates a new Paths object with all values set to empty values.
func New() *Path {
	return &Path{}
}

// InitPaths sets the default paths in the configuration based on CLI flags,
// configuration file and default values. It also tries to create the data
// path with mode 0755 and returns an error on failure.
func (paths *Path) InitPaths(cfg *Path) error {
	err := paths.initPaths(cfg)
	if err != nil {
		return err
	}

	// make sure the data path exists
	err = os.MkdirAll(paths.Data, 0755)
	if err != nil {
		return fmt.Errorf("Failed to create data path %s: %v", paths.Data, err)
	}

	return nil
}

// InitPaths sets the default paths in the configuration based on CLI flags,
// configuration file and default values. It also tries to create the data
// path with mode 0755 and returns an error on failure.
func InitPaths(cfg *Path) error {
	return Paths.InitPaths(cfg)
}

// initPaths sets the default paths in the configuration based on CLI flags,
// configuration file and default values.
func (paths *Path) initPaths(cfg *Path) error {
	paths.Home = cfg.Home
	paths.Config = cfg.Config
	paths.Data = cfg.Data
	paths.Logs = cfg.Logs

	// overwrite paths from CLI flags
	if homePath != nil && len(*homePath) > 0 {
		paths.Home = *homePath
	}
	if configPath != nil && len(*configPath) > 0 {
		paths.Config = *configPath
	}
	if dataPath != nil && len(*dataPath) > 0 {
		paths.Data = *dataPath
	}
	if logsPath != nil && len(*logsPath) > 0 {
		paths.Logs = *logsPath
	}

	// default for the home path is the binary location
	if len(paths.Home) == 0 {
		var err error
		paths.Home, err = filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			return fmt.Errorf("The absolute path to %s could not be obtained. %v",
				os.Args[0], err)
		}
	}

	// default for config path
	if len(paths.Config) == 0 {
		paths.Config = paths.Home
	}

	// default for data path
	if len(paths.Data) == 0 {
		paths.Data = filepath.Join(paths.Home, "data")
	}

	// default for logs path
	if len(paths.Logs) == 0 {
		paths.Logs = filepath.Join(paths.Home, "logs")
	}

	return nil
}

// Resolve resolves a path to a location in one of the default
// folders. For example, Resolve(Home, "test") returns an absolute
// path for "test" in the home path.
func (paths *Path) Resolve(fileType FileType, path string) string {
	// absolute paths are not changed
	if filepath.IsAbs(path) {
		return path
	}

	switch fileType {
	case Home:
		return filepath.Join(paths.Home, path)
	case Config:
		return filepath.Join(paths.Config, path)
	case Data:
		return filepath.Join(paths.Data, path)
	case Logs:
		return filepath.Join(paths.Logs, path)
	default:
		panic(fmt.Sprintf("Unknown file type: %s", fileType))
	}
}

// Resolve resolves a path to a location in one of the default
// folders. For example, Resolve(Home, "test") returns an absolute
// path for "test" in the home path.
func Resolve(fileType FileType, path string) string {
	return Paths.Resolve(fileType, path)
}

// String returns a textual representation
func (paths *Path) String() string {
	return fmt.Sprintf("Home path: [%s] Config path: [%s] Data path: [%s] Logs path: [%s]",
		paths.Home, paths.Config, paths.Data, paths.Logs)
}
