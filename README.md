# gorilla-mux-routing-example
This is an example app showing how gorilla muxers, routes, submuxers, and middleware all interact.

The intent of this is to provide an easy playground for exploring edge cases in gorilla's muxer and to help folks understand (through playing with it) how things work.

Use two terminals to play with this - one to run the server and watch its output and a second to issue `curl`s against the server.

It will listen on localhost:8080

## this branch submux_dag

this branch is to play with the idea of using submuxes in a directed acyclic graph (dag).

the idea:
* I have many endpoints
* some of them can be accessed on hostname B
* all of them can be accesesd on hostname A

I'd like to have one submux restricted to hostname A, one to hostname B. But I don't want to define the selected group that can be accessed from hostname B twice, so I'd like to have a submux that defines all the endpoints, then have the host-restricted submux let a few prefixes pass through.

```

             --> match public.domain --> only  /public/*  --
           /                                                 \
 ---> root   --> match  private.domain --------------------------> mux with all handlers



```

since I'm testing with domains, I'll need to add a `-H "Host: public.domain"` or `-H "Host: private.domain"` to the curl commands. There's middleware on each router to show the execution path.

It's interesting to note that though the public submux uses a PathPrefix to limit matches to the `/public` subtree, the main submux does _not_ automatically prepend that path to its routes in the way that it would if it were invoked with `PathPrefix("/public").Subrouter()`.


in one terminal:
```bash
$ go run main.go
```

in the other terminal:
```bash
$ curl -H "Host: public" localhost:8080/public/e1
found e1
$ curl -H "Host: public" localhost:8080/private/e2
404 page not found
$ curl -H "Host: private" localhost:8080/public/e1
found e1
$ curl -H "Host: private" localhost:8080/private/e2
found e2
```

When you issue each curl, you should also see output in the first terminal
window showing a log entry when each middleware or handler is entered and
exitted so you can see call order.

```bash
➜  go run main.go
INFO[0000] launching example app
INFO[0132] starting mRoot           <-- curl public public/e1
INFO[0132] starting mPub
INFO[0132] starting mMain
INFO[0132] starting e1
INFO[0132] ending e1
INFO[0132] leaving mMain
INFO[0132] leaving mPub
INFO[0132] leaving mRoot
                                    <-- 404 skips middleware
INFO[0146] starting mRoot           <-- curl private public/e1
INFO[0146] starting mPriv
INFO[0146] starting mMain
INFO[0146] starting e1
INFO[0146] ending e1
INFO[0146] leaving mMain
INFO[0146] leaving mPriv
INFO[0146] leaving mRoot

INFO[0154] starting mRoot           <-- curl private private/e2
INFO[0154] starting mPriv
INFO[0154] starting mMain
INFO[0154] starting e2
INFO[0154] ending e2
INFO[0154] leaving mMain
INFO[0154] leaving mPriv
INFO[0154] leaving mRoot
```
