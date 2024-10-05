# Go CLI App Kafka Retention Manager

A powerful command-line interface for managing Kafka topic retention and data archiving to Google Cloud Storage.

## Overview

This CLI application, part of the Kafka Retention Manager suite, provides functionality to manage Google Cloud Storage (GCS) buckets in relation to Kafka topic retention. It can move files between GCS buckets or delete files from a bucket, ensuring proper data retention and archiving for your Kafka-based data pipeline.


## Features

- Count files in a GCS bucket
- Move files between GCS buckets
- Delete files with a specific prefix in a GCS bucket
- Interact with Redis to compare file counts

## Prerequisites

- Go 1.16 or later
- Access to Google Cloud Storage
- Redis server (for move operation only)
- Kafka cluster (for full pipeline functionality)

## Installation

## Installation

1. Clone the repository:
   ```
   git clone https://github.com/azharizz/kafka-retention-manager-cli.git
   cd kafka-retention-manager-cli
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

3. Build the application:
   ```
   go build -o kafka-retention-manager-cli
   ```

   ```

## Usage

Run the application with the following command:

```
./kafka-retention-manager-cli
```

The application will first ask you to choose between two operations:

1. Move files
2. Delete files

After selecting an operation, the application will interactively prompt you for the necessary information.



```
./kafka-retention-manager-cli --bucket=SOURCE_BUCKET --dest-bucket=DEST_BUCKET --src-prefix=SOURCE_PREFIX --dest-prefix=DEST_PREFIX --redis-addr=REDIS_ADDRESS
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

## Pipeline Overview

The Kafka Retention Manager CLI is part of a larger data pipeline:

1. Data Ingestion: Data is ingested into Kafka topics and simultaneously counted in Redis.
2. Data Processing: Kafka consumers process the data and store it in Google Cloud Storage (GCS).
3. Data Retention: This CLI tool manages the retention of data in GCS, either by moving it to long-term storage or deleting it based on retention policies.

## Infrastructure

The `docker` folder contains the infrastructure components:

- `goproto`: Contains the data ingestion service that pushes data into Kafka and Redis.
- `kafka`: Includes a Docker Compose file for setting up the Kafka infrastructure.
- `kafka-go-gcs`: Contains the Go code for the Kafka consumer that sends data to GCS.

## Future Improvements

While the current version is functional, there are several areas for potential improvement:

- Scheduler Integration: Implement a scheduler for automated retention management tasks.
- Protocol Buffers: Utilize Protocol Buffers for more efficient data serialization.
- List Management: Implement functionality to manage lists of topics or buckets.
- Extensible Source Stack: Add support for different source systems beyond Kafka.

## Demo

For a visual demonstration of the Kafka Retention Manager CLI and its associated pipeline, check out our YouTube video:

- Description: "Kafka Retention Manager CLI: Streamline Your Data Retention Workflow"
- Link: [Kafka Retention Manager CLI Demo](https://www.youtube.com/watch?v=aGrZOSbiRoo)

[![IMAGE ALT TEXT HERE](https://img.youtube.com/vi/CEr3pQliJwE/0.jpg)](https://www.youtube.com/watch?v=aGrZOSbiRoo)

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
