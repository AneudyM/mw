package base

import (
	"log"
	"os"
	"path/filepath"
)

// consider adding a rule for the build directory. If user is not in the
// project's root directory and issues the mw build command, InRoot()
// will look for the existence of .mwrc, and src. If either .mwrc or src
// exist in the current directory then false¹ will be returned. If
//
// ¹create a custom error type instead of boolean

// InRoot checks whether the build subcommand was called from within the
// project's root directory.
func InRoot() bool {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	if _, err := os.Stat(filepath.Join(cwd, MWROOT[CFG_FILE])); os.IsNotExist(err) {
		return false
	}

	f, err := os.Stat(filepath.Join(cwd, MWROOT[DIR_SRC]))
	if os.IsNotExist(err) || !f.IsDir() {
		return false
	}
	return true
}
