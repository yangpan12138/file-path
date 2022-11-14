FROM golang:latest AS build

WORKDIR /test
COPY . /test
RUN go env -w GOPROXY=https://goproxy.cn,direct
RUN CGO_ENABLED=0 go build -v -o main .

FROM alpine AS api
RUN mkdir /app
COPY --from=build /test/main /app
WORKDIR /app
ENTRYPOINT ["./main", "-v" ,"1.3"]