package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

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
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: *todo}, err
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	_, _ = h.svc.ReadTODO(ctx, 0, 0)
	return &model.ReadTODOResponse{}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *todo}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		body := r.Body
		defer body.Close()
		var req model.CreateTODORequest
		if err := json.NewDecoder(body).Decode(&req); err != nil {
			log.Fatal(err)
			return
		}

		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		todo, err := h.Create(r.Context(), &req)
		if err != nil {
			log.Fatal(err)
			return
		}
		if err := json.NewEncoder(w).Encode(todo); err != nil {
			log.Fatal(err)
			return
		}
	case "PUT":
		body := r.Body
		defer body.Close()
		var req model.UpdateTODORequest
		if err := json.NewDecoder(body).Decode(&req); err != nil {
			log.Fatal(err)
			return
		}
		if req.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.Update(r.Context(), &req)
		if err != nil {
			log.Fatal(err)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Fatal(err)
			return
		}
	}
}
