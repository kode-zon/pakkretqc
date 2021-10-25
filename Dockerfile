FROM golang:1.17

WORKDIR /app
COPY . /app/pakkretqc

ARG PAKKRETQC_ALM_ENDPOINT=https://your.qcweb.server:9999
WORKDIR /app/pakkretqc
RUN pwd && ls -la && ./build.sh
ENTRYPOINT ["tail", "-f", "/dev/null"]