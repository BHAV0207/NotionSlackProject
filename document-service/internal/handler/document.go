package handler

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/BHAV0207/documet-service/internal/websockets"
	"github.com/BHAV0207/documet-service/pkg/models"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func (h *Handler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	title := r.FormValue("ttle")

	doc := models.Document{
		ID:    uuid.New().String(),
		Title: title,
	}

	if err := h.DB.Create(&doc).Error; err != nil {
		http.Error(w, "failed to create document", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, _ = w.Write([]byte(`{"id":"` + doc.ID + `"}`))
}

func (h *Handler) GetDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var doc models.Document
	if err := h.DB.First(&doc, "id = ?", id).Error; err != nil {
		http.Error(w, "document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"id":"` + doc.ID + `","title":"` + doc.Title + `"}`))
}

func (h *Handler) UploadSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	var doc models.Snapshot
	if err := h.DB.First(&doc, "id = ?", docID).Error; err != nil {
		http.Error(w, "document not found", http.StatusNotFound)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	snap := models.Snapshot{
		ID:         uuid.New().String(),
		DocumentID: docID,
		Data:       data,
		CreatedAt:  time.Now().UTC(),
	}

	if err := h.DB.Create(&snap).Error; err != nil {
		http.Error(w, "failed to save snapshot", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) GetLatestSnapshot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	var snap models.Snapshot
	if err := h.DB.
		Where("document_id = ?", docID).
		Order("created_at desc").
		Limit(1).
		First(&snap).Error; err != nil {
		http.Error(w, "snapshot not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.Itoa(len(snap.Data)))
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(snap.Data)
}

func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	conn, err := websockets.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		// upgrader returns error when not a websocket request; respond 400
		http.Error(w, "failed to upgrade to websocket", http.StatusBadRequest)
		return
	}

	client := &websockets.Client{
		Hub:   h.Hub,
		Conn:  conn,
		DocID: docID,
	}

	// register client and start read pump
	h.Hub.Register <- client
	go client.ReadPump()
}
