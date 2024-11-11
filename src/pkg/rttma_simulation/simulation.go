package rttma_simulation

// This package is used to parse the RTTMAS XML file and generate the simulation data after 0827.
import (
	"encoding/xml"
	"fmt" // Add this line
	"io"
	// "io/ioutil"
	"context"
	"rttmas-backend/pkg/database"

	"os"
)

// RTTMAS represents the root element of the XML.
type RTTMAS struct {
	XMLName          xml.Name         `xml:"rttmas-sim"`
	SimConfiguration SimConfiguration `xml:"sim-configuration"`
	PVBindingFacts   PVBindingFacts   `xml:"pv-binding-facts"`
	Simulation       Simulation       `xml:"simulation"`
}

// SimConfiguration represents the simulation configuration.
type SimConfiguration struct {
	StartTime             string `xml:"startTime,attr"`
	ActiveUserPercentage  string `xml:"activeUserPercentage,attr"`
	ReportIntervalSeconds string `xml:"reportIntervalInSeconds,attr"`
	SimulateGPSError      bool   `xml:"simulateGPSError,attr"`
}

// PVBindingFacts represents the vehicle binding facts.
type PVBindingFacts struct {
	Vehicles []Vehicle `xml:"veh"`
}

// Vehicle represents a vehicle with a vid and plate.
type Vehicle struct {
	Vid   string `xml:"vid,attr"`
	Plate string `xml:"plate,attr"`
}

// Simulation represents the simulation data.
type Simulation struct {
	Timesteps []Timestep `xml:"timestep"`
}

// Timestep represents a single timestep in the simulation.
type Timestep struct {
	TimeSeconds             int                      `xml:"timeSeconds,attr"`
	VehicleTrueLocations    []VehicleTrueLocation    `xml:"vehicle-true-locations>report"`
	UserLocationReports     []UserLocationReport     `xml:"user-location-reports>report"`
	PlateRecognitionReports []PlateRecognitionReport `xml:"plate-recognition-reports>report"`
}

// VehicleTrueLocation represents a true location report for a vehicle.
type VehicleTrueLocation struct {
	Timestep int     `xml:"timestep,attr"`
	VID      string  `xml:"VID,attr"`
	Lat      float64 `xml:"lat,attr"`
	Lon      float64 `xml:"lon,attr"`
}

// UserLocationReport represents a location report from a user.
type UserLocationReport struct {
	Timestep int     `xml:"timestep,attr"`
	UID      string  `xml:"UID,attr"`
	Lat      float64 `xml:"lat,attr"`
	Lon      float64 `xml:"lon,attr"`
}

// PlateRecognitionReport represents a plate recognition report.
type PlateRecognitionReport struct {
	Timestep        int     `xml:"timestep,attr"`
	ReporterUID     string  `xml:"reporterUID,attr"`
	Lat             float64 `xml:"lat,attr"`
	Lon             float64 `xml:"lon,attr"`
	PlateNumberSeen string  `xml:"plateNumberSeen,attr"`
}

// ParseRTTMAS parses the XML file into the RTTMAS struct.
func ParseRTTMAS(filePath string) (*RTTMAS, error) {
	xmlFile, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, err := io.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}

	var rttmas RTTMAS
	err = xml.Unmarshal(byteValue, &rttmas)
	if err != nil {
		return nil, err
	}

	return &rttmas, nil
}

func (r *RTTMAS) WriteToDB() error {
	// Assuming we have a database connection established

	// Store SimConfiguration
	// _, err := database.RTTMA_Collections.SimConfiguration.InsertOne(context.Background(), r.SimConfiguration)
	// if err != nil {
	// 	return fmt.Errorf("failed to insert SimConfiguration: %v", err)
	// }

	// Store TimeSteps
	for _, timestep := range r.Simulation.Timesteps {
		// Store VehicleTrueLocations
		for _, vtl := range timestep.VehicleTrueLocations {
			vtl.Timestep = timestep.TimeSeconds
			_, err := database.RTTMA_Collections.VehicleTrueLocations.InsertOne(context.Background(), vtl)
			if err != nil {
				return fmt.Errorf("failed to insert VehicleTrueLocation: %v", err)
			}
		}

		// Store UserLocationReports
		for _, ulr := range timestep.UserLocationReports {
			ulr.Timestep = timestep.TimeSeconds
			_, err := database.RTTMA_Collections.UserLocationReports.InsertOne(context.Background(), ulr)
			if err != nil {
				return fmt.Errorf("failed to insert UserLocationReport: %v", err)
			}
		}

		// Store PlateRecognitionReports
		for _, prr := range timestep.PlateRecognitionReports {
			prr.Timestep = timestep.TimeSeconds
			_, err := database.RTTMA_Collections.PlateRecognitionReports.InsertOne(context.Background(), prr)
			if err != nil {
				return fmt.Errorf("failed to insert PlateRecognitionReport: %v", err)
			}
		}
	}

	return nil
}
