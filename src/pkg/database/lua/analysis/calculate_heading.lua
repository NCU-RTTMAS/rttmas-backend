-- Function to calculate heading between two points
local function calculate_heading(lat1, lon1, lat2, lon2)
    -- Convert degrees to radians
    local rad = math.pi / 180
    lat1 = lat1 * rad
    lon1 = lon1 * rad
    lat2 = lat2 * rad
    lon2 = lon2 * rad

    -- Calculate differences
    local dLon = lon2 - lon1

    -- Calculate heading
    local y = math.sin(dLon) * math.cos(lat2)
    local x = math.cos(lat1) * math.sin(lat2) - math.sin(lat1) * math.cos(lat2) * math.cos(dLon)
    local heading = math.atan2(y, x)

    -- Convert from radians to degrees and normalize to 0-360 range
    heading = (heading * 180 / math.pi + 360) % 360

    return tostring(heading)
end

-- Retrieve geospatial positions from Redis
local key = KEYS[1]
local location1 = ARGV[1]
local location2 = ARGV[2]

-- Get latitude and longitude for both locations
local pos1 = redis.call('GEOPOS', key, location1)
local pos2 = redis.call('GEOPOS', key, location2)

-- Check if both locations exist
if not pos1[1] or not pos2[1] then
    return redis.error_reply("One or both of the locations do not exist.")
end

-- Extract longitude and latitude from the GEOPOS results
local lon1 = tonumber(pos1[1][1])
local lat1 = tonumber(pos1[1][2])
local lon2 = tonumber(pos2[1][1])
local lat2 = tonumber(pos2[1][2])

-- Calculate and return the heading
return calculate_heading(lat1, lon1, lat2, lon2)
