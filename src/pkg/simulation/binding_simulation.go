package rttmas_simulation

import (
	"context"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	rttmas_db "rttmas-backend/pkg/database"
	"rttmas-backend/pkg/utils/logger"
)

var PV_BIND_SCORE_THRESHOLD = float64(5)

// Define the structures to map to the XML structure
type SimConfiguration struct {
	StartTime             string `xml:"startTime,attr"`
	ActiveUserPercentage  string `xml:"activeUserPercentage,attr"`
	ReportIntervalSeconds string `xml:"reportIntervalInSeconds,attr"`
	SimulateGPSError      string `xml:"simulateGPSError,attr"`
}

type VehicleBindingFact struct {
	VID   string `xml:"vid,attr"`
	Plate string `xml:"plate,attr"`
}

type PVBindingFacts struct {
	VehicleFacts []VehicleBindingFact `xml:"veh"`
}

type VehicleTrueLocation struct {
	VID string  `xml:"VID,attr"`
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
}

type UserLocationReport struct {
	UID         string  `xml:"UID,attr"`
	Lat         float64 `xml:"lat,attr"`
	Lon         float64 `xml:"lon,attr"`
	AttachedVID string  `xml:"attachedVID,attr"`
}

type PlateRecognitionReport struct {
	ReporterUID     string  `xml:"reporterUID,attr"`
	Lat             float64 `xml:"lat,attr"`
	Lon             float64 `xml:"lon,attr"`
	PlateNumberSeen string  `xml:"plateNumberSeen,attr"`
	AttachedVID     string  `xml:"attachedVID,attr"`
}

type Timestep struct {
	TimeSeconds             string                   `xml:"timeSeconds,attr"`
	VehicleTrueLocations    []VehicleTrueLocation    `xml:"vehicle-true-locations>report"`
	UserLocationReports     []UserLocationReport     `xml:"user-location-reports>report"`
	PlateRecognitionReports []PlateRecognitionReport `xml:"plate-recognition-reports>report"`
}

type Simulation struct {
	Timesteps []Timestep `xml:"timestep"`
}

type xRTTMAS struct {
	XMLName          xml.Name         `xml:"rttmas-sim"`
	SimConfiguration SimConfiguration `xml:"sim-configuration"`
	PVBindingFacts   PVBindingFacts   `xml:"pv-binding-facts"`
	Simulation       Simulation       `xml:"simulation"`
}

func parseXML(filename string) (*xRTTMAS, error) {
	// Read the XML file
	xmlFile, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer xmlFile.Close()

	// Read the file contents into a byte array
	xmlData, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Unmarshal the XML data into the Go structure
	var xrttmas xRTTMAS
	err = xml.Unmarshal(xmlData, &xrttmas)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling XML: %v", err)
	}

	return &xrttmas, nil
}

func playbackTimesteps(timesteps []Timestep, geoSearchRadius int) {

	tCounter := 0

	for _, timestep := range timesteps {

		// Iterate over user location reports
		for _, report := range timestep.UserLocationReports {
			rkey := fmt.Sprintf("u_locations:%s", timestep.TimeSeconds)
			rttmas_db.RedisExecuteLuaScript("geoadd", []string{rkey}, report.Lon, report.Lat, report.UID)
			rttmas_db.RedisExecuteLuaScript("check_or_create_vid_for_uid", []string{"nil"}, report.UID)

			rawResult, _ := rttmas_db.RedisExecuteLuaScript("get_uv_convergence", []string{"nil"}, report.UID)
			if rawResult != nil {
				vid := rawResult.(string)
				if vid != "NULL" {
					rkeyForVID := fmt.Sprintf("v_locations:%s", timestep.TimeSeconds)
					rttmas_db.RedisExecuteLuaScript("geoadd", []string{rkeyForVID}, report.Lon, report.Lat, vid)
				}
			}

			rttmas_db.RedisExecuteLuaScript("adjust_uv_score", []string{"nil"}, timestep.TimeSeconds, report.UID, report.Lon, report.Lat, geoSearchRadius, 30, 50)
		}

		// Iterate over plate recognition reports
		for _, report := range timestep.PlateRecognitionReports {
			rkey := fmt.Sprintf("p_locations:%s", timestep.TimeSeconds)
			rttmas_db.RedisExecuteLuaScript("geoadd", []string{rkey}, report.Lon, report.Lat, report.PlateNumberSeen)

			rttmas_db.RedisExecuteLuaScript("adjust_pv_score", []string{"nil"}, timestep.TimeSeconds, report.PlateNumberSeen, report.ReporterUID, report.Lon, report.Lat, geoSearchRadius, 30, 20)
		}

		tCounter += 1
	}
}

