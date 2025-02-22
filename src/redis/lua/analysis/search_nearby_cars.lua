local function log(message)
    local log_key = "script_log"  -- Key for the log
    local timestamp = redis.call("TIME")  -- Get the current timestamp
    local log_message = string.format("[%d:%d] %s", timestamp[1], timestamp[2], message)
    
    -- Append the log message to a Redis list
    redis.call("RPUSH", log_key, log_message)
end

local geoKey = KEYS[1]
local radius = 5 -- meters

-- Table to hold results
local results = {}
local userData
local tempKey
local location
-- Iterate through each user in the geoKey
local users = redis.call('KEYS', geoKey) -- Assuming geoKey is a sorted set with user names

for _, user in ipairs(users) do
    -- Get the location of the user from the user_location_report
    local userKey =  user:gsub("user_location_report:", "basic_info:")
     userData = redis.call("JSON.GET", userKey  , "$.LatestTimestep")
    
    if userData then
        -- Get the geospatial member from the user's data
        tempKey = userKey:gsub("basic_info:","user_location_report:" )
        userData = userData:sub(2,-2)
        -- log(userData)
        -- log(userData:sub(2,-2))
        -- log(tostring(foo))
        -- local location = redis.call('HGET', geoKey, user) -- Assuming user's location is stored as a member in geoKey
        location = redis.call('ZSCORE', tempKey,userData) -- Assuming user's location is stored as a member in geoKey
        log(string.format("%s %s", tempKey, userData))

        if location then
            -- Perform GEOSEARCH using the member's location
            local nearbyUsers = redis.call('GEOSEARCH', geoKey, 'FROMMEMBER', user, 'BYRADIUS', radius, 'm')
            for _, nearbyUser in ipairs(nearbyUsers) do
                table.insert(results, nearbyUser)
            end
        end
    end
end

-- Return the list of nearby users found
return results
