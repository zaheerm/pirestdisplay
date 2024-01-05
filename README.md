[![Ceasefire Now](https://badge.techforpalestine.org/default)](https://techforpalestine.org/learn-more)

A simple go app to wake up the display when receiving an mqtt message. Useful for example for a rpi acting as a smart home controller and wanting for example the screen to turn on when the doorbell rings.

To build for linux arm64 (such as a raspberry pi) do:

```
make
```

To build for the local machine do:

```
make local
```

To run on the Raspberry pi, first modify the config.yaml to match your mqtt server details and then run like follows:

```
./activatedisplay_linux_arm64
```

To test it works, use an mqtt client like mqtt explorer or use the test app inside the tools directory.