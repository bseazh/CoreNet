# CoreNet Frontend (Next.js)

A Google Drive–like frontend for CoreNet APIs built with Next.js (App Router), React, Tailwind CSS.

Features
- File list via `GET /search?q=`; enrich with `GET /files/:fileId`
- Upload with chunking via `/upload/init`, `/upload/chunk`, `/upload/complete`
- Preview images and videos with `GET /files/:fileId/preview`
- Delete files via `DELETE /files/:fileId`

Assumptions
- Backend base URL is `http://localhost:8888` (configurable).
- `/upload/chunk` accepts `PUT /upload/chunk?uploadId=..&index=..` with raw body chunk.
  Adjust in `frontend/lib/api.ts` if your backend differs (e.g., headers or JSON form).

Getting Started
1. cd frontend
2. Copy `.env.local.example` to `.env.local` and adjust API URL
3. Install deps: `npm i` (or `pnpm i` / `yarn`)
4. Run dev: `npm run dev` and open http://localhost:3000

Env
- NEXT_PUBLIC_API_BASE_URL: CoreNet API base (e.g., `http://localhost:8888`)

Notes
- If `/search` requires a query, return all files for empty `q` on backend, or adjust the UI to use another list endpoint when available.
- Add auth later (middleware + tokens) as needed.

