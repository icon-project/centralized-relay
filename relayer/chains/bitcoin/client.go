package bitcoin

import (
	"os"
)

func RunApp() {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "master" {
		startMaster()
	} else {
		startSlave()
	}
}
