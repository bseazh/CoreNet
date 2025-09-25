import { put, del } from "@vercel/blob"

export async function uploadFile(file: File, path?: string): Promise<{ url: string; pathname: string }> {
  const filename = path ? `${path}/${file.name}` : file.name

  const blob = await put(filename, file, {
    access: "public",
  })

  return {
    url: blob.url,
    pathname: blob.pathname,
  }
}

export async function deleteFile(url: string): Promise<void> {
  await del(url)
}

export function formatFileSize(bytes: number): string {
  if (bytes === 0) return "0 Bytes"

  const k = 1024
  const sizes = ["Bytes", "KB", "MB", "GB", "TB"]
  const i = Math.floor(Math.log(bytes) / Math.log(k))

  return Number.parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + " " + sizes[i]
}

export function getFileIcon(mimeType: string): string {
  if (mimeType.startsWith("image/")) return "🖼️"
  if (mimeType.startsWith("video/")) return "🎥"
  if (mimeType.startsWith("audio/")) return "🎵"
  if (mimeType.includes("pdf")) return "📄"
  if (mimeType.includes("document") || mimeType.includes("word")) return "📝"
  if (mimeType.includes("spreadsheet") || mimeType.includes("excel")) return "📊"
  if (mimeType.includes("presentation") || mimeType.includes("powerpoint")) return "📽️"
  if (mimeType.includes("zip") || mimeType.includes("rar") || mimeType.includes("archive")) return "📦"
  return "📄"
}
