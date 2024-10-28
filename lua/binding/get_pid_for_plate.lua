local plateNumber = ARGV[1]

local result = redis.call(
    "FT.SEARCH", "pid_plate_binding_index",
    string.format("@plate:{%s}", plateNumber:gsub("-", "\\-")), "NOCONTENT"
)

return result