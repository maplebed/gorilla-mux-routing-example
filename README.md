# gorilla-mux-routing-example
This is an example app showing how gorilla muxers, routes, submuxers, and middleware all interact.

The intent of this is to provide an easy playground for exploring edge cases in gorilla's muxer and to help folks understand (through playing with it) how things work.

Use two terminals to play with this - one to run the server and watch its output and a second to issue `curl`s against the server.

It will listen on localhost:8080 and localhost:8081

This branch illustrates my current confusion on how middleware works, or maybe shows off a bug in how middleware is invoked.

See https://github.com/gorilla/mux/issues/429 for more conversation

in one terminal:
```bash
$ go run main.go
```

in the other terminal:
```bash
$ curl localhost:8080/s1/hello
found e1
$ curl localhost:8080/s2/hello
found e2
$ curl localhost:8080/s3/hello
found e3
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
INFO[0012] starting e3
INFO[0012] ending e3
```
Notice that middleware is not called for `/s2/hello` or `/s3/hello`, despite the endpoints being successfully run.

Now try it against port 8081 instead, and see that all middleware is run for all endpoinst.

```bash
$ curl localhost:8081/s1/hello
found e1
$ curl localhost:8081/s2/hello
found e2
$ curl localhost:8081/s3/hello
found e3
```

See in STDOUT:

```bash
INFO[0024] starting mroot
INFO[0024] starting m1
INFO[0024] starting e1
INFO[0024] ending e1
INFO[0024] leaving m1
INFO[0024] leaving mroot
INFO[0026] starting mroot
INFO[0026] starting m2
INFO[0026] starting e2
INFO[0026] ending e2
INFO[0026] leaving m2
INFO[0026] leaving mroot
INFO[0029] starting mroot
INFO[0029] starting m3
INFO[0029] starting e3
INFO[0029] ending e3
INFO[0029] leaving m3
INFO[0029] leaving mroot
```
