###########
# builder #
###########

FROM golang:1.20-buster AS builder
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    upx-ucl

WORKDIR /build
COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 go build -o ./bin/goboy \
    -ldflags='-w -s -extldflags "-static"' \
    . \
 && upx-ucl --best --ultra-brute ./bin/goboy

###########
# release #
###########

FROM golang:1.20-buster AS release
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
    git

COPY --from=builder /build/bin/goboy /bin/
WORKDIR /workdir
ENTRYPOINT ["/bin/goboy"]
