FROM golang:1.20 as build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /k8s-reloader ./cmd/main.go

FROM golang:alpine

COPY --from=build /k8s-reloader /k8s-reloader

CMD ["/k8s-reloader"]