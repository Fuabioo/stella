# stella

<img src="assets/demo.gif">

**TUI progress bar ⏱️**

The process you want to monitor just needs to expose an HTTP
endpoint returning the specified JSON and *voilà* 💅🏻✨

## Installation

### Precompiled binaries

Download a precompiled binary from [the releases](https://github.com/Fuabioo/stella/releases) 

### APT on linux

```bash
echo -e 'Package: stella\nPin: origin "apt.fuabioo.com"\nPin-Priority: 1001' | sudo tee /etc/apt/preferences.d/stella.pref > /dev/null
apt-cache policy stella
```

```bash
sudo apt update
sudo apt install stella
```

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
