package handler

import (
	"encoding/json"
	"net/http"

	"github.com/progotman/multiplexer/processor"
)

type ProcessUrlsHandler struct {
	Processor processor.UrlsProcessor
}

func (s *ProcessUrlsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	request, err := s.decodeRequest(r)
	if err != nil {
		http.Error(w, "Unprocessable Entity", http.StatusUnprocessableEntity)
		return
	}

	result := s.Processor.Process(r.Context(), request)
	s.encodeResponse(w, result)
}

func (s *ProcessUrlsHandler) decodeRequest(r *http.Request) (processor.ProcessUrlsRequest, error) {
	request := processor.ProcessUrlsRequest{}
	return request, json.NewDecoder(r.Body).Decode(&request)
}

func (s *ProcessUrlsHandler) encodeResponse(w http.ResponseWriter, result *processor.ProcessUrlsResult) {
	w.Header().Set("Content-Type", "application/json")

	encodeError := json.NewEncoder(w).Encode(result)
	if encodeError != nil {
		http.Error(w, "Error encoding data", http.StatusInternalServerError)
		return
	}
}
