"use client"
import React from 'react'

export default function Modal({ open, onClose, title, children, footer }: { open: boolean; onClose: () => void; title?: string; children: React.ReactNode; footer?: React.ReactNode }) {
  if (!open) return null
  return (
    <div className="fixed inset-0 z-50 bg-black/50 flex items-center justify-center p-4" onClick={onClose}>
      <div className="bg-white w-full max-w-md rounded shadow" onClick={(e)=>e.stopPropagation()}>
        {title && <div className="px-4 py-3 border-b text-sm font-medium">{title}</div>}
        <div className="p-4">{children}</div>
        {footer && <div className="px-4 py-3 border-t bg-gray-50">{footer}</div>}
      </div>
    </div>
  )
}

