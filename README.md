#Slurp Client

This is the client side of [slurp server](http://github.com/ebuckley/slurp_server). You will definitly need a server instance to try this out. The specifications for the client are from a distributed systems paper at [columbia](http://www.cs.columbia.edu/~roxana/teaching/DistributedSystemsF12/labs/), the exercise specifications are [here](http://www.cs.columbia.edu/~roxana/teaching/DistributedSystemsF12/labs/lab0.html).

#Usage
You will need to know the port and address of the slurp server instance, in the following example it is running on port 1337 on localhost. The server name param supports hostname or IP. The example is requesting the `examples.desktop` file from the server.

```
$ ./slurp_client localhost 1337 examples.desktop
```

