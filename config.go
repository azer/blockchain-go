package blockchain

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

var config ConfigSpec

type ConfigSpec struct {
	DBFile       string `envconfig:"DB_FILE"`
	BlocksBucket string `split_words:"true"`
	TargetBits   int    `split_words:"true"`
}

func init() {
	if err := godotenv.Load(); err != nil {
		panic(".env missing")
	}

	if err := envconfig.Process("bc", &config); err != nil {
		panic(err)
	}
}
