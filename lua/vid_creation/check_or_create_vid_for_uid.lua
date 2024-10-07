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
local uid = ARGV[1]
local convergenceThreshold = 50
local vidCreationThreshold = 5

-- Obtain or set not bind count
local rkey = string.format("uv_not_bind:%s", uid)
local notBindCount = tonumber(redis.call("HGET", rkey, "count"))
if notBindCount == nil then
    notBindCount = 1
    redis.call("HSET", rkey, "count", notBindCount)
end

-- Check if the UID is already bound to a VID
-- Increment not bind count, if false
local convergenceVID = redis.call("HGET", string.format("uv_convergence:%s", uid), "vid")
local convergenceCount = tonumber(redis.call("HGET", string.format("uv_convergence:%s", uid), "count"))
if convergenceCount ~= nil and convergenceCount >= convergenceThreshold then
    return 1
else
    notBindCount = notBindCount + 1
    redis.call("HSET", rkey, "count", notBindCount)
end

-- Generate and bind new VID if bind count exceeds VID creation threshold
if notBindCount >= vidCreationThreshold then
    -- local vid = string.format("v__%d", globalVehicleCount)

    local vid = uuid()
    
    -- Bind UV
    local newBindingKey = string.format("uv_bindings:%s", vid)
    redis.call("ZINCRBY", newBindingKey, 100, uid)

    -- Make UV converge
    local newConvergenceKey = string.format("uv_convergence:%s", uid)
    redis.call("HSET", newConvergenceKey, "vid", vid)
    redis.call("HSET", newConvergenceKey, "count", 100)

    redis.call("del", string.format("uv_not_bind:%s", uid))

    local globalVehicleCount = tonumber(redis.call("HGET", "counters", "vehicle_count")) + 1
    redis.call("HSET", "counters", "vehicle_count", globalVehicleCount)
end

return 0