import Header from "@/components/header"
import LifecycleContent from "@/components/lifecycle-content"

export default function Lifecycle() {
  return (
    <div className="min-h-screen bg-white text-gray-900 dark:bg-slate-950 dark:text-white">
      {/* Header - centered relative to viewport */}
      <div className="flex w-screen justify-center border-b border-gray-200 bg-white py-4 dark:border-slate-800 dark:bg-slate-950">
        <Header />
      </div>

      {/* Main content area */}
      <LifecycleContent />
    </div>
  )
}
