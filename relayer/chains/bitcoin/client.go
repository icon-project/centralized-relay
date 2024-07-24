package bitcoin

import "os"

func runApp() {
	goEnv := os.Getenv("GO_ENV")
	if goEnv == "master" {
		startMaster()
	} else {
		startSlave()
	}
}
