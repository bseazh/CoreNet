# TODO (CoreNet)

- [ ] Generate go-zero services from API definitions (`internal/api/*.api`) using `goctl` and scaffold logic.
- [ ] Wire adapters in `pkg/` into the service contexts (S3 storage, CoreKV, Kafka, MinerU) and fill upload workflow handlers.
- [ ] Implement persistence layer (MySQL migrations + access layer) for file metadata and user state.
- [ ] Build and test the upload pipeline `/upload/init|chunk|status|complete`, including chunk merge and Kafka `filemeta.write.v1` publishing.
- [ ] Implement OCR task queue: TaskCenter endpoints, Kafka consumer worker, MinerU integration, and storage of OCR outputs.
- [ ] Provide preview/search features (presigned URLs, OCR snippet search) and add integration tests for core flows.
- [ ] Extend observability and background jobs (metrics, retries, reconciliation) once core flows are stable.
