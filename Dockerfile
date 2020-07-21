FROM golang:1.13-alpine

#Make app directory in image
RUN mkdir /app

#copy to /app
ADD . /app

#Set directory inside /app
WORKDIR /app

#Pull dependencies
RUN go mod download

#Go build
RUN go build -o main .

#Start command
CMD ["/app/main"]

