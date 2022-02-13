FROM golang:latest
WORKDIR /
COPY / /
RUN go build -o client ./main.go 
RUN ["chmod", "+x", "./client"]  
CMD [ "./client" ]
