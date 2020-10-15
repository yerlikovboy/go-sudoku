FROM golang:1.15 AS build 

WORKDIR /go/src/gpd 

COPY . . 

RUN go get -d -v ./...
RUN go install -v ./...

FROM busybox:glibc
COPY --from=build ./go/bin/gpd /gpd

CMD ["/gpd"] 
