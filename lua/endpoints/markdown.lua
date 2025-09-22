local content = [[
# Description

This is how markdown is handled inside of the lua files.

## not everything works as expected, as the library isn't too good:

**NOTE:** I need to use a better library

- Lists
- **Bold words**
- *Italics*
- ***Bold and itialics***


> blockquotes

> multi line
>
> blockquotes

> nested blockquotes
>> this is nested

> blockquotes containing other elements
> - like tables
> - **or bold or *italic* text**

]]

return function()
  return md(content)
end
