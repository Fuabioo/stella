# stella

<img src="assets/demo.gif">

**TUI progress bar ⏱️**

The process you want to monitor just needs to expose an HTTP
endpoint returning the specified JSON and *voilà* 💅🏻✨

## Installation

TODO

## Usage

Run the demo process:

```sh
go run pkg/demo/demo.go
```

Run `stella` as a demo client consumer.

> The url is the only configuration possible right now

```sh
stella http://0.0.0.0:10301/progress
```

Logs are sent to an ephemeral logfile set to `/tmp` so the OS can clean it up.

```sh
tail -f /tmp/stella/2025-05-29.log
```

## TODO

- [ ] Implement CLI + configuration
- [ ] Add a documentation page
