# gorilla-mux-routing-example
This is an example app showing how gorilla muxers, routes, submuxers, and middleware all interact.

The intent of this is to provide an easy playground for exploring edge cases in gorilla's muxer and to help folks understand (through playing with it) how things work.

Use two terminals to play with this - one to run the server and watch its output and a second to issue `curl`s against the server.

It will listen on localhost:8080

## On this branch

time to explore 404 handlers. How do 404 handlers on submuxes affect future route matching? 404 without a default handler doesn't run middleware; does adding a custom 404 handler change that? Is it better to use a 404 handler to redirect or a wildcard path match?

### no submuxes
* `*mux.Router.NotFoundHandler` <-- takes a handler to run when a route isn't matched.
* does not run middleware when a route is not found

### not found on a submux
* request matches submux's prefix
* middleware on parent mux still runs
* middleware on submux does not run
* 404 handler on submux runs

### not found on subrouter with a path that matches later in the parent mux
* pathprefix /s1 on a subrouter
* path /s1/root defined on parent router
* if subrouter has a custom NotFoundHandler
** request for /s1/root hits subrouter's NotFoundHandler
* if subrouter has the default NotFoundHandler
** request for /s1/root hits the path defined on the parent router

in one terminal:
```bash
$ go run main.go
```

in the other terminal:
```bash
$ curl localhost:8080/s2/aoeu
found e2
$ curl localhost:8080/s2/actuallys1
found e1
$ curl localhost:8080/s2/actuallys3
found e2

```

When you issue each curl, you should also see output in the first terminal
window showing a log entry when each middleware or handler is entered and
exitted so you can see call order.

```bash
INFO[0004] starting mroot
INFO[0004] starting m1
INFO[0004] starting e1
INFO[0004] ending e1
INFO[0004] leaving m1
INFO[0004] leaving mroot
INFO[0010] starting e2
INFO[0010] ending e2
```
