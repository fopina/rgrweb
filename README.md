# rgrweb
Remote Garage Remote Web UI

This webapp is part of a (small) project detailed in [this post](https://rpg.skmobi.com/posts/0xcf53_rgr/).

It is a single "push" button page that will enable a chosen GPIO pin (specified in `--pin-out`) and a "virtual LED" to reflect the state of another GPIO pin (`--pin-in`).  
It is used to control a garage opener but it can certainly be used in many other *flip-switch* scenarios.

It uses [github.com/warthog618/gpio](https://github.com/warthog618/gpio) and [BCM GPIO mapping](https://pinout.xyz/) is used.

### Demo

![demo](https://github.com/fopina/rgrweb/raw/assets/Image.GIF)

### Usage

```
Usage of rgrweb:
  -b, --bind string         address:port to bind webserver (default "127.0.0.1:8081")
  -d, --duration duration   Time that output GPIO pin will be HIGH (default 5s)
  -i, --pin-in int          Input GPIO (feedback) - 0 means fake it
  -o, --pin-out int         Output GPIO (trigger) - 0 means fake it
      --test-input          Reading input GPIO for 5 seconds (testing)
      --test-output         Enable output GPIO for 5 seconds (testing)
```
