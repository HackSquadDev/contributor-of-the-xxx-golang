# FROM golang:1.19-buster AS builder
FROM mcr.microsoft.com/playwright:focal
RUN apt-get update \
    && apt-get install apt-utils -y \
	&& apt-get install cron -y \
	&& apt-get install vim -y 
RUN curl -OL https://golang.org/dl/go1.19.2.linux-amd64.tar.gz
RUN tar -C /usr/local -xvf go1.19.2.linux-amd64.tar.gz
RUN export PATH=$PATH:/usr/local/go/bin

WORKDIR /go/src/
COPY go.mod go.sum ./
RUN /usr/local/go/bin/go mod download && /usr/local/go/bin/go mod verify
# RUN npx playwright install-deps
COPY . .
# RUN mkdir -p images/

# COPY cron-file.sh /etc/cron.d/cronjob-container

RUN /usr/local/go/bin/go run main.go & /usr/local/go/bin/go run imageExtractor/imageExtractor.go

# RUN mkdir images/
# RUN /usr/local/go/bin/go run imageExtractor/imageExtractor.go

# Running commands for the startup of a container.
# CMD ["go run main.go &", "/bin/bash", "-c", "/script.sh && chmod 644 /etc/cron.d/cronjob-container && cron && tail -f /var/log/cron.log"]

# # Crontab file copied to cron.d directory.
# COPY ./files/cronjob /etc/cron.d/container_cronjob

# # Script file copied into container.
# COPY ./files/script.sh /script.sh

# # Giving executable permission to script file.
# RUN chmod +x /script.sh

# RUN CGO_ENABLED=0 GOOS=linux \
#     go build -o /go/bin/app .

# FROM gcr.io/distroless/static
# COPY --from=builder /go/bin/app /go/bin/app
# ENTRYPOINT ["/go/bin/app"]

