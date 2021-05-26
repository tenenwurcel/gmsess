FROM golang:1.14

RUN mkdir /projetos
WORKDIR /projetos
COPY . .

RUN go mod tidy

CMD make run