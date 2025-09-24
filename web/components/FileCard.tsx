"use client"
import React from 'react'

export type FileCardData = {
  fileId: string
  name: string
  mime?: string
}

export default function FileCard({ item, onOpen, onDelete }: { item: FileCardData; onOpen: (id: string) => void; onDelete: (id: string) => void }) {
  const ext = (item.name.split('.').pop() || '').toLowerCase()
  const isImage = (item.mime || '').startsWith('image') || ['jpg','jpeg','png','gif','webp','bmp','svg'].includes(ext)
  const isVideo = (item.mime || '').startsWith('video') || ['mp4','webm','mov','mkv','avi'].includes(ext)
  return (
    <div className="group rounded-lg border p-3 bg-white hover:shadow cursor-pointer relative">
      <div className="h-28 flex items-center justify-center rounded bg-gray-50 text-4xl">
        {isImage ? '🖼️' : isVideo ? '🎬' : '📄'}
      </div>
      <div className="mt-2 text-sm font-medium truncate" title={item.name}>{item.name}</div>
      <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition">
        <button onClick={(e)=>{e.stopPropagation(); onOpen(item.fileId)}} className="px-2 py-1 text-xs bg-blue-600 text-white rounded mr-2">预览</button>
        <button onClick={(e)=>{e.stopPropagation(); onDelete(item.fileId)}} className="px-2 py-1 text-xs bg-red-600 text-white rounded">删除</button>
      </div>
      <button onClick={()=>onOpen(item.fileId)} className="absolute inset-0" aria-label="open" />
    </div>
  )
}

