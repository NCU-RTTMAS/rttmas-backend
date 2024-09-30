-- geoadd.lua

local rkey = KEYS[1]
local longitude = ARGV[1]
local latitude = ARGV[2]
local member = ARGV[3]

-- Perform the GEOADD command
return redis.call('GEOADD', rkey, longitude, latitude, member)
