package rollups

import (
	"encoding/json"
	"fmt"
	"log"
)

type AdvanceHandlerFunc func(payload []byte, metadata Metadata) error

type Router struct {
	AdvanceHandlers map[string]AdvanceHandlerFunc
}

func NewRouter() *Router {
	return &Router{
		AdvanceHandlers: make(map[string]AdvanceHandlerFunc),
	}
}

func (r *Router) HandleAdvance(path string, handler AdvanceHandlerFunc) {
	r.AdvanceHandlers[path] = handler
}

func (r *Router) Advance(payload []byte, metadata Metadata) error {
	log.Println("Router: Advance", string(payload))
	var input Input
	if err := json.Unmarshal(payload, &input); err != nil {
		return err
	}
	handler, ok := r.AdvanceHandlers[input.Path]
	if !ok {
		return fmt.Errorf("handler: path not found: %s", input.Path)
	}
	if err := handler(input.Payload, metadata); err != nil {
		return err
	}
	return nil
}
