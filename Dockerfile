FROM golang:1.25.0 AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -v -o "groupie-tracker"

FROM alpine:latest

WORKDIR /app

EXPOSE 8080

COPY --from=build /app/groupie-tracker .

COPY --from=build /app/static ./static

CMD [ "./groupie-tracker" ]

LABEL name="groupie-tracker"
LABEL version="1.0.1"
LABEL maintainer="Olamide Ifarajimi"
LABEL maintainerEmail="lordelami@gmail.com"