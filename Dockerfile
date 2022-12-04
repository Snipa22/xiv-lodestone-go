FROM golang:1.19.1-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go get -u
COPY *.go ./
RUN go build -o /godocker


FROM scratch
WORKDIR /
COPY --from=build /godocker /godocker
EXPOSE 8080
ENTRYPOINT ["/godocker"]