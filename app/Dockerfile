FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY ./src/go.mod ./src/go.sum ./
RUN go mod download

COPY ./src .
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -ldflags="-s -w" -o /app/taskScheduler .

FROM scratch

COPY --from=builder /app/taskScheduler /app/taskScheduler

EXPOSE 8080

CMD [ "/app/taskScheduler" ]