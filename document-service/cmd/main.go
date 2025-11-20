package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/BHAV0207/documet-service/internal/handler"
	"github.com/BHAV0207/documet-service/internal/repository"
	"github.com/BHAV0207/documet-service/internal/websockets"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("no env found")
	}

	// 1) Initialize DB
	repository.Init() // this populates db.DB
	// get the gorm DB instance for handler
	d := repository.DB
	if d == nil {
		log.Fatal("database not initialized")
	}

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "7000"
	}

	hub := websockets.NewHub()
	go hub.Run()

	// 3) Create Handler (DI container)
	h := handler.NewHandler(d, hub)

	r := mux.NewRouter()

	r.HandleFunc("/documents", h.CreateDocument).Methods("POST")
	r.HandleFunc("/documents/{id}", h.GetDocument).Methods("GET")

	// Snapshots
	r.HandleFunc("/documents/{id}/snapshot", h.UploadSnapshot).Methods("POST")
	r.HandleFunc("/documents/{id}/snapshot", h.GetLatestSnapshot).Methods("GET")

	// WebSocket route -> now uses Handler.ServeWS
	r.HandleFunc("/ws/document/{id}", h.ServeWS)

	// r.HandleFunc("/ws/document/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	vars := mux.Vars(r)
	// 	DocID := vars["id"]

	// 	h.ServeWS(w, r, DocID)
	// })

	log.Println("📄 Document Collaboration Service running on", PORT)

	if err := http.ListenAndServe(":"+PORT, r); err != nil {
		log.Fatal("Server Error:", err)
	}

}
