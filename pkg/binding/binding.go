package binding

import (
	"fmt"

	rttmas_db "rttmas-backend/pkg/database"
)

func RTTMAS_InitializeBindingModule() {
	rttmas_db.RedisExecuteLuaScript("create_indices", "nil")
}

func RTTMAS_InsertUserLocationReport(
	reportTime int64,
	latitude float64, longitude float64,
	reporterUID string,
) {
	rkey := fmt.Sprintf("u_locations:%d", reportTime)
	rttmas_db.RedisGeoAdd(rkey, latitude, longitude, reporterUID)
}
