import { createContext, useContext, useState, type ReactNode } from 'react';
import type { User } from '../types';

interface AuthContextType {
    user: User | null;
    login: (user: User) => void;
    logout: () => void;
    isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(() => {
        const stored = localStorage.getItem('ats_user');
        if (stored) {
            try {
                return JSON.parse(stored);
            } catch {
                localStorage.removeItem('ats_user');
            }
        }
        return null;
    });

    const login = (userData: User) => {
        setUser(userData);
        localStorage.setItem('ats_user', JSON.stringify(userData));
        localStorage.setItem('ats_token', userData.token);
    };

    const logout = () => {
        setUser(null);
        localStorage.removeItem('ats_user');
        localStorage.removeItem('ats_token');
    };

    return (
        <AuthContext.Provider value={{ user, login, logout, isAuthenticated: !!user }}>
            {children}
        </AuthContext.Provider>
    );
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth(): AuthContextType {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error('useAuth must be used within AuthProvider');
    return ctx;
}
