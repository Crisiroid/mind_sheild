import { useState, useEffect } from 'react'
import {
    ResponsiveContainer,
    LineChart,
    Line,
    BarChart,
    Bar,
    PieChart,
    Pie,
    Cell,
    XAxis,
    YAxis,
    CartesianGrid,
    Tooltip,
    Legend,
} from 'recharts'
import { getDashboard, getUserActivity, getStressAnalytics, getCognitivePatterns } from '@/api/services/reportService'
import { getUserStats } from '@/api/services/userService'
import { PageLoading } from '@/components/ui/LoadingSpinner'
import { handleApiError } from '@/utils/errorHandler'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/Card'
import { Button } from '@/components/ui/Button'
import { Badge } from '@/components/ui/Badge'

const COLORS = ['#0e8de3', '#22c55e', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899', '#06b6d4', '#f97316']

interface DashboardData {
    total_users: number
    active_users: number
    new_users_today: number
}

interface UserActivityItem {
    date: string
    active_users: number
    new_users: number
}

interface StressAnalytic {
    category: string
    count: number
}

interface CognitivePattern {
    name: string
    value: number
}

export default function DashboardPage() {
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    const [dashboardData, setDashboardData] = useState<DashboardData | null>(null)
    const [userActivity, setUserActivity] = useState<UserActivityItem[]>([])
    const [stressAnalytics, setStressAnalytics] = useState<StressAnalytic[]>([])
    const [cognitivePatterns, setCognitivePatterns] = useState<CognitivePattern[]>([])

    const fetchData = async () => {
        setLoading(true)
        setError(null)
        try {
            const [, statsRes, activityRes, stressRes, cognitiveRes] = await Promise.all([
                getDashboard(),
                getUserStats(),
                getUserActivity(),
                getStressAnalytics(),
                getCognitivePatterns(),
            ])

            setDashboardData({
                total_users: statsRes.data.total_users,
                active_users: statsRes.data.active_users,
                new_users_today: statsRes.data.new_users_today,
            })

            setUserActivity(Array.isArray(activityRes.data) ? activityRes.data : [])
            setStressAnalytics(Array.isArray(stressRes.data) ? stressRes.data : [])
            setCognitivePatterns(Array.isArray(cognitiveRes.data) ? cognitiveRes.data : [])
        } catch (err) {
            handleApiError(err)
            setError('خطا در بارگذاری داشبورد')
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchData()
    }, [])

    if (loading) {
        return <PageLoading />
    }

    if (error) {
        return (
            <div className="animate-fade-in">
                <h1 className="text-3xl font-bold text-neutral-900 mb-6">داشبورد</h1>
                <Card variant="bordered">
                    <div className="text-center py-8">
                        <div className="w-16 h-16 mx-auto mb-4 bg-error-100 rounded-full flex items-center justify-center">
                            <svg className="w-8 h-8 text-error-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                        </div>
                        <p className="text-error-600 mb-6 font-medium">{error}</p>
                        <Button variant="danger" onClick={fetchData}>
                            تلاش مجدد
                        </Button>
                    </div>
                </Card>
            </div>
        )
    }

    const stressTotal = stressAnalytics.reduce((sum, item) => sum + item.count, 0)

    const summaryCards = [
        { title: 'کل کاربران', value: dashboardData?.total_users ?? 0, gradient: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)', icon: '👥', trend: '+12%' },
        { title: 'کاربران فعال', value: dashboardData?.active_users ?? 0, gradient: 'linear-gradient(135deg, #f093fb 0%, #f5576c 100%)', icon: '✨', trend: '+8%' },
        { title: 'کاربران جدید امروز', value: dashboardData?.new_users_today ?? 0, gradient: 'linear-gradient(135deg, #4facfe 0%, #00f2fe 100%)', icon: '🌟', trend: '+24%' },
        { title: 'تعداد تحلیل استرس', value: stressTotal, gradient: 'linear-gradient(135deg, #fa709a 0%, #fee140 100%)', icon: '📊', trend: '+15%' },
    ]

    return (
        <div className="space-y-4 md:space-y-6 lg:space-y-8 animate-fade-in">
            {/* Page Header */}
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4">
                <div>
                    <h1 className="text-2xl md:text-3xl font-bold tracking-tight text-neutral-900">داشبورد</h1>
                    <p className="mt-1 md:mt-2 text-xs md:text-sm text-neutral-600">نمای کلی از وضعیت سیستم و آمار کاربران</p>
                </div>
                <button className="px-4 md:px-5 py-2 md:py-2.5 text-xs md:text-sm font-semibold rounded-xl border border-neutral-200 text-neutral-700 bg-white transition-all duration-200 hover:bg-neutral-50 hover:border-neutral-300 hover:shadow-sm w-full sm:w-auto">
                    <div className="flex items-center justify-center gap-2">
                        <svg className="w-4 h-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                            <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4" />
                            <polyline points="7,10 12,15 17,10" />
                            <line x1="12" y1="15" x2="12" y2="3" />
                        </svg>
                        خروجی گزارش
                    </div>
                </button>
            </div>

            {/* Summary Cards */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 md:gap-6">
                {summaryCards.map((card, index) => (
                    <Card key={card.title} hover className="animate-slide-up" style={{ animationDelay: `${index * 0.05}s` }}>
                        <div className="flex items-start justify-between gap-3 md:gap-4">
                            <div className="flex-1 min-w-0">
                                <p className="text-xs md:text-sm font-medium text-neutral-600 mb-2 md:mb-3">{card.title}</p>
                                <p className="text-3xl md:text-4xl font-bold text-neutral-900 mb-2 md:mb-3">{card.value.toLocaleString()}</p>
                                <div className="flex items-center gap-2">
                                    <svg className="w-3 h-3 md:w-4 md:h-4 text-success-600 flex-shrink-0" fill="currentColor" viewBox="0 0 20 20">
                                        <path fillRule="evenodd" d="M12 7a1 1 0 110-2h5a1 1 0 011 1v5a1 1 0 11-2 0V8.414l-4.293 4.293a1 1 0 01-1.414 0L8 10.414l-4.293 4.293a1 1 0 01-1.414-1.414l5-5a1 1 0 011.414 0L11 10.586 14.586 7H12z" clipRule="evenodd" />
                                    </svg>
                                    <span className="text-xs font-semibold text-success-600">{card.trend}</span>
                                    <span className="text-xs text-neutral-500 hidden sm:inline">نسبت به دیروز</span>
                                </div>
                            </div>
                            <div
                                className="w-12 h-12 md:w-16 md:h-16 rounded-xl flex items-center justify-center text-2xl md:text-3xl shadow-lg flex-shrink-0"
                                style={{
                                    background: card.gradient,
                                }}
                            >
                                {card.icon}
                            </div>
                        </div>
                    </Card>
                ))}
            </div>

            {/* Charts Row */}
            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 md:gap-6">
                {/* User Activity Line Chart */}
                <Card hover className="animate-slide-up" style={{ animationDelay: '0.2s' }}>
                    <CardHeader>
                        <CardTitle>فعالیت کاربران</CardTitle>
                    </CardHeader>
                    <CardContent>
                        {userActivity.length > 0 ? (
                            <ResponsiveContainer width="100%" height={300}>
                                <LineChart data={userActivity}>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
                                    <XAxis dataKey="date" stroke="#94a3b8" />
                                    <YAxis stroke="#94a3b8" />
                                    <Tooltip
                                        contentStyle={{
                                            backgroundColor: '#fff',
                                            border: '1px solid #e2e8f0',
                                            borderRadius: '10px',
                                            boxShadow: '0 8px 24px rgba(0, 0, 0, 0.12)',
                                            padding: '12px'
                                        }}
                                    />
                                    <Legend />
                                    <Line type="monotone" dataKey="active_users" stroke="#667eea" strokeWidth={3} name="کاربران فعال" dot={false} />
                                    <Line type="monotone" dataKey="new_users" stroke="#f5576c" strokeWidth={3} name="کاربران جدید" dot={false} />
                                </LineChart>
                            </ResponsiveContainer>
                        ) : (
                            <div className="text-center py-16">
                                <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-neutral-100 flex items-center justify-center">
                                    <svg className="w-8 h-8 text-neutral-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                                    </svg>
                                </div>
                                <p className="text-neutral-500 font-medium">داده‌ای برای نمایش وجود ندارد</p>
                            </div>
                        )}
                    </CardContent>
                </Card>

                {/* Stress Analytics Bar Chart */}
                <Card hover className="animate-slide-up" style={{ animationDelay: '0.25s' }}>
                    <CardHeader>
                        <CardTitle>تحلیل استرس</CardTitle>
                    </CardHeader>
                    <CardContent>
                        {stressAnalytics.length > 0 ? (
                            <ResponsiveContainer width="100%" height={300}>
                                <BarChart data={stressAnalytics}>
                                    <CartesianGrid strokeDasharray="3 3" stroke="#e2e8f0" />
                                    <XAxis dataKey="category" stroke="#94a3b8" />
                                    <YAxis stroke="#94a3b8" />
                                    <Tooltip
                                        contentStyle={{
                                            backgroundColor: '#fff',
                                            border: '1px solid #e2e8f0',
                                            borderRadius: '10px',
                                            boxShadow: '0 8px 24px rgba(0, 0, 0, 0.12)',
                                            padding: '12px'
                                        }}
                                    />
                                    <Legend />
                                    <Bar dataKey="count" fill="#fa709a" name="تعداد" radius={[8, 8, 0, 0]} />
                                </BarChart>
                            </ResponsiveContainer>
                        ) : (
                            <div className="text-center py-16">
                                <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-neutral-100 flex items-center justify-center">
                                    <svg className="w-8 h-8 text-neutral-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                                    </svg>
                                </div>
                                <p className="text-neutral-500 font-medium">داده‌ای برای نمایش وجود ندارد</p>
                            </div>
                        )}
                    </CardContent>
                </Card>
            </div>

            {/* Cognitive Patterns Pie Chart */}
            <Card hover className="animate-slide-up" style={{ animationDelay: '0.3s' }}>
                <CardHeader>
                    <CardTitle>الگوهای شناختی</CardTitle>
                </CardHeader>
                <CardContent>
                    {cognitivePatterns.length > 0 ? (
                        <ResponsiveContainer width="100%" height={350}>
                            <PieChart>
                                <Pie
                                    data={cognitivePatterns}
                                    cx="50%"
                                    cy="50%"
                                    labelLine
                                    label={({ name, percent }: { name: string; percent: number }) =>
                                        `${name} ${(percent * 100).toFixed(0)}%`
                                    }
                                    outerRadius={120}
                                    dataKey="value"
                                >
                                    {cognitivePatterns.map((_entry, index) => (
                                        <Cell key={`cell-${index}`} fill={COLORS[index % COLORS.length]} />
                                    ))}
                                </Pie>
                                <Tooltip
                                    contentStyle={{
                                        backgroundColor: '#fff',
                                        border: '1px solid #e2e8f0',
                                        borderRadius: '10px',
                                        boxShadow: '0 8px 24px rgba(0, 0, 0, 0.12)',
                                        padding: '12px'
                                    }}
                                />
                                <Legend />
                            </PieChart>
                        </ResponsiveContainer>
                    ) : (
                        <div className="text-center py-16">
                            <div className="w-16 h-16 mx-auto mb-4 rounded-full bg-neutral-100 flex items-center justify-center">
                                <svg className="w-8 h-8 text-neutral-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z" />
                                </svg>
                            </div>
                            <p className="text-neutral-500 font-medium">داده‌ای برای نمایش وجود ندارد</p>
                        </div>
                    )}
                </CardContent>
            </Card>
        </div>
    )
}

