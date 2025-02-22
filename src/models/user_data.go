package models

type UserReport struct {
	Latitude  float64 `json:"latitude" bson:"latitude"`
	Longitude float64 `json:"longitude" bson:"longitude"`
	Speed     float64 `json:"speed" bson:"speed"`
	Heading   float64 `json:"heading" bson:"heading"`
}

type UserData struct {
	ID      string               `json:"id"`
	UID     string               `json:"uid" bson:"uid"`
	Reports map[int64]UserReport `json:"reports" bson:"reports"`
}
