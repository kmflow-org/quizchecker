FROM golang:1.23.0 AS build

WORKDIR /app

COPY . /app

RUN pwd && cat go.mod && ls -lathr && go mod download
RUN go test
RUN CGO_ENABLED=0 go build -o quizchecker .

FROM ubuntu
WORKDIR /app
COPY --from=build /app/quizchecker /app/quizchecker
COPY --from=build /app/config.yaml /app/config.yaml

EXPOSE 8082
CMD ["/app/quizchecker"]
