package test

import (
	"encoding/json"
	"testing"

	"github.com/jcarley/datica-users/helper/structutil"
)

func GetRawData(t *testing.T, buffer []byte) (data map[string]interface{}) {
	err := json.Unmarshal(buffer, &data)
	if err != nil {
		t.Fatalf("Failed to unmarshal buffer: ", err)
	}
	return
}

func SetRawData(t *testing.T, v interface{}) (bytes []byte) {
	bytes, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to marshal v: ", err)
	}
	return bytes
}

func DecodeStruct(t *testing.T, m interface{}, rawVal interface{}) {
	err := structutil.Decode(m, rawVal)
	if err != nil {
		t.Fatalf("Failed to decode struct: ", err)
	}
}
