# Base build image
# Use golang:1.19 and debian as `wasmd` requires 1.19+ and not compatible with alpine
# Reference:
#   + https://github.com/CosmWasm/wasmd/tree/626966841eeea6e457e4cdae32cd61a2e8800d71
#   + https://github.com/CosmWasm/wasmvm/tree/818281ce1aa783afb7f937352bbb63d44845fb49
FROM bitnami/golang:1.20 AS build_base

RUN echo 'Building Coreum Service'

# Display go version info to know minor version
RUN go version

# Force the go compiler to use modules
ENV GO111MODULE=on

WORKDIR /app

# We want to populate the module cache based on the go.{mod,sum} files.
RUN --mount=id=original-src,type=bind,target=/app/original_src \
    cd ./original_src && find . -name "go.*" -exec cp --parents {} /app \;

#This is the ‘magic’ step that will download all the dependencies that are specified in
# the go.mod and go.sum file.
# Because of how the layer caching system works in Docker, the  go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download

# Here we copy the rest of the source code
COPY ./ ./
# And compile the project
RUN GOFLAGS="-buildvcs=false" CGO_ENABLED=1 GOOS=linux go build -o bin/coreumservice ./

# Prepare the related dynamically linked lib
RUN mkdir -p /app/go-pkg && \
    find /go/pkg -type f -iname "*.so" -exec cp --parents {} /app/go-pkg \;

# In this last stage, we start from a fresh Mini Debian image, to reduce the image size and not ship the Go compiler in our production artifacts.
FROM bitnami/minideb:bullseye-amd64

WORKDIR /app

# We add certificates otherwise it cannot call AWS secrets manager
RUN apt-get update && apt-get install -y ca-certificates

# Finally we copy the statically compiled Go binary
COPY --from=build_base /app/bin/coreumservice ./bin/coreumservice

# Also copy the dynamically linked lib from the build base
COPY --from=build_base /app/go-pkg /

# Load things used during runtime
COPY ./env ./env

# Build time arguments
ARG STAGE
ENV STAGE=$STAGE

# Run the compiled binary
ENTRYPOINT ["./bin/coreumservice"]
EXPOSE 5011
