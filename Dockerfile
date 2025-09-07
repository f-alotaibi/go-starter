# Build.
FROM golang:1.25 AS build-stage
WORKDIR /app
RUN go install github.com/a-h/templ/cmd/templ@latest
COPY go.mod go.sum ./
RUN go mod download
COPY . /app
RUN touch /app/.keep
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /entrypoint

# Deploy.
FROM gcr.io/distroless/static-debian12 AS release-stage
WORKDIR /

# To keep /app owned by nonroot
COPY --from=build-stage --chown=nonroot:nonroot /app/.keep /app/.keep

COPY --from=build-stage /entrypoint /app/entrypoint
COPY --from=build-stage /app/assets /app/assets

# Move to enviroment variables in production
# COPY --from=build-stage /app/.env /app/.env

WORKDIR /app
USER nonroot:nonroot
ENTRYPOINT ["/app/entrypoint"]