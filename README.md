# This repo is in a "proof of concept" phase. Take it as a pre-alpha.

The goal is to have a flexible method to have more complex logic in NeosVR that I do not want to write and
maintain in LogiX. But if possible this should be as generic as possible so it may be of an actual use 
at some point in time.

## Websocket callable lambda functions written in lua

This app will start a HTTP server on port 8081 and provide a few endpoints:

* POST `/lambda`   Initialize a session. A session is a container for a bunch of persistent variables. As a body you can provide some initial persistent variables like `a=1|b=hello`. It returns the session token to use for all lambda calls and 
* POST `/lambda/{session token}/{function name}`   Adds a function with the given name to an internal list of lambda functions. As body provide lua code that will be executed in its own state when called
* GET/WS `/connect/{session token}`   Endpoint to connect to via websocket

The websocket connection allows you to send the name of lambda function prior setup via the HTTP call metioned above.

## lua API

The lua state is injected with a table called `neos`. There are 5 functions on that table.
* neos.persist(varName, varValue) Saves var on an internal map.
* neos.load(varName) Reads var from internal stack and returns whatever was stored
* neos.send(varName, varValue) Send the varValue as VarName to the ws client
* neos.tonumber(val) Converts the given string to a number
* neos.tostring(val) Converts given value to a string

## Example

```
POST /lambda
a=1|b=3
< abc123

POST /lambda/abc123/fn1
a = neos.tonumber(neos.load("a"))
a = a + 1
max = math.max(a, b)
neos.persist("a", a)
neos.send("max", max)

WS /connect/abc123
> fn1
< max|3
> fn1
< max|3
> fn1
< max|4
```