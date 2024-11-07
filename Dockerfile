# syntax=docker/dockerfile:1

FROM golang:1.22-alpine AS build

ARG GITHUB_TOKEN
ENV GITHUB_TOKEN=$GITHUB_TOKEN
ENV GOPRIVATE=github.com/volvo-cars

RUN sed -i 's/https/http/' /etc/apk/repositories
RUN apk --update add --no-cache git
RUN git config --global url."https://$GITHUB_TOKEN@github.com".insteadOf "https://github.com"

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o /app/go-app ./cmd/authz/main.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /app
COPY --from=build /app/go-app /app/go-app
COPY --from=build /app/iam /iam
EXPOSE 8080 8081
USER nonroot:nonroot
ENTRYPOINT ["/app/go-app"]
