FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /server .

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /server /server

USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]
