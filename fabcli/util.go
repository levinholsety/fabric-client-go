package fabcli

import (
	"encoding/json"
)

func marshalIndentJSONString(v interface{}) (data []byte) {
	data, _ = json.MarshalIndent(v, "", "    ")
	return
}

func batchExecute(functions ...func() error) (err error) {
	for _, f := range functions {
		err = f()
		if err != nil {
			return
		}
	}
	return
}
