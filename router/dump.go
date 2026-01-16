package router

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
)

// HandlerDump tbd
func HandlerDump(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("@@@ dump")
		var bodyBytes []byte
		var err error
		if r.Body != nil {
			bodyBytes, err = io.ReadAll(r.Body)
			if err != nil {
				log.Println("router", "Body reading error", err)
				return
			}
			defer r.Body.Close()
		}
		log.Println("router", "Headers", r.Header)
		if len(bodyBytes) > 0 {
			var prettyJSON bytes.Buffer
			if err = json.Indent(&prettyJSON, bodyBytes, "", "\t"); err != nil {
				log.Println("router", "JSON parse error", err)
				return
			}
			log.Println(string(prettyJSON.String()))
		} else {
			log.Println("Body: No Body Supplied")
		}
		next.ServeHTTP(w, r)
	})
}

// --- EOF
