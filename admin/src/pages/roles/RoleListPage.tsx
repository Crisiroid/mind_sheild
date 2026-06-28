import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import DataTable from '@/components/ui/DataTable'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import { listAdminRoles, deleteAdminRole } from '@/api/services/adminService'
import { showSuccess } from '@/utils/errorHandler'
import type { AdminRole } from '@/types/admin'
import type { Column } from '@/components/ui/DataTable'

export default function RoleListPage() {
  const navigate = useNavigate()
  const [data, setData] = useState<AdminRole[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [total, setTotal] = useState(0)
  const [pages, setPages] = useState(1)
  const [deleteTarget, setDeleteTarget] = useState<AdminRole | null>(null)
  const [actionLoading, setActionLoading] = useState(false)

  const fetchData = useCallback(async () => {
    setLoading(true)
    try {
      const res = await listAdminRoles({ page, page_size: pageSize })
      setData(res.data.data)
      setTotal(res.data.total)
      setPages(res.data.pages)
    } catch {
    } finally {
      setLoading(false)
    }
  }, [page, pageSize])

  useEffect(() => {
    fetchData()
  }, [fetchData])

  const handleDelete = async () => {
    if (!deleteTarget) return
    setActionLoading(true)
    try {
      await deleteAdminRole(deleteTarget.id)
      showSuccess('نقش با موفقیت حذف شد')
      setDeleteTarget(null)
      fetchData()
    } catch {
    } finally {
      setActionLoading(false)
    }
  }

  const columns: Column<AdminRole>[] = [
    { key: 'role_name', title: 'نام نقش' },
    { key: 'description', title: 'توضیحات' },
    {
      key: 'permissions',
      title: 'دسترسی‌ها',
      render: (item) => {
        try {
          const perms = JSON.parse(item.permissions)
          return Array.isArray(perms) ? perms.join(', ') : item.permissions
        } catch {
          return item.permissions || '-'
        }
      },
    },
    {
      key: 'created_at',
      title: 'تاریخ ایجاد',
      render: (item) => new Date(item.created_at).toLocaleDateString('fa-IR'),
    },
  ]

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-900">مدیریت نقش‌ها</h1>
        <button
          onClick={() => navigate('/roles/new')}
          className="px-4 py-2 bg-primary-600 text-white text-sm rounded-md hover:bg-primary-700 transition-colors"
        >
          ایجاد نقش جدید
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <DataTable<AdminRole>
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
                onClick={() => navigate(`/roles/${item.id}/edit`)}
                className="text-xs px-3 py-1 bg-blue-50 text-blue-700 rounded-md hover:bg-blue-100 transition-colors"
              >
                ویرایش
              </button>
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
        title="حذف نقش"
        message={`آیا از حذف نقش "${deleteTarget?.role_name}" اطمینان دارید؟`}
        variant="danger"
        loading={actionLoading}
      />
    </div>
  )
}

