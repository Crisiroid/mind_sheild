import { NavLink } from 'react-router-dom'
import { useState } from 'react'

const menuItems = [
    { name: 'داشبورد', icon: 'LayoutDashboard', path: '/' },
    { name: 'کاربران', icon: 'Users', path: '/users' },
    { name: 'ادمین‌ها', icon: 'Shield', path: '/admins' },
    { name: 'نقش‌ها', icon: 'Key', path: '/roles' },
    {
        name: 'داده‌ها',
        icon: 'Database',
        children: [
            { name: 'تقویم‌ها', path: '/data/calendars' },
            { name: 'احساسات', path: '/data/emotions' },
            { name: 'استرس', path: '/data/stress-events' },
            { name: 'تنش بدنی', path: '/data/body-tension' },
            { name: 'تنفس', path: '/data/breathing' },
            { name: 'بازی شناختی', path: '/data/cognitive-games' },
            { name: 'بایدهای ذهنی', path: '/data/mental-musts' },
            { name: 'افکار منفی', path: '/data/negative-thoughts' },
            { name: 'دادگاه ذهن', path: '/data/mind-court' },
            { name: 'تمرین تعارض', path: '/data/conflict-exercises' },
            { name: 'ردیاب خلق', path: '/data/mood-tracker' },
            { name: 'ارزش‌های نقش', path: '/data/roles-values' },
            { name: 'افکار آسمانی', path: '/data/sky-thoughts' },
            { name: 'تایمر ذهن‌آگاهی', path: '/data/mindful-timers' },
            { name: 'پذیرش', path: '/data/acceptance' },
            { name: 'گزارش‌های هفتگی', path: '/data/weekly-reports' },
        ],
    },
    { name: 'گزارش‌ها', icon: 'BarChart3', path: '/reports' },
    { name: 'لاگ‌ها', icon: 'FileText', path: '/logs' },
    { name: 'رسانه', icon: 'Video', path: '/media' },
    { name: 'پروفایل', icon: 'Settings', path: '/profile' },
]

const Icons: Record<string, React.FC<{ className?: string; style?: React.CSSProperties }>> = {
    LayoutDashboard: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <rect x="3" y="3" width="7" height="9" rx="1" />
            <rect x="14" y="3" width="7" height="5" rx="1" />
            <rect x="14" y="12" width="7" height="9" rx="1" />
            <rect x="3" y="16" width="7" height="5" rx="1" />
        </svg>
    ),
    Users: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2" />
            <circle cx="9" cy="7" r="4" />
            <path d="M22 21v-2a4 4 0 0 0-3-3.87" />
            <path d="M16 3.13a4 4 0 0 1 0 7.75" />
        </svg>
    ),
    Shield: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z" />
        </svg>
    ),
    Key: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <circle cx="8" cy="15" r="5" />
            <path d="m21 2-9.6 9.6" />
            <path d="m15.5 7.5 3 3L22 7l-3-3" />
        </svg>
    ),
    Database: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <ellipse cx="12" cy="5" rx="9" ry="3" />
            <path d="M3 5v14c0 1.66 4.03 3 9 3s9-1.34 9-3V5" />
            <path d="M3 12c0 1.66 4.03 3 9 3s9-1.34 9-3" />
        </svg>
    ),
    BarChart3: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M3 3v18h18" />
            <path d="M18 17V9" />
            <path d="M13 17V5" />
            <path d="M8 17v-3" />
        </svg>
    ),
    FileText: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M14.5 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V7.5L14.5 2z" />
            <polyline points="14,2 14,8 20,8" />
            <line x1="16" y1="13" x2="8" y2="13" />
            <line x1="16" y1="17" x2="8" y2="17" />
            <line x1="10" y1="9" x2="8" y2="9" />
        </svg>
    ),
    Video: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="m16 13 5.223 3.482a.5.5 0 0 0 .777-.416V7.87a.5.5 0 0 0-.752-.432L16 10.5" />
            <rect x="2" y="6" width="14" height="12" rx="2" />
        </svg>
    ),
    Settings: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z" />
            <circle cx="12" cy="12" r="3" />
        </svg>
    ),
    ChevronDown: ({ className }) => (
        <svg className={className} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
            <path d="m6 9 6 6 6-6" />
        </svg>
    ),
}

interface SidebarProps {
    onClose?: () => void
}

