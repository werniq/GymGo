package main

import (
	"encoding/json"
	"net/http"
)

// elementInArray check if integers is in array
// will be used for generating workouts, without repeated exercises
func (web *webapp) elementInArray(length int, array []int, element int) bool {
	for i := 0; i < length; i++ {
		if array[i] == element {
			return true
		}
	}
	return false
}

// elementInArray check if string element is in array
// will be used for generating workouts
func (web *webapp) stringElementInArray(length int, array []string, element string) bool {
	for i := 0; i < length; i++ {
		if array[i] == element {
			return true
		}
	}
	return false
}

// writeJSON writes arbitrary data out as JSON
func (web *webapp) writeJSON(w http.ResponseWriter, status int, data interface{}, headers ...http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		for k, v := range headers[0] {
			w.Header()[k] = v
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(out)

	return nil
}
