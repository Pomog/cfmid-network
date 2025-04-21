FROM wishartlab/cfmid:latest

# Install Go (needed to compile wrapper)
RUN apt-get update && \
    apt-get install -y golang-go && \
    rm -rf /var/lib/apt/lists/*

# Copy wrapper source, compile it
WORKDIR /cfmid
COPY wrapper.go .
RUN go build -o cfm-server wrapper.go

# Expose HTTP port for Nginx proxy
EXPOSE 5001

# Launch our wrapper
CMD ["./cfm-server"]