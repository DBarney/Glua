#GLUA - render html lua templtes in golang.

I got tired of trying to shoehorn a templating system (that works) into golangs `html/tempalte` library. So i built my own.

``` lua
local funciton Component(data)
	return div{"my name is", data.name}
end

return function(params)
	html{
		head{title="this is my page"},
		body{
			h1{"This is my small page I wrote"},
			p{
				"some content",
				"more content"
			},
			div{"page has some dynamic values",params.count},
			Component{name="daniel"}
		}
	}
end
```


Basically glua uses the ability of lua tables to contain both key values and indexed data to allow me to easily write html pages.
It supports components inplemented as functions, and can also have data passed from golang rendered into the document
