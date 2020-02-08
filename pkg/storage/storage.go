package storage

import "github.com/saromanov/golang-developer-test-task/pkg/models"

// Storage defines main interface for storage
type Storage interface {
	Insert([]models.Parking) error
	Find(*FindConfig) ([]models.Parking, error)
}

// FindConfig provides definition of the findign data
type FindConfig struct {
	GlobalID int64
	ID       int64
	ModeID   string
}
