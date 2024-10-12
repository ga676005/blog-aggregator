FROM postgres:latest

# Install curl and necessary certificates
RUN apt-get update && \
    apt-get install -y --no-install-recommends curl ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Install Go
RUN curl -O https://dl.google.com/go/go1.23.2.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.23.2.linux-amd64.tar.gz && \
    rm go1.23.2.linux-amd64.tar.gz

# Add Go binaries to PATH
ENV PATH="/usr/local/go/bin:${PATH}"

# Verify Go installation
RUN go version

# Install Goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Set up the work directory
WORKDIR /migrations

# Copy the migration SQL files from the host
COPY ../sql/schema /migrations

# Wait for PostgreSQL to be ready and run Goose migrations
CMD ["sh", "-c", "until pg_isready -h localhost -p 5432; do sleep 1; done; goose postgres \"postgres://postgres:postgres@localhost:5432/blogator\" up && tail -f /dev/null"]
