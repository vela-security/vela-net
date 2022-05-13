# net

基础网络库方法

- [net.ipv4(addr)]()
- [net.ipv6(addr)]()
- [net.ip(addr)]()
- [net.ping(addr)]()
- [net.cat(addr)]()
- [net.open(addr)]()

```lua
    local v =  net.ipv4("127.0.0.1")
    print(v) --true

    local v = net.ipv6("aa::22")
    print(v) --true

    local v = net.ip("10.0.0.1")
    print(v.ipv4)  --true
    print(v.ipv6)  --false

    local p = net.ping("127.0.0.1")
    print(p.ok)   -- true
    print(p.addr) -- 127.0.0.1
    print(p.cnt)  -- 32 
    print(p.time) -- 56
    print(p.code) -- 200
    print(p.id)   -- 123091
    print(p.warp) -- error
```
