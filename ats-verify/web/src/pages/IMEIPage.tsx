import { useState, useRef } from 'react';
import { Smartphone, Upload, Download, CheckCircle, XCircle } from 'lucide-react';
import api from '../lib/api';

interface IMEIResult {
    row: number;
    imei_code: string;
    status: 'found' | 'missing';
    matched_number?: string;
}

interface IMEIReport {
    total: number;
    found: number;
    missing: number;
    match_rate: number;
    results: IMEIResult[];
}

export default function IMEIPage() {
    const [report, setReport] = useState<IMEIReport | null>(null);
    const [loading, setLoading] = useState(false);
    const [error, setError] = useState('');
    const csvRef = useRef<HTMLInputElement>(null);
    const pdfRef = useRef<HTMLInputElement>(null);
    const [csvFile, setCsvFile] = useState<File | null>(null);
    const [pdfFile, setPdfFile] = useState<File | null>(null);

    const handleAnalyze = async () => {
        if (!csvFile || !pdfFile) return;
        setLoading(true);
        setError('');

        try {
            const formData = new FormData();
            formData.append('csv_file', csvFile);
            formData.append('pdf_file', pdfFile);
            const { data } = await api.post('/imei/analyze', formData, {
                headers: { 'Content-Type': 'multipart/form-data' },
            });
            setReport(data);
        } catch (err: unknown) {
            const msg = (err as { response?: { data?: { error?: string } } })?.response?.data?.error;
            setError(msg || 'Ошибка анализа. Убедитесь что бэкенд запущен.');
        } finally {
            setLoading(false);
        }
    };

    const results = report?.results || [];
    const matchRate = report?.match_rate || 0;

    const handleExport = () => {
        if (results.length === 0) return;
        const headers = ['Row', 'IMEI Code', 'Status', 'Matched Number'];
        const rows = results.map(r => [r.row, r.imei_code, r.status, r.matched_number || '']);
        const csv = [headers.join(','), ...rows.map(r => r.join(','))].join('\n');
        const blob = new Blob([csv], { type: 'text/csv' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a'); a.href = url; a.download = 'imei_report.csv'; a.click();
    };

    return (
        <div>
            {/* Header */}
            <div className="flex items-center justify-between mb-6">
                <div>
                    <div className="breadcrumb mb-1">
                        <span>🏠</span>
                        <span>/</span>
                        <span className="text-text-primary font-medium">IMEI Проверка</span>
                    </div>
                    <h1 className="page-title">Анализ IMEI кодов</h1>
                    <p className="page-subtitle">Проверка 14-значных IMEI из CSV против 15-значных номеров в PDF-декларации</p>
                </div>
                {results.length > 0 && (
                    <button onClick={handleExport} className="btn-primary">
                        <Download size={16} />
                        Export Report
                    </button>
                )}
            </div>

            {/* Upload Section */}
            {!report && (
                <div className="card p-6 mb-6">
                    <h3 className="text-sm font-semibold text-text-primary mb-4 flex items-center gap-2">
                        <Upload size={16} className="text-primary" />
                        Загрузка файлов
                    </h3>
                    <div className="grid grid-cols-2 gap-4 mb-4">
                        <div
                            onClick={() => csvRef.current?.click()}
                            className={`border-2 border-dashed rounded-xl p-6 text-center cursor-pointer transition-all ${csvFile ? 'border-success bg-success-light' : 'border-border hover:border-primary/40'}`}
                        >
                            <input ref={csvRef} type="file" accept=".csv" onChange={(e) => e.target.files?.[0] && setCsvFile(e.target.files[0])} className="hidden" />
                            <Smartphone size={24} className={`mx-auto mb-2 ${csvFile ? 'text-green-600' : 'text-text-muted'}`} />
                            <p className="text-sm font-medium">{csvFile ? csvFile.name : 'CSV с IMEI кодами'}</p>
                            <p className="text-xs text-text-muted mt-1">14-значные IMEI</p>
                        </div>
                        <div
                            onClick={() => pdfRef.current?.click()}
                            className={`border-2 border-dashed rounded-xl p-6 text-center cursor-pointer transition-all ${pdfFile ? 'border-success bg-success-light' : 'border-border hover:border-primary/40'}`}
                        >
                            <input ref={pdfRef} type="file" accept=".pdf" onChange={(e) => e.target.files?.[0] && setPdfFile(e.target.files[0])} className="hidden" />
                            <Upload size={24} className={`mx-auto mb-2 ${pdfFile ? 'text-green-600' : 'text-text-muted'}`} />
                            <p className="text-sm font-medium">{pdfFile ? pdfFile.name : 'PDF Декларация'}</p>
                            <p className="text-xs text-text-muted mt-1">Таможенный документ</p>
                        </div>
                    </div>

                    {error && (
                        <div className="mb-4 bg-danger-light border border-danger/20 text-danger text-sm px-4 py-3 rounded-lg">
                            {error}
                        </div>
                    )}

                    <div className="flex justify-end">
                        <button onClick={handleAnalyze} disabled={!csvFile || !pdfFile || loading} className="btn-primary disabled:opacity-50 disabled:cursor-not-allowed">
                            {loading ? 'Анализ...' : 'Запустить анализ'}
                        </button>
                    </div>
                </div>
            )}

            {/* Results */}
            {report && (
                <>
                    {/* Stats Row */}
                    <div className="grid grid-cols-4 gap-4 mb-6">
                        {/* Overall Match Rate */}
                        <div className="card p-5 flex items-center gap-4">
                            <div className="relative w-16 h-16">
                                <svg className="w-16 h-16 -rotate-90" viewBox="0 0 64 64">
                                    <circle cx="32" cy="32" r="28" fill="none" stroke="#E2E8F0" strokeWidth="4" />
                                    <circle cx="32" cy="32" r="28" fill="none" stroke="#3B82F6" strokeWidth="4"
                                        strokeDasharray={`${matchRate * 1.76} 176`} strokeLinecap="round" />
                                </svg>
                                <span className="absolute inset-0 flex items-center justify-center text-sm font-bold text-primary">{matchRate}%</span>
                            </div>
                            <div>
                                <p className="text-sm font-semibold text-text-primary">Match Rate</p>
                                <p className="text-xs text-text-muted">{report.found} out of {report.total} verified</p>
                            </div>
                        </div>

                        <div className="stat-card">
                            <p className="text-xs text-text-muted uppercase tracking-wider">Всего</p>
                            <p className="text-xl font-bold text-text-primary">{report.total}</p>
                        </div>
                        <div className="stat-card">
                            <p className="text-xs text-text-muted uppercase tracking-wider">Найдено</p>
                            <p className="text-xl font-bold text-green-600">{report.found}</p>
                        </div>
                        <div className="stat-card">
                            <p className="text-xs text-text-muted uppercase tracking-wider">Не найдено</p>
                            <p className="text-xl font-bold text-red-600">{report.missing}</p>
                        </div>
                    </div>

                    {/* New analysis button */}
                    <div className="mb-4">
                        <button onClick={() => { setReport(null); setCsvFile(null); setPdfFile(null); }} className="btn-secondary text-sm">
                            ← Новый анализ
                        </button>
                    </div>

                    {/* Matching Results Table */}
                    <div className="card overflow-hidden">
                        <div className="px-6 py-4 border-b border-border flex items-center justify-between">
                            <h3 className="text-base font-semibold text-text-primary">Matching Results</h3>
                        </div>

                        <table className="data-table">
                            <thead>
                                <tr>
                                    <th>Row #</th>
                                    <th>IMEI Code (14-digit)</th>
                                    <th>Status in PDF</th>
                                    <th>Matched 15-digit Number</th>
                                </tr>
                            </thead>
                            <tbody>
                                {results.map((r) => (
                                    <tr key={r.row}>
                                        <td className="text-text-muted">{String(r.row).padStart(3, '0')}</td>
                                        <td className="font-mono font-medium text-text-primary">{r.imei_code}</td>
                                        <td>
                                            {r.status === 'found' ? (
                                                <span className="badge-success"><CheckCircle size={12} /> Found</span>
                                            ) : (
                                                <span className="badge-danger"><XCircle size={12} /> Missing</span>
                                            )}
                                        </td>
                                        <td className="font-mono text-text-secondary">{r.matched_number || <span className="text-text-muted italic">-- Not found --</span>}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>

                        <div className="px-4 py-3 border-t border-border">
                            <span className="text-sm text-text-muted">Showing {results.length} entries</span>
                        </div>
                    </div>
                </>
            )}
        </div>
    );
}
