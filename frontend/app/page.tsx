"use client"
import React, { useEffect, useMemo, useState } from 'react'
import useSWR from 'swr'
import { deleteFile, getFile, searchFiles } from '@/lib/api'
import UploadDropzone from '@/components/UploadDropzone'
import FileCard from '@/components/FileCard'
import PreviewModal from '@/components/PreviewModal'

const fetcher = (q: string) => searchFiles(q)

export default function Page() {
  const [q, setQ] = useState('')
  const { data: items, mutate, isLoading } = useSWR(['search', q], () => fetcher(q))
  const [previewId, setPreviewId] = useState<string | null>(null)

  const grid = useMemo(() => (items || []).map(i => i.fileId), [items])

  useEffect(() => {
    // when modal closes, refresh list to reflect any changes
    return () => {}
  }, [])

  async function handleUploaded(fileId: string) {
    // ensure the newly uploaded file appears; we can refetch or insert
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
    mutate(items?.filter(i => i.fileId !== fileId), { revalidate: false })
  }

  return (
    <div className="max-w-7xl mx-auto p-6 space-y-6">
      <header className="flex items-center justify-between">
        <h1 className="text-2xl font-semibold">CoreNet Drive</h1>
        <div className="w-80">
          <input
            value={q}
            onChange={(e)=> setQ(e.target.value)}
            placeholder="搜索文件名称..."
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
        {(!items || items.length === 0) ? (
          <div className="text-gray-500 text-sm">暂无文件。上传后将显示在这里。</div>
        ) : (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {items.map(it => (
              <FileCard key={it.fileId} item={{ fileId: it.fileId, name: it.name }} onOpen={setPreviewId} onDelete={handleDelete} />
            ))}
          </div>
        )}
      </section>

      <footer className="text-xs text-gray-500">API: {process.env.NEXT_PUBLIC_API_BASE_URL || 'http://localhost:8888'}</footer>

      <PreviewModal fileId={previewId} onClose={() => setPreviewId(null)} />
    </div>
  )
}

