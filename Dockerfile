# Minimal Docker version >= 17.05

# build_image and base_image arguments are mandatory
ARG base_image=golang:1.12

# -----------------
# BINARY BUILD STAGE
# -----------------
FROM ${base_image} as build

ENV GO111MODULE=on

RUN mkdir -p /go/src

WORKDIR /go/src

COPY . .

RUN set -xe \
    && go mod vendor -v \
    && go build -ldflags "-linkmode external -extldflags -static" -a main.go

# -----------------
# IMAGE BUILD STAGE
# -----------------
FROM scratch

COPY --from=build /go/src/main /main

ENTRYPOINT ["/main"]

# Metadata
ARG base_image=scratch
ARG socks_version=not-set
ARG revision=not-set

LABEL org.kolesa-team.image.maintainer="Amangeldy Kadyl <lex0.kz@gmail.com>" \
      org.kolesa-team.image.name="go-socks-ldap" \
      org.kolesa-team.image.openresty="${socks_version}" \
      org.kolesa-team.image.revision="${revision}" \
      org.kolesa-team.image.base_image="${base_image}" \
      org.kolesa-team.image.description="Provides a go Socks5 server with LDAP image based on scratch."
