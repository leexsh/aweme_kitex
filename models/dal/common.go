package dal

import (
	"encoding/json"
	"fmt"
)

func marshal(v any) string {
	ub, err := json.Marshal(v)
	if err != nil {
		fmt.Printf("%s\n", err.Error())
	}
	return string(ub)
}

func unmarshal(data string, v any) error {
	return json.Unmarshal([]byte(data), v)
}
