import { useState, useEffect, useCallback } from 'react'
import {
  listMediaContent, getMediaContent, uploadMediaContent,
  updateMediaContent, deleteMediaContent,
} from '@/api/services/reportService'
import type { MediaListParams } from '@/api/services/reportService'
import Modal from '@/components/ui/Modal'
import ConfirmDialog from '@/components/ui/ConfirmDialog'
import DataTable from '@/components/ui/DataTable'
import { LoadingSpinner } from '@/components/ui/LoadingSpinner'
import { showSuccess, handleApiError } from '@/utils/errorHandler'
import type { Column } from '@/components/ui/DataTable'

interface MediaItem {
  id: string
  file_name: string
  file_type: string
  week_number: number
  is_active: boolean
  file_url?: string
  created_at: string
}

export default function MediaContentPage() {
  const [mediaList, setMediaList] = useState<MediaItem[]>([])
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(20)
  const [total, setTotal] = useState(0)
  const [pages, setPages] = useState(1)

  const [uploadModalOpen, setUploadModalOpen] = useState(false)
  const [editModalOpen, setEditModalOpen] = useState(false)
  const [deleteDialogOpen, setDeleteDialogOpen] = useState(false)

  const [selectedMedia, setSelectedMedia] = useState<MediaItem | null>(null)
  const [uploadFile, setUploadFile] = useState<File | null>(null)
  const [uploadWeekNumber, setUploadWeekNumber] = useState<number>(1)
  const [uploadFileType, setUploadFileType] = useState<string>('')
  const [uploadDescription, setUploadDescription] = useState<string>('')
  const [editWeekNumber, setEditWeekNumber] = useState<number>(1)
  const [editIsActive, setEditIsActive] = useState(true)
  const [submitLoading, setSubmitLoading] = useState(false)

  const fetchMedia = useCallback(async () => {
    setLoading(true)
    try {
      const params: MediaListParams = { page, page_size: pageSize }
      const response = await listMediaContent(params)
      const paginatedData = response.data
      setMediaList(paginatedData.data || [])
      setTotal(paginatedData.total || 0)
      setPages(paginatedData.pages || 1)
    } catch (err: any) {
      handleApiError(err, 'خطا در دریافت محتوای رسانه‌ای')
    } finally {
      setLoading(false)
    }
  }, [page, pageSize])

  useEffect(() => {
    fetchMedia()
  }, [fetchMedia])

  const handleUpload = async () => {
    if (!uploadFile) return
    setSubmitLoading(true)
    try {
      const formData = new FormData()
      formData.append('file', uploadFile)
      formData.append('week_number', String(uploadWeekNumber))
      if (uploadFileType) {
        formData.append('file_type', uploadFileType)
      }
      if (uploadDescription) {
        formData.append('description', uploadDescription)
      }
      await uploadMediaContent(formData)
      showSuccess('فایل با موفقیت آپلود شد')
      setUploadModalOpen(false)
      setUploadFile(null)
      setUploadWeekNumber(1)
      setUploadFileType('')
      setUploadDescription('')
      fetchMedia()
    } catch (err: any) {
      handleApiError(err, 'خطا در آپلود فایل')
    } finally {
      setSubmitLoading(false)
    }
  }

  const handleEdit = async () => {
    if (!selectedMedia) return
    setSubmitLoading(true)
    try {
      await updateMediaContent(selectedMedia.id, {
        week_number: editWeekNumber,
        is_active: editIsActive,
      })
      showSuccess('محتوای رسانه‌ای با موفقیت به‌روزرسانی شد')
      setEditModalOpen(false)
      setSelectedMedia(null)
      fetchMedia()
    } catch (err: any) {
      handleApiError(err, 'خطا در به‌روزرسانی')
    } finally {
      setSubmitLoading(false)
    }
  }

  const handleDelete = async () => {
    if (!selectedMedia) return
    setSubmitLoading(true)
    try {
      await deleteMediaContent(selectedMedia.id)
      showSuccess('محتوای رسانه‌ای با موفقیت حذف شد')
      setDeleteDialogOpen(false)
      setSelectedMedia(null)
      fetchMedia()
    } catch (err: any) {
      handleApiError(err, 'خطا در حذف')
    } finally {
      setSubmitLoading(false)
    }
  }

  const openEditModal = async (item: MediaItem) => {
    try {
      const response = await getMediaContent(item.id)
      const media = response.data
      setSelectedMedia(media)
      setEditWeekNumber(media.week_number || 1)
      setEditIsActive(media.is_active ?? true)
      setEditModalOpen(true)
    } catch (err: any) {
      handleApiError(err, 'خطا در دریافت اطلاعات')
    }
  }

  const openDeleteDialog = (item: MediaItem) => {
    setSelectedMedia(item)
    setDeleteDialogOpen(true)
  }

  const columns: Column<MediaItem>[] = [
    {
      key: 'file_name',
      title: 'نام فایل',
      render: (item: MediaItem) => (
        <span className="text-sm text-gray-700 font-medium">{item.file_name}</span>
      ),
    },
    {
      key: 'file_type',
      title: 'نوع فایل',
      render: (item: MediaItem) => (
        <span className="text-xs bg-gray-100 text-gray-700 px-2 py-1 rounded">{item.file_type}</span>
      ),
    },
    {
      key: 'week_number',
      title: 'شماره هفته',
    },
    {
      key: 'is_active',
      title: 'فعال',
      render: (item: MediaItem) => (
        <span className={`text-xs font-medium px-2 py-1 rounded ${item.is_active ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
          }`}>
          {item.is_active ? 'فعال' : 'غیرفعال'}
        </span>
      ),
    },
    {
      key: 'created_at',
      title: 'تاریخ ایجاد',
      render: (item: MediaItem) => (
        <span className="text-xs text-gray-500">{new Date(item.created_at).toLocaleDateString('fa-IR')}</span>
      ),
    },
  ]

  const renderActions = (item: MediaItem) => (
    <div className="flex gap-2">
      {item.file_url && (
        <a
          href={item.file_url}
          target="_blank"
          rel="noopener noreferrer"
          className="text-xs px-3 py-1 bg-green-50 text-green-700 rounded-md hover:bg-green-100 transition-colors"
        >
          پیش‌نمایش
        </a>
      )}
      <button
        onClick={() => openEditModal(item)}
        className="text-xs px-3 py-1 bg-blue-50 text-blue-700 rounded-md hover:bg-blue-100 transition-colors"
      >
        ویرایش
      </button>
      <button
        onClick={() => openDeleteDialog(item)}
        className="text-xs px-3 py-1 bg-red-50 text-red-700 rounded-md hover:bg-red-100 transition-colors"
      >
        حذف
      </button>
    </div>
  )

  if (loading && mediaList.length === 0) {
    return (
      <div>
        <h1 className="text-3xl font-bold text-gray-900 mb-6">محتوای رسانه‌ای</h1>
        <div className="py-20">
          <LoadingSpinner size="lg" />
        </div>
      </div>
    )
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold text-gray-900">محتوای رسانه‌ای</h1>
        <button
          onClick={() => setUploadModalOpen(true)}
          className="px-4 py-2 bg-primary-600 text-white rounded-md hover:bg-primary-700 text-sm"
        >
          آپلود فایل جدید
        </button>
      </div>

      <div className="bg-white rounded-lg shadow">
        <DataTable
          columns={columns}
          data={mediaList}
          loading={loading}
          actions={renderActions}
          page={page}
          pageSize={pageSize}
          total={total}
          pages={pages}
          onPageChange={setPage}
          onPageSizeChange={size => { setPageSize(size); setPage(1) }}
          emptyMessage="محتوای رسانه‌ای یافت نشد"
        />
      </div>

      <Modal
        isOpen={uploadModalOpen}
        onClose={() => { setUploadModalOpen(false); setUploadFile(null); setUploadWeekNumber(1); setUploadFileType(''); setUploadDescription('') }}
        title="آپلود فایل جدید"
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-2">انتخاب فایل</label>
            <input
              type="file"
              onChange={e => setUploadFile(e.target.files?.[0] || null)}
              className="w-full text-sm text-gray-600 file:mr-4 file:py-2 file:px-4 file:rounded-md file:border-0 file:text-sm file:font-medium file:bg-primary-50 file:text-primary-700 hover:file:bg-primary-100"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">شماره هفته *</label>
            <input
              type="number"
              min={1}
              max={52}
              value={uploadWeekNumber}
              onChange={e => setUploadWeekNumber(Number(e.target.value))}
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
              placeholder="۱ تا ۵۲"
            />
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">نوع فایل</label>
            <select
              value={uploadFileType}
              onChange={e => setUploadFileType(e.target.value)}
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            >
              <option value="">خودکار (بر اساس پسوند)</option>
              <option value="audio">صوتی (Audio)</option>
              <option value="video">ویدیویی (Video)</option>
              <option value="image">تصویر (Image)</option>
              <option value="document">سند (Document)</option>
            </select>
            <p className="text-xs text-gray-500 mt-1">در صورت خالی بودن، نوع فایل به صورت خودکار تشخیص داده می‌شود</p>
          </div>
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">توضیحات</label>
            <textarea
              value={uploadDescription}
              onChange={e => setUploadDescription(e.target.value)}
              rows={3}
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent resize-none"
              placeholder="توضیحات اختیاری درباره فایل..."
            />
          </div>
          <div className="flex justify-end gap-3">
            <button
              onClick={() => { setUploadModalOpen(false); setUploadFile(null); setUploadWeekNumber(1); setUploadFileType(''); setUploadDescription('') }}
              className="px-4 py-2 text-sm text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
            >
              انصراف
            </button>
            <button
              onClick={handleUpload}
              disabled={!uploadFile || !uploadWeekNumber || submitLoading}
              className="px-4 py-2 text-sm text-white bg-primary-600 rounded-md hover:bg-primary-700 disabled:opacity-50"
            >
              {submitLoading ? 'در حال آپلود...' : 'آپلود'}
            </button>
          </div>
        </div>
      </Modal>

      <Modal
        isOpen={editModalOpen}
        onClose={() => { setEditModalOpen(false); setSelectedMedia(null) }}
        title="ویرایش محتوای رسانه‌ای"
      >
        <div className="space-y-4">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">شماره هفته</label>
            <input
              type="number"
              min={1}
              value={editWeekNumber}
              onChange={e => setEditWeekNumber(Number(e.target.value))}
              className="w-full px-3 py-2 text-sm border border-gray-300 rounded-md focus:ring-2 focus:ring-primary-500 focus:border-transparent"
            />
          </div>
          <div className="flex items-center gap-3">
            <label className="text-sm font-medium text-gray-700">فعال</label>
            <input
              type="checkbox"
              checked={editIsActive}
              onChange={e => setEditIsActive(e.target.checked)}
              className="h-4 w-4 text-primary-600 border-gray-300 rounded focus:ring-primary-500"
            />
          </div>
          <div className="flex justify-end gap-3">
            <button
              onClick={() => { setEditModalOpen(false); setSelectedMedia(null) }}
              className="px-4 py-2 text-sm text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200"
            >
              انصراف
            </button>
            <button
              onClick={handleEdit}
              disabled={submitLoading}
              className="px-4 py-2 text-sm text-white bg-primary-600 rounded-md hover:bg-primary-700 disabled:opacity-50"
            >
              {submitLoading ? 'در حال ذخیره...' : 'ذخیره'}
            </button>
          </div>
        </div>
      </Modal>

      <ConfirmDialog
        isOpen={deleteDialogOpen}
        onConfirm={handleDelete}
        onCancel={() => { setDeleteDialogOpen(false); setSelectedMedia(null) }}
        title="حذف محتوای رسانه‌ای"
        message={selectedMedia ? `آیا از حذف "${selectedMedia.file_name}" اطمینان دارید؟` : ''}
        confirmText="حذف"
        cancelText="انصراف"
        variant="danger"
        loading={submitLoading}
      />
    </div>
  )
}

