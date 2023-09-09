package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/marintailor/rcstate/cmd/server/handlers/env"
	"github.com/marintailor/rcstate/cmd/server/handlers/vm"
)

func NewServer(port string) {
	router := http.NewServeMux()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "root path")
	})

	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{\"status\":\"ok\"}")
	})

	router.HandleFunc("/v1/env/down", env.Down)
	router.HandleFunc("/v1/env/show", env.Show)
	router.HandleFunc("/v1/env/up", env.Up)

	router.HandleFunc("/v1/vm/list", vm.List)
	router.HandleFunc("/v1/vm/start", vm.Start)
	router.HandleFunc("/v1/vm/status", vm.Status)
	router.HandleFunc("/v1/vm/stop", vm.Stop)

	log.Println("Listening on localhost:8080")

	log.Fatal(http.ListenAndServe(":"+port, router))
}
