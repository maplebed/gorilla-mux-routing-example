# gorilla-mux-routing-example
This is an example app showing how gorilla muxers, routes, submuxers, and middleware all interact.

The intent of this is to provide an easy playground for exploring edge cases in gorilla's muxer and to help folks understand (through playing with it) how things work.

Use two terminals to play with this - one to run the server and watch its output and a second to issue `curl`s against the server.

It will listen on localhost:8080


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
