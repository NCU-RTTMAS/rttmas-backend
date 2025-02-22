package models

type PlateReport struct {
	Latitude    float64 `json:"latitude" bson:"latitude"`
	Longitude   float64 `json:"longitude" bson:"longitude"`
	ReporterUID string  `json:"reporter_uid" bson:"reporter_uid"`
}

type PlateData struct {
	ID          string                  `json:"id"`
	PlateNumber string                  `json:"plate_number" bson:"plate_number"`
	Reports     map[int64][]PlateReport `json:"reports" bson:"reports"`
}
