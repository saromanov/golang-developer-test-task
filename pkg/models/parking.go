package models

// Parking provides data for taxi parking
type Parking struct {
	GlobalID         int64  `json:"global_id"`
	ID               int64  `json:"id"`
	SystemObjectID   string `json:"system_object_id"`
	Name             string `json:"name"`
	Mode             string `json:"mode"`
	AdmArea          string `json:"AdmArea"`
	District         string `json:"District"`
	Address          string `json:"Address"`
	LongitudeWGS84   string `json:"Longitude_WGS84"`
	LatitudeWGS84    string `json:"Latitude_WGS84"`
	CarCapacity      int64  `json:"CarCapacity"`
	IDEN             int64  `json:"ID_en"`
	NameEn           string `json:"Name_en"`
	AdmAreaEn        string `json:"AdmArea_en"`
	DistrictEn       string `json:"District_en"`
	AddressEn        string `json:"Address_en"`
	LongitudeWGS84En string `json:"Longitude_WGS84_en"`
	LatitudeWGS84En  string `json:"Latitude_WGS84_en"`
	CarCapacityEn    int64  `json:"CarCapacity_en"`
	ModeEn           string `json:"Mode_en"`
}
