# CLI App for GCS Bucket Management and Redis Interaction

This CLI application provides functionality to manage Google Cloud Storage (GCS) buckets and interact with Redis. It can count files in a bucket, move them between buckets, and optionally delete files.

## Features

- Count files in a GCS bucket
- Move files between GCS buckets
- Delete files with a specific prefix in a GCS bucket
- Interact with Redis to compare file counts

## Prerequisites

- Go 1.16 or later
- Access to Google Cloud Storage
- Redis server

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/yourusername/cli-app.git
   cd cli-app
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Build the application:
   ```
   go build -o cli-app
   ```

## Usage

Run the application with the following command:

```
./cli-app --bucket=SOURCE_BUCKET --dest-bucket=DEST_BUCKET --src-prefix=SOURCE_PREFIX --dest-prefix=DEST_PREFIX --redis-addr=REDIS_ADDRESS
```

Replace the placeholders with your actual values:

- `SOURCE_BUCKET`: The name of the source GCS bucket
- `DEST_BUCKET`: The name of the destination GCS bucket
- `SOURCE_PREFIX`: The prefix for files in the source bucket
- `DEST_PREFIX`: The prefix for files in the destination bucket
- `REDIS_ADDRESS`: The address of your Redis server (default: localhost:6379)

## Configuration

You can use a configuration file instead of command-line flags. Create a `.cli-app.yaml` file in your home directory or use the `--config` flag to specify a different location.

Example configuration file:

```yaml
bucket: "source-bucket-name"
dest-bucket: "destination-bucket-name"
src-prefix: "source/prefix/"
dest-prefix: "destination/prefix/"
redis-addr: "localhost:6379"
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
