local VID = ARGV[1]

local rkey = string.format("uv_bindings:%s", VID)

local result = redis.call(
    "ZREVRANGE", rkey, 0, 0, "WITHSCORES"
)

return result