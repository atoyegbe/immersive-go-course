package main

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"net/http"
	"os"

	"golang.org/x/time/rate"
)

func main() {

	AUTH_USERNAME := os.Getenv("AUTH_USERNAME")
	AUTH_PASSWORD := os.Getenv("AUTH_PASSWORD")
	limiter := rate.NewLimiter(100, 30)

	http.HandleFunc("/limited", func(w http.ResponseWriter, r *http.Request) {
		if limiter.Allow() {
			w.Write([]byte("Not rate limited!"))
			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})

	http.HandleFunc("/authenticated", func(w http.ResponseWriter, r *http.Request) {
		user, password, ok := r.BasicAuth()
		if ok {
			if AUTH_USERNAME != user || AUTH_PASSWORD != password {
				http.Error(w, "Unable to authenticate user", http.StatusUnauthorized)
				return
			}

			w.WriteHeader(http.StatusAccepted)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
		}

	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			queryParams := r.URL.Query()
			fmt.Fprint(w, "<!DOCTYPE html> <html> <em>Hello, world</em><p>Query parameters: <ul>")
			for key, values := range queryParams {
				for _, value := range values {
					res := html.EscapeString(value)
					fmt.Fprintf(w, "<li>%s: [%s]</li>", key, res)
				}
			}
			fmt.Fprint(w, "</ul> </html>")
		}

		if r.Method == "POST" {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "Error parsing form data", http.StatusInternalServerError)
				return
			}

			w.Header().Add("Content-Type", "text/html")
			var buffer bytes.Buffer
			buffer.Write(body)
			stringValue := buffer.String()
			fmt.Fprint(w, html.EscapeString(stringValue))
		}

		w.WriteHeader(http.StatusAccepted)
	})

	http.HandleFunc("/200", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	})

	http.Handle("/404", http.NotFoundHandler())

	http.HandleFunc("/500", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Internal service error", http.StatusInternalServerError)
	})

	http.ListenAndServe(":8000", nil)
}
