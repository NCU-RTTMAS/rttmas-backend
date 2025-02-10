-- UUID Generator (for new VID)
local random = math.random
local function uuid()
    local template ='v__xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'
    return string.gsub(template, '[xy]', function (c)
        local v = (c == 'x') and random(0, 0xf) or random(8, 0xb)
        return string.format('%x', v)
    end)
end


-- Input arguments
local plate = ARGV[1]
local convergenceThreshold = 10
local vidCreationThreshold = 3

-- Obtain or set not bind count
local rkey = string.format("pv_not_bind:%s", plate)
local notBindCount = tonumber(redis.call("HGET", rkey, "count"))
if notBindCount == nil then
    notBindCount = 0
    redis.call("HSET", rkey, "count", notBindCount)
end

-- Check if the UID is already bound to a VID
-- Increment not bind count, if false
local convergenceVID = redis.call("HGET", string.format("pv_convergence:%s", plate), "vid")
local convergenceCount = tonumber(redis.call("HGET", string.format("pv_convergence:%s", plate), "count"))
if convergenceCount ~= nil and convergenceCount >= convergenceThreshold then
    return 1
else
    notBindCount = notBindCount + 1
    redis.call("HSET", rkey, "count", notBindCount)
end


-- Generate and bind new VID if not bind count exceeds VID creation threshold
if notBindCount >= vidCreationThreshold then
    -- local globalVehicleCount = tonumber(redis.call("HGET", "counters", "vehicle_count")) + 1
    -- local vid = string.format("v__%d", globalVehicleCount)

    local vid = uuid()
    
    -- Bind PV
    local newBindingKey = string.format("pv_bindings:%s", vid)
    redis.call("ZINCRBY", newBindingKey, 1000, plate)

    -- Make PV converge
    local newConvergenceKey = string.format("pv_convergence:%s", plate)
    redis.call("HSET", newConvergenceKey, "vid", vid)
    redis.call("HSET", newConvergenceKey, "count", 100)

    local globalVehicleCount = tonumber(redis.call("HGET", "counters", "vehicle_count")) + 1
    redis.call("HSET", "counters", "vehicle_count", globalVehicleCount)
end

return 0