FROM golang:1.21-alpine as builder

WORKDIR /src

COPY . .

RUN go mod tidy 

RUN go build -o statusbot main.go


FROM alpine

WORKDIR /app

COPY --from=builder /src/statusbot .

CMD ./statusbot -t ${STATUSBOT_TOKEN}

