ARG APP_NAME="tasmota_backupper"
ARG USER=tasmota_backupper

# Build App itself
FROM golang:1.18 AS builder
COPY . /go/src/tasmota_backup
WORKDIR /go/src/tasmota_backup
RUN GOOS=linux go build -ldflags "-linkmode external -extldflags -static" -a -installsuffix cgo -v -o tasmota_backupper .


# Serve with proper permissions and etc
FROM alpine
ARG USER
ARG APP_NAME
ENV HOME /home/$USER
ENV GIN_MODE release
RUN apk add --no-cache tzdata ca-certificates && update-ca-certificates
RUN mkdir -p /app
WORKDIR /app
COPY --from=builder /go/src/tasmota_backup/tasmota_backupper tasmota_backupper
RUN chmod +x tasmota_backupper
RUN addgroup -S $USER && adduser -S $USER -G $USER
RUN chown -R $USER:$USER /app
USER $USER
EXPOSE 8080
CMD ./tasmota_backupper





