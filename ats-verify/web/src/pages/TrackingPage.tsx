import { useState } from 'react';
import { MapPin, Search, Truck, CheckCircle } from 'lucide-react';

interface TrackingStep {
    date: string;
    status: string;
    location: string;
    done: boolean;
}

const MOCK_STEPS: TrackingStep[] = [
    { date: '2026-02-15 14:30', status: 'Вручено получателю', location: 'Алматы, KZ', done: true },
    { date: '2026-02-14 09:15', status: 'Прибыло в пункт выдачи', location: 'Алматы, KZ', done: true },
    { date: '2026-02-12 18:00', status: 'В пути', location: 'Караганда → Алматы', done: true },
    { date: '2026-02-10 12:00', status: 'Прибыло в страну', location: 'Хоргос, KZ', done: true },
    { date: '2026-02-08 06:00', status: 'Отправлено', location: 'Гуанчжоу, CN', done: true },
];

export default function TrackingPage() {
    const [trackNumber, setTrackNumber] = useState('');
    const [loading, setLoading] = useState(false);
    const [steps, setSteps] = useState<TrackingStep[]>([]);
    const [carrier, setCarrier] = useState('');

    const handleTrack = async () => {
        if (!trackNumber.trim()) return;
        setLoading(true);
        await new Promise(r => setTimeout(r, 800));
        setSteps(MOCK_STEPS);
        setCarrier('Казпочта');
        setLoading(false);
    };

    return (
        <div>
            <div className="page-header">
                <h1 className="page-title">Отслеживание посылки</h1>
                <p className="page-subtitle">Проверьте статус доставки через Казпочту или СДЭК</p>
            </div>

            {/* Search */}
            <div className="card p-6 mb-6">
                <div className="flex gap-3">
                    <div className="relative flex-1">
                        <Search size={16} className="absolute left-3 top-1/2 -translate-y-1/2 text-text-muted" />
                        <input
                            type="text"
                            value={trackNumber}
                            onChange={(e) => setTrackNumber(e.target.value)}
                            onKeyDown={(e) => e.key === 'Enter' && handleTrack()}
                            className="input pl-10"
                            placeholder="Введите трек-номер (например, KZ1234567890)"
                        />
                    </div>
                    <button onClick={handleTrack} disabled={loading || !trackNumber.trim()} className="btn-primary disabled:opacity-50 disabled:cursor-not-allowed">
                        <Truck size={16} />
                        {loading ? 'Поиск...' : 'Отследить'}
                    </button>
                </div>
            </div>

            {/* Timeline */}
            {steps.length > 0 && (
                <div className="card p-6">
                    <div className="flex items-center justify-between mb-6">
                        <h3 className="text-base font-semibold text-text-primary">История доставки</h3>
                        <span className="badge-info">{carrier}</span>
                    </div>

                    <div className="space-y-0">
                        {steps.map((step, i) => (
                            <div key={i} className="flex gap-4">
                                {/* Timeline line */}
                                <div className="flex flex-col items-center">
                                    <div className={`w-8 h-8 rounded-full flex items-center justify-center ${i === 0 ? 'bg-success text-white' : 'bg-primary-light text-primary'
                                        }`}>
                                        {i === 0 ? <CheckCircle size={16} /> : <MapPin size={14} />}
                                    </div>
                                    {i < steps.length - 1 && (
                                        <div className="w-px h-12 bg-border" />
                                    )}
                                </div>

                                {/* Content */}
                                <div className={`pb-6 ${i < steps.length - 1 ? '' : 'pb-0'}`}>
                                    <p className="text-sm font-medium text-text-primary">{step.status}</p>
                                    <p className="text-xs text-text-secondary mt-0.5">{step.location}</p>
                                    <p className="text-xs text-text-muted mt-0.5">{step.date}</p>
                                </div>
                            </div>
                        ))}
                    </div>
                </div>
            )}
        </div>
    );
}
