"use client"
import React, { useMemo, useState } from 'react'
import useSWR from 'swr'
import { deleteFile, getFile, searchFiles, createOCR, getJob } from '@/lib/api'
import UploadDropzone from '@/components/UploadDropzone'
import FileCard from '@/components/FileCard'
import PreviewModal from '@/components/PreviewModal'
import Link from 'next/link'
import { useParams } from 'next/navigation'

export default function FolderPage() {
  const params = useParams() as { folderId: string }
  const folderId = params.folderId || 'root'
  const [q, setQ] = useState('')
  const { data: items, mutate, isLoading } = useSWR(['search', q], () => searchFiles(q))
  const [previewId, setPreviewId] = useState<string | null>(null)

  const files = useMemo(() => items || [], [items])

  async function handleUploaded(fileId: string) {
    try {
      const meta = await getFile(fileId)
      await mutate(async (prev) => {
        const next = [...(prev || [])]
        next.unshift({ fileId: meta.fileId, name: meta.name })
        return next
      }, { revalidate: false })
    } catch {
      mutate()
    }
  }

  async function handleDelete(fileId: string) {
    await deleteFile(fileId)
    mutate(files?.filter(i => i.fileId !== fileId), { revalidate: false })
  }

  async function handleOCR(fileId: string) {
    const { jobId } = await createOCR(fileId)
    // naive poll
    let tries = 0
    const timer = setInterval(async () => {
      tries++
      const job = await getJob(jobId)
      if (job.status === 'done' || job.status === 'failed' || tries > 30) {
        clearInterval(timer)
        // refresh results if backend search indexes OCR
        mutate()
      }
    }, 2000)
  }

  return (
    <div className="max-w-7xl mx-auto p-6 space-y-6">
      <header className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link href="/f/root" className="text-xl font-semibold">CoreNet Drive</Link>
          <span className="text-gray-400">/</span>
          <span className="text-gray-700">{folderId}</span>
        </div>
        <div className="w-80">
          <input
            value={q}
            onChange={(e)=> setQ(e.target.value)}
            placeholder="搜索文件名或内容（需OCR后）"
            className="w-full border rounded px-3 py-2 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        </div>
      </header>

      <UploadDropzone onUploaded={handleUploaded} />

      <section>
        <div className="flex items-center justify-between mb-2">
          <h2 className="text-lg font-medium">文件</h2>
          {isLoading && <span className="text-sm text-gray-500">加载中...</span>}
        </div>
        {(!files || files.length === 0) ? (
          <div className="text-gray-500 text-sm">暂无文件。上传后将显示在这里。</div>
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {files.map(it => (
              <FileCard key={it.fileId} item={{ fileId: it.fileId, name: it.name }} onOpen={setPreviewId} onDelete={handleDelete} onOCR={handleOCR} />
            ))}
          </div>
        )}
      </section>

      <footer className="text-xs text-gray-500">API: {process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8888'} · 在 /login 可设置 Token</footer>

      <PreviewModal fileId={previewId} onClose={() => setPreviewId(null)} />
    </div>
  )
}

