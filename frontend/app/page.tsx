import Link from 'next/link'

export default function Home() {
  return (
    <main className="min-h-screen flex items-center justify-center p-6">
      <div className="max-w-xl w-full border rounded-lg p-8 bg-white shadow-sm text-center space-y-4">
        <h1 className="text-2xl font-semibold">欢迎使用 CoreNet 前端（v0/shadcn 模板）</h1>
        <p className="text-gray-600">你可以直接使用 v0 生成的界面，或进入网盘功能。</p>
        <div className="flex items-center justify-center gap-3">
          <Link href="/f/root" className="px-4 py-2 rounded bg-blue-600 text-white">打开网盘</Link>
          <Link href="/login" className="px-4 py-2 rounded border">登录</Link>
        </div>
      </div>
    </main>
  )
}
