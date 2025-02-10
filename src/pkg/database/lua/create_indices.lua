redis.call(
    "FT.CREATE", "uv_binding_index", "ON", "HASH", "PREFIX", "1",
    "uv_bindings:v__", "SCHEMA", "uid", "TEXT"
)

redis.call(
    "FT.CREATE", "pv_binding_index", "ON", "HASH", "PREFIX", "1",
    "pv_bindings:v__", "SCHEMA", "pid", "TEXT"
)

redis.call(
    "FT.CREATE", "pid_plate_binding_index", "ON", "HASH", "PREFIX", "1",
    "pid-plate-bindings:p__", "SCHEMA", "plate", "TAG"
)

redis.call("HSET", "counters", "pid_plate_bindings", 0)
redis.call("HSET", "counters", "vehicle_count", 0)

return 0