import React from 'react'

export default function Input(props: React.InputHTMLAttributes<HTMLInputElement>) {
  const base = 'w-full border rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500'
  return <input {...props} className={`${base} ${props.className || ''}`} />
}

