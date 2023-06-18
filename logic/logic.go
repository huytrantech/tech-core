package logic

import "encoding/json"

func ParseStringToStruct(str string, model interface{}) error {
	if err := json.Unmarshal([]byte(str), &model); err != nil {
		return err
	}
	return nil
}
