package main

import (
	"fmt"
	"github.com/huytq/tech-core/read_configuration"
)

func main() {
	configProvider := read_configuration.NewConfigProvider(read_configuration.ConfigEnv{
		FileName: "env-yaml",
		FileType: read_configuration.YAML,
	})
	fmt.Println(configProvider.GetValue("KEY_DATA"))
}
