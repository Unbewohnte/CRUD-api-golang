package logs

import (
	"log"
	"os"
	"path/filepath"
)

// Create a logfile in the same directory as executable and set output
// of `log` package to it
func SetUp() error {
	executablePath, err := os.Executable()
	if err != nil {
		return err
	}
	exeDir := filepath.Dir(executablePath)

	logfile, err := os.Create(filepath.Join(exeDir, "logs.log"))
	if err != nil {
		return err
	}

	log.SetOutput(logfile)

	return nil
}
