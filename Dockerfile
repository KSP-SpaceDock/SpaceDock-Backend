FROM golang:1.9-alpine

# Install curl, git
RUN apk add --no-cache curl, git

# Configure DB
ENV DB_DIALECT="postgres" \
        DB_CONNECTION_STRING="postgresql://spacedock:spacedock@host/spacedock?sslmode=disable" \
        REDIS_STRING="redis://localhost:6379/0"

# SpaceDock Backend repository location
ENV SDB_REPOSITORY="github.com/KSP-SpaceDock/SpaceDock-Backend"

# Make the source code path
RUN mkdir -p /go/src/$SDB_REPOSITORY

# Add all source code
ADD .  /go/src/$SDB_REPOSITORY

# Change working dir to app location
WORKDIR /go/src/$SDB_REPOSITORY

# Add plugins.txt
RUN touch build/plugins.txt

# Install Glide for package management
RUN curl https://glide.sh/get | sh

# Fetch plugins
RUN chmod +x build/fetch_plugins.sh
RUN sh build/fetch_plugins.sh

# Setup configuration file
RUN cp config/config.example.yml config/config.yml \
        && sed -i 's!postgres\b!$DB_DIALECT!g' config/config.yml \
        && sed -i 's!postgresql://user:password@host/dbname?sslmode=disable!$DB_CONNECTION_STRING!g' config/config.yml \
        && sed -i 's!redis://localhost:6379/0!$REDIS_STRING!g' config/config.yml

# Install dependencies
RUN glide install

# Build the app
RUN go build -v -o sdb build_sdb.go

RUN mkdir /output/ \
        && mkdir /output/config \
        && cp sdb /output/sdb \
        && cp config/config.yml /output/config.yml



# Setup actual app
FROM alpine
WORKDIR /app
COPY --from=builder /output/sdb /app/sdb
COPY --from=builder /output/config.yml /app/config/config.yml

ENTRYPOINT ['./app/sdb']

EXPOSE 5000