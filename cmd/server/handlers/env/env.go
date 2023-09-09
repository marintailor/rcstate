package env

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/marintailor/rcstate/cmd/api/env"
)

// Down is a handler function to stop all resources in the environment(s).
func Down(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("read body:", err)
		}

		cfg := env.Config{}
		if err := cfg.GetConfig(body); err != nil {
			log.Println("get config:", err)
		}

		json, err := cfg.Down()
		if err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(json)); err != nil {
			log.Printf("write to response: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("{\"error\":\"method not allowed\"}")); err != nil {
			log.Printf("write to response: %v", err)
		}
	}
}

// Show is a handler function to show all resources in the environment(s).
func Show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("read body:", err)
		}

		cfg := env.Config{}
		if err := cfg.GetConfig(body); err != nil {
			log.Println("get config:", err)
		}

		json, err := cfg.Show()
		if err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		if json == "null" {
			if cfg.Name != "" {
				json = fmt.Sprintf("environment %q with label %q not found", cfg.Name, cfg.Label)
			}
			if cfg.All {
				json = fmt.Sprintf("no environment with label %q was found", cfg.Label)
			}
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(json)); err != nil {
			log.Printf("write to response: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("{\"error\":\"method not allowed\"}")); err != nil {
			log.Printf("write to response: %v", err)
		}
	}
}

// Down is a handler function to start all resources in the environment(s).
func Up(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("read body:", err)
		}

		cfg := env.Config{}
		if err := cfg.GetConfig(body); err != nil {
			log.Println("get config:", err)
		}

		json, err := cfg.Up()
		if err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(json)); err != nil {
			log.Printf("write to response: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("{\"error\":\"method not allowed\"}")); err != nil {
			log.Printf("write to response: %v", err)
		}
	}
}
