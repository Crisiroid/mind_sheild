import { useState, useEffect } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { getAdminProfile, updateAdminProfile } from '@/api/services/adminService'
import { changeAdminPassword } from '@/api/services/reportService'
import { useAuthStore } from '@/store/auth'
import { showSuccess, handleApiError } from '@/utils/errorHandler'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'

const profileSchema = z.object({
  email: z.string().email('ایمیل نامعتبر است').min(1, 'ایمیل الزامی است'),
  full_name: z.string().min(1, 'نام کامل الزامی است'),
})

const passwordSchema = z.object({
  old_password: z.string().min(1, 'رمز عبور فعلی الزامی است'),
  new_password: z.string().min(6, 'رمز عبور جدید باید حداقل ۶ کاراکتر باشد'),
  confirm_password: z.string().min(1, 'تکرار رمز عبور الزامی است'),
}).refine(data => data.new_password === data.confirm_password, {
  message: 'رمز عبور و تکرار آن یکسان نیستند',
  path: ['confirm_password'],
})

type ProfileFormData = z.infer<typeof profileSchema>
type PasswordFormData = z.infer<typeof passwordSchema>

export default function ProfilePage() {
  const [profileLoading, setProfileLoading] = useState(true)
  const [profileSubmitting, setProfileSubmitting] = useState(false)
  const [passwordSubmitting, setPasswordSubmitting] = useState(false)
  const updateProfile = useAuthStore(state => state.updateProfile)

  const profileForm = useForm<ProfileFormData>({
    resolver: zodResolver(profileSchema),
  })

  const passwordForm = useForm<PasswordFormData>({
    resolver: zodResolver(passwordSchema),
  })

  useEffect(() => {
    const loadProfile = async () => {
      setProfileLoading(true)
      try {
        const response = await getAdminProfile()
        const user = response.data
        profileForm.reset({
          email: user.email || '',
          full_name: user.full_name || '',
        })
      } catch (err: any) {
        handleApiError(err, 'خطا در دریافت پروفایل')
      } finally {
        setProfileLoading(false)
      }
    }
    loadProfile()
  }, [profileForm])

  const onProfileSubmit = async (data: ProfileFormData) => {
    setProfileSubmitting(true)
    try {
      await updateAdminProfile(data)
      updateProfile(data)
      showSuccess('پروفایل با موفقیت به‌روزرسانی شد')
    } catch (err: any) {
      handleApiError(err, 'خطا در به‌روزرسانی پروفایل')
    } finally {
      setProfileSubmitting(false)
    }
  }

  const onPasswordSubmit = async (data: PasswordFormData) => {
    setPasswordSubmitting(true)
    try {
      await changeAdminPassword({
        old_password: data.old_password,
        new_password: data.new_password,
      })
      showSuccess('رمز عبور با موفقیت تغییر کرد')
      passwordForm.reset()
    } catch (err: any) {
      handleApiError(err, 'خطا در تغییر رمز عبور')
    } finally {
      setPasswordSubmitting(false)
    }
  }

  if (profileLoading) {
    return (
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-6">تنظیمات پروفایل</h1>
        <div className="py-20">
          <LoadingSpinner size="lg" />
        </div>
      </div>
    )
  }

  return (
    <div>
      <h1 className="text-3xl font-bold text-gray-900 mb-6">تنظیمات پروفایل</h1>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">ویرایش پروفایل</h2>
          <form onSubmit={profileForm.handleSubmit(onProfileSubmit)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">ایمیل</label>
              <input
                {...profileForm.register('email')}
                type="email"
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                placeholder="example@mail.com"
              />
              {profileForm.formState.errors.email && (
                <p className="mt-1 text-sm text-red-600">{profileForm.formState.errors.email.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">نام کامل</label>
              <input
                {...profileForm.register('full_name')}
                type="text"
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                placeholder="نام کامل خود را وارد کنید"
              />
              {profileForm.formState.errors.full_name && (
                <p className="mt-1 text-sm text-red-600">{profileForm.formState.errors.full_name.message}</p>
              )}
            </div>

            <button
              type="submit"
              disabled={profileSubmitting}
              className="w-full bg-primary-600 text-white py-2 px-4 rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {profileSubmitting ? 'در حال ذخیره...' : 'ذخیره تغییرات'}
            </button>
          </form>
        </div>

        <div className="bg-white rounded-lg shadow p-6">
          <h2 className="text-lg font-semibold text-gray-900 mb-4">تغییر رمز عبور</h2>
          <form onSubmit={passwordForm.handleSubmit(onPasswordSubmit)} className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">رمز عبور فعلی</label>
              <input
                {...passwordForm.register('old_password')}
                type="password"
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                placeholder="رمز عبور فعلی"
              />
              {passwordForm.formState.errors.old_password && (
                <p className="mt-1 text-sm text-red-600">{passwordForm.formState.errors.old_password.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">رمز عبور جدید</label>
              <input
                {...passwordForm.register('new_password')}
                type="password"
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                placeholder="رمز عبور جدید"
              />
              {passwordForm.formState.errors.new_password && (
                <p className="mt-1 text-sm text-red-600">{passwordForm.formState.errors.new_password.message}</p>
              )}
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">تکرار رمز عبور جدید</label>
              <input
                {...passwordForm.register('confirm_password')}
                type="password"
                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                placeholder="تکرار رمز عبور جدید"
              />
              {passwordForm.formState.errors.confirm_password && (
                <p className="mt-1 text-sm text-red-600">{passwordForm.formState.errors.confirm_password.message}</p>
              )}
            </div>

            <button
              type="submit"
              disabled={passwordSubmitting}
              className="w-full bg-primary-600 text-white py-2 px-4 rounded-md hover:bg-primary-700 focus:outline-none focus:ring-2 focus:ring-primary-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
            >
              {passwordSubmitting ? 'در حال تغییر...' : 'تغییر رمز عبور'}
            </button>
          </form>
        </div>
      </div>
    </div>
  )
}

