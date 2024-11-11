package binding

import (
	rttmas_db "rttmas-backend/pkg/database"
)

var RTTMAS_UV_BINDING_GEO_SEARCH_RADIUS = 45
var RTTMAS_UV_BINDING_PATH_SIMILARITY_WINDOW_IN_SECONDS = 30
var RTTMAS_UV_BINDING_CONVERGENCE_THRESHOLD = 50

var RTTMAS_PV_BINDING_GEO_SEARCH_RADIUS = 45
var RTTMAS_PV_BINDING_PATH_SIMILARITY_WINDOW_IN_SECONDS = 30
var RTTMAS_PV_BINDING_CONVERGENCE_THRESHOLD = 50

func RTTMAS_InitializeBindingModule() {
	rttmas_db.RedisExecuteLuaScript("create_indices", []string{"nil"})
}
