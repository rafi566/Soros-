package api

import "time"

type Source struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Destination struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Connection struct {
	ID            string `json:"id"`
	SourceID      string `json:"sourceId"`
	DestinationID string `json:"destinationId"`
	Status        string `json:"status"`
}

type Fanout struct {
	ID             string   `json:"id"`
	SourceID       string   `json:"sourceId"`
	DestinationIDs []string `json:"destinationIds"`
	Status         string   `json:"status"`
}

type Health struct {
	Status string `json:"status"`
}

type Job struct {
	ID         string    `json:"id"`
	Status     string    `json:"status"`
	Progress   int       `json:"progress"`
	StartedAt  time.Time `json:"startedAt"`
	FinishedAt time.Time `json:"finishedAt,omitempty"`
	Error      string    `json:"error,omitempty"`
	SourceID   string    `json:"sourceId,omitempty"`
	DestIDs    []string  `json:"destinationIds,omitempty"`
}
