package binding

import (
	"fmt"
	rttmas_db "rttmas-backend/pkg/database"
)

func RTTMAS_InsertPlateReport(
	reportTime int64,
	latitude float64, longitude float64,
	reportedPlate string, reporterUID string,
) {
	rkey := fmt.Sprintf("p_locations:%d", reportTime)
	rttmas_db.RedisGeoAdd(rkey, latitude, longitude, reportedPlate)
}

func RTTMAS_AdjustPVScore(
	reportTime int64,
	latitude float64, longitude float64,
	reportedPlate string, reporterUID string,
) {
	rttmas_db.RedisExecuteLuaScript(
		"adjust_pv_score", []string{"nil"},
		reportTime, reportedPlate, reporterUID,
		longitude, latitude,
		RTTMAS_PV_BINDING_GEO_SEARCH_RADIUS,
		RTTMAS_PV_BINDING_PATH_SIMILARITY_WINDOW_IN_SECONDS,
		RTTMAS_PV_BINDING_CONVERGENCE_THRESHOLD,
	)
}

func RTTMAS_OnPlateReport(
	reportTime int64,
	latitude float64, longitude float64,
	reportedPlate string, reporterUID string,
) {
	RTTMAS_InsertPlateReport(reportTime, latitude, longitude, reportedPlate, reporterUID)
	RTTMAS_AdjustPVScore(reportTime, latitude, longitude, reportedPlate, reporterUID)
}
