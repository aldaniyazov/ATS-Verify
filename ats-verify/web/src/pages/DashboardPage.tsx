import { useAuth } from '../hooks/useAuth';
import { Link } from 'react-router-dom';
import {
    Package,
    Upload,
    TrendingUp,
    ShieldAlert,
    Smartphone,
    Search,
    ArrowUpRight,
} from 'lucide-react';
import type { UserRole } from '../types';

interface StatCard {
    label: string;
    value: string;
    change?: string;
    changeType?: 'up' | 'down' | 'neutral';
    icon: React.ReactNode;
    roles: UserRole[];
}

const STATS: StatCard[] = [
    { label: 'Всего посылок', value: '12,847', change: '+5.2%', changeType: 'up', icon: <Package size={18} />, roles: ['admin', 'customs_staff'] },
    { label: 'Загружено сегодня', value: '342', change: '+12', changeType: 'up', icon: <Upload size={18} />, roles: ['admin', 'marketplace_staff'] },
    { label: 'Использованных', value: '8,231', icon: <TrendingUp size={18} />, roles: ['admin', 'customs_staff'] },
    { label: 'Красный риск', value: '23', icon: <ShieldAlert size={18} />, roles: ['admin', 'customs_staff'] },
    { label: 'IMEI проверок', value: '156', icon: <Smartphone size={18} />, roles: ['admin', 'customs_staff', 'paid_user'] },
    { label: 'Поисков трека', value: '89', icon: <Search size={18} />, roles: ['admin', 'ats_staff'] },
];

interface QuickAction {
    label: string;
    to: string;
    icon: React.ReactNode;
    roles: UserRole[];
}

const ACTIONS: QuickAction[] = [
    { label: 'Загрузить CSV', to: '/upload', icon: <Upload size={16} />, roles: ['marketplace_staff'] },
    { label: 'Поиск трек-номера', to: '/track', icon: <Search size={16} />, roles: ['ats_staff', 'admin'] },
    { label: 'Проверка IMEI', to: '/imei', icon: <Smartphone size={16} />, roles: ['customs_staff', 'paid_user'] },
    { label: 'Управление рисками', to: '/risks', icon: <ShieldAlert size={16} />, roles: ['admin', 'customs_staff'] },
];

function getGreeting(): string {
    const h = new Date().getHours();
    if (h < 6) return 'Доброй ночи';
    if (h < 12) return 'Доброе утро';
    if (h < 18) return 'Добрый день';
    return 'Добрый вечер';
}

export default function DashboardPage() {
    const { user } = useAuth();
    if (!user) return null;

    const userStats = STATS.filter((s) => s.roles.includes(user.role));
    const userActions = ACTIONS.filter((a) => a.roles.includes(user.role));

    return (
        <div>
            {/* Header */}
            <div className="page-header">
                <h1 className="page-title">{getGreeting()}, {user.username}</h1>
                <p className="page-subtitle">Обзор системы ATS-Verify</p>
            </div>

            {/* Stat Cards */}
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-8">
                {userStats.map((stat) => (
                    <div key={stat.label} className="stat-card">
                        <div className="flex items-center justify-between">
                            <div className="w-9 h-9 rounded-lg bg-primary-50 flex items-center justify-center text-primary">
                                {stat.icon}
                            </div>
                            {stat.change && (
                                <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${stat.changeType === 'up' ? 'bg-success-light text-green-700' : 'bg-danger-light text-red-700'
                                    }`}>
                                    {stat.change}
                                </span>
                            )}
                        </div>
                        <div>
                            <p className="stat-value">{stat.value}</p>
                            <p className="stat-label">{stat.label}</p>
                        </div>
                    </div>
                ))}
            </div>

            {/* Quick Actions */}
            {userActions.length > 0 && (
                <div className="card p-6">
                    <h2 className="text-base font-semibold text-text-primary mb-4">Быстрые действия</h2>
                    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
                        {userActions.map((action) => (
                            <Link
                                key={action.to}
                                to={action.to}
                                className="flex items-center justify-between px-4 py-3 rounded-lg border border-border hover:border-primary/30 hover:bg-primary-50 transition-all group"
                            >
                                <div className="flex items-center gap-3">
                                    <span className="text-text-muted group-hover:text-primary transition-colors">{action.icon}</span>
                                    <span className="text-sm font-medium text-text-secondary group-hover:text-text-primary transition-colors">{action.label}</span>
                                </div>
                                <ArrowUpRight size={14} className="text-text-muted group-hover:text-primary transition-colors" />
                            </Link>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
