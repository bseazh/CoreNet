"use client"
import React, { useState } from 'react'
import { useRouter } from 'next/navigation'

export default function LoginPage() {
  const [token, setToken] = useState('')
  const router = useRouter()
  return (
    <div className="min-h-screen flex items-center justify-center p-6">
      <div className="w-full max-w-sm bg-white rounded border p-6 space-y-4">
        <h1 className="text-lg font-semibold">登录</h1>
        <div className="space-y-2">
          <label className="text-sm text-gray-600">API Token（用于 Authorization Bearer）</label>
          <input className="w-full border rounded px-3 py-2 text-sm" value={token} onChange={(e)=>setToken(e.target.value)} placeholder="粘贴你的后端 Token" />
        </div>
        <button className="px-3 py-2 bg-blue-600 text-white rounded" onClick={()=>{ localStorage.setItem('token', token); router.push('/f/root') }}>保存并进入</button>
      </div>
    </div>
  )
}

