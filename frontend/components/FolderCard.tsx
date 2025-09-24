"use client"
import React from 'react'

export type Folder = { id: string; name: string }

export default function FolderCard({ folder, onOpen, onDropFiles }: { folder: Folder; onOpen: (id: string) => void; onDropFiles: (folderId: string, draggedIds: string[]) => void }) {
  return (
    <div
      className="group rounded-lg border p-3 bg-white hover:shadow cursor-pointer relative"
      onDoubleClick={() => onOpen(folder.id)}
      onDragOver={(e) => { e.preventDefault() }}
      onDrop={(e) => { e.preventDefault(); const ids = (e.dataTransfer.getData('text/plain')||'').split(',').filter(Boolean); if (ids.length) onDropFiles(folder.id, ids) }}
    >
      <div className="h-28 flex items-center justify-center rounded bg-yellow-50 text-4xl">📁</div>
      <div className="mt-2 text-sm font-medium truncate" title={folder.name}>{folder.name}</div>
    </div>
  )
}

