package internal

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strings"
)

var envVarPrefix = "USERSRV"

type RawAppConfig struct {
	Server   RawAppServerConfig
	Database RawAppDatabaseConfig
}

type RawAppDatabaseConfig struct {
	Host     string
	Port     uint16
	Name     string
	User     string
	Password string
	SslMode  string
}

type RawAppServerConfig struct {
	Port uint16
}

func LoadAppConfig(deployment string) *RawAppConfig {
	v := viper.New()
	applyEnvVariables(v)
	v.SetConfigName("config-" + deployment)
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs/usersrv")
	err := v.ReadInConfig()
	if err != nil {
		logrus.Fatalln("error reading configuration file:", err)
	}
	var rawCfg RawAppConfig
	err = v.UnmarshalExact(&rawCfg)
	if err != nil {
		logrus.Fatalln("error parsing configuration file:", err)
	}
	return &rawCfg
}

func applyEnvVariables(v *viper.Viper) {
	v.SetEnvPrefix(envVarPrefix)
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
}
