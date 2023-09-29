FROM golang:1.19-alpine as builder
RUN apk update && apk add git bash build-base && rm -rf /var/cache/apk/* \
  && mkdir -p /github.com/simagix/bond && cd /github.com/simagix \
  && git clone --depth 1 https://github.com/simagix/bond.git
WORKDIR /github.com/simagix/bond
RUN ./build.sh
FROM alpine
LABEL Ken Chen <ken.chen@simagix.com>
RUN addgroup -S simagix && adduser -S simagix -G simagix
USER simagix
WORKDIR /home/simagix
COPY --from=builder /github.com/simagix/bond/dist/bond /bond
CMD ["/bond", "--version"]
