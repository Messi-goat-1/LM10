package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
	"context"
	"encoding/json"
	"fmt"
)

// PCAPAnalyzeHandler routes PCAP analysis events to the PCAP service.
//
// NOTE: This handler is a thin bridge between the event layer
// and the analysis/service layer.
// It does not contain business logic.
type PCAPAnalyzeHandler struct {
	// pcapService performs the actual PCAP analysis work.
	pcapService *services.PCAPService
}

// NewPCAPAnalyzeHandler creates a new PCAPAnalyzeHandler.
//
// NOTE: The service is injected to keep the handler decoupled.
// TODO: Add nil validation for the service.
func NewPCAPAnalyzeHandler(s *services.PCAPService) *PCAPAnalyzeHandler {
	return &PCAPAnalyzeHandler{pcapService: s}
}

// Handle receives a PCAPAnalyzeEvent and triggers analysis.
//
// NOTE: This is usually called by the dispatcher when routingKey is "pcap.analyze".
// FIXME: No validation is done on FileID/FilePath here.
func (h *PCAPAnalyzeHandler) Handle(data []byte) error {
	var event events.PCAPAnalyzeEvent

	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	ctx := context.Background()

	go func() {
		if err := h.pcapService.Analyze(ctx, event.FileID, event.FilePath); err != nil {
			fmt.Printf(
				"Background analysis failed for %s: %v\n",
				event.FileID,
				err,
			)
		}
	}()

	return nil
}

/*
pcap.analyze event
   ↓
EventDispatcher.Dispatch
   ↓
PCAPAnalyzeHandler.Handle
   ↓
services.PCAPService.Analyze
   ↓
analysis.AnalyzePCAP (pcap_analyzer.go)

*/
