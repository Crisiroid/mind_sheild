import { useState, useEffect, useCallback } from 'react'
import { listSystemLogs } from '@/api/services/adminService'
import type { LogListParams } from '@/api/services/adminService'
import type { SystemLog } from '@/types/admin'
import DataTable from '@/components/ui/DataTable'
import FilterBar from '@/components/ui/FilterBar'
import Pagination from '@/components/ui/Pagination'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { handleApiError } from '@/utils/errorHandler'
import type { Column } from '@/components/ui/DataTable'

const SEVERITY_COLORS: Record<string, string> = {
  error: 'bg-red-100 text-red-800',
  warning: 'bg-yellow-100 text-yellow-800',
  info: 'bg-blue-100 text-blue-800',
}

const LOG_TYPES = ['system', 'auth', 'user', 'error', 'security']
const SEVERITIES = ['error', 'warning', 'info']

export default function LogsPage() {
  const [logs, setLogs] = useState<SystemLog[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [total, setTotal] = useState(0)
  const [pages, setPages] = useState(1)
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')
  const [logType, setLogType] = useState('')
  const [severity, setSeverity] = useState('')

  const fetchLogs = useCallback(async () => {
    setLoading(true)
    try {
      const params: LogListParams = { page, page_size: pageSize }
      if (dateFrom) params.date_from = dateFrom
      if (dateTo) params.date_to = dateTo
      if (logType) params.log_type = logType
      if (severity) params.severity = severity

      const response = await listSystemLogs(params)
      const paginatedData = response.data
      setLogs(paginatedData.data || [])
      setTotal(paginatedData.total || 0)
      setPages(paginatedData.pages || 1)
    } catch (err: any) {
      handleApiError(err, 'خطا در دریافت لاگ‌ها')
    } finally {
      setLoading(false)
    }
  }, [page, pageSize, dateFrom, dateTo, logType, severity])

  useEffect(() => {
    fetchLogs()
  }, [fetchLogs])

  const handleReset = () => {
    setDateFrom('')
    setDateTo('')
    setLogType('')
    setSeverity('')
    setPage(1)
  }

  const columns: Column<SystemLog>[] = [
    {
      key: 'log_type',
      title: 'نوع لاگ',
      render: (item: SystemLog) => (
        <span className="text-xs font-medium text-gray-700 bg-gray-100 px-2 py-1 rounded">
          {item.log_type}
        </span>
      ),
    },
    {
      key: 'log_message',
      title: 'پیام',
      render: (item: SystemLog) => (
        <span className="text-sm text-gray-700 max-w-xs truncate block" title={item.log_message}>
          {item.log_message}
        </span>
      ),
    },
    {
      key: 'user_id',
      title: 'شناسه کاربر',
      render: (item: SystemLog) => (
        <span className="text-xs text-gray-500 font-mono">{item.user_id || '-'}</span>
      ),
    },
    {
      key: 'severity',
      title: 'شدت',
      render: (item: SystemLog) => (
        <span className={`text-xs font-medium px-2 py-1 rounded ${SEVERITY_COLORS[item.severity] || 'bg-gray-100 text-gray-800'}`}>
          {item.severity === 'error' ? 'خطا' : item.severity === 'warning' ? 'هشدار' : item.severity === 'info' ? 'اطلاعات' : item.severity}
        </span>
      ),
    },
    {
      key: 'created_at',
      title: 'تاریخ ایجاد',
      render: (item: SystemLog) => (
        <span className="text-xs text-gray-500">{new Date(item.created_at).toLocaleDateString('fa-IR')}</span>
      ),
    },
  ]

  const extraFilters = (
    <>
      <div>
        <label className="block text-xs font-medium text-gray-600 mb-1">نوع لاگ</label>
        <select
          value={logType}
          onChange={e => { setLogType(e.target.value); setPage(1) }}
          className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
        >
          <option value="">همه</option>
          {LOG_TYPES.map(t => (
            <option key={t} value={t}>{t}</option>
          ))}
        </select>
      </div>
      <div>
        <label className="block text-xs font-medium text-gray-600 mb-1">شدت</label>
        <select
          value={severity}
          onChange={e => { setSeverity(e.target.value); setPage(1) }}
          className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
        >
          <option value="">همه</option>
          {SEVERITIES.map(s => (
            <option key={s} value={s}>
              {s === 'error' ? 'خطا' : s === 'warning' ? 'هشدار' : 'اطلاعات'}
            </option>
          ))}
        </select>
      </div>
    </>
  )

  if (loading && logs.length === 0) {
    return (
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-6">لاگ‌های سیستم</h1>
        <div className="py-20">
          <LoadingSpinner size="lg" />
        </div>
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">لاگ‌های سیستم</h1>

      <FilterBar
        dateFrom={dateFrom}
        dateTo={dateTo}
        onDateFromChange={d => { setDateFrom(d); setPage(1) }}
        onDateToChange={d => { setDateTo(d); setPage(1) }}
        extraFilters={extraFilters}
        onReset={handleReset}
      />

      <div className="bg-white rounded-lg shadow">
        <DataTable
          columns={columns}
          data={logs}
          loading={loading}
          emptyMessage="لاگی یافت نشد"
        />

        {total > 0 && (
          <div className="px-4 pb-4">
            <Pagination
              page={page}
              pageSize={pageSize}
              total={total}
              pages={pages}
              onPageChange={setPage}
              onPageSizeChange={size => { setPageSize(size); setPage(1) }}
            />
          </div>
        )}
      </div>
    </div>
  )
}

