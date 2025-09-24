export type FileItem = {
  fileId: string
  name: string
  size?: number
  mime?: string
  snippet?: string
}

const API_BASE = process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8888'

function authHeaders(): HeadersInit {
  if (typeof window === 'undefined') return {}
  const token = localStorage.getItem('token')
  return token ? { Authorization: `Bearer ${token}` } : {}
}

function to<T>(res: Response): Promise<T> {
  if (!res.ok) {
    return res.text().then(t => Promise.reject(new Error(t || res.statusText)))
  }
  return res.json() as Promise<T>
}

export async function searchFiles(q: string): Promise<FileItem[]> {
  const url = new URL('/search', API_BASE)
  if (q) url.searchParams.set('q', q)
  const data = await to<{ items: { fileId: string; name: string; snippet?: string }[] }>(
    await fetch(url.toString(), { cache: 'no-store', headers: { ...authHeaders() } })
  )
  return data.items || []
}

export async function getFile(fileId: string): Promise<{ fileId: string; name: string; size: number; mime: string; sha1: string; version: number }> {
  const url = `${API_BASE}/files/${encodeURIComponent(fileId)}`
  return to(await fetch(url, { cache: 'no-store', headers: { ...authHeaders() } }))
}

export async function deleteFile(fileId: string): Promise<void> {
  const url = `${API_BASE}/files/${encodeURIComponent(fileId)}`
  const res = await fetch(url, { method: 'DELETE', headers: { ...authHeaders() } })
  if (!res.ok) throw new Error(await res.text())
}

export async function getPreviewUrl(fileId: string): Promise<string> {
  const url = `${API_BASE}/files/${encodeURIComponent(fileId)}/preview`
  const data = await to<{ url: string }>(await fetch(url, { cache: 'no-store', headers: { ...authHeaders() } }))
  return data.url
}

export async function initUpload(input: { name: string; size: number; mime: string }): Promise<{ uploadId: string; chunkSize: number }> {
  const url = `${API_BASE}/upload/init`
  return to(
    await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeaders() }, body: JSON.stringify(input) })
  )
}

export async function putChunk(uploadId: string, index: number, chunk: Blob): Promise<void> {
  // Backend specifics are TBD; using query params for now
  const url = new URL('/upload/chunk', API_BASE)
  url.searchParams.set('uploadId', uploadId)
  url.searchParams.set('index', String(index))
  const res = await fetch(url.toString(), { method: 'PUT', body: chunk, headers: { ...authHeaders() } })
  if (!res.ok) throw new Error(await res.text())
}

export async function completeUpload(uploadId: string): Promise<{ fileId: string }> {
  const url = `${API_BASE}/upload/complete`
  return to(
    await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeaders() }, body: JSON.stringify({ uploadId }) })
  )
}

// Folders & hierarchy (backend alignment required)
export type NodeType = 'file' | 'folder'
export type DriveNode = { id: string; type: NodeType; name: string; mime?: string; size?: number; updatedAt?: string }

export async function listChildren(folderId: string): Promise<DriveNode[]> {
  const url = `${API_BASE}/folders/${encodeURIComponent(folderId)}/children`
  const res = await fetch(url, { headers: { ...authHeaders() }, cache: 'no-store' })
  if (res.ok) return res.json()
  // Fallback: show files only via search
  const files = await searchFiles('')
  return files.map(f => ({ id: f.fileId, type: 'file' as const, name: f.name }))
}

export async function createFolder(parentId: string, name: string): Promise<{ id: string }> {
  const url = `${API_BASE}/folders`
  return to(await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeaders() }, body: JSON.stringify({ parentId, name }) }))
}

export async function renameNode(id: string, type: NodeType, name: string): Promise<void> {
  const url = type === 'folder' ? `${API_BASE}/folders/${encodeURIComponent(id)}` : `${API_BASE}/files/${encodeURIComponent(id)}`
  const res = await fetch(url, { method: 'PATCH', headers: { 'Content-Type': 'application/json', ...authHeaders() }, body: JSON.stringify({ name }) })
  if (!res.ok) throw new Error(await res.text())
}

export async function moveNodes(ids: string[], targetFolderId: string): Promise<void> {
  const url = `${API_BASE}/nodes/move`
  const res = await fetch(url, { method: 'POST', headers: { 'Content-Type': 'application/json', ...authHeaders() }, body: JSON.stringify({ ids, targetFolderId }) })
  if (!res.ok) throw new Error(await res.text())
}

