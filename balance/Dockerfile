FROM golang:rc as build

WORKDIR /build

COPY . .
RUN go install -installsuffix cgo ./balance

FROM alpine

COPY --from=build /go/bin .

ENTRYPOINT ["svc"]