import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { PageLoading } from '@/components/ui/LoadingSpinner'
import { getAdminUser, createAdminUser, updateAdminUser, listAdminRoles } from '@/api/services/adminService'
import { showSuccess } from '@/utils/errorHandler'
import type { AdminRole } from '@/types/admin'

const formSchema = z.object({
  username: z.string().min(3, 'نام کاربری حداقل ۳ کاراکتر'),
  email: z.string().email('ایمیل معتبر وارد کنید'),
  password: z.string().optional(),
  full_name: z.string().min(1, 'نام کامل الزامی است'),
  role_id: z.string().min(1, 'نقش الزامی است'),
  is_active: z.boolean(),
})

type FormValues = z.infer<typeof formSchema>

export default function AdminFormPage() {
  const { id } = useParams<{ id: string }>()
  const navigate = useNavigate()
  const isEditing = !!id
  const [loading, setLoading] = useState(isEditing)
  const [submitting, setSubmitting] = useState(false)
  const [roles, setRoles] = useState<AdminRole[]>([])

  const {
    register,
    handleSubmit,
    formState: { errors },
    reset,
    setError,
  } = useForm<FormValues>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      username: '',
      email: '',
      password: '',
      full_name: '',
      role_id: '',
      is_active: true,
    },
  })

  useEffect(() => {
    listAdminRoles({ page: 1, page_size: 100 })
      .then(res => setRoles(res.data.data))
      .catch(() => { })
  }, [])

  useEffect(() => {
    if (isEditing && id) {
      setLoading(true)
      getAdminUser(id)
        .then(res => {
          const admin = res.data
          reset({
            username: admin.username,
            email: admin.email,
            password: '',
            full_name: admin.full_name,
            role_id: admin.role_id,
            is_active: admin.is_active,
          })
        })
        .catch(() => navigate('/admins'))
        .finally(() => setLoading(false))
    }
  }, [id, isEditing, reset, navigate])

  const onSubmit = async (values: FormValues) => {
    if (!isEditing && (!values.password || values.password.length < 6)) {
      setError('password', { message: 'رمز عبور حداقل ۶ کاراکتر' })
      return
    }

    setSubmitting(true)
    try {
      if (isEditing && id) {
        await updateAdminUser(id, {
          email: values.email,
          full_name: values.full_name,
          role_id: values.role_id,
          is_active: values.is_active,
        })
        showSuccess('ادمین با موفقیت به‌روزرسانی شد')
      } else {
        await createAdminUser({
          username: values.username,
          email: values.email,
          password: values.password!,
          full_name: values.full_name,
          role_id: values.role_id,
          is_active: values.is_active,
        })
        showSuccess('ادمین با موفقیت ایجاد شد')
      }
      navigate('/admins')
    } catch {
    } finally {
      setSubmitting(false)
    }
  }

  if (loading) return <PageLoading />

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">
        {isEditing ? 'ویرایش ادمین' : 'ایجاد ادمین جدید'}
      </h1>
      <div className="bg-white rounded-lg shadow p-6 max-w-2xl">
        <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">نام کاربری</label>
            <input
              {...register('username')}
              disabled={isEditing}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm disabled:bg-gray-100 focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            {errors.username && <p className="text-red-500 text-xs mt-1">{errors.username.message}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">ایمیل</label>
            <input
              {...register('email')}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            {errors.email && <p className="text-red-500 text-xs mt-1">{errors.email.message}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {isEditing ? 'رمز عبور (خالی بگذارید اگر تغییری نمی‌خواهید)' : 'رمز عبور'}
            </label>
            <input
              type="password"
              {...register('password')}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            {errors.password && <p className="text-red-500 text-xs mt-1">{errors.password.message}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">نام کامل</label>
            <input
              {...register('full_name')}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
            {errors.full_name && <p className="text-red-500 text-xs mt-1">{errors.full_name.message}</p>}
          </div>

          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">نقش</label>
            <select
              {...register('role_id')}
              className="w-full px-3 py-2 border border-gray-300 rounded-md text-sm focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            >
              <option value="">انتخاب نقش</option>
              {roles.map(role => (
                <option key={role.id} value={role.id}>{role.role_name}</option>
              ))}
            </select>
            {errors.role_id && <p className="text-red-500 text-xs mt-1">{errors.role_id.message}</p>}
          </div>

          <div className="flex items-center gap-2">
            <input
              type="checkbox"
              {...register('is_active')}
              id="is_active"
              className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
            />
            <label htmlFor="is_active" className="text-sm text-gray-700">فعال</label>
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
              onClick={() => navigate('/admins')}
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

