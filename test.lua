local mt = {}

function vector(x, y)
    local v = {x = x, y = y}
    setmetatable(v, mt)
    return v
end

mt.__add = function(v1, v2)
    return vector(v1.x+v2.x, v1.y+v2.y)
end

mt.__sub = function(v1, v2)
    return vector(v1.x-v2.x, v1.y-v2.y)
end

mt.__mul = function(v1, n)
    return vector(v1.x*n, v1.y*n)
end

mt.__div = function(v1, n)
    return vector(v1.x/n, v2.y/n)
end

mt.__len = function(v)
    return (v.x*v.x+v.y*v.y)^0.5
end

mt.__eq = function(v1, v2)
    return v1.x == v2.x and v1.y == v2.y
end

mt.__index = function(v, k)
    if k == 'print' then
        return function()
            print("[" .. v.x .. "," .. v.y .. "]")
        end
    end
end

mt.__call = function(v)
    print("[" .. v.x .. "," .. v.y .. "]")
end

local v1 = vector(1, 2)
v1:print()

local v2 = vector(3, 4)
v2:print()

local v3 = v1*2
v3:print()

local v4 = v1+v2
v4:print()

print(#v2)
print(v1 == v2)
print(v2 == vector(3, 4))
v4()

print("------------------")
local t = {a = 1, b = 2, c = 3}
for k, v in pairs(t) do
    print(k, v)
end
print("------------------")
local t2 = {4, 3, 2, 1}
for i, v in ipairs(t2) do
    print(i, v)
end
print("------------------")

local function div0(a, b)
    if b == 0 then
        error("DIV BY ZERO")
    else
        return a/b
    end
end

local function div1(a, b) return div0(a, b) end
local function div2(a, b) return div1(a, b) end

local ok, result = pcall(div2, 4, 2)
print(ok, result, type(result))
ok, result = pcall(div2, 5, 0)
print(ok, result, type(result))
ok, result = pcall(div2, {}, {})
print(ok, result, type(result))

tb = {1, 2, 3, 4, 5, 6, 7}
for i, v in ipairs(tb) do
    if v%2 == 0 then
        print(i, v)
    end
end

print(type(_G))
print(os.time())
print(package)
for k, v in pairs(package) do
    print(k, v)
    if type(v) == 'table' then
        for kk, vv in pairs(v) do
            print('\t', kk, vv)
        end
    end
end

local function foo (a)
    print("foo", a)
    return coroutine.yield(2*a)
end
  
local co = coroutine.create(function (a,b)
    print("co-body", a, b)
    local r = foo(a+1)
    print("co-body", r)
    local r, s = coroutine.yield(a+b, a-b)
    print("co-body", r, s)
    return b, "end"
end)

print("------------------")
print("main", coroutine.resume(co, 1, 10))
print("main", coroutine.resume(co, "r"))
print("main", coroutine.resume(co, "x", "y"))
print("main", coroutine.resume(co, "x", "y"))

--[[
co-body 1 10
foo 2
main true 4
co-body r
main true 11 -9
co-body x y
main true 10 end
main false xxx
]]

