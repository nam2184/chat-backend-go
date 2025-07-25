package r2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestR2Bucket(t *testing.T) {
  r2Client := Init("stream")
 
  listObjectsOutput, err := r2Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
    Bucket: &r2Client.Bucket,
  })
  
  if err != nil {
    t.Fatal(err)
  }
  for _, object := range listObjectsOutput.Contents {
    obj, _ := json.MarshalIndent(object, "", "\t")
    fmt.Println(string(obj))
  }

  listBucketsOutput, err := r2Client.ListBuckets(context.TODO(), &s3.ListBucketsInput{})
  if err != nil {
    t.Fatal(err)
  }

  for _, object := range listBucketsOutput.Buckets {
    obj, _ := json.MarshalIndent(object, "", "\t")
    fmt.Println(string(obj))
  }
}

func TestR2MultipartUpload(t *testing.T) {
	// Initialize mock service
  r2Client := Init("stream")
  // Generate a large file in memory (e.g., 150 MB)
	const fileSize = 150 * 1024 * 1024 // 150 MB
	largeFile := bytes.Repeat([]byte("A"), int(fileSize)) // Repeating 'A' to fill the buffer

	// Call the multipart upload function
	ctx := context.TODO()
	err := r2Client.R2MultipartUpload(ctx, largeFile, "test-large-file.mp4", "video/mp4", fileSize)
	if err != nil {
		t.Fatalf("Multipart upload failed: %v", err)
	}

	t.Log("Multipart upload test passed successfully.")
}
