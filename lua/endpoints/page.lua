return function(params)
  return html{
    head{
      title{params.title or "default title"}
    },
    body{
      params
    }
  }
end
