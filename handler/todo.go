package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.

func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {

	res, err:= h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		log.Println("Error creating TODO: ", err)
		return nil, err
	}
	return &model.CreateTODOResponse{TODO: res},nil
}
// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	res, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		log.Println("Error reading TODO: ", err)
		return nil, err
	}

	return &model.ReadTODOResponse{TODOs: res}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	res, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		log.Println("Error updating TODO: ", err)
		return nil, err
	}
	return &model.UpdateTODOResponse{TODO: res}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		{req := &model.CreateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			log.Println("Error decoding request: ", err)
			http.Error(w, "error", http.StatusBadRequest)
			return
			}
	
		if req.Subject == "" {
			log.Println("Invalid subject")
			http.Error(w, "error", http.StatusBadRequest)
			return
			}
		res, err := h.Create(r.Context(), req)
		if err != nil {
			log.Println("Error creating TODO: ", err)
			http.Error(w, "error", http.StatusBadRequest)
			return
			}
		
		jsonRes, err := json.Marshal(res)
		w.Write(jsonRes)
		if err != nil {
			log.Println("Error encoding TODO for logging: ", err)
			}
		}
	case http.MethodPut:
		{
			req := &model.UpdateTODORequest{}
			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
				log.Println("Error decoding request: ", err)
				http.Error(w, "error", http.StatusBadRequest)
				return
			}
			if req.Subject == ""  || req.ID == 0{
				log.Println("Invalid subject or ID")
				http.Error(w, "error", http.StatusBadRequest)
				return
			}
			res, err := h.Update(r.Context(), req)
			if err != nil {
				log.Println("Error updating TODO: ", err)
				http.Error(w, "error", http.StatusNotFound)
				return
			}
			jsonRes, err := json.Marshal(res)
			w.Write(jsonRes)
			if err != nil {
				log.Println("Error encoding TODO for logging: ", err)
			}
		}
	case http.MethodGet:
		{
			prevIDStr := r.URL.Query().Get("prev_id")
			sizeStr := r.URL.Query().Get("size")
			var prevID int64
			var err error
			if prevIDStr != "" {
				prevID, err = strconv.ParseInt(prevIDStr, 10, 64)
				if err != nil {
					log.Println("Error parsing prev_id: ", err)
					http.Error(w, "error", http.StatusBadRequest)
					return
				}
			}

			var size int64 = 5
			if sizeStr != "" {
				size, err = strconv.ParseInt(sizeStr, 10, 64)
				if err != nil {
					log.Println("Error parsing size: ", err)
					http.Error(w, "error", http.StatusBadRequest)
					return
				}
			}

			req := &model.ReadTODORequest{
				PrevID: prevID,
				Size: size,
			}
			res, err := h.Read(r.Context(), req)
			if err != nil {
				log.Println("Error updating TODO: ", err)
				http.Error(w, "error", http.StatusNotFound)
				return
			}
			jsonRes, err := json.Marshal(res)
			w.Write(jsonRes)
			if err != nil {
				log.Println("Error encoding TODO for logging: ", err)
			}
		}
		

	default:
		{
			log.Println("Invalid method")
			http.Error(w, "error", http.StatusBadRequest)
			return
		}

	}
	
}