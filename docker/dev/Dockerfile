FROM golang:1.26-alpine3.23 AS download

RUN go install github.com/cespare/reflex@v0.3.2

FROM golang:1.26-alpine3.23 AS final

WORKDIR /app
COPY --from=download /go/bin/reflex /bin/
COPY ./reflex.conf /
CMD ["reflex", "--decoration=none", "--config=/reflex.conf"]
