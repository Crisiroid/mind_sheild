import { useState, useEffect, useCallback } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import DataTable from '@/components/ui/DataTable'
import Modal from '@/components/ui/Modal'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import { listUsers, createUser, updateUser, deleteUser } from '@/api/services/userService'
import { handleApiError, showSuccess } from '@/utils/errorHandler'
import type { User } from '@/types/user'
import type { Column } from '@/components/ui/DataTable'

const createUserSchema = z.object({
    phone_number: z
        .string()
        .min(10, 'شماره تلفن باید حداقل ۱۰ رقم باشد')
        .regex(/^[0-9]+$/, 'شماره تلفن باید فقط شامل ارقام باشد'),
    android_version: z.string().optional(),
    app_version: z.string().optional(),
})

const editUserSchema = z.object({
    phone_number: z
        .string()
        .min(10, 'شماره تلفن باید حداقل ۱۰ رقم باشد')
        .regex(/^[0-9]+$/, 'شماره تلفن باید فقط شامل ارقام باشد')
        .optional(),
    cloud_sync_enabled: z.boolean().optional(),
    do_not_disturb_enabled: z.boolean().optional(),
    dnd_start_time: z.string().optional(),
    dnd_end_time: z.string().optional(),
})

type CreateUserFormData = z.infer<typeof createUserSchema>
type EditUserFormData = z.infer<typeof editUserSchema>

