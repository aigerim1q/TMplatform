import Header from '@/components/header';
import LifecycleContent from '@/components/lifecycle-content';

export default function Lifecycle() {
  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header - centered relative to viewport */}
      <div className="flex w-screen justify-center border-b border-gray-200 bg-gray-50 py-4">
        <Header />
      </div>

      {/* Main content area */}
      <LifecycleContent />
    </div>
  );
}
