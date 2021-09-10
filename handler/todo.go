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
	todo, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	return &model.CreateTODOResponse{TODO: *todo}, err
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	todo, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	return &model.ReadTODOResponse{TODOs: todo}, err
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	todo, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	return &model.UpdateTODOResponse{TODO: *todo}, err
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	err := h.svc.DeleteTODO(ctx, req.IDs)
	return &model.DeleteTODOResponse{}, err
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
	case "GET":
		params := r.URL.Query()
		PrevId := params.Get("prev_id")
		Size := params.Get("size")
		var req model.ReadTODORequest
		var err error
		if PrevId == "" {
			req.PrevID = 0
		} else {
			PrevIdToInt, err := strconv.ParseInt(PrevId, 10, 64)
			if err != nil {
				log.Fatal(err)
				return
			}
			req.PrevID = PrevIdToInt
		}

		if Size == "" {
			req.Size = 0
		} else {
			SizeToInt, err := strconv.ParseInt(Size, 10, 64)
			if err != nil {
				log.Fatal(err)
				return
			}
			req.Size = SizeToInt
		}
		res, err := h.Read(r.Context(), &req)
		if err != nil {
			log.Fatal(err)
			return
		}

		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Fatal(err)
			return
		}
	case "DELETE":
		body := r.Body
		defer body.Close()
		var req model.DeleteTODORequest
		if err := json.NewDecoder(body).Decode(&req); err != nil {
			log.Fatal(err)
			return
		}

		if len(req.IDs) == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		res, err := h.Delete(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		id
		if err = json.NewEncoder(w).Encode(&res); err != nil {
			log.Fatal(err)
			return
		}
	}
}
