# Milestones & Tasks

M0 Infra
- [ ] docker-compose up: mysql/minio/kafka/mineru-mock
- [ ] create bucket 'gopan' in MinIO console

M1 Adapters
- [ ] pkg/storage MinIO
- [ ] pkg/kv CoreKV
- [ ] pkg/queue Kafka

M2 Upload
- [ ] /upload/init|chunk|status|complete
- [ ] merge & S3 put; Kafka filemeta.write.v1
- [ ] second-pass (sha1->fileId) DB+KV

M3 FileMeta & Preview
- [ ] file lookup KV->DB
- [ ] preview presigned URL (image/pdf)
- [ ] migrations apply

M4 OCR
- [ ] taskcenter create/get
- [ ] ocrworker consume -> MinerU -> S3 -> KV snippets
- [ ] search MVP

M5 Video Preview (optional)
- [ ] transcodeworker ffmpeg->HLS->S3
- [ ] preview returns m3u8

M6 Observability
- [ ] metrics/logs/traces
- [ ] backfill/retry/reconciliation jobs
