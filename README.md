# screenbot
A simple program that takes screenshots of a whiteboard/blackboard and sends them to a Discord bot via TCP whenever the board gets erased.

## Build & run

Dependencies:
 - [Golang](https://golang.org/)

In order to build the program, run `build.sh` and it will do all the work.  
To run it without building, use the command `go run .`

Run the program as follows (once compiled):
```
./screenbot <IP>:<port>
```
where `<address>` and `<port>` are the address and port of the server to send the images to.
