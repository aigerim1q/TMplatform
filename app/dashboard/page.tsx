import Header from '@/components/header';
import Sidebar from '@/components/sidebar';
import DashboardContent from '@/components/dashboard-content';

export default function Dashboard() {
  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header - centered relative to viewport */}
      <div className="flex w-screen justify-center border-b border-gray-200 bg-gray-50 py-4">
        <Header />
      </div>

      {/* Dashboard content area with sidebar */}
      <div className="flex">
        <Sidebar />
        <DashboardContent />
      </div>
    </div>
  );
}
