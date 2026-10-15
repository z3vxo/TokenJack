package server

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (s *Server) HandleListCampaigns(w http.ResponseWriter, r *http.Request) {
	campaigns, err := s.DB.GetCampaigns()
	if err != nil {
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "database failed"})
		return
	}

	s.WriteJson(w, http.StatusOK, CampaignsListResp{
		Total:     len(campaigns),
		Campaigns: campaigns,
	})
}

func (s *Server) HandleNewCampaign(w http.ResponseWriter, r *http.Request) {
	var req CampaingNewReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed decoding json"})
		return
	}

	if req.TenantId == "" || req.Name == "" || len(req.Resource.Scopes) == 0 {
		s.WriteJson(w, http.StatusBadRequest, map[string]string{"error": "invalid json body"})
		return
	}

	const maxAttempts = 10
	var err error
	var UrlString string

	for range maxAttempts {
		UrlString = genString()
		err := s.DB.InsertNewCampaign(UrlString, req.Name, req.TenantId, req.Resource.Name, req.Resource.Scopes)

		if err == nil {
			break
		}
		if !strings.Contains(err.Error(), "UNIQUE constraint failed") {
			s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "database failed"})
			return
		}

	}
	if err != nil {
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed generating unique string"})
		return
	}

	s.WriteJson(w, http.StatusOK, map[string]string{"id": UrlString})
}
