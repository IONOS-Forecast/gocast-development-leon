package utils

import (
	"log"

	"github.com/IONOS-Forecast/gocast-development-leon/Gocast/pkg/model"
	"github.com/jessevdk/go-flags"
	"github.com/joho/godotenv"
)

var Options model.Options

var (
	FdbUser,
	FdbPass,
	FdbDB,
	FdbAddress string
)

func DotEnvironment() model.Options {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	err = godotenv.Load(".env.dev")
	if err != nil {
		log.Fatal(err)
	}
	_, err = flags.Parse(&Options)
	if err != nil {
		log.Fatal(err)
	}
	FdbUser = Options.FdbUser
	FdbPass = Options.FdbPassword
	FdbDB = Options.FdbDatabase
	FdbAddress = Options.FdbAddress
	return Options
}
