package testutils

import "os"

func AcceptanceEnabed() bool {
	return os.Getenv("ACC_ENABLE") == "1"
}
