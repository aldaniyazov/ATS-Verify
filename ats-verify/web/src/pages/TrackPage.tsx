import { useState } from 'react';
import { Search, Download, Play } from 'lucide-react';
import api from '../lib/api';

interface TrackResult {
    track_number: string;
    found: boolean;
    parcel?: {
        is_used: boolean;
        marketplace: string;
        product_name: string;
        updated_at: string;
    };
}

export default function TrackPage() {
    const [input, setInput] = useState('');
    const [loading, setLoading] = useState(false);
    const [results, setResults] = useState<TrackResult[]>([]);
    const [error, setError] = useState('');

    const handleSearch = async () => {
        if (!input.trim()) return;
        setLoading(true);
        setError('');

        const tracks = input.split('\n').map(t => t.trim()).filter(Boolean);
        if (tracks.length === 0) { setLoading(false); return; }

        try {
            const { data } = await api.post('/track/bulk', { tracks });
            setResults(data.results || []);
        } catch {
            setError('Ошибка проверки. Убедитесь что бэкенд запущен.');
            setResults([]);
        } finally {
            setLoading(false);
        }
    };

    const foundCount = results.filter(r => r.found).length;
    const notFoundCount = results.filter(r => !r.found).length;

    const handleExport = () => {
        if (results.length === 0) return;
        const headers = ['Трек-номер', 'Найден', 'Использован', 'Маркетплейс', 'Товар'];
        const rows = results.map(r => [
            r.track_number,
            r.found ? 'Да' : 'Нет',
            r.found ? (r.parcel?.is_used ? 'Да' : 'Нет') : '--',
            r.parcel?.marketplace || '--',
            r.parcel?.product_name || '--',
        ]);
        const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
        const blob = new Blob([csv], { type: 'text/csv' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a'); a.href = url; a.download = 'track_results.csv'; a.click();
    };

    return (
        <div>
            <div className="flex items-center justify-between mb-6">
                <div>
                    <div className="breadcrumb mb-1">
                        <span>🏠</span>
                        <span>/</span>
                        <span className="text-text-primary font-medium">Поиск трека</span>
                    </div>
                    <h1 className="page-title">Массовая проверка трек-номеров</h1>
                    <p className="page-subtitle">Проверьте трек-номера по базе данных ATS в реальном времени</p>
                </div>
                <div className="badge-success">
                    API Status: Online
                </div>
            </div>

            {/* Input Area */}
            <div className="card p-6 mb-6">
                <div className="flex items-center justify-between mb-3">
                    <h3 className="text-sm font-semibold text-text-primary flex items-center gap-2">
                        📋 Входные данные
                    </h3>
                    <span className="text-xs text-text-muted">{input.split('\n').filter(l => l.trim()).length} / 500 lines</span>
                </div>

                <textarea
                    value={input}
                    onChange={(e) => setInput(e.target.value)}
                    className="input font-mono text-sm h-40 resize-none"
                    placeholder="Вставьте трек-номера, по одному на строку (например, TRK-9283402)..."
                />

                {error && (
                    <div className="mt-2 bg-danger-light border border-danger/20 text-danger text-sm px-4 py-3 rounded-lg">
                        {error}
                    </div>
                )}

                <div className="flex items-center justify-between mt-3">
                    <button onClick={() => { setInput(''); setResults([]); }} className="text-sm text-primary hover:text-primary-dark cursor-pointer underline">
                        Очистить
                    </button>
                    <button
                        onClick={handleSearch}
                        disabled={!input.trim() || loading}
                        className="btn-primary disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                        <Play size={16} />
                        {loading ? 'Проверка...' : 'Запустить проверку'}
                    </button>
                </div>
            </div>

            {/* Results */}
            {results.length > 0 && (
                <div className="card overflow-hidden">
                    {/* Table Header */}
                    <div className="px-6 py-4 border-b border-border flex items-center justify-between">
                        <div className="flex items-center gap-3">
                            <h3 className="text-base font-semibold text-text-primary">Результаты</h3>
                            <span className="badge-success">{foundCount} Найдено</span>
                            {notFoundCount > 0 && <span className="badge-danger">{notFoundCount} Не найдено</span>}
                        </div>
                        <div className="flex items-center gap-2">
                            <div className="relative">
                                <Search size={14} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
                                <input className="input pl-9 py-1.5 text-xs w-48" placeholder="Фильтр..." />
                            </div>
                            <button onClick={handleExport} className="btn-secondary text-xs py-1.5">
                                <Download size={14} />
                                Export CSV
                            </button>
                        </div>
                    </div>

                    {/* Table */}
                    <table className="data-table">
                        <thead>
                            <tr>
                                <th>Трек-номер</th>
                                <th>Статус обнаружения</th>
                                <th>Статус использования</th>
                                <th>Маркетплейс</th>
                                <th>Товар</th>
                            </tr>
                        </thead>
                        <tbody>
                            {results.map((r, i) => (
                                <tr key={i} className={!r.found ? 'bg-danger-light/30' : ''}>
                                    <td className="font-mono font-medium text-text-primary">{r.track_number}</td>
                                    <td>
                                        {r.found ? (
                                            <span className="flex items-center gap-1.5 text-green-600">
                                                <span className="dot-success" /> Exists
                                            </span>
                                        ) : (
                                            <span className="badge-danger">● Not Found</span>
                                        )}
                                    </td>
                                    <td>
                                        {r.found ? (
                                            <span className={r.parcel?.is_used ? 'badge-success' : 'badge-info'}>
                                                ● {r.parcel?.is_used ? 'Used' : 'Unused'}
                                            </span>
                                        ) : (
                                            <span className="text-text-muted">--</span>
                                        )}
                                    </td>
                                    <td>{r.parcel?.marketplace || '--'}</td>
                                    <td>{r.parcel?.product_name || '--'}</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>

                    {/* Footer */}
                    <div className="px-4 py-3 border-t border-border flex items-center justify-between">
                        <span className="text-sm text-text-muted">Showing 1 to {results.length} of {results.length} results</span>
                        <div className="pagination">
                            <button>‹</button>
                            <button>›</button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
}
