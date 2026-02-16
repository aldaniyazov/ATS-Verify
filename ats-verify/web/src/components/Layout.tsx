import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import Sidebar from './Sidebar';

export default function Layout() {
    const { user } = useAuth();

    if (!user) return <Navigate to="/login" replace />;

    return (
        <div className="flex min-h-screen bg-bg">
            <Sidebar />
            <main className="flex-1 p-8 overflow-y-auto">
                <Outlet />
            </main>
        </div>
    );
}
