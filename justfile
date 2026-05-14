set windows-shell := ["powershell.exe", "-NoLogo", "-Command"]
set shell := ["bash", "-c"]

APP := "owlet"
BUILD_FLAGS := "-trimpath"
LD_FLAGS := "-s -w"

default:
    just --list

build:
    go build {{BUILD_FLAGS}} -ldflags="{{LD_FLAGS}}" -o {{APP}} .

[env('GOOS', 'windows'), env('GOARCH', 'amd64')]
build-windows:
    go build {{BUILD_FLAGS}} -ldflags="{{LD_FLAGS}}" -o {{APP}}-windows-amd64.exe .

[env('GOOS', 'linux'), env('GOARCH', 'amd64')]
build-linux:
    go build {{BUILD_FLAGS}} -ldflags="{{LD_FLAGS}}" -o {{APP}}-linux-amd64 .

[env('GOOS', 'darwin'), env('GOARCH', 'arm64')]
build-macos:
    go build {{BUILD_FLAGS}} -ldflags="{{LD_FLAGS}}" -o {{APP}}-darwin-arm64 .

build-all: build-windows build-linux build-macos