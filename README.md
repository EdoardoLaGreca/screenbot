# screenbot
A simple program that takes screenshots of a whiteboard/blackboard and sends them to a Discord bot via TCP whenever the board gets erased.

Note that TCP messages sent by this bot are structured as follows:
```
bytes:   0       2             10   10 + image length
         +-------+--------------+-------+
content: | "BOT" | image length | image |
         +-------+--------------+-------+
```
Where:
 - The string `BOT` is used to identify that the TCP package is from this bot and not from other internet bots sending requests at random addresses on the internet.
 - The image length is encoded into a 64 bit unsigned big-endian integer.

## Build & run

Dependencies:
 - [Golang](https://golang.org/)

In order to build the program, run `build.sh`, it will do all the work and place the binary file into the `target/` directory.

Run the program as follows (in the `target/` directory, once compiled):
```
./screenbot <IP>:<port>
```
where `<IP>` and `<port>` are the address and port of the server to send the images to.

When it starts, it will ask for the screen coordinates to take the images from, to ease this process you can use a tool called `xdotool`. By running the following command you can get the screen coordinates of the mouse.
```
xdotool getmouselocation
```
Note that this tool works for sure on Unix-like systems running Xorg while I don't know if it works as well on Wayland or WSL.

If the bot fails to send one or more images, it will store them in the `offline/` directory and send them once it will reach the remote host again. It also stores all the images (both sent and not sent) in the `imgs/` directory.
