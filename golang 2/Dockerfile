# docker build -t golang-server .

FROM golang

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY *.go ./

# RUN GO GET -t
RUN go get -t
RUN go build -o main .

EXPOSE 5000

CMD ["go", "run", "main.go"]