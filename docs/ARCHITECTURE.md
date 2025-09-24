# Architecture Overview

- go-zero microservices generated from .api
- pkg/* provide adapters for infra
- CoreKV holds hot/state data (chunk status, hash->fileId, task states, hot metadata)
- MySQL holds authoritative metadata
- MinIO stores objects (orig/preview/ocr)
- Kafka for async writes & jobs
- MinerU OCR as external HTTP service
- Optional FFmpeg HLS worker

Data flows:
1) Upload -> KV chunk status -> merge -> S3 -> Kafka(filemeta) -> DB
2) OCR -> TaskCenter -> Kafka(ocr) -> Worker -> MinerU -> S3 -> KV index
3) Preview -> presigned URL
