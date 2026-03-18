FROM golang:1.26-alpine3.23 AS download

RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.1 &&\
    go install github.com/pressly/goose/v3/cmd/goose@v3.27.0

FROM golang:1.26-alpine3.23 AS final

WORKDIR /workspace
COPY --from=download /go/bin/* /bin/
RUN apk update &&\
    apk add npm --no-cache &&\
    npm i -g @redocly/cli@2.18.2
