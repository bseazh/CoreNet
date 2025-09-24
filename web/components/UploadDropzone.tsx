"use client"
import React, { useCallback, useRef, useState } from 'react'
import { completeUpload, initUpload, putChunk } from '@/lib/api'

type Props = {
  onUploaded?: (fileId: string) => void
}

export default function UploadDropzone({ onUploaded }: Props) {
  const [isDragging, setDragging] = useState(false)
  const inputRef = useRef<HTMLInputElement>(null)
  const [progress, setProgress] = useState<{ name: string; pct: number } | null>(null)

  const onDrop = useCallback(async (files: FileList | null) => {
    if (!files || files.length === 0) return
    const file = files[0]
    const init = await initUpload({ name: file.name, size: file.size, mime: file.type || 'application/octet-stream' })
    const { uploadId, chunkSize } = init
    const size = chunkSize || 5 * 1024 * 1024
    let index = 0
    for (let offset = 0; offset < file.size; offset += size) {
      const chunk = file.slice(offset, Math.min(offset + size, file.size))
      await putChunk(uploadId, index++, chunk)
      setProgress({ name: file.name, pct: Math.round((Math.min(offset + size, file.size) / file.size) * 100) })
    }
    const done = await completeUpload(uploadId)
    setProgress(null)
    onUploaded?.(done.fileId)
  }, [onUploaded])

  return (
    <div
      onDragOver={(e) => { e.preventDefault(); setDragging(true) }}
      onDragLeave={() => setDragging(false)}
      onDrop={(e) => { e.preventDefault(); setDragging(false); onDrop(e.dataTransfer.files) }}
      className={`border-2 border-dashed rounded-lg p-6 text-center cursor-pointer transition ${isDragging ? 'border-blue-500 bg-blue-50' : 'border-gray-300'}`}
      onClick={() => inputRef.current?.click()}
    >
      <input ref={inputRef} type="file" className="hidden" onChange={(e) => onDrop(e.target.files)} />
      <p className="text-gray-600">拖拽文件到此处，或点击选择上传</p>
      {progress && (
        <div className="mt-3 text-sm text-gray-700">
          {progress.name} · {progress.pct}%
          <div className="w-full bg-gray-200 rounded h-2 mt-1">
            <div className="bg-blue-600 h-2 rounded" style={{ width: `${progress.pct}%` }} />
          </div>
        </div>
      )}
    </div>
  )
}

