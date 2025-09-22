# GLUA - render html lua templates in golang.

Problem:

I was using `html/template` to try and do component based html templates. It worked ok, and was native, but you can't pass the name of a second template to be rendered inside of the first template. Meaning that you can't have a generic container template that renders custom content. No Wrappers. No Containers. No easy way to do components that contain other components.


This does not work inside of a template:

```
{{template .name .values}}
```

https://go.dev/play/p/FawHmxpH4SB

## Why not use an existing pure go library?

Oh Please no.

```golang

# probably some other libraries will look like this
func BuildPage(id string) *html.Node {
    return html.Html(
        html.Body(
            html.Div(
                html.P("some text").ID(id)))
    }
```

Golang is statically typed. So everything needs a type added to it. I really don't want to have to update that for every argument I pass around these components.

Dont get me wrong, I like golang. This whole project is meant to be embedded in a go program. But I want it to be easier, not harder, then just writing the html directly.

## What does lua look like instead?

Since I have used lua before, and I know it is really good at building DSLs and also at being embedded inside of other programs, I decided it might be a good fit for this small project.

This is for prototyping html sites, so I don't really care too much about performance. That being said, this is probably pretty fast. It also probably could be much faster.

Here is what a basic lua html template might look like:

``` lua
local funciton Component(data)
	return div{"my name is", data.name}
end

return function(params)
	return html{
		head{title="this is my page"},
		body{
			h1{"This is my small page I wrote"},
			p{
				"some content",
				"more content"
			},
			div{"page has some dynamic values", params.count},
			Component{name="daniel"}
		}
	}
end
```

And this is the html that it would generate:

```html
<html>
  <head>
    <title>this is my page</title>
  </head>
  <body>
    <h1>This is my small page I wrote</h1>
    <p>Some content more content</p>
    <div>page has some dynamic values 5</div>
    <div>my name is daniel</div>
  </body>
</html>
```

# Aritecture

The actual code is divided into two sections, the golang code, and the lua code.

## Golang

The golang code is responsible for managing the different lua runtime globals, and setting them up with the correct environment to run in. Golang handles the final writing to the `io.Writer` and converts the golang structues to lua tables.

## Lua

Lua is where most of the interesting code actually live currently. This is where the global enviroment is setup with the correct lua metatable to support the `html{}` syntax without needing to define every possible element as its own function.

Lua is also where the code that takes the lua tables and renders them correctly to a witer lives. This code should probably be rewitten in golang at some point, as it will probably work faster. Its not a high priority though as it works right now.

## Markdown

Markdown currently lives as lua code downloaded into the project. I want to convert this over to golang, and maybe expose a function in lua that can be called, but I haven't bothered to do that yet as the current one works. Writing markdown inside of a lua string is not a good developer experience, so I am looking at better ways to do that.

The library I am using is from https://github.com/mpeterv/markdown. For ease of use it is copied into this repository, as the installation instructions describe.

# How does this work?

## Optional Parenthesis

Parenthesis in lua are optional when making a function call. `print("string")` and `print "string"` are identical. This removes a lot of the extra fluff that would clutter up this implementation in a different language.

## Tables combine arrays and maps

Tables in lua can be used as arrays and maps. `myTable = {"first", "second"}` is a table with two elements that can be accessed like an array: `myTable[1]`. Tables also can contain elements indexed by a key. `myMap = {key="value"}` is a table that has a single entry that can be accessed with the correct key: `myTable.key` or `myTable["key"]` are the same.

Lua also lets you do both at the same time! Which is a super powerful way to combine both index specific values, and keyed values into a single table: `someTable = {"first","second"", third="value"}`.

I use this ability to set both the children of the html node, and the parameters of the html node with a single table. For example this creates a div with an id and text elements: `div{"Am I made of chesse or not?", id="Moon"}` renders to `<div id="Moon">Am I made of cheese or not?</div>`.

## Lookup tables can be used to handle undefined functions

Lua allows an enrvionment to be set when calling a function. Any function that has not been defined upto that point in the envrionment can be dynaically handled. By using this functionality I do not have to define a function for the html tags. I can just setup the envrionment table correctly, and then it just works.

```lua

-- we need to store off the original Global envrionment so that we can access it later
local Global = _G
-- and we need a replacement
local newGlobal = {}

-- setup the replacement with a metatable that will do the lookup for us
setmetatable(newGlobal,{
	__index= function(_, name)
        -- if the function has been defined globally, return it.
        -- otherwise return our new function
		return Global[n] or function(params)
            -- This is where you put the custom code that this dynamic function should run
		end
	end,
})

-- set the new Global environment up
setfenv(0, newGlobal)

```

## Examples

There is an example runner in the `/example` folder. This can be run like so: `go run ./example/main.go`. It will print out all the examples that can be run and once can be passed as the first argument to reder the output to stdout: `go run ./example/main.go simple`

Current examples are:

- **Simple** | Just to show how this works
- **Component** | Shows how to import a reusable component and use it
- **Optional** | Shows a quirk with null values in a lua table and how to get around them
- **Markdown** | Shows how to use the built in markdown converter to render a markdown string to html
