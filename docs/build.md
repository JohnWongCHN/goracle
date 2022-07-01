# How to build binary

## 跨平台编译

使用 `xgo` 实现跨平台编译

```bash
xgo -out ./bin/goracle --targets=linux/386,linux/amd64 -ldflags '-s -w'  .
```

## 指定平台便宜

无跨平台编译的需求时，可以直接使用 `go build` 需要注意的是 `CGO_ENABLE` 变量需设为 `1`（默认情况是开启的)

```bash
# with centos 6.10 and go version 1.15
go build -o bin/goracle-linux-amd64-el6 github.com/JohnWongCHN/goracle
```
