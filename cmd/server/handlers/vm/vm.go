package vm

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/marintailor/rcstate/cmd/api/vm"
)

// List is a handler function to list all virtual machines.
func List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("read body:", err)
		}

		cfg := vm.Config{}
		if err := cfg.GetConfig(body); err != nil {
			log.Println("get config:", err)
		}

		vm, err := vm.NewVirtualMachine(cfg.Project, cfg.Zone)
		if err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		list := vm.List()
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte(list)); err != nil {
			log.Printf("write to response: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("{\"error\":\"method not allowed\"}")); err != nil {
			log.Printf("could not write to response: %v", err)
		}
	}
}

// Start is a handler function to start a virtual machine.
func Start(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("read body:", err)
		}

		cfg := vm.Config{}
		if err := cfg.GetConfig(body); err != nil {
			log.Println("get config:", err)
		}

		vm, err := vm.NewVirtualMachine(cfg.Project, cfg.Zone)
		if err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		if err := vm.Start(cfg.Name); err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		if cfg.DNS.RecordName != "" {
			dnsRecord := fmt.Sprintf("%s.%s", cfg.DNS.RecordName, cfg.DNS.Domain)
			cfg.Record(dnsRecord)
		}

		if cfg.Script.CMD != "" {
			cfg.ExecuteScript()
		}

		w.WriteHeader(http.StatusOK)
		msg := "{ \"status\": \"success\"}"
		if _, err := w.Write([]byte(msg)); err != nil {
			log.Printf("write to response: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("{\"error\":\"method not allowed\"}")); err != nil {
			log.Printf("could not write to response: %v", err)
		}
	}
}

// Status is a handler function that returns the status of a virtual machine.
func Status(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("read body:", err)
		}

		cfg := vm.Config{}
		if err := cfg.GetConfig(body); err != nil {
			log.Println("get config:", err)
		}

		vm, err := vm.NewVirtualMachine(cfg.Project, cfg.Zone)
		if err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		status, err := vm.Status(cfg.Name)
		if err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		w.WriteHeader(http.StatusOK)
		msg := fmt.Sprintf("{ \"status\": \"%s\"}", status)
		if _, err := w.Write([]byte(msg)); err != nil {
			log.Printf("write to response: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("{\"error\":\"method not allowed\"}")); err != nil {
			log.Printf("could not write to response: %v", err)
		}
	}
}

// Stop is a handler function to stop a virtual machine.
func Stop(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodPost:
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("read body:", err)
		}

		cfg := vm.Config{}
		if err := cfg.GetConfig(body); err != nil {
			log.Println("get config:", err)
		}

		vm, err := vm.NewVirtualMachine(cfg.Project, cfg.Zone)
		if err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		if err := vm.Stop(cfg.Name); err != nil {
			msg := fmt.Sprintf("{ \"error\": \"%s\"}", err)
			w.WriteHeader(http.StatusInternalServerError)
			if _, err := w.Write([]byte(msg)); err != nil {
				log.Printf("write to response: %v", err)
			}
		}

		w.WriteHeader(http.StatusOK)
		msg := "{ \"status\": \"success\"}"
		if _, err := w.Write([]byte(msg)); err != nil {
			log.Printf("write to response: %v", err)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		if _, err := w.Write([]byte("{\"error\":\"method not allowed\"}")); err != nil {
			log.Printf("could not write to response: %v", err)
		}
	}
}
