import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { useAuthStore } from '@/store/auth'
import { toast } from 'sonner'
import { AxiosError } from 'axios'

const loginSchema = z.object({
    username: z.string().min(1, 'نام کاربری الزامی است'),
    password: z.string().min(1, 'رمز عبور الزامی است'),
})

type LoginFormData = z.infer<typeof loginSchema>

const getLoginErrorMessage = (error: unknown): string => {
    const axiosError = error as AxiosError<{ message?: string; error?: string }>
    if (axiosError.response?.data?.message) return axiosError.response.data.message
    if (axiosError.response?.data?.error) return axiosError.response.data.error
    if (axiosError.response?.status === 401) return 'نام کاربری یا رمز عبور اشتباه است'
    if (axiosError.response?.status === 403) return 'حساب شما مسدود شده است'
    if (axiosError.response?.status === 429) return 'تعداد تلاش‌ها بیش از حد است، لطفاً بعداً دوباره تلاش کنید'
    if (axiosError.code === 'ECONNABORTED') return 'اتصال به سرور با تأخیر مواجه شد'
    if (!axiosError.response) return 'اتصال به سرور برقرار نشد، اینترنت خود را بررسی کنید'
    return 'خطایی رخ داده است، لطفاً دوباره تلاش کنید'
}

export default function LoginPage() {
    const navigate = useNavigate()
    const { login, isLoading } = useAuthStore()
    const [showPassword, setShowPassword] = useState(false)

    const {
        register,
        handleSubmit,
        formState: { errors },
    } = useForm<LoginFormData>({
        resolver: zodResolver(loginSchema),
    })

    const onSubmit = async (data: LoginFormData) => {
        try {
            await login(data)
            navigate('/')
        } catch (error) {
            const msg = getLoginErrorMessage(error)
            toast.error(msg)
        }
    }

    return (
        <div className="login-page">
            {/* Background decoration */}
            <div className="login-bg">
                <div className="login-bg-orb login-bg-orb-1" />
                <div className="login-bg-orb login-bg-orb-2" />
                <div className="login-bg-orb login-bg-orb-3" />
            </div>

            <div className="login-container">
                {/* Left: Branding panel */}
                <div className="login-brand">
                    <div className="login-brand-content">
                        <div className="login-brand-icon">
                            <svg width="40" height="40" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                                <path d="M12 2a7 7 0 0 0-7 7c0 3 2 5.5 4 7.5L12 22l3-5.5c2-2 4-4.5 4-7.5a7 7 0 0 0-7-7z" />
                                <circle cx="9" cy="10" r="1" fill="currentColor" />
                                <circle cx="15" cy="10" r="1" fill="currentColor" />
                                <path d="M9.5 14c.83.67 2.17.67 3 0" />
                            </svg>
                        </div>
                        <h1 className="login-brand-title">پنل مدیریت</h1>
                        <p className="login-brand-subtitle">Psychology App</p>

                        <div className="login-brand-features">
                            <div className="login-feature">
                                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                                    <polyline points="20 6 9 17 4 12" />
                                </svg>
                                <span>مدیریت کاربران و ادمین‌ها</span>
                            </div>
                            <div className="login-feature">
                                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                                    <polyline points="20 6 9 17 4 12" />
                                </svg>
                                <span>گزارش‌ها و تحلیل‌های جامع</span>
                            </div>
                            <div className="login-feature">
                                <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                                    <polyline points="20 6 9 17 4 12" />
                                </svg>
                                <span>مدیریت داده‌های روانشناختی</span>
                            </div>
                        </div>
                    </div>
                </div>

                {/* Right: Login form */}
                <div className="login-form-wrapper">
                    <div className="login-form-card">
                        {/* Mobile logo */}
                        <div className="login-mobile-logo">
                            <div className="login-mobile-icon">
                                <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round">
                                    <path d="M12 2a7 7 0 0 0-7 7c0 3 2 5.5 4 7.5L12 22l3-5.5c2-2 4-4.5 4-7.5a7 7 0 0 0-7-7z" />
                                    <circle cx="9" cy="10" r="1" fill="currentColor" />
                                    <circle cx="15" cy="10" r="1" fill="currentColor" />
                                    <path d="M9.5 14c.83.67 2.17.67 3 0" />
                                </svg>
                            </div>
                            <h1 className="login-mobile-title">پنل مدیریت</h1>
                        </div>

                        <div className="login-form-header">
                            <h2 className="login-form-title">ورود به حساب</h2>
                            <p className="login-form-desc">برای ادامه، اطلاعات خود را وارد کنید</p>
                        </div>

                        <form onSubmit={handleSubmit(onSubmit)} className="login-form">
                            {/* Username */}
                            <div className="login-field">
                                <label className="login-label">نام کاربری</label>
                                <div className="login-input-wrap">
                                    <div className="login-input-icon" style={{ color: errors.username ? '#ef4444' : undefined }}>
                                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
                                            <circle cx="12" cy="7" r="4" />
                                        </svg>
                                    </div>
                                    <input
                                        {...register('username')}
                                        type="text"
                                        className={`login-input ${errors.username ? 'login-input-error' : ''}`}
                                        placeholder="username"
                                        autoComplete="username"
                                    />
                                </div>
                                {errors.username && (
                                    <p className="login-error-msg">
                                        <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                                            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z" />
                                        </svg>
                                        {errors.username.message}
                                    </p>
                                )}
                            </div>

                            {/* Password */}
                            <div className="login-field">
                                <label className="login-label">رمز عبور</label>
                                <div className="login-input-wrap">
                                    <div className="login-input-icon" style={{ color: errors.password ? '#ef4444' : undefined }}>
                                        <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                            <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                                            <path d="M7 11V7a5 5 0 0 1 10 0v4" />
                                        </svg>
                                    </div>
                                    <input
                                        {...register('password')}
                                        type={showPassword ? 'text' : 'password'}
                                        className={`login-input ${errors.password ? 'login-input-error' : ''}`}
                                        placeholder="••••••••"
                                        autoComplete="current-password"
                                    />
                                    <button
                                        type="button"
                                        className="login-toggle-pwd"
                                        onClick={() => setShowPassword(!showPassword)}
                                    >
                                        {showPassword ? (
                                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                <path d="M17.94 17.94A10.07 10.07 0 0 1 12 20c-7 0-11-8-11-8a18.45 18.45 0 0 1 5.06-5.94M9.9 4.24A9.12 9.12 0 0 1 12 4c7 0 11 8 11 8a18.5 18.5 0 0 1-2.16 3.19m-6.72-1.07a3 3 0 1 1-4.24-4.24" />
                                                <line x1="1" y1="1" x2="23" y2="23" />
                                            </svg>
                                        ) : (
                                            <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                                                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
                                                <circle cx="12" cy="12" r="3" />
                                            </svg>
                                        )}
                                    </button>
                                </div>
                                {errors.password && (
                                    <p className="login-error-msg">
                                        <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor">
                                            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm1 15h-2v-2h2v2zm0-4h-2V7h2v6z" />
                                        </svg>
                                        {errors.password.message}
                                    </p>
                                )}
                            </div>

                            {/* Submit */}
                            <button
                                type="submit"
                                disabled={isLoading}
                                className="login-submit"
                            >
                                {isLoading ? (
                                    <span className="login-submit-loading">
                                        <svg className="login-spinner" width="20" height="20" viewBox="0 0 24 24">
                                            <circle cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4" fill="none" opacity="0.25" />
                                            <path d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z" fill="currentColor" opacity="0.75" />
                                        </svg>
                                        در حال ورود...
                                    </span>
                                ) : (
                                    'ورود به پنل مدیریت'
                                )}
                            </button>
                        </form>

                        <div className="login-footer">
                            <p>پنل مدیریت سپر روان • نسخه ۱.۰</p>
                        </div>
                    </div>
                </div>
            </div>

            <style>{`
                /* ===== Login Page Styles ===== */
                .login-page {
                    min-height: 100vh;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    background: #f0f0f5;
                    position: relative;
                    overflow: hidden;
                }

                /* Background orbs */
                .login-bg {
                    position: fixed;
                    inset: 0;
                    pointer-events: none;
                    z-index: 0;
                }
                .login-bg-orb {
                    position: absolute;
                    border-radius: 50%;
                    filter: blur(80px);
                    opacity: 0.35;
                }
                .login-bg-orb-1 {
                    width: 500px; height: 500px;
                    background: #7c3aed;
                    top: -150px; right: -100px;
                    animation: orbFloat1 12s ease-in-out infinite;
                }
                .login-bg-orb-2 {
                    width: 400px; height: 400px;
                    background: #6366f1;
                    bottom: -100px; left: -80px;
                    animation: orbFloat2 10s ease-in-out infinite;
                }
                .login-bg-orb-3 {
                    width: 300px; height: 300px;
                    background: #a78bfa;
                    top: 50%; left: 50%;
                    transform: translate(-50%, -50%);
                    animation: orbFloat3 14s ease-in-out infinite;
                }

                @keyframes orbFloat1 {
                    0%, 100% { transform: translate(0, 0) scale(1); }
                    50% { transform: translate(-40px, 30px) scale(1.05); }
                }
                @keyframes orbFloat2 {
                    0%, 100% { transform: translate(0, 0) scale(1); }
                    50% { transform: translate(30px, -40px) scale(1.08); }
                }
                @keyframes orbFloat3 {
                    0%, 100% { transform: translate(-50%, -50%) scale(1); }
                    50% { transform: translate(-50%, -55%) scale(0.95); }
                }

                /* Main container */
                .login-container {
                    position: relative;
                    z-index: 1;
                    display: flex;
                    width: 100%;
                    max-width: 1100px;
                    min-height: 600px;
                    margin: 24px;
                    border-radius: 24px;
                    overflow: hidden;
                    box-shadow:
                        0 0 0 1px rgba(255,255,255,0.1),
                        0 25px 50px -12px rgba(0,0,0,0.15),
                        0 0 80px -20px rgba(124,58,237,0.15);
                }

                /* Brand panel */
                .login-brand {
                    flex: 1;
                    background: linear-gradient(135deg, #7c3aed 0%, #6366f1 50%, #8b5cf6 100%);
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    padding: 48px;
                    position: relative;
                    overflow: hidden;
                }
                .login-brand::before {
                    content: '';
                    position: absolute;
                    inset: 0;
                    background:
                        radial-gradient(circle at 20% 80%, rgba(255,255,255,0.08) 0%, transparent 50%),
                        radial-gradient(circle at 80% 20%, rgba(255,255,255,0.06) 0%, transparent 50%);
                }
                .login-brand-content {
                    position: relative;
                    z-index: 1;
                    text-align: center;
                    color: #fff;
                }
                .login-brand-icon {
                    width: 72px; height: 72px;
                    margin: 0 auto 24px;
                    background: rgba(255,255,255,0.15);
                    backdrop-filter: blur(12px);
                    border-radius: 20px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    color: #fff;
                    border: 1px solid rgba(255,255,255,0.2);
                }
                .login-brand-title {
                    font-size: 2.25rem;
                    font-weight: 800;
                    margin: 0 0 8px;
                    letter-spacing: -0.02em;
                }
                .login-brand-subtitle {
                    font-size: 1rem;
                    opacity: 0.75;
                    margin: 0 0 40px;
                    font-weight: 400;
                }
                .login-brand-features {
                    display: flex;
                    flex-direction: column;
                    gap: 16px;
                    text-align: right;
                }
                .login-feature {
                    display: flex;
                    align-items: center;
                    gap: 12px;
                    font-size: 0.95rem;
                    opacity: 0.85;
                }
                .login-feature svg {
                    flex-shrink: 0;
                    color: rgba(255,255,255,0.9);
                }

                /* Form panel */
                .login-form-wrapper {
                    width: 440px;
                    flex-shrink: 0;
                    background: #fff;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    padding: 48px 40px;
                }
                .login-form-card {
                    width: 100%;
                    max-width: 360px;
                }

                /* Mobile logo */
                .login-mobile-logo {
                    display: none;
                    text-align: center;
                    margin-bottom: 32px;
                }
                .login-mobile-icon {
                    width: 56px; height: 56px;
                    margin: 0 auto 12px;
                    background: linear-gradient(135deg, #7c3aed, #6366f1);
                    border-radius: 16px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    color: #fff;
                }
                .login-mobile-title {
                    font-size: 1.5rem;
                    font-weight: 800;
                    color: #1e1b4b;
                    margin: 0;
                }

                /* Form header */
                .login-form-header {
                    text-align: center;
                    margin-bottom: 32px;
                }
                .login-form-title {
                    font-size: 1.5rem;
                    font-weight: 800;
                    color: #1e1b4b;
                    margin: 0 0 8px;
                    letter-spacing: -0.02em;
                }
                .login-form-desc {
                    font-size: 0.9rem;
                    color: #9ca3af;
                    margin: 0;
                }

                /* Form fields */
                .login-form {
                    display: flex;
                    flex-direction: column;
                    gap: 20px;
                }
                .login-field {
                    display: flex;
                    flex-direction: column;
                    gap: 8px;
                }
                .login-label {
                    font-size: 0.85rem;
                    font-weight: 600;
                    color: #374151;
                }
                .login-input-wrap {
                    position: relative;
                    display: flex;
                    align-items: center;
                }
                .login-input-icon {
                    position: absolute;
                    right: 14px;
                    color: #9ca3af;
                    pointer-events: none;
                    transition: color 0.2s;
                }
                .login-input {
                    width: 100%;
                    padding: 12px 44px 12px 14px;
                    border: 2px solid #e5e7eb;
                    border-radius: 12px;
                    font-size: 0.95rem;
                    font-family: inherit;
                    background: #fafafa;
                    color: #1e1b4b;
                    transition: all 0.2s ease;
                    outline: none;
                }
                .login-input::placeholder {
                    color: #c4b5fd;
                    font-size: 0.85rem;
                }
                .login-input:focus {
                    border-color: #7c3aed;
                    background: #fff;
                    box-shadow: 0 0 0 4px rgba(124,58,237,0.1);
                }

                .login-input-error {
                    border-color: #fca5a5 !important;
                    background: #fef2f2 !important;
                }
                .login-input-error:focus {
                    box-shadow: 0 0 0 4px rgba(239,68,68,0.1) !important;
                    border-color: #ef4444 !important;
                }
                .login-toggle-pwd {
                    position: absolute;
                    left: 14px;
                    background: none;
                    border: none;
                    cursor: pointer;
                    color: #9ca3af;
                    padding: 4px;
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    transition: color 0.2s;
                }
                .login-toggle-pwd:hover {
                    color: #6b7280;
                }
                .login-error-msg {
                    display: flex;
                    align-items: center;
                    gap: 6px;
                    font-size: 0.8rem;
                    color: #ef4444;
                    margin: 0;
                }

                /* Submit button */
                .login-submit {
                    width: 100%;
                    padding: 14px;
                    border: none;
                    border-radius: 12px;
                    font-size: 1rem;
                    font-weight: 700;
                    font-family: inherit;
                    color: #fff;
                    background: linear-gradient(135deg, #7c3aed 0%, #6366f1 100%);
                    cursor: pointer;
                    transition: all 0.25s ease;
                    box-shadow: 0 4px 15px -3px rgba(124,58,237,0.4);
                    margin-top: 8px;
                }
                .login-submit:hover:not(:disabled) {
                    transform: translateY(-2px);
                    box-shadow: 0 8px 25px -5px rgba(124,58,237,0.5);
                }
                .login-submit:active:not(:disabled) {
                    transform: translateY(0);
                }
                .login-submit:disabled {
                    opacity: 0.6;
                    cursor: not-allowed;
                }
                .login-submit-loading {
                    display: flex;
                    align-items: center;
                    justify-content: center;
                    gap: 10px;
                }
                .login-spinner {
                    animation: spin 0.8s linear infinite;
                }
                @keyframes spin {
                    to { transform: rotate(360deg); }
                }

                /* Footer */
                .login-footer {
                    margin-top: 28px;
                    padding-top: 20px;
                    border-top: 1px solid #f3f4f6;
                    text-align: center;
                }
                .login-footer p {
                    font-size: 0.75rem;
                    color: #c4b5fd;
                    margin: 0;
                }

                /* Responsive */
                @media (max-width: 1024px) {
                    .login-brand {
                        display: none;
                    }
                    .login-form-wrapper {
                        width: 100%;
                    }
                    .login-mobile-logo {
                        display: block;
                    }
                }
                @media (max-width: 480px) {
                    .login-container {
                        margin: 12px;
                        border-radius: 20px;
                        min-height: auto;
                    }
                    .login-form-wrapper {
                        padding: 32px 24px;
                    }
                }
            `}</style>
        </div>
    )
}
