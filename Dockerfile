FROM golang

WORKDIR /go/src/filas-backend
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 8080
CMD ["filas-backend"]
