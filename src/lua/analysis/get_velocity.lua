-- Calculate the distance between two geospatial points
local function distance(member1, member2)
    return redis.call('GEODIST', KEYS[1], member1, member2, 'm') -- distance in meters
end

-- Main function to get velocity
local function get_velocity(key, terms)
    local total_path = 0
    local total_timespan = 0

    -- Fetch the geospatial locations from the ZSET based on the provided timespan
    local data = redis.call('ZRANGE', key, -terms, -1, 'WITHSCORES')
    -- Ensure the timespan does not exceed available locations
    if #data < terms * 2 then
        return nil -- Not enough data
    end
    local foo = 0 -- temp
    -- Iterate over the timespan and calculate the total path distance
    for t = 1, terms - 1 do
        local time1 = tonumber(data[(t - 1) * 2 + 1])   -- The location (odd indices in the data array)
        -- local member1 = tonumber(data[(t - 1) * 2 + 2]) -- The timestamp (even indices in the data array)
        local time2 = tonumber(data[t * 2 + 1])   -- The next location
        -- local member2 = tonumber(data[t * 2 + 2]) -- The next timestamp
        total_path = total_path + distance(time1, time2)
        total_timespan = total_timespan  + math.abs( time2 - time1)
    end

    -- Return the average velocity
    return  3.6 * total_path / total_timespan 
end

-- Call the function with the provided ZSET key and timespan
return get_velocity(KEYS[1], tonumber(ARGV[1]))
