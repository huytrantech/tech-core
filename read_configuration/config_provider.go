package read_configuration

import (
	"github.com/spf13/viper"
)

type ConfigFileType string

const (
	YAML ConfigFileType = "yaml"
	ENV  ConfigFileType = "env"
	JSON ConfigFileType = "json"
)

type ConfigEnv struct {
	FileName string
	FileType ConfigFileType
}

type IConfigProvider interface {
	GetValue(key string) string
}

type configProvider struct {
	viper *viper.Viper
}

func NewConfigProvider(config ConfigEnv) IConfigProvider {

	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigType(string(config.FileType))
	v.SetConfigName(config.FileName)
	if err := v.ReadInConfig(); err != nil {
		panic(err)
	}
	return &configProvider{viper: v}
}

func (configProvider *configProvider) GetValue(key string) string {
	return configProvider.viper.GetString(key)
}
