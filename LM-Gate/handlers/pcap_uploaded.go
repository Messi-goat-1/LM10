package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
)

type PCAPAnalyzeHandler struct {
	pcapService *services.PCAPService
}

func NewPCAPAnalyzeHandler(s *services.PCAPService) *PCAPAnalyzeHandler {
	return &PCAPAnalyzeHandler{pcapService: s}
}

func (h *PCAPAnalyzeHandler) Handle(event events.PCAPAnalyzeEvent) error {
	return h.pcapService.Analyze(event.FileID, event.FilePath)
}
