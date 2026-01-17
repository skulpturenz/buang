package app

import (
	"encoding/json"
	"net/http"
)

func WriteJson(w http.ResponseWriter, ret any, code int) error {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(ret)
	if err != nil {
		return err
	}

	return nil
}

func WriteError(w http.ResponseWriter, errs map[string]string, code int) error {
	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(errs)
	if err != nil {
		return err
	}

	return nil
}
