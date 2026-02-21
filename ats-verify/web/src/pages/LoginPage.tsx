import { useState, type FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import api from '../lib/api';
import { Lock, User, CheckCircle2, Zap } from 'lucide-react';

export default function LoginPage() {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const { login } = useAuth();
    const navigate = useNavigate();

    const handleSubmit = async (e: FormEvent) => {
        e.preventDefault();
        setError('');
        setLoading(true);
        try {
            const { data } = await api.post('/auth/login', { username, password });
            login({ id: data.user.id, username: data.user.username, role: data.user.role, marketplace_prefix: data.user.marketplace_prefix, token: data.token });
            navigate('/');
        } catch (err: unknown) {
            const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error;
            setError(msg || 'Ошибка авторизации');
        } finally {
            setLoading(false);
        }
    };

    return (
        <div className="min-h-screen flex items-center justify-center bg-bg px-4">
            <div className="w-full max-w-md">
                {/* Logo */}
                <div className="text-center mb-8">
                    <div className="inline-flex items-center justify-center w-14 h-14 rounded-2xl bg-primary mb-4">
                        <CheckCircle2 size={28} className="text-white" />
                    </div>
                    <h1 className="text-3xl font-bold text-text-primary">ATS-Verify</h1>
                    <p className="text-text-muted mt-1 text-sm">Платформа верификации</p>
                </div>

                {/* Login Card */}
                <div className="card p-8">
                    <h2 className="text-lg font-semibold text-text-primary mb-6">Вход в систему</h2>

                    <form onSubmit={handleSubmit} className="space-y-4">
                        {error && (
                            <div className="bg-danger-light border border-danger/20 text-danger text-sm px-4 py-3 rounded-lg">
                                {error}
                            </div>
                        )}

                        <div>
                            <label className="text-sm font-medium text-text-primary block mb-1.5">Логин</label>
                            <div className="relative">
                                <User size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
                                <input
                                    type="text"
                                    value={username}
                                    onChange={(e) => setUsername(e.target.value)}
                                    className="input"
                                    style={{ paddingLeft: '2.5rem' }}
                                    placeholder="Введите логин"
                                    required
                                />
                            </div>
                        </div>

                        <div>
                            <label className="text-sm font-medium text-text-primary block mb-1.5">Пароль</label>
                            <div className="relative">
                                <Lock size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
                                <input
                                    type="password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    className="input"
                                    style={{ paddingLeft: '2.5rem' }}
                                    placeholder="Введите пароль"
                                    required
                                />
                            </div>
                        </div>

                        <button type="submit" disabled={loading} className="btn-primary w-full justify-center py-2.5">
                            {loading ? 'Вход...' : 'Войти в систему'}
                        </button>
                    </form>

                    {/* Demo Login Shortcuts */}
                    <div className="mt-6 pt-6 border-t border-border">
                        <div className="flex items-center gap-2 mb-4 text-xs font-semibold text-text-muted uppercase tracking-wider">
                            <Zap size={14} className="text-warning" />
                            Демо-вход
                        </div>
                        <div className="flex flex-wrap gap-2">
                            <button
                                type="button"
                                onClick={() => { setUsername('admin'); setPassword('admin'); }}
                                className="px-3 py-1.5 text-xs font-medium rounded-md bg-purple-50 text-purple-700 hover:bg-purple-100 transition-colors cursor-pointer"
                            >
                                Admin
                            </button>
                            <button
                                type="button"
                                onClick={() => { setUsername('ats_staff'); setPassword('ats_staff'); }}
                                className="px-3 py-1.5 text-xs font-medium rounded-md bg-blue-50 text-blue-700 hover:bg-blue-100 transition-colors cursor-pointer"
                            >
                                ATS
                            </button>
                            <button
                                type="button"
                                onClick={() => { setUsername('customs_staff'); setPassword('customs_staff'); }}
                                className="px-3 py-1.5 text-xs font-medium rounded-md bg-cyan-50 text-cyan-700 hover:bg-cyan-100 transition-colors cursor-pointer"
                            >
                                Таможня
                            </button>
                            <button
                                type="button"
                                onClick={() => { setUsername('marketplace_wb'); setPassword('marketplace_wb'); }}
                                className="px-3 py-1.5 text-xs font-medium rounded-md bg-orange-50 text-orange-700 hover:bg-orange-100 transition-colors cursor-pointer"
                            >
                                Маркетплейс
                            </button>
                        </div>
                    </div>
                </div>

                <p className="text-center text-text-muted text-xs mt-6">
                    © 2026 ATS-Verify. Все права защищены.
                </p>
            </div>
        </div>
    );
}
