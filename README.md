# This repo is in a "proof of concept" phase. Take it as a pre-alpha.

The goal is to have a flexible method to have more complex logic in NeosVR that I do not want to write and
maintain in LogiX. But if possible this should be as generic as possible so it may be of an actual use 
at some point in time.

## Websocket callable lambda functions written in lua

This app will start a HTTP server on port 8081 and provide a few endpoints:

* POST `/lambda/{namespace}`   Initialize a common namespace. A namespace is a container for a bunch of persistent variables. As a body provide a valid JSON array of strings. Each string will be setup as an empty persistent variable... Why? Why not I guess.
* POST `/lambda/{namespace}/{function name}`   Adds a function with the given name to an internal list of lambda functions. As body provide lua code that will be executed in its own state when called
* GET/WS `/connect/{namespace}`   Endpoint to connect to via websocket

The websocket connection allows you to send the name of lambda function prior setup via the HTTP call metioned above.

## lua API

The lua state is injected with a table called `neos`. There are 3 functions on that table.
* neos.persist(varName, varValue) Saves var on an internal map.
* neos.load(varName) Reads var from internal stack and returns whatever was stored
* neos.update(varName) Send the current value of var to the websocket connection

## Example

```
POST /lambda/test
["a", "b"]

POST /lambda/test/init
neos.persist("a", "1")
neos.persist("b", "2")

POST /lambda/test/fn1
local a = tonumber(neos.load("a"))
local a = a + 1
neos.persist("a", a)
neos.update("a")

WS /connect/test
> init
> fn1
< a|2
> fn1
< a|3
```