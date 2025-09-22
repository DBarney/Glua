-- this is how to reuse components to avoid writing the same html over and over again
local page = require("page")

return function()
  -- NOTE: page looks just like a standard html node, but is actually a custom function
  -- that accepts a table with parameters and arguments. I like this style, but if you
  -- prefer a standard function call: `page("title", {"children"})`, then do it that way
  -- its your code, write it how you prefer.
  return page {
    title="This is my title",
    "some content goes here",
    "and here."
  }
end
