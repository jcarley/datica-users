package jsonutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Values struct {
	data map[string]interface{}
}

func NewValues() *Values {
	return &Values{
		data: make(map[string]interface{}),
	}
}

func (this *Values) Add(key string, value interface{}) error {
	if key == "" {
	}

	if value == nil {
	}

	delete(this.data, key)
	this.Set(key, value)

	return nil
}

func (this *Values) Set(key string, value interface{}) error {
	if key == "" {
	}

	if value == nil {
	}

	this.data[key] = value

	return nil
}

func (this *Values) Del(key string, value interface{}) error {
	if key == "" {
	}

	if value == nil {
	}

	delete(this.data, key)

	return nil
}

func (this *Values) Data() map[string]interface{} {
	return this.data
}

func (this *Values) Encode() ([]byte, error) {
	return EncodeJSON(this.data)
}

func (this *Values) EncodeToWriter(writer io.Writer) error {
	return EncodeJSONToWriter(writer, this.data)
}

func Decode(data []byte) (*Values, error) {
	var output map[string]interface{}
	if err := DecodeJSON(data, &output); err != nil {
		panic(err)
	}

	v := &Values{
		data: output,
	}

	return v, nil
}

// Encodes/Marshals the given object into JSON
func EncodeJSON(v interface{}) ([]byte, error) {
	var buf bytes.Buffer
	if err := EncodeJSONToWriter(&buf, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func EncodeJSONToString(v interface{}) (string, error) {
	data, err := EncodeJSON(v)
	return string(data), err
}

func EncodeJSONToWriter(writer io.Writer, v interface{}) error {
	encoder := json.NewEncoder(writer)
	return encoder.Encode(v)
}

// Decodes/Unmarshals the given JSON into a desired object
func DecodeJSON(data []byte, out interface{}) error {
	if data == nil {
		return fmt.Errorf("'data' being decoded is nil")
	}
	if out == nil {
		return fmt.Errorf("output parameter 'out' is nil")
	}

	return DecodeJSONFromReader(bytes.NewReader(data), out)
}

// Decodes/Unmarshals the given io.Reader pointing to a JSON, into a desired object
func DecodeJSONFromReader(r io.Reader, out interface{}) error {
	if r == nil {
		return fmt.Errorf("'io.Reader' being decoded is nil")
	}
	if out == nil {
		return fmt.Errorf("output parameter 'out' is nil")
	}

	dec := json.NewDecoder(r)

	// While decoding JSON values, intepret the integer values as `json.Number`s instead of `float64`.
	dec.UseNumber()

	// Since 'out' is an interface representing a pointer, pass it to the decoder without an '&'
	return dec.Decode(out)
}
