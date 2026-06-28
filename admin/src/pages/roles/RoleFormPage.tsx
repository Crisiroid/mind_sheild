import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { PageLoading } from '@/components/ui/LoadingSpinner'
import { getAdminRole, createAdminRole, updateAdminRole } from '@/api/services/adminService'
import { showSuccess } from '@/utils/errorHandler'

const formSchema = z.object({
  role_name: z.string().min(2, 'نام نقش حداقل ۲ کاراکتر'),
  description: z.string().min(1, 'توضیحات الزامی است'),
  permissions: z.string().min(1, 'دسترسی‌ها الزامی است'),
})

type FormValues = z.infer<typeof formSchema>

export default function RoleFormPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const isEditing = !!id
  const [loading, setLoading] = useState(isEditing)
  const [submitting, setSubmitting] = useState(false)

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      role_name: '',
      description: '',
      permissions: '',
    },
  })

  useEffect(() => {
    if (isEditing && id) {
      setLoading(true)
      getAdminRole(id)
        .then(res => {
          const role = res.data
          reset({
            role_name: role.role_name,
            description: role.description,
            permissions: role.permissions,
          })
        })
        .catch(() => navigate('/roles'))
        .finally(() => setLoading(false))
    }
  }, [id, isEditing, reset, navigate])

  const onSubmit = async (values: FormValues) => {
    setSubmitting(true)
    try {
      if (isEditing && id) {
        await updateAdminRole(id, {
          description: values.description,
          permissions: values.permissions,
        })
        showSuccess('نقش با موفقیت به‌روزرسانی شد')
      } else {
        await createAdminRole({
          role_name: values.role_name,
          description: values.description,
          permissions: values.permissions,
        })
        showSuccess('نقش با موفقیت ایجاد شد')
      }
      navigate('/roles')
    } catch {
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) return <PageLoading />

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">
        {isEditing ? 'ویرایش نقش' : 'ایجاد نقش جدید'}
      </h1>
      <div className="bg-white rounded-lg shadow p-6 max-w-2xl">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">نام نقش</label>
            <input
              {...register('role_name')}
              disabled={isEditing}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm disabled:bg-gray-100 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            {errors.role_name && <p className="text-red-500 text-xs mt-1">{errors.role_name.message}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">توضیحات</label>
            <textarea
              {...register('description')}
              rows={3}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            {errors.description && <p className="text-red-500 text-xs mt-1">{errors.description.message}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">دسترسی‌ها (JSON)</label>
            <textarea
              {...register('permissions')}
              rows={5}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm font-mono focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              placeholder='["permission1", "permission2"]'
            />
            {errors.permissions && <p className="text-red-500 text-xs mt-1">{errors.permissions.message}</p>}
          </div>

          <div className="flex gap-3 pt-4 border-t">
            <button
              type="submit"
              disabled={submitting}
              className="px-4 py-2 bg-primary-600 text-white text-sm rounded-md hover:bg-primary-700 disabled:opacity-50 transition-colors"
            >
              {submitting ? 'در حال ذخیره...' : (isEditing ? 'به‌روزرسانی' : 'ایجاد')}
            </button>
            <button
              type="button"
              onClick={() => navigate('/roles')}
              className="px-4 py-2 bg-gray-100 text-gray-700 text-sm rounded-md hover:bg-gray-200 transition-colors"
            >
              انصراف
            </button>
          </div>
        </form>
      </div>
    </div>
  )
}

