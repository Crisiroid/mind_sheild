import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import DataTable from '@/components/ui/DataTable'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import FilterBar from '@/components/ui/FilterBar'
import { listAdminUsers, deleteAdminUser, deactivateAdminUser, listAdminRoles } from '@/api/services/adminService'
import { showSuccess } from '@/utils/errorHandler'
import type { AdminUser, AdminRole } from '@/types/admin'
import type { Column } from '@/components/ui/DataTable'

export default function AdminListPage() {
  const navigate = useNavigate()
  const [data, setData] = useState<AdminUser[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [total, setTotal] = useState(0)
  const [pages, setPages] = useState(1)
  const [search, setSearch] = useState('')
  const [roleFilter, setRoleFilter] = useState('')
  const [activeFilter, setActiveFilter] = useState('')
  const [roles, setRoles] = useState<AdminRole[]>([])
  const [deleteTarget, setDeleteTarget] = useState<AdminUser | null>(null)
  const [deactivateTarget, setDeactivateTarget] = useState<AdminUser | null>(null)
  const [actionLoading, setActionLoading] = useState(false)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const params: Record<string, unknown> = { page, page_size: pageSize }
      if (search) params.search = search
      if (roleFilter) params.role_id = roleFilter
      if (activeFilter) params.is_active = activeFilter === 'true'
      const res = await listAdminUsers(params)
      setData(res.data.data)
      setTotal(res.data.total)
      setPages(res.data.pages)
    } catch {
    } finally {
      setLoading(false)
    }
  }, [page, pageSize, search, roleFilter, activeFilter])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  useEffect(() => {
    listAdminRoles({ page: 1, page_size: 100 })
      .then(res => setRoles(res.data.data))
      .catch(() => { })
  }, [])

  const handleDelete = async () => {
    if (!deleteTarget) return
    setActionLoading(true)
    try {
      await deleteAdminUser(deleteTarget.id)
      showSuccess('ادمین با موفقیت حذف شد')
      setDeleteTarget(null)
      fetchData()
    } catch {
    } finally {
      setActionLoading(false)
    }
  }

  const handleDeactivate = async () => {
    if (!deactivateTarget) return
    setActionLoading(true)
    try {
      await deactivateAdminUser(deactivateTarget.id)
      showSuccess('ادمین با موفقیت غیرفعال شد')
      setDeactivateTarget(null)
      fetchData()
    } catch {
    } finally {
      setActionLoading(false)
    }
  }

  const handleReset = () => {
    setSearch('')
    setRoleFilter('')
    setActiveFilter('')
    setPage(1)
  }

  const columns: Column<AdminUser>[] = [
    { key: 'username', title: 'نام کاربری' },
    { key: 'full_name', title: 'نام کامل' },
    { key: 'email', title: 'ایمیل' },
    {
      key: 'role_id',
      title: 'نقش',
      render: (item) => {
        const role = roles.find(r => r.id === item.role_id)
        return role?.role_name || '-'
      },
    },
    {
      key: 'is_active',
      title: 'وضعیت',
      render: (item) => (
        <span className={`px-2 py-1 text-xs rounded-full ${item.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}`}>
          {item.is_active ? 'فعال' : 'غیرفعال'}
        </span>
      ),
    },
    {
      key: 'last_login',
      title: 'آخرین ورود',
      render: (item) => (item.last_login ? new Date(item.last_login).toLocaleDateString('fa-IR') : '-'),
    },
  ]

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-900">مدیریت ادمین‌ها</h1>
        <button
          onClick={() => navigate('/admins/new')}
          className="px-4 py-2 bg-primary-600 text-white text-sm rounded-md hover:bg-primary-700 transition-colors"
        >
          ایجاد ادمین جدید
        </button>
      </div>

      <FilterBar
        search={search}
        onSearchChange={(val) => { setSearch(val); setPage(1) }}
        onReset={handleReset}
        extraFilters={
          <>
            <div>
              <label className="block text-xs font-medium text-gray-600 mb-1">نقش</label>
              <select
                value={roleFilter}
                onChange={e => { setRoleFilter(e.target.value); setPage(1) }}
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
                <option value="">همه نقش‌ها</option>
                {roles.map(role => (
                  <option key={role.id} value={role.id}>{role.role_name}</option>
                ))}
              </select>
            </div>
            <div>
              <label className="block text-xs font-medium text-gray-600 mb-1">وضعیت</label>
              <select
                value={activeFilter}
                onChange={e => { setActiveFilter(e.target.value); setPage(1) }}
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
                <option value="">همه</option>
                <option value="true">فعال</option>
                <option value="false">غیرفعال</option>
              </select>
            </div>
          </>
        }
      />

      <div className="bg-white rounded-lg shadow">
        <DataTable<AdminUser>
          columns={columns}
          data={data}
          loading={loading}
          page={page}
          pageSize={pageSize}
          total={total}
          pages={pages}
          onPageChange={setPage}
          onPageSizeChange={setPageSize}
          actions={(item) => (
            <div className="flex gap-2">
              <button
                onClick={() => navigate(`/admins/${item.id}/edit`)}
                className="text-xs px-3 py-1 bg-blue-50 text-blue-700 rounded-md hover:bg-blue-100 transition-colors"
              >
                ویرایش
              </button>
              {item.is_active && (
                <button
                  onClick={() => setDeactivateTarget(item)}
                  className="text-xs px-3 py-1 bg-yellow-50 text-yellow-700 rounded-md hover:bg-yellow-100 transition-colors"
                >
                  غیرفعال
                </button>
              )}
              <button
                onClick={() => setDeleteTarget(item)}
                className="text-xs px-3 py-1 bg-red-50 text-red-700 rounded-md hover:bg-red-100 transition-colors"
              >
                حذف
              </button>
            </div>
          )}
        />
      </div>

      <ConfirmDialog
        isOpen={!!deleteTarget}
        onConfirm={handleDelete}
        onCancel={() => setDeleteTarget(null)}
        title="حذف ادمین"
        message={`آیا از حذف ادمین "${deleteTarget?.username}" اطمینان دارید؟`}
        variant="danger"
        loading={actionLoading}
      />

      <ConfirmDialog
        isOpen={!!deactivateTarget}
        onConfirm={handleDeactivate}
        onCancel={() => setDeactivateTarget(null)}
        title="غیرفعال کردن ادمین"
        message={`آیا از غیرفعال کردن ادمین "${deactivateTarget?.username}" اطمینان دارید؟`}
        variant="warning"
        loading={actionLoading}
      />
    </div>
  )
}

