package models

// Parking provides data for taxi parking
type Parking struct {
	GlobalID       int64  `json:"global_id"`
	ID             int64  `json:"id"`
	SystemObjectID string `json:"system_object_id"`
	Name           string `json:"name"`
	Mode           string `json:"mode"`
}
