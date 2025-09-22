-- page is a wrapper around other content. This is one way to easily allow children to be passed in
return function(params)
  
  -- we clear out parameters that are used specifically for this node, so that they do not
  -- get rendered into the body node
  local t = params.title or "default title"
  params.title = null

  return html{
    head{
      title{t}
    },
    -- params is directly passed into the body allowing children to be rendered
    -- could also do this like `body params` or `body{params} if no other children need to be
    -- specified before or after the children in params
    body{
      h1{"some thing"},
      div{params},
      footer{"example footer"}
    }
  }
end
