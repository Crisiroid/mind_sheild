import { useState, useEffect, useCallback } from 'react'
import {
  LineChart, Line, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend,
  ResponsiveContainer, PieChart, Pie, Cell,
} from 'recharts'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import {
  getUserActivity, getStressAnalytics, getBodyTensionReport,
  getCognitivePatterns, getMoodTrends, getEngagement, getWeeklyProgress,
} from '@/api/services/reportService'
import { handleApiError } from '@/utils/errorHandler'

const TABS = [
  'فعالیت کاربران',
  'تحلیل استرس',
  'تنش بدنی',
  'الگوهای شناختی',
  'روند خلق',
  'تعامل',
  'پیشرفت هفتگی',
]

const TAB_APIS = [
  getUserActivity,
  getStressAnalytics,
  getBodyTensionReport,
  getCognitivePatterns,
  getMoodTrends,
  getEngagement,
  getWeeklyProgress,
] as const

const PIE_COLORS = ['#10B981', '#F59E0B', '#EF4444', '#3B82F6', '#8B5CF6', '#EC4899']

export default function ReportsPage() {
  const [activeTab, setActiveTab] = useState(0)
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')
  const [data, setData] = useState<any>(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchData = useCallback(async (apiFn: (params?: { date_from?: string; date_to?: string }) => Promise<any>) => {
    setLoading(true)
    setError(null)
    try {
      const params: { date_from?: string; date_to?: string } = {}
      if (dateFrom) params.date_from = dateFrom
      if (dateTo) params.date_to = dateTo
      const response = await apiFn(params)
      setData(response.data)
    } catch (err: any) {
      handleApiError(err)
      setError('خطا در دریافت اطلاعات')
    } finally {
      setLoading(false)
    }
  }, [dateFrom, dateTo])

  useEffect(() => {
    fetchData(TAB_APIS[activeTab])
  }, [activeTab, fetchData])

  const renderUserActivity = () => {
    const chartData = Array.isArray(data) ? data : data?.activities || data?.data || []
    return (
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">نمودار فعالیت کاربران</h3>
        <ResponsiveContainer width="100%" height={350}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line type="monotone" dataKey="count" stroke="#3B82F6" name="تعداد فعالیت" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    )
  }

  const renderStressAnalytics = () => {
    const chartData = Array.isArray(data) ? data : data?.levels || data?.data || []
    return (
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">تحلیل سطح استرس</h3>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          <ResponsiveContainer width="100%" height={350}>
            <BarChart data={chartData}>
              <CartesianGrid strokeDasharray="3 3" />
              <XAxis dataKey="level" />
              <YAxis />
              <Tooltip />
              <Legend />
              <Bar dataKey="count" fill="#EF4444" name="تعداد" radius={[4, 4, 0, 0]} />
            </BarChart>
          </ResponsiveContainer>
          <ResponsiveContainer width="100%" height={350}>
            <PieChart>
              <Pie
                data={chartData}
                dataKey="count"
                nameKey="level"
                cx="50%"
                cy="50%"
                outerRadius={100}
                label
              >
                {chartData.map((_: any, idx: number) => (
                  <Cell key={`cell-${idx}`} fill={PIE_COLORS[idx % PIE_COLORS.length]} />
                ))}
              </Pie>
              <Tooltip />
              <Legend />
            </PieChart>
          </ResponsiveContainer>
        </div>
      </div>
    )
  }

  const renderBodyTension = () => {
    const items = Array.isArray(data) ? data : data?.tensions || data?.data || []
    return (
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">گزارش تنش بدنی</h3>
        {items.length > 0 ? (
          <div className="bg-gray-50 rounded-lg p-4">
            <pre className="text-sm text-gray-700 whitespace-pre-wrap">{JSON.stringify(items, null, 2)}</pre>
          </div>
        ) : (
          <p className="text-gray-500">داده‌ای برای نمایش وجود ندارد</p>
        )}
      </div>
    )
  }

  const renderCognitivePatterns = () => {
    const items = Array.isArray(data) ? data : data?.patterns || data?.data || []
    return (
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">الگوهای شناختی</h3>
        {items.length > 0 ? (
          <div className="bg-gray-50 rounded-lg p-4">
            <pre className="text-sm text-gray-700 whitespace-pre-wrap">{JSON.stringify(items, null, 2)}</pre>
          </div>
        ) : (
          <p className="text-gray-500">داده‌ای برای نمایش وجود ندارد</p>
        )}
      </div>
    )
  }

  const renderMoodTrends = () => {
    const chartData = Array.isArray(data) ? data : data?.trends || data?.data || []
    return (
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">روند خلق کاربران</h3>
        <ResponsiveContainer width="100%" height={350}>
          <LineChart data={chartData}>
            <CartesianGrid strokeDasharray="3 3" />
            <XAxis dataKey="date" />
            <YAxis />
            <Tooltip />
            <Legend />
            <Line type="monotone" dataKey="positive" stroke="#10B981" name="مثبت" strokeWidth={2} />
            <Line type="monotone" dataKey="neutral" stroke="#F59E0B" name="خنثی" strokeWidth={2} />
            <Line type="monotone" dataKey="negative" stroke="#EF4444" name="منفی" strokeWidth={2} />
          </LineChart>
        </ResponsiveContainer>
      </div>
    )
  }

  const renderEngagement = () => {
    const items = Array.isArray(data) ? data : data?.engagement || data?.data || []
    return (
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">تعامل کاربران</h3>
        {items.length > 0 ? (
          <div className="bg-gray-50 rounded-lg p-4">
            <pre className="text-sm text-gray-700 whitespace-pre-wrap">{JSON.stringify(items, null, 2)}</pre>
          </div>
        ) : (
          <p className="text-gray-500">داده‌ای برای نمایش وجود ندارد</p>
        )}
      </div>
    )
  }

  const renderWeeklyProgress = () => {
    const items = Array.isArray(data) ? data : data?.progress || data?.data || []
    return (
      <div>
        <h3 className="text-lg font-semibold text-gray-900 mb-4">پیشرفت هفتگی</h3>
        {items.length > 0 ? (
          <div className="bg-gray-50 rounded-lg p-4">
            <pre className="text-sm text-gray-700 whitespace-pre-wrap">{JSON.stringify(items, null, 2)}</pre>
          </div>
        ) : (
          <p className="text-gray-500">داده‌ای برای نمایش وجود ندارد</p>
        )}
      </div>
    )
  }

  const renderContent = () => {
    if (loading) {
      return (
        <div className="py-20">
          <LoadingSpinner size="lg" />
        </div>
      )
    }

    if (error) {
      return (
        <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
          <p className="text-red-600">{error}</p>
          <button
            onClick={() => fetchData(TAB_APIS[activeTab])}
            className="mt-3 px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 text-sm"
          >
            تلاش مجدد
          </button>
        </div>
      )
    }

    if (!data) {
      return (
        <div className="bg-white rounded-lg shadow p-6 text-center">
          <p className="text-gray-500">داده‌ای برای نمایش وجود ندارد</p>
        </div>
      )
    }

    switch (activeTab) {
      case 0: return renderUserActivity()
      case 1: return renderStressAnalytics()
      case 2: return renderBodyTension()
      case 3: return renderCognitivePatterns()
      case 4: return renderMoodTrends()
      case 5: return renderEngagement()
      case 6: return renderWeeklyProgress()
      default: return null
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">گزارش‌ها</h1>

      <div className="bg-white rounded-lg shadow p-4 mb-6">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 items-end">
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">از تاریخ</label>
            <input
              type="date"
              value={dateFrom}
              onChange={e => setDateFrom(e.target.value)}
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>
          <div>
            <label className="block text-xs font-medium text-gray-600 mb-1">تا تاریخ</label>
            <input
              type="date"
              value={dateTo}
              onChange={e => setDateTo(e.target.value)}
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>
          <button
            onClick={() => fetchData(TAB_APIS[activeTab])}
            className="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 text-sm"
          >
            اعمال فیلتر
          </button>
          <button
            onClick={() => {
              setDateFrom('')
              setDateTo('')
            }}
            className="px-4 py-2 text-sm text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50"
          >
            پاک کردن
          </button>
        </div>
      </div>

      <div className="bg-white rounded-lg shadow mb-6">
        <div className="border-b border-gray-200 overflow-x-auto">
          <nav className="flex whitespace-nowrap">
            {TABS.map((tab, index) => (
              <button
                key={tab}
                onClick={() => setActiveTab(index)}
                className={`px-4 py-3 text-sm font-medium border-b-2 transition-colors ${activeTab === index
                  ? 'border-primary-600 text-primary-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
                  }`}
              >
                {tab}
              </button>
            ))}
          </nav>
        </div>
        <div className="p-6">
          {renderContent()}
        </div>
      </div>
    </div>
  )
}

