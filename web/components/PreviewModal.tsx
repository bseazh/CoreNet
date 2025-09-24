"use client"
import React, { useEffect, useState } from 'react'
import { getFile, getPreviewUrl } from '@/lib/api'

export default function PreviewModal({ fileId, onClose }: { fileId: string | null; onClose: () => void }) {
  const [url, setUrl] = useState<string>('')
  const [meta, setMeta] = useState<{ name: string; mime?: string } | null>(null)

  useEffect(() => {
    let mounted = true
    async function load() {
      if (!fileId) return
      const m = await getFile(fileId)
      const u = await getPreviewUrl(fileId)
      if (!mounted) return
      setMeta({ name: m.name, mime: m.mime })
      setUrl(u)
    }
    load()
    return () => { mounted = false }
  }, [fileId])

  if (!fileId) return null
  const isImage = (meta?.mime || '').startsWith('image')
  const isVideo = (meta?.mime || '').startsWith('video')

  return (
    <div className="fixed inset-0 z-50 bg-black/70 flex items-center justify-center p-4" onClick={onClose}>
      <div className="bg-white max-w-5xl w-full rounded shadow-lg overflow-hidden" onClick={(e)=>e.stopPropagation()}>
        <div className="px-4 py-2 border-b flex items-center justify-between">
          <div className="font-medium truncate">{meta?.name}</div>
          <button onClick={onClose} className="text-gray-600 hover:text-black">关闭</button>
        </div>
        <div className="p-4 max-h-[75vh] overflow-auto flex items-center justify-center bg-gray-50">
          {isImage && url && (<img src={url} alt={meta?.name} className="max-h-[70vh] object-contain" />)}
          {isVideo && url && (<video src={url} controls className="max-h-[70vh]" />)}
          {!isImage && !isVideo && url && (
            <iframe src={url} className="w-full h-[70vh] bg-white" />
          )}
        </div>
      </div>
    </div>
  )
}

