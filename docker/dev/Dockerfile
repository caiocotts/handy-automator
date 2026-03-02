FROM golang:1.26-alpine3.23

RUN go install github.com/cespare/reflex@v0.3.2
RUN go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.5.1
RUN go install github.com/pressly/goose/v3/cmd/goose@v3.27.0
RUN apk update && apk add npm
RUN npm i -g @redocly/cli@2.18.2

COPY ./reflex.conf /

CMD ["reflex", "--decoration=none", "--config=/reflex.conf"]
