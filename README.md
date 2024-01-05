[![Ceasefire Now](https://badge.techforpalestine.org/default)](https://techforpalestine.org/learn-more)

A simple go web server to wake up the display. Useful for example for a rpi acting as a smart home controller and wanting for example the screen to turn on when the doorbell rings.

To build for linux arm64 (such as a raspberry pi) do:

```
make
```

To build for the local machine do:

```
make local
```

To run on the Raspberry pi for example:

```
./activatedisplay_linux_arm64
```

To test it works, do from another machine:

```
curl http://<ip_of_pi>:8080/activate-display
```