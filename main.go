/*
Copyright © 2022 Vijay
*/
package main

import (
	"os"

	"github.com/evijayan2/vault-secret-finder/cmd"
	"github.com/rs/zerolog"
)

func main() {

	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	zerolog.TimeFieldFormat = ""
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	cmd.Execute(logger)
}
