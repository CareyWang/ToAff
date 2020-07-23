FROM golang:1.14-alpine AS dependencies
WORKDIR /app
RUN go env -w GO111MODULE="on"

COPY go.sum go.mod ./
RUN go mod tidy

FROM dependencies as build
WORKDIR /app
COPY main.go ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o toaff main.go

FROM scratch
WORKDIR /app
COPY --from=build /app/toaff ./
EXPOSE 8006
ENTRYPOINT ["/app/toaff"]
