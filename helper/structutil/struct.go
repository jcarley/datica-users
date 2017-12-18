package structutil

import "github.com/mitchellh/mapstructure"

func Decode(m interface{}, rawVal interface{}) error {
	md := &mapstructure.Metadata{}
	d, _ := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Metadata: md,
		Result:   m,
		TagName:  "json",
	})
	return d.Decode(rawVal)
}
