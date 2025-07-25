package r2

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var (
  accessKey string = os.Getenv("CLOUDFARE_ACCESS_KEY") 
  secretKey string = os.Getenv("CLOUDFARE_SEC_KEY")
  accountID string = os.Getenv("CLOUDFARE_ID")
)

type Bucket string

const (
  Stream Bucket = "stream"
)

// S3Service wraps the AWS S3 client and configuration for R2.
type R2Service struct {
  *s3.Client
  Bucket   string
}

func Init(bucket string) *R2Service {
  cfg, err := config.LoadDefaultConfig(context.TODO(),
    config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
    config.WithRegion("auto"),
  )
  if err != nil {
    log.Fatal(err)
  }
  client := s3.NewFromConfig(cfg, func(o *s3.Options) {
      o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
  })
  
  return &R2Service{
    client,
    bucket,
  }
  
}

func InitStream() *R2Service {
  cfg, err := config.LoadDefaultConfig(context.TODO(),
    config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")),
    config.WithRegion("auto"),
  )
  if err != nil {
    log.Fatal(err)
  }
  
  client := s3.NewFromConfig(cfg, func(o *s3.Options) {
      o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountID))
  })
  
  return &R2Service{
    client,
    string(Stream),
  }
}


func (r *R2Service) R2MultipartUpload(ctx context.Context, body []byte, uploadPath string, contType string, fileSize int64) error {
  input := &s3.CreateMultipartUploadInput{
    Bucket:      aws.String(r.Bucket),
    Key:         aws.String(uploadPath),
    ContentType: aws.String(contType),
  }
  resp, err := r.CreateMultipartUpload(ctx, input)
  if err != nil {
    return err
  }

  var curr, partLength int64
  var remaining = fileSize
  var completedParts []types.CompletedPart
  const maxPartSize int64 = int64(50 * 1024 * 1024)
  partNumber := 1
  for curr = 0; remaining != 0; curr += partLength {
    if remaining < maxPartSize {
       partLength = remaining
    } else {
       partLength = maxPartSize
    }

    partInput := &s3.UploadPartInput{
      Body:       bytes.NewReader(body[curr : curr+partLength]),
      Bucket:     resp.Bucket,
      Key:        resp.Key,
      PartNumber: aws.Int32(int32(partNumber)),
      UploadId:   resp.UploadId,
    }
  
    uploadResult, err := r.UploadPart(ctx, partInput)
    if err != nil {
     aboInput := &s3.AbortMultipartUploadInput{
      Bucket:   resp.Bucket,
      Key:      resp.Key,
      UploadId: resp.UploadId,
     }
     _, aboErr := r.AbortMultipartUpload(ctx, aboInput)
     if aboErr != nil {
      return aboErr
     }
     return err
    }

    completedParts = append(completedParts, types.CompletedPart{
     ETag:       uploadResult.ETag,
     PartNumber: aws.Int32(int32(partNumber)),
    })
    remaining -= partLength
    partNumber++
  }

  compInput := &s3.CompleteMultipartUploadInput{
    Bucket:   resp.Bucket,
    Key:      resp.Key,
    UploadId: resp.UploadId,
    MultipartUpload: &types.CompletedMultipartUpload{
       Parts: completedParts,
     },
  }
  _, compErr := r.CompleteMultipartUpload(ctx, compInput)
  if err != nil {
    return compErr
  }

  return nil
}
