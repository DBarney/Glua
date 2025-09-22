-- just the most simple example showing how nexting of tables and dynamic function names work
-- this renders out the html exactly how it is described
-- NOTE: I return the function as each module should only have a single node that it is responsible
-- for rendering. Other functions can be present in this file, but that is upto you to decide
-- when that is needed.
return function()
  return html{
    head{
      title{"My Site"}
    },
    body{
      h1{"Welcome to my basic site!"},
      img{src="/some.png"}
    }
  }
end
