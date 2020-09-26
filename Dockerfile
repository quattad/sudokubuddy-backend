## DOCKER BUILD CMDS


## docker build -t sudokubuddy-backend:latest .

## Run  
# docker run -p 8080:8080 --network="test-network" sudokubuddy-backend:latest

## For dev environment, connect to local db via 192.168.1.99:3306 where <host_ip>:<port>
## Make HTTP request to 192.168.99.100:8080 where <docker_machine_ip>:<forwarded_port>

## Reference original image for go
FROM golang:latest

## Install application globally
RUN go get -u github.com/quattad/sudokubuddy-backend/src/main

## Configure environment variables
ENV PORT=8080
ENV DB_DRIVER=mysql
ENV DB_USER=root
ENV DB_PASSWORD=Sudokubuddy123!
ENV DB_NAME=tcp(172.18.0.2:3306)/sudokubuddy
# ENV DB_NAME=tcp(sudokubuddy-db.cspmexpcmnes.us-east-2.rds.amazonaws.com:3306)/sudokubuddy
ENV API_SECRET=WhHYReR5BzKkmUeoc5gIZOTmXNjLivlZ
ENV ROOT_URL=github.com/quattad/sudokubuddy-backend

## Add maintainer info
LABEL maintainer="Jonathan Quah <quahjieren@gmail.com>"

## Set current working directory in container
WORKDIR /app

## Copy go mod and go sum files
COPY go.mod go.sum ./

## Download dependencies
RUN go mod download

## Copy source from current to working directory in container
COPY . .

## Run go app
RUN cd src/main && go build -o main .

## Expose port 8000
EXPOSE 8080

## Run executable
CMD ["./src/main/main"]
