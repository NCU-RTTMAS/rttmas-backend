local UID = ARGV[1]
local convergenceThreshold = 0

local convergenceVID = redis.call("HGET", string.format("uv_convergence:%s", UID), "vid")
local convergenceCount = tonumber(redis.call("HGET", string.format("uv_convergence:%s", UID), "count"))

if convergenceCount ~= nil and convergenceCount >= convergenceThreshold then
    return convergenceVID
else
    return "NULL"
end
