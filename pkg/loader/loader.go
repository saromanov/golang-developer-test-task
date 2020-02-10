package loader

import "github.com/saromanov/golang-developer-test-task/pkg/models"

// Loader defines struct for loading data
type Loader interface {
	Load() ([]models.Parking, error)
}
