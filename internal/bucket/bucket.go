package bucket

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

func CountFilesInBucket(bucketName, srcPrefix string) (int, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return 0, fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	query := &storage.Query{Prefix: srcPrefix}

	count := 0

	it := bucket.Objects(ctx, query)
	for {
		_, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, fmt.Errorf("iterating object: %v", err)
		}
		count++
	}

	return count, nil
}

func MoveFilesToBucket(scrBucketName, destBucketName, srcPrefix, destPrefix string) error {
	ctx := context.Background()

	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage NewClient: %v", err)
	}
	defer client.Close()

	srcBucket := client.Bucket(scrBucketName)
	destBucket := client.Bucket(destBucketName)

	query := &storage.Query{Prefix: srcPrefix}

	it := srcBucket.Objects(ctx, query)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("iterating objects: %v", err)
		}

		destObjectName := destPrefix + objAttrs.Name[len(srcPrefix):]

		_, err = destBucket.Object(destObjectName).CopierFrom(srcBucket.Object(objAttrs.Name)).Run(ctx)
		if err != nil {
			return fmt.Errorf("copying object %s to %s: %v", objAttrs.Name, destObjectName, err)
		}

		err = srcBucket.Object(objAttrs.Name).Delete(ctx)
		if err != nil {
			return fmt.Errorf("deleting original object %s: %v", objAttrs.Name, err)
		}

		fmt.Printf("Moved %s to %s\n", objAttrs.Name, destObjectName)
	}
	return nil
}

func DeleteFilesWithPrefix(bucketName, prefix string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("storage.NewClient: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	query := &storage.Query{Prefix: prefix}

	it := bucket.Objects(ctx, query)
	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("iterating objects: %v", err)
		}

		err = bucket.Object(objAttrs.Name).Delete(ctx)
		if err != nil {
			return fmt.Errorf("deleting object %s: %v", objAttrs.Name, err)
		}

		fmt.Printf("Deleted %s\n", objAttrs.Name)
	}

	return nil
}
