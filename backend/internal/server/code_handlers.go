package server

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
)

type GenCodeReq struct {
	Id string `json:"id"`
}
type GenCodeResp struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Url     string `json:"url"`
}

func (s *Server) GenerateCode(w http.ResponseWriter, r *http.Request) {
	var req GenCodeReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed decoding json"})
		return
	}

	CodeInfo, err := s.DB.GetCodeInfo(req.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s.WriteJson(w, http.StatusNotFound, map[string]string{"error": "no entry found for provided id"})
			return
		}
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "database error"})
		return
	}

	code, url, err := s.Codes.InsertAndBeginPoll(CodeInfo)
	if err != nil {
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed inserting and polling code"})
		return
	}

	s.DB.UpdateClick(req.Id)
	data := &GenCodeResp{
		Message: "polling",
		Code:    code,
		Url:     url,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(data)
	return

}
