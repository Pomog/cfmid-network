# ──────────────────────────
# Stage 1: Build the Go wrapper (modules OFF)
# ──────────────────────────
FROM golang:1.24 AS builder
WORKDIR /cfmid

# Disable modules so we can compile without go.mod/go.sum
ENV GO111MODULE=off

# Copy only the wrapper source
COPY wrapper.go ./

# Build the binary directly
RUN go build -o cfm-server wrapper.go

# ──────────────────────────
# Stage 2: Final image with CFM-ID
# ──────────────────────────
FROM wishartlab/cfmid:latest
WORKDIR /cfmid

# Copy the compiled binary from the builder stage
COPY --from=builder /cfmid/cfm-server /cfmid/cfm-server

# Expose the HTTP port
EXPOSE 5001

# Run the server
CMD ["./cfm-server"]