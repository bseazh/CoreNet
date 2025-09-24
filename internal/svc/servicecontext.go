package svc

import (
	"time"

	"corenet/internal/config"
	"corenet/pkg/kv"
	"corenet/pkg/mineru"
	"corenet/pkg/queue"
	"corenet/pkg/storage"

	"github.com/zeromicro/go-zero/core/logx"
)

type ServiceContext struct {
	Config        config.Config
	Storage       storage.S3Client
	KV            kv.KV
	KafkaProducer queue.Producer
	MinerU        mineru.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	ctx := &ServiceContext{
		Config: c,
	}

	if c.S3.Endpoint != "" {
		s3Client, err := storage.NewMinioClient(
			c.S3.Endpoint,
			c.S3.AccessKey,
			c.S3.SecretKey,
			c.S3.Bucket,
			c.S3.UseSSL,
		)
		if err != nil {
			logx.Must(err)
		}
		ctx.Storage = s3Client
	}

	if c.CoreKV.DataDir != "" {
		kvStore, err := kv.NewCoreKV(c.CoreKV.DataDir)
		if err != nil {
			logx.Must(err)
		}
		ctx.KV = kvStore
	}

	if len(c.Kafka.Brokers) > 0 {
		producer, err := queue.NewSaramaProducer(c.Kafka.Brokers, nil)
		if err != nil {
			logx.Must(err)
		}
		ctx.KafkaProducer = producer
	}

	if c.MinerU.Endpoint != "" {
		timeout := time.Duration(c.MinerU.TimeoutSec) * time.Second
		if timeout <= 0 {
			timeout = 60 * time.Second
		}
		client, err := mineru.NewHTTPClient(c.MinerU.Endpoint, timeout)
		if err != nil {
			logx.Must(err)
		}
		ctx.MinerU = client
	}

	return ctx
}
