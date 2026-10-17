package tokens

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/z3vxo/TokenJack/internal/api"
	"github.com/z3vxo/TokenJack/internal/auth"
	"github.com/z3vxo/TokenJack/internal/db"
)

type Handler struct {
	db   *db.DB
	auth *auth.Auth
}

type RefreshReq struct {
	TokenId string `json:"token_id"`
}

func New(d *db.DB, a *auth.Auth) *Handler {
	return &Handler{db: d, auth: a}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(h.auth.AuthWare())

	// id is unique string value in db
	r.Get("/list/{id}", h.ListTokens)

	r.Post("/delete/{id}/{token_id}", h.DeleteToken)
	r.Post("/delete/{id}/all", h.DeleteAllTokens)

	r.Post("/refresh/{id}", h.RefreshToken)
	return r
}

type ListTokensResp struct {
	Total  int        `json:"total"`
	Tokens []db.Token `json:"tokens"`
}

func redact(token string) string {
	if len(token) <= 8 {
		return token
	}

	return token[:4] + strings.Repeat("*", len(token)-8) + token[len(token)-4:]
}

func (h *Handler) ListTokens(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tokens, err := h.db.ListTokens(id)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		api.Error(w, http.StatusInternalServerError, "failed retriving tokens")
		return
	}

	if len(tokens) == 0 {
		api.WriteJSON(w, http.StatusOK, map[string]any{
			"total":  0,
			"tokens": map[string]any{},
		})
		return
	}

	respTokens := make([]db.Token, 0, len(tokens))
	i := 1
	for _, t := range tokens {
		respTokens = append(respTokens, db.Token{
			Id:             i,
			AccessToken:    redact(t.AccessToken),
			AccessTokenLen: len(t.AccessToken),
			RefreshToken:   redact(t.RefreshToken),
			RfreshTokenLen: len(t.RefreshToken),
			Created_at:     t.Created_at,
			TenantId:       t.TenantId,
		})
		i += 1
	}

	data := &ListTokensResp{
		Total:  len(respTokens),
		Tokens: respTokens,
	}

	api.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) DeleteToken(w http.ResponseWriter, r *http.Request) {
	unqiueID := chi.URLParam(r, "id")
	tokenID := chi.URLParam(r, "token_id")

	if err := h.db.DeleteToken(unqiueID, tokenID); err != nil {
		api.Error(w, http.StatusInternalServerError, "failed deleting")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) DeleteAllTokens(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.db.DeleteAllTokens(id); err != nil {
		api.Error(w, http.StatusInternalServerError, "failed deleting")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req RefreshReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.Error(w, http.StatusInternalServerError, "failed decoding json")
		return
	}

	info, err := h.db.GetToken(id, req.TokenId)
	if err != nil {
		api.Error(w, http.StatusInternalServerError, "failed retreving token")
		return

	}

	_, _, err = h.http_RefreshToken(info.Refresh, info.TenantId)

	w.WriteHeader(http.StatusNoContent)
}
