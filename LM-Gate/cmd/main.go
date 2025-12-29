package main

import (
	"os"

	lmgate "LM-Gate"
)

func main() {
	if err := lmgate.Execute(); err != nil {
		os.Exit(1)
	}
}
