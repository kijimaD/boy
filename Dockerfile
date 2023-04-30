###########
# builder #
###########

FROM golang:1.20-buster AS builder
RUN apt-get update \
    && apt-get install -yq \
    upx-ucl \
    xorg-dev \
    libx11-dev \
    libxrandr-dev \
    libxinerama-dev \
    libxcursor-dev \
    libxi-dev \
    libopenal-dev  \
    libasound2-dev \
    libgl1-mesa-dev

WORKDIR /build
COPY . .

RUN GO111MODULE=on go build -o ./bin/goboy \
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
