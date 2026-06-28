import { useState, useEffect, useCallback } from 'react'
import DataTable from '@/components/ui/DataTable'
import FilterBar from '@/components/ui/FilterBar'
import { listCalendars } from '@/api/services/dataService'
import type { Column } from '@/components/ui/DataTable'
import type { ListParams } from '@/api/services/dataService'

export default function CalendarsPage() {
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
      const response = await listCalendars(params)
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
    { key: 'event_date', title: 'تاریخ رویداد', render: (item) => new Date(item.event_date).toLocaleDateString('fa-IR') },
    { key: 'event_type', title: 'نوع رویداد' },
    { key: 'description', title: 'توضیحات' },
    { key: 'created_at', title: 'تاریخ ایجاد', render: (item) => new Date(item.created_at).toLocaleDateString('fa-IR') },
  ]

  return (
    <div className="space-y-4 md:space-y-6">
      <h1 className="text-2xl md:text-3xl font-bold text-gray-900">تقویم‌ها</h1>
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

