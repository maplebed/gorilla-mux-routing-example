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
  * request for /s1/root hits subrouter's NotFoundHandler
* if subrouter has the default NotFoundHandler
  * request for /s1/root hits the path defined on the parent router

in one terminal:
```bash
$ go run main.go
```

in the other terminal:
```bash➜  for i in /s1/aoeu /s1/root /s2/aoeu /s2/root ; do curl "localhost:8080$i" ; done
responding notFoundSub
responding notFoundSub
responding notFound
found e2root
```

When you issue each curl, you should also see output in the first terminal
window showing a log entry when each middleware or handler is entered and
exitted so you can see call order.

```bash
INFO[0006] starting notFoundSub
INFO[0006] ending notFoundSub
INFO[0006] starting notFoundSub
INFO[0006] ending notFoundSub
INFO[0006] starting notFound
INFO[0006] ending notFound
INFO[0006] starting e2root
INFO[0006] ending e2root
```

This illustrates a call to
* `/s1/aoeu` is answered by the s1 submux NotFoundHandler
* `/s1/root` is answered by the s1 submux NotFoundHandler (despite having another route defined)
* `/s2/aoeu` is answered by the s2 submux root mux NotFoundHandler
* `/s2/root` is answered by the specific route handler
