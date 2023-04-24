package config

import (
	"log"
	"os"
	"strings"

	"github.com/danztran/sample-admission-validating-webhook/pkg/server"
	"github.com/spf13/viper"
)

var Values Config

type Config struct {
	Server server.Config `mapstructure:"server"`
}

func init() {
	config := viper.New()
	config.SetConfigName("sample-admission-validating-webhook-config") // config file name
	if configPath, ok := os.LookupEnv("ISTIO_GUARD_CONFIG"); ok {
		config.AddConfigPath(configPath)
	}
	config.AddConfigPath(".")
	config.AddConfigPath("./config/")
	config.AddConfigPath("../config/")
	config.AddConfigPath("../../config/")
	config.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	config.AutomaticEnv()

	err := config.ReadInConfig()
	if err != nil {
		log.Fatalf("error read config / %s", err)
	}

	err = config.Unmarshal(&Values)
	if err != nil {
		log.Fatalf("error parse config / %s", err)
	}
}
