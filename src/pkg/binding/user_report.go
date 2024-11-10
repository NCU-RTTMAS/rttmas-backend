package binding

import (
	"fmt"
	rttmas_db "rttmas-backend/pkg/database"
)

func RTTMAS_InsertUserLocationReport(
	reportTime int64,
	latitude float64, longitude float64,
	reporterUID string,
) {
	rkey := fmt.Sprintf("u_locations:%d", reportTime)
	rttmas_db.RedisGeoAdd(rkey, latitude, longitude, reporterUID)
}

func RTTMAS_AdjustUVScore(
	reportTime int64,
	latitude float64, longitude float64,
	reporterUID string,
) {
	rttmas_db.RedisExecuteLuaScript(
		"binding/adjust_uv_score", "nil",
		reportTime, reporterUID,
		longitude, latitude,
		RTTMAS_UV_BINDING_GEO_SEARCH_RADIUS,
		RTTMAS_UV_BINDING_PATH_SIMILARITY_WINDOW_IN_SECONDS,
		RTTMAS_UV_BINDING_CONVERGENCE_THRESHOLD,
	)
}

func RTTMAS_OnUserLocationReport(
	reportTime int64,
	latitude float64, longitude float64,
	reporterUID string,
) {
	RTTMAS_InsertUserLocationReport(reportTime, latitude, longitude, reporterUID)
	RTTMAS_AdjustUVScore(reportTime, latitude, longitude, reporterUID)
}
