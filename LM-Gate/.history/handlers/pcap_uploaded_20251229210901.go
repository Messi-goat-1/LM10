package handlers

import (
	"LM-Gate/events"
	"LM-Gate/services"
	"context"
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
func (h *PCAPAnalyzeHandler) Handle(event events.PCAPAnalyzeEvent) error {
	ctx := context.Background()

	// Start analysis in a background goroutine
	go func() {
		err := h.pcapService.Analyze(ctx, event.FileID, event.FilePath)
		if err != nil {
			// Log the error (Analysis failed in background)
			fmt.Printf("Background analysis failed for %s: %v\n", event.FileID, err)
		}
	}()

	return nil // Return immediately to acknowledge the event
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
