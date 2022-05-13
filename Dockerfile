FROM golang:alpine3.15

LABEL maintainer = "Rofviv <royvillarroel94@gmail.com>"
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build

CMD [ "./maps.patio.com" ]