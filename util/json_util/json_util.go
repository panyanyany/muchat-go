package json_util

import (
	"encoding/json"
	"fmt"
)

func StructToMap(obj interface{}) (m map[string]interface{}, err error) {
	var bs []byte
	m = make(map[string]interface{})
	bs, err = json.Marshal(obj)
	if err != nil {
		err = fmt.Errorf("json.Marshal: %w", err)
		return
	}

	err = json.Unmarshal(bs, &m)
	if err != nil {
		err = fmt.Errorf("json.Unmarshal: %w", err)
		return
	}

	return
}
func MustStructToMap(obj interface{}) (m map[string]interface{}) {
	var err error
	m, err = StructToMap(obj)
	if err != nil {
		panic(err)
	}
	return
}
