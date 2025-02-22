package binding

import (
	"fmt"
	rttmas_models "rttmas-backend/models"
	rttmas_redis "rttmas-backend/redis"
	rttmas_service "rttmas-backend/services"
	"rttmas-backend/utils/logger"
)

func RTTMAS_InsertUserLocationReport(
	reportTime int64,
	latitude float64, longitude float64,
	reporterUID string,
) {
	rkey := fmt.Sprintf("u_locations:%d", reportTime)
	rttmas_redis.RedisGeoAdd(rkey, latitude, longitude, reporterUID)
}

func RTTMAS_AdjustUVScore(
	reportTime int64,
	latitude float64, longitude float64,
	reporterUID string,
) {
	_, err := rttmas_redis.RedisExecuteLuaScript(
		"adjust_uv_score", []string{"nil"},
		reportTime, reporterUID,
		longitude, latitude,
		RTTMAS_UV_BINDING_GEO_SEARCH_RADIUS,
		RTTMAS_UV_BINDING_PATH_SIMILARITY_WINDOW_IN_SECONDS,
		RTTMAS_UV_BINDING_CONVERGENCE_THRESHOLD,
	)
	if err != nil {
		logger.Error("AdjustUVScore ERR:", err)
	}

}

func RTTMAS_UpdateUVBinding(
	reportTime int64,
	latitude float64, longitude float64,
	reporterUID string,
) {
	rttmas_redis.RedisExecuteLuaScript("check_or_create_vid_for_uid", []string{"nil"}, reporterUID)

	rawResult, _ := rttmas_redis.RedisExecuteLuaScript("get_uv_convergence", []string{"nil"}, reporterUID)
	if rawResult != nil {
		vid := rawResult.(string)
		if vid != "NULL" {
			rkeyForVID := fmt.Sprintf("v_locations:%d", reportTime)
			rttmas_redis.RedisExecuteLuaScript("geoadd", []string{rkeyForVID}, longitude, latitude, vid)
		}
	}
}

func RTTMAS_StoreUserReport(
	reportTime int64,
	latitude float64, longitude float64,
	reporterUID string,
) {
	mongoReportRecord := rttmas_models.UserReport{
		Latitude:  latitude,
		Longitude: longitude,
		Speed:     0,
		Heading:   0,
	}

	rttmas_service.StoreUserReportToMongoDB(reporterUID, reportTime, mongoReportRecord)
}

func RTTMAS_OnUserLocationReport(
	reportTime int64,
	latitude float64, longitude float64,
	reporterUID string,
) {
	RTTMAS_InsertUserLocationReport(reportTime, latitude, longitude, reporterUID)
	RTTMAS_UpdateUVBinding(reportTime, latitude, longitude, reporterUID)
	RTTMAS_AdjustUVScore(reportTime, latitude, longitude, reporterUID)
	RTTMAS_StoreUserReport(reportTime, latitude, longitude, reporterUID)
}

// func RTTMAS_InsertFoo (){
// 	rkey := fmt.Sprintf("u_locations:%s", timestep.TimeSeconds)
// 			rttmas_redis.RedisExecuteLuaScript("geoadd", []string{rkey}, report.Lon, report.Lat, report.UID)
// 			rttmas_redis.RedisExecuteLuaScript("check_or_create_vid_for_uid", []string{"nil"}, report.UID)

// 			rawResult, _ := rttmas_redis.RedisExecuteLuaScript("get_uv_convergence", []string{"nil"}, report.UID)
// 			if rawResult != nil {
// 				vid := rawResult.(string)
// 				if vid != "NULL" {
// 					rkeyForVID := fmt.Sprintf("v_locations:%s", timestep.TimeSeconds)
// 					rttmas_redis.RedisExecuteLuaScript("geoadd", []string{rkeyForVID}, report.Lon, report.Lat, vid)
// 				}
// 			}

// 			rttmas_redis.RedisExecuteLuaScript("adjust_uv_score", []string{"nil"}, timestep.TimeSeconds, report.UID, report.Lon, report.Lat, geoSearchRadius, 30, 50)

// 			var mongoReportRecord rttmas_models.UserReport
// 			mongoReportRecord.Latitude = report.Lat
// 			mongoReportRecord.Longitude = report.Lon
// 			mongoReportRecord.Speed = 0
// 			mongoReportRecord.Heading = 0

// 			reportTime, _ := strconv.ParseInt(timestep.TimeSeconds, 10, 64)

// 			rttmas_service.StoreUserReportToMongoDB(report.UID, reportTime, mongoReportRecord)
// }
