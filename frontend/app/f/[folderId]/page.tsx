"use client"
import React, { useCallback, useMemo, useState } from 'react'
import useSWR from 'swr'
import Button from '@/components/ui/Button'
import Input from '@/components/ui/Input'
import Checkbox from '@/components/ui/Checkbox'
import UploadDropzone from '@/components/UploadDropzone'
import FileCard from '@/components/FileCard'
import FolderCard from '@/components/FolderCard'
import PreviewModal from '@/components/PreviewModal'
import Modal from '@/components/ui/Modal'
import { createFolder, deleteFile, DriveNode, listChildren, moveNodes, renameNode } from '@/lib/api'
import Link from 'next/link'
import { useParams } from 'next/navigation'

export default function FolderPage() {
  const params = useParams() as { folderId: string }
  const folderId = params.folderId || 'root'
  const { data, mutate, isLoading } = useSWR(['children', folderId], () => listChildren(folderId))
  const [view, setView] = useState<'grid'|'list'>('grid')
  const [selected, setSelected] = useState<Record<string, boolean>>({})
  const [previewId, setPreviewId] = useState<string | null>(null)
  const [renaming, setRenaming] = useState<{ id: string; type: 'file'|'folder'; name: string }|null>(null)
  const [newFolderOpen, setNewFolderOpen] = useState(false)
  const [newFolderName, setNewFolderName] = useState('新建文件夹')

  const items = data || []
  const folders = items.filter(x => x.type === 'folder')
  const files = items.filter(x => x.type === 'file')

  const anySelected = useMemo(() => Object.values(selected).some(Boolean), [selected])
  const selectedIds = useMemo(() => Object.entries(selected).filter(([,v]) => v).map(([k])=>k), [selected])

  const toggleSel = useCallback((id: string) => setSelected(s => ({ ...s, [id]: !s[id] })), [])

  async function onDropToFolder(targetId: string, draggedIds: string[]) {
    await moveNodes(draggedIds, targetId)
    mutate()
  }

  async function bulkDelete() {
    await Promise.all(selectedIds.map(id => deleteFile(id).catch(()=>{})))
    setSelected({})
    mutate()
  }

  function onDragStartFile(e: React.DragEvent, id: string) {
    const ids = selected[id] ? selectedIds : [id]
    e.dataTransfer.setData('text/plain', ids.join(','))
  }

  async function onCreateFolder() {
    await createFolder(folderId, newFolderName || '新建文件夹')
    setNewFolderOpen(false)
    setNewFolderName('新建文件夹')
    mutate()
  }

  async function onRenameConfirm() {
    if (!renaming) return
    await renameNode(renaming.id, renaming.type, renaming.name)
    setRenaming(null)
    mutate()
  }

  return (
    <div className="max-w-7xl mx-auto p-6 space-y-6">
      <header className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link href="/f/root" className="text-xl font-semibold">CoreNet Drive</Link>
          <span className="text-gray-400">/</span>
          <span className="text-gray-700">{folderId}</span>
        </div>
        <div className="flex items-center gap-2">
          <Button onClick={()=>setNewFolderOpen(true)} variant="secondary">新建文件夹</Button>
          <Button onClick={()=>setView(view==='grid'?'list':'grid')} variant="ghost">{view==='grid'?'列表视图':'网格视图'}</Button>
          <Link href="/login" className="text-sm text-gray-600 hover:underline">登录</Link>
        </div>
      </header>

      <UploadDropzone onUploaded={() => mutate()} />

      <div className="flex items-center gap-2">
        <Input placeholder="搜索（当前原型使用全局 /search，可后续接入 folder 过滤）" />
        {anySelected && (
          <>
            <Button variant="primary" onClick={()=>{/* placeholder for bulk move via modal */}}>移动</Button>
            <Button variant="ghost" onClick={bulkDelete}>删除</Button>
          </>
        )}
      </div>

      {isLoading && <div className="text-sm text-gray-500">加载中...</div>}

      {/* Folders */}
      {folders.length > 0 && (
        <section>
          <div className="mb-2 text-sm text-gray-600">文件夹</div>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {folders.map(f => (
              <FolderCard key={f.id} folder={{ id: f.id, name: f.name }} onOpen={(id)=>window.location.href=`/f/${id}`} onDropFiles={onDropToFolder} />
            ))}
          </div>
        </section>
      )}

      {/* Files */}
      <section>
        <div className="mb-2 text-sm text-gray-600">文件</div>
        {files.length === 0 ? (
          <div className="text-gray-500 text-sm">暂无文件</div>
        ) : view === 'grid' ? (
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
            {files.map(it => (
              <div key={it.id} draggable onDragStart={(e)=>onDragStartFile(e, it.id)}>
                <div className="flex items-center gap-2 mb-1">
                  <Checkbox checked={!!selected[it.id]} onChange={()=>toggleSel(it.id)} />
                  <button className="text-xs text-blue-700" onClick={()=>setRenaming({ id: it.id, type:'file', name: it.name })}>重命名</button>
                </div>
                <FileCard item={{ fileId: it.id, name: it.name }} onOpen={setPreviewId} onDelete={async(id)=>{ await deleteFile(id); mutate() }} />
              </div>
            ))}
          </div>
        ) : (
          <div className="rounded border bg-white overflow-hidden">
            <div className="grid grid-cols-12 text-xs px-3 py-2 bg-gray-50 border-b">
              <div className="col-span-6">名称</div>
              <div className="col-span-3">大小</div>
              <div className="col-span-3 text-right">操作</div>
            </div>
            {files.map(it => (
              <div key={it.id} className="grid grid-cols-12 items-center px-3 py-2 border-b last:border-0" draggable onDragStart={(e)=>onDragStartFile(e, it.id)}>
                <div className="col-span-6 flex items-center gap-2">
                  <Checkbox checked={!!selected[it.id]} onChange={()=>toggleSel(it.id)} />
                  <span className="truncate">{it.name}</span>
                </div>
                <div className="col-span-3 text-sm text-gray-500">{it.size ? (Math.round(it.size/1024)+' KB') : '-'}</div>
                <div className="col-span-3 text-right space-x-2">
                  <button className="text-blue-700 text-sm" onClick={()=>setPreviewId(it.id)}>预览</button>
                  <button className="text-gray-700 text-sm" onClick={()=>setRenaming({ id: it.id, type:'file', name: it.name })}>重命名</button>
                  <button className="text-red-600 text-sm" onClick={async()=>{ await deleteFile(it.id); mutate() }}>删除</button>
                </div>
              </div>
            ))}
          </div>
        )}
      </section>

      <PreviewModal fileId={previewId} onClose={()=>setPreviewId(null)} />

      {/* New Folder */}
      <Modal open={newFolderOpen} onClose={()=>setNewFolderOpen(false)} title="新建文件夹"
        footer={<div className="flex justify-end gap-2"><Button variant="ghost" onClick={()=>setNewFolderOpen(false)}>取消</Button><Button variant="primary" onClick={onCreateFolder}>创建</Button></div>}>
        <Input value={newFolderName} onChange={(e)=>setNewFolderName(e.target.value)} />
      </Modal>

      {/* Rename */}
      <Modal open={!!renaming} onClose={()=>setRenaming(null)} title="重命名"
        footer={<div className="flex justify-end gap-2"><Button variant="ghost" onClick={()=>setRenaming(null)}>取消</Button><Button variant="primary" onClick={onRenameConfirm}>保存</Button></div>}>
        {renaming && (
          <Input value={renaming.name} onChange={(e)=>setRenaming({...renaming, name: e.target.value})} />
        )}
      </Modal>
    </div>
  )
}

