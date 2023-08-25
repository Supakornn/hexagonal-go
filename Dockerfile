FROM golang:1.20-buster AS build

WORKDIR /app

COPY . ./
RUN go mod download

RUN CGO_ENABLED=0 go build -o /bin/app

## Deploy
FROM gcr.io/distroless/static-debian11

COPY --from=build /bin/app /bin
COPY .env.prod /bin
# COPY /assets /bin/assets

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT [ "/bin/app", "/bin/.env.prod" ]