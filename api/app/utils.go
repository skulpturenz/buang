package app

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func WriteText(w http.ResponseWriter, ret any, code int) error {
	switch ret.(type) {
	case string:
		break
	default:
		panic("res not a string")
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)

	_, err := fmt.Fprintf(w, "%s", ret.(string))
	if err != nil {
		return err
	}

	return nil
}

func WriteJson(w http.ResponseWriter, ret any, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(ret)
	if err != nil {
		return err
	}

	return nil
}

func WriteError(w http.ResponseWriter, errs map[string]string, code int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	err := json.NewEncoder(w).Encode(errs)
	if err != nil {
		return err
	}

	return nil
}
