import Header from '@/components/header';
import DashboardContent from '@/components/dashboard-content';

export default function Dashboard() {
  return (
    <div className="min-h-screen bg-background">
      {/* Header - centered */}
      <div className="flex justify-center pt-6">
        <Header />
      </div>

      {/* Dashboard content */}
      <DashboardContent />
    </div>
  );
}
