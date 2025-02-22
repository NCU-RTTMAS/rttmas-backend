local allVIDs = redis.call(
    "SMEMBERS",
    "all_global_vids"
)

return allVIDs
