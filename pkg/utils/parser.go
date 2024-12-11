package utils

import (
	"encoding/json"
	"net/http"
)

func Decode(r *http.Request, arg any) error {
	err := json.NewDecoder(r.Body).Decode(arg)
	if err != nil {
		return err
	}
	return nil
}

func Encode(r http.ResponseWriter, arg any) error {
	err := json.NewEncoder(r).Encode(arg)
	if err != nil {
		return err
	}
	return nil
}
