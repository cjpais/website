FROM golang:latest 
RUN mkdir /app 
ADD . /app/ 
WORKDIR /app 

RUN go get golang.org/x/crypto/...
RUN go get golang.org/x/crypto/acme/autocert
RUN go get github.com/gorilla/sessions
RUN go get github.com/nfnt/resize

RUN go build -o main . 
CMD ["/app/main", "-clear"]

EXPOSE 8080
