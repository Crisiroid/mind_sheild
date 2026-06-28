import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { exportData } from '@/api/services/reportService'
import { exportUsers } from '@/api/services/userService'
import { showSuccess, handleApiError } from '@/utils/errorHandler'

const exportSchema = z.object({
  exportType: z.enum(['reports', 'users']),
  dateFrom: z.string().optional(),
  dateTo: z.string().optional(),
  userId: z.string().optional(),
})

type ExportFormData = z.infer<typeof exportSchema>

export default function ExportPage() {
  const [loading, setLoading] = useState(false)
  const [result, setResult] = useState<any>(null)

  const {
    register,
    handleSubmit,
    watch,
    formState: { errors },
  } = useForm<ExportFormData>({
    resolver: zodResolver(exportSchema),
    defaultValues: { exportType: 'reports' },
  })

  const exportType = watch('exportType')

  const onSubmit = async (formData: ExportFormData) => {
    setLoading(true)
    setResult(null)
    try {
      let response: any
      if (formData.exportType === 'users') {
        response = await exportUsers({
          date_from: formData.dateFrom || undefined,
          date_to: formData.dateTo || undefined,
        })
      } else {
        response = await exportData({
          date_from: formData.dateFrom || undefined,
          date_to: formData.dateTo || undefined,
          user_id: formData.userId || undefined,
        })
      }
      setResult(response.data)
      showSuccess('خروجی با موفقیت دریافت شد')
    } catch (err: any) {
      handleApiError(err, 'خطا در دریافت خروجی')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">خروجی گرفتن</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">فرم خروجی</h2>
          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">نوع خروجی</label>
              <select
                {...register('exportType')}
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              >
                <option value="reports">گزارش عمومی</option>
                <option value="users">اطلاعات کاربران</option>
              </select>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">از تاریخ</label>
              <input
                {...register('dateFrom')}
                type="date"
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">تا تاریخ</label>
              <input
                {...register('dateTo')}
                type="date"
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              />
            </div>

            {exportType === 'reports' && (
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">شناسه کاربر (اختیاری)</label>
                <input
                  {...register('userId')}
                  type="text"
                  placeholder="شناسه کاربر..."
                  className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                />
              </div>
            )}

            {errors.userId && (
              <p className="text-sm text-red-600">{errors.userId.message}</p>
            )}

            <button
              type="submit"
              disabled={loading}
              className="w-full bg-primary-600 text-white py-2 px-4 rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {loading ? 'در حال دریافت...' : 'دریافت خروجی'}
            </button>
          </form>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">پیش‌نمایش خروجی</h2>
          {result ? (
            <div className="bg-gray-50 rounded-lg p-4 max-h-[500px] overflow-auto">
              <pre className="text-sm text-gray-700 whitespace-pre-wrap">
                {JSON.stringify(result, null, 2)}
              </pre>
            </div>
          ) : (
            <div className="bg-gray-50 rounded-lg p-6 text-center">
              <p className="text-gray-500">
                پس از دریافت خروجی، نتیجه در اینجا نمایش داده می‌شود
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  )
}

