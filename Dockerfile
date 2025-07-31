# Using Debian Bookworm as it's more stable than Ubuntu
FROM debian:bookworm

# Update, upgrade and install dependencies
RUN apt-get update -y && apt-get full-upgrade -y && apt-get install -y \
    curl \
    wget \
    gnupg \
    ca-certificates \
    git \
    build-essential \
    postgresql-client \
    && rm -rf /var/lib/apt/lists/*

# setting go version
ENV GO_VERSION=1.24.5

# download go
RUN curl -OL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz && \
    rm go${GO_VERSION}.linux-amd64.tar.gz

# Add Go to PATH
ENV PATH="/usr/local/go/bin:${PATH}"

WORKDIR /app

COPY . .

RUN go mod tidy && go build -o main .

CMD ["./main"]
