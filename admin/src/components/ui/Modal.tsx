import { useEffect, useRef } from 'react'

interface ModalProps {
    isOpen: boolean
    onClose: () => void
    title: string
    children: React.ReactNode
    maxWidth?: string
}

export default function Modal({ isOpen, onClose, title, children, maxWidth = 'max-w-2xl' }: ModalProps) {
    const overlayRef = useRef<HTMLDivElement>(null)

    useEffect(() => {
        const handleEscape = (e: KeyboardEvent) => {
            if (e.key === 'Escape') onClose()
        }
        if (isOpen) {
            document.addEventListener('keydown', handleEscape)
            document.body.style.overflow = 'hidden'
        }
        return () => {
            document.removeEventListener('keydown', handleEscape)
            document.body.style.overflow = 'unset'
        }
    }, [isOpen, onClose])

    if (!isOpen) return null

    return (
        <div
            ref={overlayRef}
            className="fixed inset-0 z-50 flex items-center justify-center bg-black/50 p-3 md:p-4"
            onClick={e => e.target === overlayRef.current && onClose()}
        >
            <div className={`bg-white rounded-lg shadow-xl w-full ${maxWidth} max-h-[85vh] md:max-h-[90vh] flex flex-col animate-fade-in`}>
                <div className="flex items-center justify-between px-4 md:px-6 py-3 md:py-4 border-b">
                    <h2 className="text-base md:text-lg font-semibold text-gray-900">{title}</h2>
                    <button onClick={onClose} className="text-gray-400 hover:text-gray-600 text-xl md:text-2xl leading-none p-1">&times;</button>
                </div>
                <div className="flex-1 overflow-y-auto px-4 md:px-6 py-3 md:py-4">
                    {children}
                </div>
            </div>
        </div>
    )
}
