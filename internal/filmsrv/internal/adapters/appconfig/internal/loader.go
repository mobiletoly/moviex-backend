package internal

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var envVarPrefix = "FILMSRV"

type RawAppConfig struct {
	Database struct {
		Host     string
		Port     uint16
		Name     string
		User     string
		Password string
		SslMode  string
	}
}

func LoadAppConfig() RawAppConfig {
	v := viper.New()
	applyEnvVariables(v)
	v.SetConfigName("config-local")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs/filmsrv")
	err := v.ReadInConfig()
	if err != nil {
		logrus.Fatalln("error reading configuration file:", err)
	}
	var rawCfg RawAppConfig
	err = v.UnmarshalExact(&rawCfg)
	if err != nil {
		logrus.Fatalln("error parsing configuration file:", err)
	}
	return rawCfg
}

func applyEnvVariables(v *viper.Viper) {
	v.SetEnvPrefix(envVarPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}
