local VID = ARGV[1]

local rkey = string.format("pv_bindings:%s", VID)

local result = redis.call("ZCOUNT", rkey, "0", "+inf")

return result
