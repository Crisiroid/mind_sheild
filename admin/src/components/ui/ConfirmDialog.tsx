interface ConfirmDialogProps {
    isOpen: boolean
    onConfirm: () => void
    onCancel: () => void
    title: string
    message: string
    confirmText?: string
    cancelText?: string
    variant?: 'danger' | 'warning'
    loading?: boolean
}

export default function ConfirmDialog({
    isOpen,
    onConfirm,
    onCancel,
    title,
    message,
    confirmText = 'تأیید',
    cancelText = 'انصراف',
    variant = 'danger',
    loading = false,
}: ConfirmDialogProps) {
    if (!isOpen) return null

    return (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-4">
            <div className="bg-white rounded-lg shadow-xl w-full max-w-md p-6 animate-fade-in">
                <h3 className="text-lg font-semibold text-gray-900 mb-2">{title}</h3>
                <p className="text-gray-600 mb-6">{message}</p>
                <div className="flex justify-end gap-3">
                    <button
                        onClick={onCancel}
                        disabled={loading}
                        className="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 rounded-md hover:bg-gray-200 disabled:opacity-50"
                    >
                        {cancelText}
                    </button>
                    <button
                        onClick={onConfirm}
                        disabled={loading}
                        className={`px-4 py-2 text-sm font-medium text-white rounded-md disabled:opacity-50 ${variant === 'danger' ? 'bg-red-600 hover:bg-red-700' : 'bg-yellow-600 hover:bg-yellow-700'
                            }`}
                    >
                        {loading ? 'در حال انجام...' : confirmText}
                    </button>
                </div>
            </div>
        </div>
    )
}
