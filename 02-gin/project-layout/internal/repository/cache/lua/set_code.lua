-- 发送到的 key，也就是 code:业务:手机号码
-- phone_code:login:152xxxxxxxx
local key = KEYS[1]
-- 验证次数，我们一个验证码，最多重复三次，这个记录还可以验证几次
-- phone_code:login:152xxxxxxxx:cnt
local cntKey = key..":cnt"
-- 你的验证码 123456
local val= ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))

-- -1 是 key 存在，但是没有过期时间
if ttl == -1 then
    -- key 存在，但是没有过期时间
    -- 系统错误，你的同事手贱，手动设置了这个 key，但是没给过期时间
    return -2
-- -2 是 key 不存在，ttl < 540 是发了一个验证码，已经超过一分钟了，可以重新发送
elseif ttl == -2 or ttl < 540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("set", cntKey, 3)
    redis.call("expire", cntKey, 600)
    -- 完美，符合预期
    return 0
else
    -- 发送太频繁, 已经发送了一个验证码，但是还不到一分钟
    return -1
end

