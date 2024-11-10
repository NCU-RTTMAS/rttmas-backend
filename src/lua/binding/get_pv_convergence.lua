local plate = ARGV[1]
local convergenceThreshold = 50

local convergenceVID = redis.call("HGET", string.format("pv_convergence:%s", plate), "vid")
local convergenceCount = tonumber(redis.call("HGET", string.format("pv_convergence:%s", plate), "count"))

if convergenceCount ~= nil and convergenceCount >= convergenceThreshold then
    return convergenceVID
else
    return "NULL"
end
