import { useNavigate } from 'react-router-dom'
import { useAuthStore } from '@/store/auth'

interface HeaderProps {
    onMenuToggle: () => void
}

export default function Header({ onMenuToggle }: HeaderProps) {
    const navigate = useNavigate()
    const { adminUser, logout } = useAuthStore()

    const handleLogout = () => {
        logout()
        navigate('/login')
    }

    return (
        <header className="sticky top-0 z-30 px-4 md:px-6 lg:px-8 py-4" style={{ background: '#ffffff', borderBottom: '1px solid #e2e8f0', boxShadow: '0 1px 3px 0 rgb(0 0 0 / 0.1)' }}>
            <div className="flex items-center justify-between">
                {/* Mobile menu button and user info */}
                <div className="flex items-center gap-3 flex-1 min-w-0">
                    <button
                        onClick={onMenuToggle}
                        className="lg:hidden p-2 rounded-lg hover:bg-neutral-100 transition-colors"
                        style={{ color: '#64748b' }}
                    >
                        <svg className="w-6 h-6" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="3" y1="12" x2="21" y2="12" />
                            <line x1="3" y1="6" x2="21" y2="6" />
                            <line x1="3" y1="18" x2="21" y2="18" />
                        </svg>
                    </button>

                    {/* User avatar icon */}
                    <div className="hidden sm:flex items-center justify-center w-10 h-10 rounded-full flex-shrink-0" style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}>
                        <svg className="w-5 h-5 text-white" viewBox="0 0 24 24" fill="currentColor">
                            <path d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z" />
                        </svg>
                    </div>

                    {/* User info */}
                    <div className="min-w-0 flex-1">
                        <div className="flex items-center gap-2">
                            <h2 className="text-sm sm:text-base md:text-lg font-semibold truncate" style={{ color: '#1e293b' }}>
                                خوش آمدید، <span className="font-bold">{adminUser?.full_name || adminUser?.username}</span>
                            </h2>
                        </div>
                        <p className="text-xs sm:text-sm mt-0.5 truncate hidden sm:block" style={{ color: '#64748b' }}>
                            {adminUser?.email}
                        </p>
                        <p className="text-xs mt-0.5 truncate sm:hidden" style={{ color: '#94a3b8' }}>
                            {adminUser?.email}
                        </p>
                    </div>
                </div>

                <div className="flex items-center gap-2 md:gap-4">
                    {/* Profile button - hidden on mobile */}
                    <button
                        onClick={() => navigate('/profile')}
                        className="px-3 md:px-4 py-2 text-xs md:text-sm font-semibold rounded-lg transition-all duration-200 border hidden md:flex"
                        style={{
                            background: '#ffffff',
                            color: '#475569',
                            borderColor: '#e2e8f0'
                        }}
                        onMouseEnter={(e) => {
                            e.currentTarget.style.background = '#f8fafc'
                            e.currentTarget.style.borderColor = '#cbd5e1'
                        }}
                        onMouseLeave={(e) => {
                            e.currentTarget.style.background = '#ffffff'
                            e.currentTarget.style.borderColor = '#e2e8f0'
                        }}
                    >
                        <div className="flex items-center gap-2">
                            <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z" />
                                <circle cx="12" cy="12" r="3" />
                            </svg>
                            <span className="hidden lg:inline">تنظیمات پروفایل</span>
                        </div>
                    </button>
                    <button
                        onClick={handleLogout}
                        className="px-3 md:px-4 py-2 text-xs md:text-sm font-semibold rounded-lg transition-all duration-200 shadow-sm hover:shadow-md btn-lift"
                        style={{
                            background: 'linear-gradient(135deg, #f87171 0%, #ef4444 100%)',
                            color: '#fff',
                        }}
                    >
                        <div className="flex items-center gap-2">
                            <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4" />
                                <polyline points="16,17 21,12 16,7" />
                                <line x1="21" y1="12" x2="9" y2="12" />
                            </svg>
                            <span className="hidden sm:inline">خروج</span>
                        </div>
                    </button>
                </div>
            </div>
        </header>
    )
}
