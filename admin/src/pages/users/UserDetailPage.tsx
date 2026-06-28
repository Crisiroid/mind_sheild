import { useState, useEffect } from 'react'
import { useParams, useNavigate } from 'react-router-dom'
import { getUser, updateUser, deleteUser } from '@/api/services/userService'
import { PageLoading } from '@/components/ui/LoadingSpinner'
import Modal from '@/components/ui/Modal'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import { handleApiError, showSuccess } from '@/utils/errorHandler'
import type { User } from '@/types/user'

export default function UserDetailPage() {
    const { id } = useParams<{ id: string }>()
    const navigate = useNavigate()

    const [user, setUser] = useState<User | null>(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState<string | null>(null)

    const [editModalOpen, setEditModalOpen] = useState(false)
    const [editFormData, setEditFormData] = useState({
        phone_number: '',
        cloud_sync_enabled: false,
        do_not_disturb_enabled: false,
        dnd_start_time: '',
        dnd_end_time: '',
    })
    const [submitting, setSubmitting] = useState(false)

    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
    const [deleting, setDeleting] = useState(false)

    const fetchUser = async () => {
        if (!id) return
        setLoading(true)
        setError(null)
        try {
            const res = await getUser(id)
            setUser(res.data)
            setEditFormData({
                phone_number: res.data.phone_number,
                cloud_sync_enabled: res.data.cloud_sync_enabled,
                do_not_disturb_enabled: res.data.do_not_disturb_enabled,
                dnd_start_time: res.data.dnd_start_time || '',
                dnd_end_time: res.data.dnd_end_time || '',
            })
        } catch (err) {
            handleApiError(err)
            setError('خطا در بارگذاری اطلاعات کاربر')
        } finally {
            setLoading(false)
        }
    }

    useEffect(() => {
        fetchUser()
    }, [id])

    const handleEdit = async () => {
        if (!user) return
        setSubmitting(true)
        try {
            const res = await updateUser(user.id, {
                phone_number: editFormData.phone_number || undefined,
                cloud_sync_enabled: editFormData.cloud_sync_enabled,
                do_not_disturb_enabled: editFormData.do_not_disturb_enabled,
                dnd_start_time: editFormData.dnd_start_time || undefined,
                dnd_end_time: editFormData.dnd_end_time || undefined,
            })
            setUser(res.data)
            setEditModalOpen(false)
            showSuccess('کاربر با موفقیت ویرایش شد')
        } catch (err) {
            handleApiError(err)
        } finally {
            setSubmitting(false)
        }
    }

    const handleDelete = async () => {
        if (!user) return
        setDeleting(true)
        try {
            await deleteUser(user.id)
            showSuccess('کاربر با موفقیت حذف شد')
            navigate('/users')
        } catch (err) {
            handleApiError(err)
        } finally {
            setDeleting(false)
        }
    }

    if (loading) return <PageLoading />

    if (error || !user) {
        return (
            <div>
                <h1 className="text-3xl font-bold text-gray-900 mb-6">جزئیات کاربر</h1>
                <div className="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
                    <p className="text-red-600 mb-4">{error || 'کاربر یافت نشد'}</p>
                    <button
                        onClick={() => navigate('/users')}
                        className="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700 transition-colors"
                    >
                        بازگشت به لیست
                    </button>
                </div>
            </div>
        )
    }

    const InfoRow = ({ label, value }: { label: string; value: string | number | boolean | null | undefined }) => {
        let displayValue = '-'
        if (typeof value === 'boolean') {
            displayValue = value ? 'فعال' : 'غیرفعال'
        } else if (value !== null && value !== undefined && value !== '') {
            displayValue = String(value)
        }

        return (
            <div className="flex items-center justify-between py-3 border-b border-gray-100 last:border-b-0">
                <span className="text-sm font-medium text-gray-600">{label}</span>
                <span className="text-sm text-gray-900">{displayValue}</span>
            </div>
        )
    }

    return (
        <div>
            <div className="flex items-center justify-between mb-6">
                <div className="flex items-center gap-4">
                    <button
                        onClick={() => navigate('/users')}
                        className="px-3 py-2 text-sm text-gray-600 bg-gray-100 rounded-md hover:bg-gray-200 transition-colors"
                    >
                        ← بازگشت به لیست
                    </button>
                    <h1 className="text-3xl font-bold text-gray-900">جزئیات کاربر</h1>
                </div>
                <div className="flex gap-3">
                    <button
                        onClick={() => setEditModalOpen(true)}
                        className="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 transition-colors text-sm"
                    >
                        ویرایش کاربر
                    </button>
                    <button
                        onClick={() => setDeleteDialogOpen(true)}
                        className="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 transition-colors text-sm"
                    >
                        حذف کاربر
                    </button>
                </div>
            </div>

            <div className="bg-white rounded-lg shadow">
                <div className="px-6 py-4 border-b border-gray-200">
                    <h2 className="text-lg font-semibold text-gray-900">{user.phone_number}</h2>
                    <p className="text-sm text-gray-500 mt-1">شناسه: {user.id}</p>
                </div>

                <div className="px-6 py-4 grid grid-cols-1 md:grid-cols-2 gap-x-8">
                    <div>
                        <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">اطلاعات پایه</h3>
                        <InfoRow label="شماره تلفن" value={user.phone_number} />
                        <InfoRow
                            label="تاریخ ثبت‌نام"
                            value={new Date(user.registration_date).toLocaleDateString('fa-IR')}
                        />
                        <InfoRow
                            label="آخرین ورود"
                            value={new Date(user.last_login).toLocaleDateString('fa-IR')}
                        />
                        <InfoRow label="تعداد ورود" value={user.login_count} />
                        <InfoRow label="وضعیت توافق" value={user.agreement_accepted ? 'بلی' : 'خیر'} />
                        <InfoRow
                            label="تاریخ توافق"
                            value={
                                user.agreement_accepted_at
                                    ? new Date(user.agreement_accepted_at).toLocaleDateString('fa-IR')
                                    : '-'
                            }
                        />
                    </div>
                    <div>
                        <h3 className="text-xs font-semibold text-gray-500 uppercase tracking-wider mb-2">تنظیمات و دستگاه</h3>
                        <InfoRow label="همگام‌سازی ابری" value={user.cloud_sync_enabled} />
                        <InfoRow label="حالت مزاحم نشوید" value={user.do_not_disturb_enabled} />
                        <InfoRow label="زمان شروع DND" value={user.dnd_start_time || '-'} />
                        <InfoRow label="زمان پایان DND" value={user.dnd_end_time || '-'} />
                        <InfoRow label="نسخه اپلیکیشن" value={user.app_version || '-'} />
                        <InfoRow label="نسخه اندروید" value={user.android_version || '-'} />
                    </div>
                </div>

                <div className="px-6 py-4 border-t border-gray-200 bg-gray-50 rounded-b-lg">
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-x-8">
                        <InfoRow
                            label="تاریخ ایجاد"
                            value={new Date(user.created_at).toLocaleDateString('fa-IR')}
                        />
                        <InfoRow
                            label="آخرین به‌روزرسانی"
                            value={new Date(user.updated_at).toLocaleDateString('fa-IR')}
                        />
                    </div>
                </div>
            </div>

            <Modal isOpen={editModalOpen} onClose={() => setEditModalOpen(false)} title="ویرایش کاربر">
                <div className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">شماره تلفن</label>
                        <input
                            type="text"
                            value={editFormData.phone_number}
                            onChange={(e) => setEditFormData((prev) => ({ ...prev, phone_number: e.target.value }))}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        />
                    </div>
                    <div className="flex items-center gap-3">
                        <input
                            type="checkbox"
                            id="detail-cloud-sync"
                            checked={editFormData.cloud_sync_enabled}
                            onChange={(e) =>
                                setEditFormData((prev) => ({ ...prev, cloud_sync_enabled: e.target.checked }))
                            }
                            className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                        />
                        <label htmlFor="detail-cloud-sync" className="text-sm text-gray-700">همگام‌سازی ابری</label>
                    </div>
                    <div className="flex items-center gap-3">
                        <input
                            type="checkbox"
                            id="detail-dnd"
                            checked={editFormData.do_not_disturb_enabled}
                            onChange={(e) =>
                                setEditFormData((prev) => ({ ...prev, do_not_disturb_enabled: e.target.checked }))
                            }
                            className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                        />
                        <label htmlFor="detail-dnd" className="text-sm text-gray-700">حالت مزاحم نشوید</label>
                    </div>
                    <div className="grid grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">زمان شروع DND</label>
                            <input
                                type="time"
                                value={editFormData.dnd_start_time}
                                onChange={(e) =>
                                    setEditFormData((prev) => ({ ...prev, dnd_start_time: e.target.value }))
                                }
                                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">زمان پایان DND</label>
                            <input
                                type="time"
                                value={editFormData.dnd_end_time}
                                onChange={(e) =>
                                    setEditFormData((prev) => ({ ...prev, dnd_end_time: e.target.value }))
                                }
                                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                            />
                        </div>
                    </div>
                    <div className="flex justify-end gap-3 pt-2">
                        <button
                            type="button"
                            onClick={() => setEditModalOpen(false)}
                            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
                        >
                            انصراف
                        </button>
                        <button
                            type="button"
                            onClick={handleEdit}
                            disabled={submitting}
                            className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700 disabled:opacity-50"
                        >
                            {submitting ? 'در حال ویرایش...' : 'ذخیره تغییرات'}
                        </button>
                    </div>
                </div>
            </Modal>

            <ConfirmDialog
                isOpen={deleteDialogOpen}
                onConfirm={handleDelete}
                onCancel={() => setDeleteDialogOpen(false)}
                title="حذف کاربر"
                message={`آیا از حذف کاربر ${user.phone_number} اطمینان دارید؟ این عمل قابل بازگشت نیست.`}
                confirmText="حذف"
                cancelText="انصراف"
                variant="danger"
                loading={deleting}
            />
        </div>
    )
}

