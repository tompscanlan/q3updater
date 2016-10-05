FROM scratch
#FROM golang:1.6
MAINTAINER Tom Scanlan <tscanlan@vmware.com>

EXPOSE 9999

# Add the microservice
ADD q3updater /q3updater

CMD ["/q3updater", "--port", "9999"]
