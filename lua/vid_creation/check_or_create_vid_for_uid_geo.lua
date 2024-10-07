-- UUID Generator (for new VID)
local random = math.random
local function uuid()
    local template ='xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[xy]', function (c)
        local v = (c == 'x') and random(0, 0xf) or random(8, 0xb)
        return string.format('%x', v)
    end)
end


-- Input arguments
local timestamp = ARGV[1]
local uid = ARGV[2]
local longitude = ARGV[3]
local latitude = ARGV[4]
local vLocationKey = ARGV[5]
local geoSearchRadius = ARGV[6] -- in meters

local convergenceThreshold = 50
local vidCreationThreshold = 3


-- Get surrounding vehicles
local surroundingVehicleCount = 0
local surroundingVehicles = redis.call(
    "GEORADIUS",
    string.format("v_locations:%s", timestamp),
    longitude, latitude,
    geoSearchRadius, "m", "COUNT", 3, "ASC"
)

for Index, possibleVID in pairs(surroundingVehicles) do
    surroundingVehicleCount = surroundingVehicleCount + 1
end


-- Obtain or set not bind count
local rkey = string.format("uv_not_bind:%s", uid)
local notBindCount = tonumber(redis.call("HGET", rkey, "count"))
if notBindCount == nil then
    notBindCount = 1
    redis.call("HSET", rkey, "count", notBindCount)
end

if (surroundingVehicleCount == 0) then
    notBindCount = notBindCount + 1
    redis.call("HSET", rkey, "count", notBindCount)
end

-- Generate and bind new VID if not bind count exceeds VID creation threshold
if notBindCount >= vidCreationThreshold then
    -- local globalVehicleCount = tonumber(redis.call("HGET", "counters", "vehicle_count")) + 1
    -- local vid = string.format("v__%d", globalVehicleCount)

    local vid = uuid()
    
    -- Bind UV
    local newBindingKey = string.format("uv_bindings:%s", vid)
    redis.call("ZINCRBY", newBindingKey, 1000, uid)

    -- Make UV converge
    local newConvergenceKey = string.format("uv_convergence:%s", uid)
    redis.call("HSET", newConvergenceKey, "vid", vid)
    redis.call("HSET", newConvergenceKey, "count", 100)

    -- redis.call('GEOADD', string.format("v_locations:%s", vid), longitude, latitude, vid)

    redis.call("del", rkey)
end

-- Check uv-binding
for Index, possibleVID in pairs(surroundingVehicles) do
    local rkey = string.format("uv_bindings:%s", VID)

    local result = redis.call(
        "ZREVRANGE", rkey, 0, 0
    )

    if result ~= nil and result[0] == uid then
        redis.call('GEOADD', string.format("v_locations:%s", possibleVID), longitude, latitude, possibleVID)
    end
end




return 0