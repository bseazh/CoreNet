import React from 'react'

export default function Checkbox({ className='', ...rest }: React.InputHTMLAttributes<HTMLInputElement>) {
  return (
    <input type="checkbox" className={`h-4 w-4 rounded border-gray-300 text-blue-600 focus:ring-blue-500 ${className}`} {...rest} />
  )}

