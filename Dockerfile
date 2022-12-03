FROM golang:1.19-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN go build -o /godocker


FROM scratch
WORKDIR /
COPY --from=build /godocker /godocker
EXPOSE 8080
ENTRYPOINT ["/godocker"]