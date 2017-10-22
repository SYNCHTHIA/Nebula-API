FROM alpine:latest

WORKDIR /app

# Install Package
RUN set -x && \
    mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2 && \

    # Add Certificates
    apk --no-cache add ca-certificates && \

    # Create Artifact Folder
    mkdir -p /app

# COPY Bin
COPY nebula /app

CMD ["/app/nebula"]