export default function UserListPage() {
    const navigate = useNavigate()

    const [users, setUsers] = useState<User[]>([])
    const [total, setTotal] = useState(0)
    const [pages, setPages] = useState(1)

    const [loading, setLoading] = useState(true)
    const [page, setPage] = useState(1)
    const [pageSize] = useState(20)
    const [search, setSearch] = useState('')
    const [dateFrom, setDateFrom] = useState('')
    const [dateTo, setDateTo] = useState('')
    const [agreementFilter, setAgreementFilter] = useState('all')

    const [createModalOpen, setCreateModalOpen] = useState(false)
    const [editModalOpen, setEditModalOpen] = useState(false)
    const [editingUser, setEditingUser] = useState<User | null>(null)
    const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)
    const [deletingUser, setDeletingUser] = useState<User | null>(null)
    const [submitting, setSubmitting] = useState(false)

    const createForm = useForm<CreateUserFormData>({
        resolver: zodResolver(createUserSchema),
        defaultValues: { phone_number: '', android_version: '', app_version: '' },
    })

    const editForm = useForm<EditUserFormData>({
        resolver: zodResolver(editUserSchema),
        defaultValues: {
            phone_number: '',
            cloud_sync_enabled: false,
            do_not_disturb_enabled: false,
            dnd_start_time: '',
            dnd_end_time: '',
        },
    })

    const fetchUsers = useCallback(async () => {
        setLoading(true)
        try {
            const params: Record<string, unknown> = { page, page_size: pageSize }
            if (search.trim()) params.search = search.trim()
            if (dateFrom) params.date_from = dateFrom
            if (dateTo) params.date_to = dateTo
            if (agreementFilter !== 'all') params.agreement_accepted = agreementFilter === 'yes'

            const res = await listUsers(params as any)
            setUsers(res.data.data)
            setTotal(res.data.total)
            setPages(res.data.pages)
        } catch (err) {
            handleApiError(err)
        } finally {
            setLoading(false)
        }
    }, [page, pageSize, search, dateFrom, dateTo, agreementFilter])

    useEffect(() => {
        fetchUsers()
    }, [fetchUsers])

    const handleCreate = async (data: CreateUserFormData) => {
        setSubmitting(true)
        try {
            await createUser({
                phone_number: data.phone_number,
                android_version: data.android_version || undefined,
                app_version: data.app_version || undefined,
            })
            showSuccess('کاربر با موفقیت ایجاد شد')
            setCreateModalOpen(false)
            createForm.reset()
            setPage(1)
            fetchUsers()
        } catch (err) {
            handleApiError(err)
        } finally {
            setSubmitting(false)
        }
    }

    const openEditModal = (user: User) => {
        setEditingUser(user)
        editForm.reset({
            phone_number: user.phone_number,
            cloud_sync_enabled: user.cloud_sync_enabled,
            do_not_disturb_enabled: user.do_not_disturb_enabled,
            dnd_start_time: user.dnd_start_time || '',
            dnd_end_time: user.dnd_end_time || '',
        })
        setEditModalOpen(true)
    }

    const handleEdit = async (data: EditUserFormData) => {
        if (!editingUser) return
        setSubmitting(true)
        try {
            await updateUser(editingUser.id, {
                phone_number: data.phone_number || undefined,
                cloud_sync_enabled: data.cloud_sync_enabled,
                do_not_disturb_enabled: data.do_not_disturb_enabled,
                dnd_start_time: data.dnd_start_time || undefined,
                dnd_end_time: data.dnd_end_time || undefined,
            })
            showSuccess('کاربر با موفقیت ویرایش شد')
            setEditModalOpen(false)
            setEditingUser(null)
            fetchUsers()
        } catch (err) {
            handleApiError(err)
        } finally {
            setSubmitting(false)
        }
    }

    const openDeleteDialog = (user: User) => {
        setDeletingUser(user)
        setDeleteDialogOpen(true)
    }

    const handleDelete = async () => {
        if (!deletingUser) return
        setSubmitting(true)
        try {
            await deleteUser(deletingUser.id)
            showSuccess('کاربر با موفقیت حذف شد')
            setDeleteDialogOpen(false)
            setDeletingUser(null)
            fetchUsers()
        } catch (err) {
            handleApiError(err)
        } finally {
            setSubmitting(false)
        }
    }

    const resetFilters = () => {
        setSearch('')
        setDateFrom('')
        setDateTo('')
        setAgreementFilter('all')
        setPage(1)
    }

    const columns: Column<User>[] = [
        { key: 'phone_number', title: 'شماره تلفن' },
        {
            key: 'registration_date',
            title: 'تاریخ ثبت‌نام',
            render: (item: User) => (
                <span dir="ltr" className="text-xs">{new Date(item.registration_date).toLocaleDateString('fa-IR')}</span>
            ),
        },
        {
            key: 'last_login',
            title: 'آخرین ورود',
            render: (item: User) => (
                <span dir="ltr" className="text-xs">{new Date(item.last_login).toLocaleDateString('fa-IR')}</span>
            ),
        },
        { key: 'login_count', title: 'تعداد ورود' },
        {
            key: 'agreement_accepted',
            title: 'توافق',
            render: (item: User) =>
                item.agreement_accepted ? (
                    <span className="inline-flex px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">بلی</span>
                ) : (
                    <span className="inline-flex px-2 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">خیر</span>
                ),
        },
    ]

    return (
        <div className="space-y-4 md:space-y-6">
            <div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-3 md:gap-4">
                <h1 className="text-2xl md:text-3xl font-bold text-gray-900">مدیریت کاربران</h1>
                <button
                    onClick={() => {
                        createForm.reset({ phone_number: '', android_version: '', app_version: '' })
                        setCreateModalOpen(true)
                    }}
                    className="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 transition-colors text-sm w-full sm:w-auto"
                >
                    + ایجاد کاربر جدید
                </button>
            </div>

            <div className="bg-white rounded-lg shadow p-3 md:p-4">
                <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-5 gap-3 md:gap-4">
                    <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">جستجوی شماره تلفن</label>
                        <input
                            type="text"
                            value={search}
                            onChange={(e) => { setSearch(e.target.value); setPage(1) }}
                            placeholder="جستجو..."
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">از تاریخ</label>
                        <input
                            type="date"
                            value={dateFrom}
                            onChange={(e) => { setDateFrom(e.target.value); setPage(1) }}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">تا تاریخ</label>
                        <input
                            type="date"
                            value={dateTo}
                            onChange={(e) => { setDateTo(e.target.value); setPage(1) }}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        />
                    </div>
                    <div>
                        <label className="block text-xs font-medium text-gray-600 mb-1">وضعیت توافق</label>
                        <select
                            value={agreementFilter}
                            onChange={(e) => { setAgreementFilter(e.target.value); setPage(1) }}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent bg-white"
                        >
                            <option value="all">همه</option>
                            <option value="yes">قبول شده</option>
                            <option value="no">رد شده</option>
                        </select>
                    </div>
                    <div className="flex items-end">
                        <button
                            onClick={resetFilters}
                            className="w-full px-3 py-2 text-sm text-gray-600 border border-gray-300 rounded-md hover:bg-gray-50 transition-colors"
                        >
                            پاک کردن فیلترها
                        </button>
                    </div>
                </div>
            </div>

            <div className="bg-white rounded-lg shadow overflow-hidden">
                <DataTable<User>
                    columns={columns}
                    data={users}
                    loading={loading}
                    page={page}
                    pageSize={pageSize}
                    total={total}
                    pages={pages}
                    onPageChange={setPage}
                    onRowClick={(item) => navigate(`/users/${item.id}`)}
                    actions={(item: User) => (
                        <div className="flex gap-2">
                            <button
                                onClick={(e) => { e.stopPropagation(); openEditModal(item) }}
                                className="px-2 md:px-3 py-1 text-xs bg-blue-50 text-blue-700 rounded-md hover:bg-blue-100 transition-colors"
                            >
                                ویرایش
                            </button>
                            <button
                                onClick={(e) => { e.stopPropagation(); openDeleteDialog(item) }}
                                className="px-2 md:px-3 py-1 text-xs bg-red-50 text-red-700 rounded-md hover:bg-red-100 transition-colors"
                            >
                                حذف
                            </button>
                        </div>
                    )}
                />
            </div>

            <Modal isOpen={createModalOpen} onClose={() => setCreateModalOpen(false)} title="ایجاد کاربر جدید" maxWidth="max-w-lg">
                <form onSubmit={createForm.handleSubmit(handleCreate)} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">شماره تلفن *</label>
                        <input
                            {...createForm.register('phone_number')}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                            placeholder="مثال: 09123456789"
                        />
                        {createForm.formState.errors.phone_number && (
                            <p className="mt-1 text-xs text-red-600">{createForm.formState.errors.phone_number.message}</p>
                        )}
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">نسخه اپلیکیشن</label>
                        <input
                            {...createForm.register('app_version')}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                            placeholder="مثال: 1.0.0"
                        />
                    </div>
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">نسخه اندروید</label>
                        <input
                            {...createForm.register('android_version')}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                            placeholder="مثال: 13"
                        />
                    </div>
                    <div className="flex flex-col sm:flex-row justify-end gap-3 pt-2">
                        <button
                            type="button"
                            onClick={() => setCreateModalOpen(false)}
                            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 w-full sm:w-auto"
                        >
                            انصراف
                        </button>
                        <button
                            type="submit"
                            disabled={submitting}
                            className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700 disabled:opacity-50 w-full sm:w-auto"
                        >
                            {submitting ? 'در حال ایجاد...' : 'ایجاد کاربر'}
                        </button>
                    </div>
                </form>
            </Modal>

            <Modal isOpen={editModalOpen} onClose={() => { setEditModalOpen(false); setEditingUser(null) }} title="ویرایش کاربر" maxWidth="max-w-lg">
                <form onSubmit={editForm.handleSubmit(handleEdit)} className="space-y-4">
                    <div>
                        <label className="block text-sm font-medium text-gray-700 mb-1">شماره تلفن</label>
                        <input
                            {...editForm.register('phone_number')}
                            className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                        />
                        {editForm.formState.errors.phone_number && (
                            <p className="mt-1 text-xs text-red-600">{editForm.formState.errors.phone_number.message}</p>
                        )}
                    </div>
                    <div className="flex items-center gap-3">
                        <input
                            type="checkbox"
                            id="cloud_sync"
                            {...editForm.register('cloud_sync_enabled')}
                            className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                        />
                        <label htmlFor="cloud_sync" className="text-sm text-gray-700">همگام‌سازی ابری</label>
                    </div>
                    <div className="flex items-center gap-3">
                        <input
                            type="checkbox"
                            id="dnd"
                            {...editForm.register('do_not_disturb_enabled')}
                            className="rounded border-gray-300 text-primary-600 focus:ring-primary-500"
                        />
                        <label htmlFor="dnd" className="text-sm text-gray-700">حالت مزاحم نشوید</label>
                    </div>
                    <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">زمان شروع DND</label>
                            <input
                                type="time"
                                {...editForm.register('dnd_start_time')}
                                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                            />
                        </div>
                        <div>
                            <label className="block text-sm font-medium text-gray-700 mb-1">زمان پایان DND</label>
                            <input
                                type="time"
                                {...editForm.register('dnd_end_time')}
                                className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
                            />
                        </div>
                    </div>
                    <div className="flex flex-col sm:flex-row justify-end gap-3 pt-2">
                        <button
                            type="button"
                            onClick={() => { setEditModalOpen(false); setEditingUser(null) }}
                            className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 w-full sm:w-auto"
                        >
                            انصراف
                        </button>
                        <button
                            type="submit"
                            disabled={submitting}
                            className="px-4 py-2 text-sm font-medium text-white bg-primary-600 rounded-md hover:bg-primary-700 disabled:opacity-50 w-full sm:w-auto"
                        >
                            {submitting ? 'در حال ویرایش...' : 'ویرایش'}
                        </button>
                    </div>
                </form>
            </Modal>

            <ConfirmDialog
                isOpen={deleteDialogOpen}
                onConfirm={handleDelete}
                onCancel={() => { setDeleteDialogOpen(false); setDeletingUser(null) }}
                title="حذف کاربر"
                message={
                    deletingUser
                        ? `آیا از حذف کاربر ${deletingUser.phone_number} اطمینان دارید؟ این عمل قابل بازگشت نیست.`
                        : 'آیا از حذف این کاربر اطمینان دارید؟'
                }
                confirmText="حذف"
                cancelText="انصراف"
                variant="danger"
                loading={submitting}
            />
        </div>
    )
}

