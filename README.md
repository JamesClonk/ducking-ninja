ducking-ninja
=========

a remote command execution webservice/app, written in [Go](https://golang.org/).

![Screenshot](https://github.com/JamesClonk/ducking-ninja/raw/master/assets/images/screenshot.png "Screenshot")

### Installation

1. Install [Go](https://golang.org/)
2. Make sure `$PATH` contains `$GOPATH/bin`
3. `go get github.com/nitrous-io/goop`
4. `git clone https://github.com/JamesClonk/ducking-ninja.git`
5. `cd ducking-ninja`
6. `goop install`
6. `goop go build`

### Usage

1. `vi auth.json`
2. `vi commands.json`
3. `./ducking-ninja` or `PORT=3456 ./ducking-ninja`

### Dockerize

1. `docker build -t ducking-ninja .`
2. `docker run -i -t -p 3333:3333 ducking-ninja`

=============
*disclaimer: name suggested by github ;)*
