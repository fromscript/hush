// components/chat/loading-skeleton.tsx
'use client'

export function LoadingSkeleton() {
  return (
    <div className="animate-pulse space-y-4 bg-white p-6 rounded-xl shadow-sm">
      <div className="h-4 bg-gray-200 rounded w-3/4"/>
      <div className="h-4 bg-gray-200 rounded"/>
      <div className="h-4 bg-gray-200 rounded w-5/6"/>
      <div className="h-10 bg-gray-200 rounded-md mt-4"/>
    </div>
  )
}
