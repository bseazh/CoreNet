"use client"
import React, { useMemo, useState } from 'react'
import { Search, Upload, FolderPlus, LayoutGrid, List, Trash } from 'lucide-react'
import useSWR from 'swr'
import { deleteFile, getFile, searchFiles, createOCR, getJob } from '@/lib/api'
import UploadDropzone from '@/components/UploadDropzone'
import FileCard from '@/components/FileCard'
import PreviewModal from '@/components/PreviewModal'
import Button from '@/components/ui/button'
import Input from '@/components/ui/input'
import Checkbox from '@/components/ui/checkbox'
import Modal from '@/components/ui/modal'
import { Toaster, toast } from 'sonner'

export default function DrivePage() {
  const [q, setQ] = useState('')
  const { data: items, mutate, isLoading } = useSWR(['search', q], () => searchFiles(q))
  const [previewId, setPreviewId] = useState<string | null>(null)
  const [view, setView] = useState<'grid'|'list'>('grid')
  const [selected, setSelected] = useState<Record<string, boolean>>({})
  const [newFolderOpen, setNewFolderOpen] = useState(false)
  const [newFolderName, setNewFolderName] = useState('新建文件夹')

  const files = useMemo(() => items || [], [items])
  const anySelected = useMemo(() => Object.values(selected).some(Boolean), [selected])
  const selectedIds = useMemo(() => Object.entries(selected).filter(([,v]) => v).map(([k])=>k), [selected])

  async function handleUploaded(fileId: string) {
    try {
      const meta = await getFile(fileId)
      await mutate(async (prev) => {
        const next = [...(prev || [])]
        next.unshift({ fileId: meta.fileId, name: meta.name })
        return next
      }, { revalidate: false })
      toast.success('上传完成')
    } catch {
      mutate(); toast.success('上传完成，列表已刷新')
    }
  }

  async function handleDelete(fileId: string) {
    await deleteFile(fileId)
    mutate(files?.filter(i => i.fileId !== fileId), { revalidate: false })
    setSelected(s => { const n={...s}; delete n[fileId]; return n })
    toast.success('已删除')
  }

  async function handleOCR(fileId: string) {
    const { jobId } = await createOCR(fileId)
    toast('已提交 OCR 任务', { description: `Job ${jobId}` })
    let tries = 0
    const timer = setInterval(async () => {
      tries++
      const job = await getJob(jobId)
      if (job.status === 'done') { toast.success('OCR 完成'); clearInterval(timer); mutate() }
      if (job.status === 'failed' || tries > 30) { toast.error('OCR 失败或超时'); clearInterval(timer) }
    }, 2000)
  }

  function toggle(id: string) { setSelected(s => ({ ...s, [id]: !s[id] })) }

  return (
    <div className="min-h-screen grid grid-cols-12">
      <Toaster richColors />
      {/* Sidebar */}
      <aside className="col-span-2 border-r bg-gray-50 p-4 space-y-2">
        <div className="font-semibold mb-2 text-gray-700">我的云端硬盘</div>
        <nav className="space-y-1 text-sm">
          <a className="block px-3 py-2 rounded hover:bg-white">我的文件</a>
          <a className="block px-3 py-2 rounded hover:bg-white">最近</a>
          <a className="block px-3 py-2 rounded hover:bg-white">回收站</a>
        </nav>
      </aside>

      {/* Main */}
      <main className="col-span-10 p-6 space-y-4">
        {/* Top bar */}
        <div className="flex items-center gap-3">
          <div className="flex-1 relative">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-4 w-4 text-gray-400" />
            <Input value={q} onChange={(e)=>setQ(e.target.value)} placeholder="搜索文件名或内容（需 OCR 后）" className="pl-9" />
          </div>
          <Button variant="primary" onClick={()=>setNewFolderOpen(true)} title="新建文件夹"><FolderPlus className="h-4 w-4 mr-2"/>新建</Button>
          <Button variant="ghost" onClick={()=>setView(view==='grid'?'list':'grid')} title="切换视图">
            {view==='grid'? <List className="h-4 w-4"/> : <LayoutGrid className="h-4 w-4"/>}
          </Button>
        </div>

        {/* Uploader */}
        <div className="border rounded-lg p-4">
          <div className="flex items-center gap-2 mb-3 text-sm text-gray-600"><Upload className="h-4 w-4"/> 上传文件</div>
          <UploadDropzone onUploaded={handleUploaded} />
        </div>

        {/* Actions for selection */}
        {anySelected && (
          <div className="flex items-center gap-2 text-sm">
            <span className="text-gray-600">已选 {selectedIds.length} 项</span>
            <Button variant="danger" onClick={()=>{ selectedIds.forEach(id=>handleDelete(id)) }}><Trash className="h-4 w-4 mr-1"/> 删除</Button>
          </div>
        )}

        {/* Files */}
        <section>
          <div className="flex items-center justify-between mb-2">
            <h2 className="text-lg font-medium">文件</h2>
            {isLoading && <span className="text-sm text-gray-500">加载中...</span>}
          </div>
          {(!files || files.length === 0) ? (
            <div className="text-gray-500 text-sm">暂无文件。上传后将显示在这里。</div>
          ) : view==='grid' ? (
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-4">
              {files.map(it => (
                <div key={it.fileId}>
                  <div className="flex items-center gap-2 mb-1">
                    <Checkbox checked={!!selected[it.fileId]} onChange={()=>toggle(it.fileId)} />
                  </div>
                  <FileCard item={{ fileId: it.fileId, name: it.name }} onOpen={setPreviewId} onDelete={handleDelete} onOCR={handleOCR} />
                </div>
              ))}
            </div>
          ) : (
            <div className="rounded border bg-white overflow-hidden">
              <div className="grid grid-cols-12 text-xs px-3 py-2 bg-gray-50 border-b">
                <div className="col-span-7">名称</div>
                <div className="col-span-2">类型</div>
                <div className="col-span-3 text-right">操作</div>
              </div>
              {files.map(it => (
                <div key={it.fileId} className="grid grid-cols-12 items-center px-3 py-2 border-b last:border-0">
                  <div className="col-span-7 flex items-center gap-2">
                    <Checkbox checked={!!selected[it.fileId]} onChange={()=>toggle(it.fileId)} />
                    <span className="truncate">{it.name}</span>
                  </div>
                  <div className="col-span-2 text-sm text-gray-500">文件</div>
                  <div className="col-span-3 text-right space-x-2">
                    <button className="text-blue-700 text-sm" onClick={()=>setPreviewId(it.fileId)}>预览</button>
                    <button className="text-amber-700 text-sm" onClick={()=>handleOCR(it.fileId)}>OCR</button>
                    <button className="text-red-600 text-sm" onClick={()=>handleDelete(it.fileId)}>删除</button>
                  </div>
                </div>
              ))}
            </div>
          )}
        </section>

        <PreviewModal fileId={previewId} onClose={()=>setPreviewId(null)} />

        {/* New Folder (待后端支持) */}
        <Modal open={newFolderOpen} onClose={()=>setNewFolderOpen(false)} title="新建文件夹"
          footer={<div className="flex justify-end gap-2"><Button variant="ghost" onClick={()=>setNewFolderOpen(false)}>取消</Button><Button variant="primary" onClick={()=>{ setNewFolderOpen(false); toast.info('新建文件夹待后端接口支持'); }}>创建</Button></div>}>
          <input value={newFolderName} onChange={(e)=>setNewFolderName(e.target.value)} className="w-full border rounded px-3 py-2 text-sm" />
        </Modal>
      </main>
    </div>
  )
}

