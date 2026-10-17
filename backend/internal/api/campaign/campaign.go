package campaign

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/z3vxo/TokenJack/internal/api"
	"github.com/z3vxo/TokenJack/internal/auth"
	"github.com/z3vxo/TokenJack/internal/db"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

type Handler struct {
	DB   *db.DB
	Auth *auth.Auth
}

func New(d *db.DB, a *auth.Auth) *Handler {
	return &Handler{DB: d, Auth: a}
}

func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(h.Auth.AuthWare())
	r.Get("/list", h.List)
	r.Post("/new", h.Create)
	r.Get("/meta/{id}", h.Meta)
	return r
}

type NewCampaignReq struct {
	TenantId string         `json:"tenant_id"`
	Name     string         `json:"name"`
	Resource ResourceScopes `json:"resources"`
}

type ResourceScopes struct {
	Name   string     `json:"name"`
	Scopes []db.Scope `json:"scopes"`
}

type ListResp struct {
	Total     int           `json:"total"`
	Campaigns []db.Campaign `json:"campaigns"`
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	campaigns, err := h.DB.GetCampaigns()
	if err != nil {
		api.Error(w, http.StatusInternalServerError, "database failed")
		return
	}

	api.WriteJSON(w, http.StatusOK, ListResp{
		Total:     len(campaigns),
		Campaigns: campaigns,
	})
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req NewCampaignReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.Error(w, http.StatusInternalServerError, "failed decoding json")
		return
	}

	if req.TenantId == "" || req.Name == "" || len(req.Resource.Scopes) == 0 {
		api.Error(w, http.StatusBadRequest, "invalid json body")
		return
	}

	const maxAttempts = 10
	var err error
	var unique string

	for range maxAttempts {
		unique = genString()
		err = h.DB.InsertNewCampaign(unique, req.Name, req.TenantId, req.Resource.Name, req.Resource.Scopes)
		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
			api.Error(w, http.StatusInternalServerError, "database failed")
			return
		}
	}
	if err != nil {
		api.Error(w, http.StatusInternalServerError, "failed generating unique string")
		return
	}

	api.WriteJSON(w, http.StatusOK, map[string]string{"id": unique})
}

func (h *Handler) Meta(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	data, err := h.DB.GetCampaignMeta(id)
	if err != nil {
		api.Error(w, http.StatusInternalServerError, "database failed")
		return
	}

	api.WriteJSON(w, http.StatusOK, data)
}

func genString() string {
	b := make([]byte, 14)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		b[i] = alphabet[n.Int64()]
	}
	return string(b)
}
