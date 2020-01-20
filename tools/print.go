package tools

import (
	"encoding/json"
	"fmt"
)

func PrettyPrint(i interface{}) {
	res, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", string(res))
}
