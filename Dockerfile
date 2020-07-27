FROM registry.access.redhat.com/ubi8/ubi:latest

COPY bin/internal-service-monitor /
WORKDIR /
