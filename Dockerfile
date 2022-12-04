FROM golang:1.19.1-alpine AS build
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . .
RUN go build -o /godocker main.go


FROM scratch
WORKDIR /
COPY --from=build /godocker /godocker
EXPOSE 8080
ENTRYPOINT ["/godocker"]