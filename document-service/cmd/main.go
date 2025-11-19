package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BHAV0207/documet-service/internal/websockets"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("no env found")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "7000"
	}

	hub := websockets.NewHub()
	go hub.Run()

	r := mux.NewRouter()

	r.HandleFunc("/ws/document/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		DocID := vars["id"]

		websockets.ServerWs(hub, w, r, DocID)
	})

	log.Println("📄 Document Collaboration Service running on", PORT)

	if err := http.ListenAndServe(":"+PORT, r); err != nil {
		log.Fatal("Server Error:", err)
	}

}
