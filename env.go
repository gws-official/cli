package cli

import "os"

func lookupEnv(name string) (string, bool) {
	return os.LookupEnv(name)
}