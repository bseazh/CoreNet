# Gopan × CoreKV × MinerU — Starter Monorepo

This repo scaffolds a **go-zero** microservice project that stitches together:
- **Gopan**-style net-disk (upload/download/metadata) flows
- **CoreKV** as a hot/state layer (chunk status, second-pass index, task states, hot metadata)
- **MinerU** OCR (async jobs producing Markdown/JSON)
- (Optional) **FFmpeg** HLS video preview pipeline

> ✅ The skeleton favors **pluggable adapters** (pkg/*) and **API-first** development (.api files).  
> ✅ Use `goctl` to generate service code from `internal/api/*.api`, then fill logics with pkg adapters.  
> ✅ `deploy/docker-compose.yml` brings up infra (MinIO, MySQL, Kafka).

## Quickstart

1) **Prereqs**
- Go 1.21+
- Node (if you build any dashboard)
- Docker & Docker Compose
- `go install github.com/zeromicro/go-zero/tools/goctl@latest`

2) **Infra up**
```bash
cd deploy
docker compose up -d
# MinIO  : http://localhost:9001 (admin/minioadmin)
# MySQL  : 127.0.0.1:3306 (root/rootpass)
# Kafka  : 127.0.0.1:9092
```

3) **Generate services from API files**
```bash
# Example for upload service:
goctl api go -api internal/api/upload.api -dir .

# Repeat for: filemeta, preview, search, taskcenter
```

4) **Fill implementations**
- Use adapters in `pkg/kv`, `pkg/storage`, `pkg/queue`, `pkg/mineru`, `pkg/ffmpeg`.
- See comments in files for TODOs and suggested contracts.

5) **Run**
```bash
make dev    # or run each service with `go run` after generation
```

## Layout

```
repo/
├── cmd/                      # (optional) hand-written mains (you can also rely on goctl outputs)
├── internal/
│   ├── api/                  # .api definitions for goctl
│   ├── config/               # yaml configs (per service)
│   ├── logic/                # (to be generated & edited)
│   ├── svc/                  # (to be generated)
│   └── types/                # (to be generated)
├── pkg/
│   ├── kv/                   # CoreKV adapter
│   ├── storage/              # MinIO/S3 wrapper
│   ├── queue/                # Kafka wrapper
│   ├── mineru/               # MinerU HTTP client
│   ├── ffmpeg/               # HLS transcoder helper
│   └── util/                 # errs, retry, helpers
├── migrations/               # SQL DDL
├── deploy/                   # docker-compose & configs
├── docs/                     # PRD, Arch, Tasks
└── Makefile
```

## Notes

- **CoreKV** here is modeled as an **embedded** KV. You can swap with Redis if you prefer, but keep the same interface `pkg/kv`.
- **Search** is lightweight in MVP. For production full-text, plug Meilisearch/ES and replace `pkg/search`.
- **MinerU** is treated as an external HTTP endpoint; a mock is provided in compose to help wiring.
- **HLS** is optional; add `transcodeworker` only when you’re ready.

Happy building!
