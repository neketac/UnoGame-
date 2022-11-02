package help

import (
	"encoding/json"
	"net/http"
)

func RenderAndWrite(w http.ResponseWriter, read interface{}) {
	js, err := json.MarshalIndent(read, "", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