func AnalyzeUVBindingAccuracy(xrttmas *xRTTMAS) {

	totalUserCount := 500
	nullCount := 0
	correctCount := 0

	for i := 0; i < totalUserCount; i++ {
		VID := fmt.Sprintf("v__%d", i)

		rawResult, _ := rttmas_db.RedisExecuteLuaScript("get_most_probable_uid_for_vid", []string{"nil"}, VID)
		resultArr := rawResult.([]interface{})

		if len(resultArr) == 0 {
			nullCount++
			continue
		}

		predictedUID := resultArr[0].(string)

		factUID := fmt.Sprintf("u__%d", i)
		isMatch := predictedUID == factUID

		if isMatch {
			correctCount++
		}

		logger.Info(fmt.Sprintf("VID: %s  ==>  Predicted UID: %s  ==>  %t", VID, predictedUID, isMatch))
	}

	logger.Info(fmt.Sprintf("Correct: %d / %d", correctCount, totalUserCount-nullCount))
	logger.Info(fmt.Sprintf("NULL Count: %d", nullCount))
}

func AnalyzePVBindingAccuracy(xrttmas *xRTTMAS) {

	totalSeenVehicleCount := 0
	nullCount := 0
	correctCount := 0

	for _, fact := range xrttmas.PVBindingFacts.VehicleFacts {
		rawResult, _ := rttmas_db.RedisExecuteLuaScript("get_most_probable_plate_for_vid", []string{"nil"}, fact.VID)
		resultArr := rawResult.([]interface{})

		if len(resultArr) == 0 {
			continue
		}

		predictedPlate := resultArr[0].(string)
		predictedScore, _ := strconv.ParseFloat(resultArr[1].(string), 64)

		if predictedScore < PV_BIND_SCORE_THRESHOLD {
			predictedPlate = "NULL"
		}

		isMatch := predictedPlate == fact.Plate

		if isMatch {
			correctCount++
		}

		// Check if this vehicle has been seen
		rawResult, _ = rttmas_db.RedisExecuteLuaScript("get_vehicle_seen_count", []string{"nil"}, fact.VID)
		seenCount := rawResult.(int64)
		if seenCount > 0 {
			if !isMatch {
				if predictedPlate == "NULL" {
					nullCount++
				}
				logger.Info(fmt.Sprintf("%s:  %s == %s  ->  %t", fact.VID, predictedPlate, fact.Plate, isMatch))
			}
			totalSeenVehicleCount++
		}
	}

	logger.Info(fmt.Sprintf("Correct: %d / %d", correctCount, totalSeenVehicleCount))
	logger.Info(fmt.Sprintf("NULL Count: %d", nullCount))
}

func AnalyzeVIDCreation(xrttmas *xRTTMAS) {

	totalSeenVehicleCount := 0
	correctCount := 0

	for _, fact := range xrttmas.PVBindingFacts.VehicleFacts {
		totalSeenVehicleCount++
		UID := strings.ReplaceAll(fact.VID, "v", "u")

		rawResult, _ := rttmas_db.RedisExecuteLuaScript("analyze_plate_for_uid", []string{"nil"}, UID)
		if rawResult == nil {
			continue
		}
		resultArr := rawResult.([]interface{})

		if len(resultArr) == 0 {
			continue
		}

		predictedPlate := resultArr[0].(string)

		isMatch := predictedPlate == fact.Plate

		if isMatch {
			correctCount++
		}

		logger.Info(fmt.Sprintf("UID: %s  :  Predicted: %s  ==>  Real: %s", UID, predictedPlate, fact.Plate))
	}

	logger.Info(fmt.Sprintf("Correct: %d / %d", correctCount, totalSeenVehicleCount))
}

