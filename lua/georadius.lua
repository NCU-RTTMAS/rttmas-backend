local rkey = KEYS[1]
local longitude = ARGV[1]
local latitude = ARGV[2]
local radius = ARGV[3]
local unit = ARGV[4]
local result = redis.call('GEORADIUS', rkey, longitude, latitude, radius, unit)
return result