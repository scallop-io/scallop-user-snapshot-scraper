package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func PopulateModelFromBody(body io.ReadCloser, model interface{}) error {
	bytes, err := ioutil.ReadAll(io.LimitReader(body, 1048576))
	if err != nil {
		return err
	}
	if err := body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, model); err != nil {
		return err
	}
	return nil
}