func AnalyzeBindingsForDynamicVIDCreation(xrttmas *xRTTMAS) {
	allVIDs, _ := rttmas_db.RedisExecuteLuaScript("analyze_all_vid_bindings", []string{"nil"})
	if allVIDs == nil {
		logger.Debug("None")
	}

	predictions := make(map[string]string)

	// Extract predictions from DB
	for _, _vid := range allVIDs.([]interface{}) {
		VID := _vid.(string)
		logger.Info(VID)

		// Plate Prediction
		rawResult1, _ := rttmas_db.RedisExecuteLuaScript("get_most_probable_plate_for_vid", []string{"nil"}, VID)
		resultArr1 := rawResult1.([]interface{})
		if len(resultArr1) == 0 {
			continue
		}
		predictedPlate := resultArr1[0].(string)
		predictedScore, _ := strconv.ParseFloat(resultArr1[1].(string), 64)
		if predictedScore < PV_BIND_SCORE_THRESHOLD {
			predictedPlate = "NULL"
		}

		// UID Prediction
		rawResult2, _ := rttmas_db.RedisExecuteLuaScript("get_most_probable_uid_for_vid", []string{"nil"}, VID)
		resultArr2 := rawResult2.([]interface{})
		if len(resultArr2) == 0 {
			continue
		}
		predictedUID := resultArr2[0].(string)

		predictions[predictedUID] = predictedPlate
	}

	totalCount := len(predictions)
	correctCount := 0

	// Compare predictions with facts
	for _, fact := range xrttmas.PVBindingFacts.VehicleFacts {
		UID := strings.ReplaceAll(fact.VID, "v", "u")

		predictedPlate := predictions[UID]

		isMatch := predictedPlate == fact.Plate

		if isMatch {
			correctCount++
		}

		logger.Info(fmt.Sprintf("%s  ->  %s == %s  ->  %t", UID, predictedPlate, fact.Plate, isMatch))
	}

	logger.Info(fmt.Sprintf("Correct: %d / %d", correctCount, totalCount))
	logger.Info(fmt.Sprintf("Accuracy: %.02f%%", float64(correctCount)/float64(totalCount)*100))
}

func AnalysisExperiment() {
	// Replace with the actual path to your XML file
	// filename := "pkg/simulation/sumo-scenarios/output_20240915_1915_taipei_nogpserror_forwardonly.xml"
	// filename := "pkg/simulation/sumo-scenarios/output_20240915_1935_nogpserror_forwardonly.xml"
	filename := "pkg/simulation/sumo-scenarios/output_20240915_1830_taipei_nogpserror.xml"
	// filename := "pkg/simulation/sumo-scenarios/output_20240915_1830_taipei_withgpserror.xml"

	// filename := "pkg/simulation/sumo-scenarios/output_20240915_1751_withgpserror.xml"
	// filename := "pkg/simulation/sumo-scenarios/output_20240915_1751_nogpserror.xml"

	// filename := "pkg/simulation/sumo-scenarios/output_20240905_1615_nogpserror.xml"
	// filename := "pkg/simulation/sumo-scenarios/output_20240911_1210_withgpserror.xml"

	xrttmas, err := parseXML(filename)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	analysisOnly := false

	if !analysisOnly {
		radius := 45
		logger.Info(fmt.Sprintf("Radius: %d", radius))

		rttmas_db.GetRedis().FlushAll(context.TODO())

		rttmas_db.RedisExecuteLuaScript("create_indices", []string{"nil"})

		playbackTimesteps(xrttmas.Simulation.Timesteps, radius)
	}

	AnalyzeBindingsForDynamicVIDCreation(xrttmas)
}
