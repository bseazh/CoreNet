package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	MySQL  MySQLConfig
	S3     S3Config
	CoreKV CoreKVConfig
	Kafka  KafkaConfig
	MinerU MinerUConfig
	FFmpeg FFmpegConfig
}

type MySQLConfig struct {
	DSN string
}

type S3Config struct {
	Endpoint           string
	Bucket             string
	AccessKey          string
	SecretKey          string
	UseSSL             bool
	SignedURLExpireSec int
}

type CoreKVConfig struct {
	DataDir         string
	CacheTTLSeconds int
}

type KafkaConfig struct {
	Brokers []string
	Topics  KafkaTopics
}

type KafkaTopics struct {
	FileMetaWrite string
	OCRJob        string
	TranscodeJob  string
}

type MinerUConfig struct {
	Endpoint   string
	TimeoutSec int
}

type FFmpegConfig struct {
	BinPath string
	HLS     HLSConfig
}

type HLSConfig struct {
	SegmentSec    int
	ListSize      int
	SingleVariant bool
}
