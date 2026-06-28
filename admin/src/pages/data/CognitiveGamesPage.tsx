import { useState, useEffect, useCallback } from 'react'
import DataTable from '@/components/ui/DataTable'
import FilterBar from '@/components/ui/FilterBar'
import { listCognitiveGames } from '@/api/services/dataService'
import type { Column } from '@/components/ui/DataTable'
import type { ListParams } from '@/api/services/dataService'

export default function CognitiveGamesPage() {
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
      const response = await listCognitiveGames(params)
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
    { key: 'game_date', title: 'تاریخ بازی', render: (item) => new Date(item.game_date).toLocaleDateString('fa-IR') },
    { key: 'scenario_type', title: 'نوع سناریو' },
    {
      key: 'score',
      title: 'امتیاز',
      render: (item) => item.score != null ? item.score : '-',
    },
    {
      key: 'is_correct',
      title: 'درست',
      render: (item) =>
        item.is_correct === true ? <span className="text-green-600 text-lg">✓</span>
          : item.is_correct === false ? <span className="text-red-500 text-lg">✗</span>
            : '-',
    },
    { key: 'time_taken_seconds', title: 'زمان (ثانیه)', render: (item) => item.time_taken_seconds != null ? `${item.time_taken_seconds}s` : '-' },
    { key: 'created_at', title: 'تاریخ ایجاد', render: (item) => new Date(item.created_at).toLocaleDateString('fa-IR') },
  ]

  return (
    <div className="space-y-4 md:space-y-6">
      <h1 className="text-2xl md:text-3xl font-bold text-gray-900">بازی‌های شناختی</h1>
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

