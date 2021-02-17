
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
