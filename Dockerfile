FROM golang:1.4-onbuild

MAINTAINER JamesClonk

RUN go get github.com/JamesClonk/ducking-ninja
EXPOSE 3333
CMD ["ducking-ninja"]
