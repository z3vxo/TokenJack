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

	if err := s.Codes.InsertAndBeginPoll(CodeInfo); err != nil {
		s.WriteJson(w, http.StatusInternalServerError, map[string]string{"error": "failed inserting and polling code"})
		return
	}

	s.WriteJson(w, http.StatusOK, map[string]string{"status": "polling"})
	return

}
