local page = require("page")

return function()
  return page {
    title="This is my title",
    "some content goes here",
    "and here."
  }
end
