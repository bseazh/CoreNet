export type FileItem = {
  fileId: string
  name: string
  size?: number
  mime?: string
  snippet?: string
}

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8888'

function to<T>(res: Response): Promise<T> {
  if (!res.ok) {
    return res.text().then(t => Promise.reject(new Error(t || res.statusText)))
  }
  return res.json() as Promise<T>
}

export async function searchFiles(q: string): Promise<FileItem[]> {
  const url = new URL('/search', API_BASE)
  if (q) url.searchParams.set('q', q)
  const data = await to<{ items: { fileId: string; name: string; snippet?: string }[] }>(await fetch(url.toString(), { cache: 'no-store' }))
  return data.items || []
}

export async function getFile(fileId: string): Promise<{ fileId: string; name: string; size: number; mime: string; sha1: string; version: number }> {
  const url = `${API_BASE}/files/${encodeURIComponent(fileId)}`
  return to(await fetch(url, { cache: 'no-store' }))
}

export async function deleteFile(fileId: string): Promise<void> {
  const url = `${API_BASE}/files/${encodeURIComponent(fileId)}`
  const res = await fetch(url, { method: 'DELETE' })
  if (!res.ok) throw new Error(await res.text())
}

export async function getPreviewUrl(fileId: string): Promise<string> {
  const url = `${API_BASE}/files/${encodeURIComponent(fileId)}/preview`
  const data = await to<{ url: string }>(await fetch(url, { cache: 'no-store' }))
  return data.url
}

export async function initUpload(input: { name: string; size: number; mime: string }): Promise<{ uploadId: string; chunkSize: number }> {
  const url = `${API_BASE}/upload/init`
  return to(await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(input) }))
}

export async function putChunk(uploadId: string, index: number, chunk: Blob): Promise<void> {
  // Backend specifics are TBD; using query params for now
  const url = new URL('/upload/chunk', API_BASE)
  url.searchParams.set('uploadId', uploadId)
  url.searchParams.set('index', String(index))
  const res = await fetch(url.toString(), { method: 'PUT', body: chunk })
  if (!res.ok) throw new Error(await res.text())
}

export async function completeUpload(uploadId: string): Promise<{ fileId: string }> {
  const url = `${API_BASE}/upload/complete`
  return to(await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ uploadId }) }))
}

