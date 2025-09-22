-- This is to show one of the gotchas around null values in tables.
-- if a table contains a null, then that is the end of the array like properties.
-- {"a",null,"b"} when used with ipairs will only be able to access "a" as the second
-- element is null.


-- This shows how to return something, even a blank table, to preserve the table order
local function elem(show)
  if show then
    return div{"table was returned"}
  else
    return {}
  end
end

local function nullElem(show)
  if show then
    return div{"optional table was returned"}
  end
end

return function()
  return div{
    elem(true),
    elem(false),
    elem(true),
    elem(false),
    -- other option is to always put the optional content with in its own table
    -- this preserves the table iteration order
    {null},
    {nullElem(false)},
    {nullElem(true)}
  }
end
