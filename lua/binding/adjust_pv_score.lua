-- Inputs
local timestamp = ARGV[1]
local plateNumber = ARGV[2]
local reporterUID = ARGV[3]
local reporterLongitude = ARGV[4]
local reporterLatitude = ARGV[5]

local geoSearchRadius = ARGV[6] -- in meters
local pathSimilarityBacktraceWindow = ARGV[7] -- in seconds
local convergenceThreshold = tonumber(ARGV[8]) -- in seconds


-- Find or generate PID for plate number
-- local pid = 1


-- Search similar VIDs
-- local similarVIDs = redis.call(
--     "FT.SEARCH", "pv_index",
--     "*=>[KNN 5 @vector $query_vec AS dist]",
--     "PARAMS", "2", "query_vec", queryVector,
--     "SORTBY", "dist",
--     "RETURN", "2", "vid", "dist",
--     "DIALECT", "2"
-- )


-- Check convergence. Skip if converged
local convergenceVID = redis.call("HGET", string.format("pv_convergence:%s", plateNumber), "vid")
local convergenceCount = tonumber(redis.call("HGET", string.format("pv_convergence:%s", plateNumber), "count"))
if convergenceCount ~= nil and convergenceCount >= convergenceThreshold then
    return 1
end


-- Vehicle Density Weighting
local surroundingVehicleCount = 0
local surroundingVehicles = redis.call(
    "GEORADIUS",
    string.format("v_locations:%s", timestamp),
    reporterLongitude, reporterLatitude,
    geoSearchRadius, "m", "COUNT", 3, "ASC"
)

for Index, possibleVID in pairs(surroundingVehicles) do
    surroundingVehicleCount = surroundingVehicleCount + 1
end

if (surroundingVehicleCount == 0) then
    return -1
end


for Index, possibleVID in pairs(surroundingVehicles) do

    -- Define score
    local score = 5 / surroundingVehicleCount


    -- Double-check:  Path Similarity Search
    local pathSimilarityFrequencyTable = {}
    for t = timestamp, timestamp-pathSimilarityBacktraceWindow, -1 do
        local rkeyPLocationsForTimeT = string.format("p_locations:%s", t)
        local rkeyVLocationsForTimeT = string.format("v_locations:%s", t)

        local queryCoords = redis.call(
            "GEOPOS", rkeyPLocationsForTimeT, plateNumber
        )
        -- Check if the member was found
        if #queryCoords > 0 and queryCoords[1] ~= false then
            -- Extract longitude and latitude from the GEOPOS result
            local longitude = queryCoords[1][1]
            local latitude = queryCoords[1][2]
            
            local nearbyVehiclesAtTimeT = redis.call(
                "GEORADIUS", rkeyVLocationsForTimeT,
                longitude, latitude,
                geoSearchRadius, "m", "COUNT", 3, "ASC"
            )

            for _, memberID in pairs(nearbyVehiclesAtTimeT) do
                -- local memberID = nearbyVehiclesAtTimeT[i]
                
                -- If the memberID is already in the table, increment its count, else initialize it to 1
                if pathSimilarityFrequencyTable[memberID] then
                    pathSimilarityFrequencyTable[memberID] = pathSimilarityFrequencyTable[memberID] + 1
                else
                    pathSimilarityFrequencyTable[memberID] = 1
                end
            end
        end
    end

    -- Step to find the member with the highest frequency
    local maxMember = nil
    local maxFrequency = 0

    for memberID, count in pairs(pathSimilarityFrequencyTable) do
        if count > maxFrequency then
            maxFrequency = count
            maxMember = memberID
        end
    end

    if maxMember == possibleVID then
        score = score * 2
    end


    -- Adjust Score
    local bindingKeyForMostProbablePID = string.format("pv_bindings:%s", possibleVID)

    redis.call("ZINCRBY", bindingKeyForMostProbablePID, score, plateNumber)
end

-- Get reporter's VID
local reporterVIDs = redis.call(
    "FT.SEARCH", "uv_binding_index",
    string.format("@uid:%s", reporterUID), "NOCONTENT"
)
for Index, Value in pairs(reporterVIDs) do
    if Index > 1 then
        local VID = tostring(Value):gsub("uv_bindings:", "")
        local bindingKeyForReporter = string.format("pv_bindings:%s", VID)
        redis.call("ZINCRBY", bindingKeyForReporter, -100, plateNumber)
    end
end

local mostProbableVID = surroundingVehicles[1]

-- Update convergence counter
if convergenceVID == mostProbableVID then
    redis.call("HSET", string.format("pv_convergence:%s", plateNumber), "count", convergenceCount+1)
else
    redis.call("HSET", string.format("pv_convergence:%s", plateNumber), "vid", mostProbableVID)
    redis.call("HSET", string.format("pv_convergence:%s", plateNumber), "count", 1)
end


return 0