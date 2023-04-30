# goboy

reference:

https://github.com/bokuweb/gopher-boy

## install

```
$ go install github.com/kijimaD/goboy@main
```

## run

```
$ go run main.go roms/helloworld/hello.gb
```

## development

```
sudo apt install libgl1-mesa-dev
export PKG_CONFIG_PATH=/usr/lib/x86_64-linux-gnu/pkgconfig
export DISPLAY=$(cat /etc/resolv.conf | grep nameserver | awk '{print $2}'):0
```
