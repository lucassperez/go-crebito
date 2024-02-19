### BASE ###

FROM golang:1.22 AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download && go mod verify

### DEV ###

FROM base AS dev

RUN go install github.com/cosmtrek/air@latest

RUN cat <<EOF >> /etc/bash.bashrc
alias ls='ls --color=auto -F --group-directories-first'
PS1='${debian_chroot:+($debian_chroot)}\u@\h:\e[35m\w\e[0;1m \\$\e[0m '
EOF

COPY . .

CMD air --build.cmd='go build -buildvcs=false -o ./tmp/main .' --build.bin='./tmp/main'

### PROD ###

FROM base AS prod

COPY . .
RUN rm -f /app/tmp/*

RUN go build -v -o /usr/local/bin/go-crebito

CMD ["/usr/local/bin/go-crebito"]