export default function Sidebar({ onClose }: SidebarProps) {
    const [expandedMenu, setExpandedMenu] = useState<string | null>(null)
    const ChevronDown = Icons.ChevronDown

    const handleNavLinkClick = () => {
        if (onClose) {
            onClose()
        }
    }

    return (
        <aside className="w-72 flex flex-col overflow-hidden h-full" style={{ background: '#1e293b', borderLeft: '1px solid #334155' }}>
            { }
            <div className="px-6 py-5" style={{ borderBottom: '1px solid #334155' }}>
                <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                        <div
                            className="w-11 h-11 rounded-lg flex items-center justify-center shadow-md"
                            style={{ background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)' }}
                        >
                            <svg className="w-6 h-6 text-white" viewBox="0 0 24 24" fill="currentColor">
                                <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z" />
                            </svg>
                        </div>
                        <div>
                            <h1 className="text-lg font-bold" style={{ color: '#f1f5f9' }}>پنل مدیریت</h1>
                            <p className="text-xs mt-0.5" style={{ color: '#94a3b8' }}>سپر روان</p>
                        </div>
                    </div>
                    { }
                    <button
                        onClick={onClose}
                        className="lg:hidden p-2 rounded-lg hover:bg-slate-700 transition-colors"
                        style={{ color: '#cbd5e1' }}
                    >
                        <svg className="w-5 h-5" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <line x1="18" y1="6" x2="6" y2="18" />
                            <line x1="6" y1="6" x2="18" y2="18" />
                        </svg>
                    </button>
                </div>
            </div>

            { }
            <nav className="flex-1 px-3 py-4 space-y-1 overflow-y-auto">
                {menuItems.map((item) => {
                    const IconComponent = Icons[item.icon]

                    return (
                        <div key={item.name}>
                            {item.children ? (
                                <div>
                                    <button
                                        onClick={() =>
                                            setExpandedMenu(expandedMenu === item.name ? null : item.name)
                                        }
                                        className="w-full flex items-center justify-between px-3 py-2.5 rounded-lg transition-all duration-200 group"
                                        style={{ color: '#cbd5e1' }}
                                        onMouseEnter={(e) => e.currentTarget.style.background = '#334155'}
                                        onMouseLeave={(e) => e.currentTarget.style.background = 'transparent'}
                                    >
                                        <div className="flex items-center gap-3">
                                            {IconComponent && <IconComponent className="w-5 h-5" />}
                                            <span className="font-medium text-sm">{item.name}</span>
                                        </div>
                                        {Icons.ChevronDown && <Icons.ChevronDown
                                            className={`w-4 h-4 transition-transform duration-200 ${expandedMenu === item.name ? 'rotate-180' : ''
                                                }`}
                                        />}
                                    </button>

                                    {expandedMenu === item.name && (
                                        <div className="mr-4 mt-1 space-y-0.5 animate-slide-in">
                                            {item.children.map((child) => (
                                                <NavLink
                                                    key={child.path}
                                                    to={child.path}
                                                    onClick={handleNavLinkClick}
                                                    className={({ isActive }) =>
                                                        `block px-3 py-2 rounded-md text-sm transition-all duration-200 ${isActive
                                                            ? 'text-white font-semibold'
                                                            : 'hover:bg-opacity-50'
                                                        }`
                                                    }
                                                    style={({ isActive }) => ({
                                                        background: isActive ? 'rgba(102, 126, 234, 0.15)' : 'transparent',
                                                        color: isActive ? '#e2e8f0' : '#94a3b8',
                                                        borderRight: isActive ? '3px solid #667eea' : '3px solid transparent',
                                                    })}
                                                    onMouseEnter={(e) => {
                                                        if (!e.currentTarget.classList.contains('active')) {
                                                            e.currentTarget.style.background = '#334155'
                                                        }
                                                    }}
                                                    onMouseLeave={(e) => {
                                                        if (!e.currentTarget.classList.contains('active')) {
                                                            e.currentTarget.style.background = 'transparent'
                                                        }
                                                    }}
                                                >
                                                    {child.name}
                                                </NavLink>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            ) : (
                                <NavLink
                                    to={item.path!}
                                    onClick={handleNavLinkClick}
                                    className={({ isActive }) =>
                                        `flex items-center gap-3 px-3 py-2.5 rounded-lg transition-all duration-200 group ${isActive
                                            ? 'text-white font-semibold'
                                            : ''
                                        }`
                                    }
                                    style={({ isActive }) => ({
                                        background: isActive ? 'rgba(102, 126, 234, 0.15)' : 'transparent',
                                        color: isActive ? '#e2e8f0' : '#cbd5e1',
                                        borderRight: isActive ? '3px solid #667eea' : '3px solid transparent',
                                    })}
                                    onMouseEnter={(e) => {
                                        if (!e.currentTarget.classList.contains('active')) {
                                            e.currentTarget.style.background = '#334155'
                                        }
                                    }}
                                    onMouseLeave={(e) => {
                                        if (!e.currentTarget.classList.contains('active')) {
                                            e.currentTarget.style.background = 'transparent'
                                        }
                                    }}
                                >
                                    {IconComponent && <IconComponent className="w-5 h-5" />}
                                    <span className="font-medium text-sm">{item.name}</span>
                                </NavLink>
                            )}
                        </div>
                    )
                })}
            </nav>
        </aside>
    )
}
