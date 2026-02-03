import Header from '@/components/header';
import LifecycleContent from '@/components/lifecycle-content';

export default function Lifecycle() {
  return (
    <div className="min-h-screen bg-background">
      {/* Header - centered */}
      <div className="flex justify-center pt-6">
        <Header />
      </div>

      {/* Main content area */}
      <LifecycleContent />
    </div>
  );
}
