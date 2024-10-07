local uid = ARGV[1]

local convergenceVID = redis.call("HGET", string.format("uv_convergence:%s", uid), "vid")

if convergenceVID == nil then
    return "NULL"
end

local rkey = string.format("pv_bindings:%s", tostring(convergenceVID))

local result = redis.call(
    "ZREVRANGE", rkey, 0, 0, "WITHSCORES"
)

return result