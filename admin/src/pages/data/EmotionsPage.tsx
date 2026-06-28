import { useState, useEffect, useCallback } from 'react'
import DataTable from '@/components/ui/DataTable'
import FilterBar from '@/components/ui/FilterBar'
import { listEmotionInteractions } from '@/api/services/dataService'
import type { Column } from '@/components/ui/DataTable'
import type { ListParams } from '@/api/services/dataService'

const emotionEmojis: Record<string, string> = {
  happy: '😊', sad: '😢', angry: '😡', anxious: '😰',
  calm: '😌', excited: '🤩', tired: '😴', stressed: '😫',
}

export default function EmotionsPage() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [total, setTotal] = useState(0)
  const [pages, setPages] = useState(1)
  const [search, setSearch] = useState('')
  const [dateFrom, setDateFrom] = useState('')
  const [dateTo, setDateTo] = useState('')
  const [userId, setUserId] = useState('')

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const params: ListParams = { page, page_size: pageSize }
      if (search) params.search = search
      if (dateFrom) params.date_from = dateFrom
      if (dateTo) params.date_to = dateTo
      if (userId) params.user_id = userId
      const response = await listEmotionInteractions(params)
      setData(response.data.data)
      setTotal(response.data.total)
      setPages(response.data.pages)
    } catch (error) {
    } finally {
      setLoading(false)
    }
  }, [page, pageSize, search, dateFrom, dateTo, userId])

  useEffect(() => { fetchData() }, [fetchData])

  const handleReset = () => {
    setSearch('')
    setDateFrom('')
    setDateTo('')
    setUserId('')
    setPage(1)
  }

  const columns: Column[] = [
    { key: 'id', title: 'شناسه' },
    { key: 'user_id', title: 'کاربر' },
    {
      key: 'emotion_type',
      title: 'نوع احساس',
      render: (item) => `${emotionEmojis[item.emotion_type] || '❓'} ${item.emotion_type}`,
    },
    {
      key: 'intensity',
      title: 'شدت',
      render: (item) => (
        <div className="flex items-center gap-2">
          <div className="w-20 h-2 bg-gray-200 rounded-full overflow-hidden">
            <div className="h-full bg-primary-500 rounded-full transition-all" style={{ width: `${Math.min(item.intensity, 100)}%` }} />
          </div>
          <span className="text-xs text-gray-600">{item.intensity}</span>
        </div>
      ),
    },
    { key: 'trigger', title: 'علت' },
    { key: 'created_at', title: 'تاریخ ایجاد', render: (item) => new Date(item.created_at).toLocaleDateString('fa-IR') },
  ]

  return (
    <div className="space-y-4 md:space-y-6">
      <h1 className="text-2xl md:text-3xl font-bold text-gray-900">احساسات</h1>
      <FilterBar
        search={search}
        onSearchChange={setSearch}
        dateFrom={dateFrom}
        onDateFromChange={setDateFrom}
        dateTo={dateTo}
        onDateToChange={setDateTo}
        userId={userId}
        onUserIdChange={setUserId}
        onReset={handleReset}
      />
      <div className="bg-white rounded-lg shadow overflow-hidden">
        <DataTable
          columns={columns}
          data={data}
          loading={loading}
          page={page}
          pageSize={pageSize}
          total={total}
          pages={pages}
          onPageChange={setPage}
          onPageSizeChange={setPageSize}
        />
      </div>
    </div>
  )
}

