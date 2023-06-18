package read_configuration

import (
	"encoding/json"
	"gopkg.in/yaml.v3"
	"io"
	"os"
)

func ReadFileType(fileName string, fileType ConfigFileType, data interface{}) error {

	switch fileType {
	case JSON:
		return readFileJSON(fileName, &data)
	case YAML:
		return readFileYAML(fileName, &data)
	}

	return nil
}
func readFileJSON(fileName string, data interface{}) error {

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(byteValue, &data); err != nil {
		return err
	}

	return nil
}

func readFileYAML(fileName string, data interface{}) error {

	jsonFile, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	if err = yaml.Unmarshal(byteValue, &data); err != nil {
		return err
	}

	return nil
}
