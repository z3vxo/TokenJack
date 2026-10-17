package code

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/z3vxo/TokenJack/internal/api"
	"github.com/z3vxo/TokenJack/internal/codes"
	"github.com/z3vxo/TokenJack/internal/db"
)

type Handler struct {
	DB    *db.DB
	Codes *codes.Codes
}

func New(d *db.DB, c *codes.Codes) *Handler {
	return &Handler{DB: d, Codes: c}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/generate", h.Generate)
	return r
}

type GenCodeReq struct {
	Id string `json:"id"`
}

type GenCodeResp struct {
	Code string `json:"code"`
	Url  string `json:"url"`
}

func (h *Handler) Generate(w http.ResponseWriter, r *http.Request) {
	var req GenCodeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.Error(w, http.StatusInternalServerError, "failed decoding json")
		return
	}

	info, err := h.DB.GetCodeInfo(req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			api.Error(w, http.StatusNotFound, "no entry found for provided id")
			return
		}
		api.Error(w, http.StatusInternalServerError, "database error")
		return
	}

	code, url, err := h.Codes.InsertAndBeginPoll(info)
	if err != nil {
		api.Error(w, http.StatusInternalServerError, "failed inserting and polling code")
		return
	}

	h.DB.UpdateClick(req.Id)

	api.WriteJSON(w, http.StatusOK, &GenCodeResp{
		Code: code,
		Url:  url,
	})
}
