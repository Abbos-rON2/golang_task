FROM golang:latest
WORKDIR /
COPY / /
RUN go build -o server ./main.go 
RUN ["chmod", "+x", "./server"]  
CMD [ "./server" ]
