import { useState, type FormEvent } from 'react'
import { Link } from 'react-router-dom'
import { supabase } from '../../../lib/supabaseClient'
import { api } from '../../../lib/api'
import { getPasswordIssues, isValidEmail, isValidFullName } from '../../../lib/validation'

const PASSWORD_RULE_HINT = '8+ characters, with uppercase, lowercase, a number, and a symbol'

type FieldErrors = {
  fullName?: string
  email?: string
  password?: string
  confirmPassword?: string
  agreeToTerms?: string
}

export function SignUpPage() {
  const [fullName, setFullName] = useState('')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [agreeToTerms, setAgreeToTerms] = useState(false)
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({})
  const [formError, setFormError] = useState<string | null>(null)
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [signUpComplete, setSignUpComplete] = useState(false)

  function validate(): FieldErrors {
    const errors: FieldErrors = {}

    if (!isValidFullName(fullName)) {
      errors.fullName = 'Enter your full name (letters, spaces, hyphens, and apostrophes only)'
    }
    if (!isValidEmail(email)) {
      errors.email = 'Enter a valid email address'
    }
    const passwordIssues = getPasswordIssues(password)
    if (passwordIssues.length > 0) {
      errors.password = `Password must have: ${passwordIssues.join(', ').toLowerCase()}`
    }
    if (password !== confirmPassword) {
      errors.confirmPassword = 'Passwords do not match'
    }
    if (!agreeToTerms) {
      errors.agreeToTerms = 'You must agree to the Terms of Service and Privacy Policy'
    }

    return errors
  }

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    setFormError(null)

    const errors = validate()
    setFieldErrors(errors)
    if (Object.keys(errors).length > 0) return

    setIsSubmitting(true)
    try {
      await api.post('/api/v1/auth/signup', {
        email,
        password,
        full_name: fullName,
        agree_to_terms: agreeToTerms,
      })
      setSignUpComplete(true)
    } catch (err) {
      setFormError(err instanceof Error ? err.message : 'Unable to create account, please try again')
    } finally {
      setIsSubmitting(false)
    }
  }

  async function handleGoogleSignUp() {
    setFormError(null)
    if (!supabase) {
      setFormError('Authentication is not configured.')
      return
    }
    const { error: oauthError } = await supabase.auth.signInWithOAuth({ provider: 'google' })
    if (oauthError) setFormError(oauthError.message)
  }

  async function handleAppleSignUp() {
    setFormError(null)
    if (!supabase) {
      setFormError('Authentication is not configured.')
      return
    }
    const { error: oauthError } = await supabase.auth.signInWithOAuth({ provider: 'apple' })
    if (oauthError) setFormError(oauthError.message)
  }

  return (
    <div className="flex min-h-screen w-full items-center justify-center bg-[#FAFAFB] p-6 font-sans">
      <div className="relative flex min-h-[900px] w-full max-w-[1050px] flex-col overflow-hidden rounded-[28px] bg-linear-to-br from-[#EAF2FF] via-[#F8FAFF] to-[#FFF8ED] shadow-[0_30px_80px_-20px_rgba(31,31,31,0.10)]">
        <div className="pointer-events-none absolute -top-32 -left-32 h-[480px] w-[480px] rounded-full bg-linear-to-br from-indigo-100/50 to-transparent blur-3xl" />
        <div className="pointer-events-none absolute -right-32 -bottom-32 h-[520px] w-[520px] rounded-full bg-linear-to-br from-amber-100/60 via-orange-50/40 to-transparent blur-3xl" />

        <header className="relative z-10 flex items-center justify-between px-12 pt-10">
          <div className="flex items-center gap-2">
            <FinEduLogo />
            <span className="text-lg font-semibold text-[#1F1F1F]">FinEdu</span>
          </div>
          <Link
            to="/"
            className="flex items-center gap-1.5 text-sm font-medium text-[#6D7280] transition hover:text-[#1F1F1F]"
          >
            Back to Homepage
            <ArrowRightIcon />
          </Link>
        </header>

        <div className="relative z-10 flex flex-1 items-center">
          {/* LEFT SECTION */}
          <div className="hidden w-[45%] flex-col justify-center px-12 lg:flex">
            <h1 className="text-[42px] leading-[1.1] font-bold tracking-tight text-[#1F1F1F]">
              Start your journey!
            </h1>
            <p className="mt-4 max-w-xs text-base leading-relaxed text-[#6D7280]">
              Create an account to unlock personalized market insights and start learning to
              invest with confidence.
            </p>
            <LoginIllustration />
          </div>

          {/* RIGHT SECTION */}
          <div className="flex w-full items-center justify-center p-8 lg:w-[55%]">
            <div className="w-full max-w-[520px] rounded-[24px] border border-[#EAEAEA] bg-white p-10 shadow-[0_20px_50px_-20px_rgba(31,31,31,0.08)]">
              <h2 className="text-[28px] font-bold text-[#1F1F1F]">Create your account</h2>
              <p className="mt-2 text-base text-[#6D7280]">
                Sign up to start tracking markets and building your knowledge
              </p>

              {signUpComplete ? (
                <p className="mt-8 rounded-xl border border-[#EAEAEA] bg-[#FAFAFB] p-4 text-sm text-[#1F1F1F]">
                  Check your email to verify your account before logging in.
                </p>
              ) : (
                <form className="mt-8 space-y-6" onSubmit={handleSubmit} noValidate>
                  <div>
                    <label htmlFor="fullName" className="mb-2 block text-sm font-medium text-[#1F1F1F]">
                      Full name
                    </label>
                    <input
                      id="fullName"
                      type="text"
                      autoComplete="name"
                      placeholder="Jane Doe"
                      value={fullName}
                      onChange={(e) => setFullName(e.target.value)}
                      className="w-full rounded-xl border border-[#EAEAEA] px-4 py-3 text-sm text-[#1F1F1F] placeholder:text-[#9CA0A8] transition focus:border-[#6C7CFF] focus:ring-2 focus:ring-[#DDE7FF] focus:outline-none"
                    />
                    {fieldErrors.fullName && (
                      <p className="mt-1.5 text-xs text-red-600">{fieldErrors.fullName}</p>
                    )}
                  </div>

                  <div>
                    <label htmlFor="email" className="mb-2 block text-sm font-medium text-[#1F1F1F]">
                      Email
                    </label>
                    <input
                      id="email"
                      type="email"
                      autoComplete="email"
                      placeholder="you@example.com"
                      value={email}
                      onChange={(e) => setEmail(e.target.value)}
                      className="w-full rounded-xl border border-[#EAEAEA] px-4 py-3 text-sm text-[#1F1F1F] placeholder:text-[#9CA0A8] transition focus:border-[#6C7CFF] focus:ring-2 focus:ring-[#DDE7FF] focus:outline-none"
                    />
                    {fieldErrors.email && <p className="mt-1.5 text-xs text-red-600">{fieldErrors.email}</p>}
                  </div>

                  <div>
                    <label htmlFor="password" className="mb-2 block text-sm font-medium text-[#1F1F1F]">
                      Password
                    </label>
                    <input
                      id="password"
                      type="password"
                      autoComplete="new-password"
                      placeholder="Create a password"
                      value={password}
                      onChange={(e) => setPassword(e.target.value)}
                      className="w-full rounded-xl border border-[#EAEAEA] px-4 py-3 text-sm text-[#1F1F1F] placeholder:text-[#9CA0A8] transition focus:border-[#6C7CFF] focus:ring-2 focus:ring-[#DDE7FF] focus:outline-none"
                    />
                    <p className="mt-1.5 text-xs text-[#6D7280]">{PASSWORD_RULE_HINT}</p>
                    {fieldErrors.password && (
                      <p className="mt-1.5 text-xs text-red-600">{fieldErrors.password}</p>
                    )}
                  </div>

                  <div>
                    <label
                      htmlFor="confirmPassword"
                      className="mb-2 block text-sm font-medium text-[#1F1F1F]"
                    >
                      Confirm password
                    </label>
                    <input
                      id="confirmPassword"
                      type="password"
                      autoComplete="new-password"
                      placeholder="Re-enter your password"
                      value={confirmPassword}
                      onChange={(e) => setConfirmPassword(e.target.value)}
                      className="w-full rounded-xl border border-[#EAEAEA] px-4 py-3 text-sm text-[#1F1F1F] placeholder:text-[#9CA0A8] transition focus:border-[#6C7CFF] focus:ring-2 focus:ring-[#DDE7FF] focus:outline-none"
                    />
                    {fieldErrors.confirmPassword && (
                      <p className="mt-1.5 text-xs text-red-600">{fieldErrors.confirmPassword}</p>
                    )}
                  </div>

                  <div>
                    <label className="flex items-start gap-2 text-sm text-[#6D7280]">
                      <input
                        type="checkbox"
                        checked={agreeToTerms}
                        onChange={(e) => setAgreeToTerms(e.target.checked)}
                        className="mt-0.5 h-4 w-4 rounded border-[#EAEAEA] text-[#6C7CFF] focus:ring-[#6C7CFF]"
                      />
                      <span>
                        I agree to the{' '}
                        <a href="#" className="font-medium text-[#6C7CFF] hover:text-[#5768f0]">
                          Terms of Service
                        </a>{' '}
                        and{' '}
                        <a href="#" className="font-medium text-[#6C7CFF] hover:text-[#5768f0]">
                          Privacy Policy
                        </a>
                      </span>
                    </label>
                    {fieldErrors.agreeToTerms && (
                      <p className="mt-1.5 text-xs text-red-600">{fieldErrors.agreeToTerms}</p>
                    )}
                  </div>

                  {formError && <p className="text-sm text-red-600">{formError}</p>}

                  <button
                    type="submit"
                    disabled={isSubmitting}
                    className="h-[52px] w-full rounded-xl bg-[#1F1F1F] text-sm font-semibold text-white transition hover:bg-black disabled:opacity-60"
                  >
                    {isSubmitting ? 'Creating account…' : 'Create account'}
                  </button>
                </form>
              )}

              {!signUpComplete && (
                <>
                  <div className="my-8 flex items-center gap-3">
                    <div className="h-px flex-1 bg-[#EAEAEA]" />
                    <span className="text-[13px] text-[#6D7280]">or continue with</span>
                    <div className="h-px flex-1 bg-[#EAEAEA]" />
                  </div>

                  <div className="space-y-3">
                    <button
                      type="button"
                      onClick={handleGoogleSignUp}
                      className="flex h-[52px] w-full items-center justify-center gap-2 rounded-xl border border-[#EAEAEA] text-sm font-medium text-[#1F1F1F] transition hover:bg-[#FAFAFB]"
                    >
                      <GoogleIcon />
                      Continue with Google
                    </button>

                    <button
                      type="button"
                      onClick={handleAppleSignUp}
                      className="flex h-[52px] w-full items-center justify-center gap-2 rounded-xl border border-[#EAEAEA] text-sm font-medium text-[#1F1F1F] transition hover:bg-[#FAFAFB]"
                    >
                      <AppleIcon />
                      Continue with Apple
                    </button>
                  </div>
                </>
              )}

              <p className="mt-8 text-center text-sm text-[#6D7280]">
                Already have an account?{' '}
                <Link to="/login" className="font-medium text-[#6C7CFF] hover:text-[#5768f0]">
                  Log in
                </Link>
              </p>
            </div>
          </div>
        </div>

        <footer className="relative z-10 flex items-center justify-between px-12 pb-10 text-[13px] text-[#6D7280]">
          <div className="flex items-center gap-2">
            <FinEduLogo />
            <span className="font-medium text-[#1F1F1F]">FinEdu</span>
          </div>
          <span>© {new Date().getFullYear()} FinEdu. All rights reserved.</span>
          <div className="flex gap-4">
            <a href="#" className="transition hover:text-[#1F1F1F]">
              Privacy Policy
            </a>
            <a href="#" className="transition hover:text-[#1F1F1F]">
              Terms of Service
            </a>
            <a href="#" className="transition hover:text-[#1F1F1F]">
              Contact Us
            </a>
          </div>
        </footer>
      </div>
    </div>
  )
}

function FinEduLogo() {
  return (
    <svg width="28" height="28" viewBox="0 0 28 28" fill="none" xmlns="http://www.w3.org/2000/svg">
      <rect width="28" height="28" rx="8" fill="#1F1F1F" />
      <path
        d="M8 18.5L11.8 14.2L14.6 16.8L20 10.5"
        stroke="white"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
      <path
        d="M16.2 10.5H20V14.3"
        stroke="white"
        strokeWidth="1.8"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function ArrowRightIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        d="M3.5 8H12.5M12.5 8L8.5 4M12.5 8L8.5 12"
        stroke="currentColor"
        strokeWidth="1.5"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  )
}

function GoogleIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" xmlns="http://www.w3.org/2000/svg">
      <path
        fill="#4285F4"
        d="M17.64 9.2c0-.64-.06-1.25-.16-1.84H9v3.48h4.84a4.14 4.14 0 0 1-1.8 2.72v2.26h2.9c1.7-1.57 2.7-3.88 2.7-6.62z"
      />
      <path
        fill="#34A853"
        d="M9 18c2.43 0 4.47-.8 5.96-2.18l-2.9-2.26c-.8.54-1.84.86-3.06.86-2.35 0-4.34-1.59-5.05-3.72H.95v2.33A9 9 0 0 0 9 18z"
      />
      <path
        fill="#FBBC05"
        d="M3.95 10.7A5.4 5.4 0 0 1 3.67 9c0-.59.1-1.17.28-1.7V4.97H.95A9 9 0 0 0 0 9c0 1.45.35 2.83.95 4.03l3-2.33z"
      />
      <path
        fill="#EA4335"
        d="M9 3.58c1.32 0 2.51.46 3.44 1.35l2.58-2.58C13.46.89 11.43 0 9 0A9 9 0 0 0 .95 4.97l3 2.33C4.66 5.17 6.65 3.58 9 3.58z"
      />
    </svg>
  )
}

function AppleIcon() {
  return (
    <svg width="18" height="18" viewBox="0 0 18 18" fill="none" xmlns="http://www.w3.org/2000/svg">
      <path
        fill="#1F1F1F"
        d="M13.53 9.56c-.02-1.86 1.52-2.75 1.59-2.8-.87-1.27-2.22-1.45-2.7-1.47-1.15-.12-2.24.68-2.82.68-.58 0-1.48-.66-2.43-.64-1.25.02-2.4.72-3.04 1.84-1.3 2.25-.33 5.57.93 7.4.62.89 1.36 1.9 2.32 1.86.93-.04 1.28-.6 2.4-.6 1.12 0 1.44.6 2.42.58 1-.02 1.63-.9 2.24-1.8.7-1.03.99-2.03 1.01-2.08-.02-.01-1.9-.73-1.92-2.97z"
      />
      <path
        fill="#1F1F1F"
        d="M11.75 3.87c.52-.63.87-1.5.77-2.37-.75.03-1.65.5-2.19 1.12-.48.55-.9 1.44-.79 2.28.83.06 1.68-.42 2.21-1.03z"
      />
    </svg>
  )
}

const CHART_POINTS: [number, number][] = [
  [10, 218],
  [78, 132],
  [146, 188],
  [214, 84],
  [282, 140],
  [350, 48],
  [410, 14],
]

const PULSE_INDICES = new Set([1, 3, 5, 6])

function LoginIllustration() {
  const linePath = CHART_POINTS.map(([x, y], i) => `${i === 0 ? 'M' : 'L'}${x} ${y}`).join(' ')
  const [lastX, lastY] = CHART_POINTS[CHART_POINTS.length - 1]
  const [firstX] = CHART_POINTS[0]
  const areaPath = `${linePath} L${lastX} 250 L${firstX} 250 Z`

  return (
    <div className="pointer-events-none relative mt-7.5 h-78 w-full min-w-86.25 select-none">
      {/* big, vivid ambient glow behind the chart, background stays fully transparent */}
      <div className="absolute inset-x-0 top-3 h-54 rounded-full bg-linear-to-br from-indigo-200/50 via-blue-100/30 to-amber-200/40 blur-3xl" />
      <div className="absolute top-12 right-3 h-30 w-30 rounded-full bg-amber-200/40 blur-3xl" />

      <svg
        className="absolute top-10.5 left-1/2 h-auto w-[96%] max-w-90 -translate-x-1/2"
        viewBox="0 0 420 260"
        fill="none"
        xmlns="http://www.w3.org/2000/svg"
      >
        <defs>
          <linearGradient id="lineStrokeSignUp" x1="0" y1="0" x2="1" y2="0">
            <stop offset="0" stopColor="#60A5FA" />
            <stop offset="0.5" stopColor="#6C7CFF" />
            <stop offset="1" stopColor="#F5B942" />
          </linearGradient>
          <linearGradient id="areaFillSignUp" x1="0" y1="0" x2="0" y2="1">
            <stop offset="0" stopColor="#6C7CFF" stopOpacity="0.4" />
            <stop offset="1" stopColor="#6C7CFF" stopOpacity="0" />
          </linearGradient>
          <filter id="pointGlowSignUp" x="-150%" y="-150%" width="400%" height="400%">
            <feGaussianBlur stdDeviation="7" />
          </filter>
          <filter id="lineGlowSignUp" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur stdDeviation="6" />
          </filter>
        </defs>

        {/* faint baseline grid, kept minimal so the chart reads as transparent */}
        <g stroke="#C7D2FE" strokeOpacity="0.35" strokeWidth="1" strokeDasharray="3 6">
          <line x1="0" y1="60" x2="420" y2="60" />
          <line x1="0" y1="140" x2="420" y2="140" />
          <line x1="0" y1="220" x2="420" y2="220" />
        </g>

        <path d={areaPath} fill="url(#areaFillSignUp)" />

        {/* neon glow trail beneath the crisp trend line */}
        <path
          d={linePath}
          stroke="url(#lineStrokeSignUp)"
          strokeWidth="7"
          strokeLinecap="round"
          strokeLinejoin="round"
          opacity="0.45"
          filter="url(#lineGlowSignUp)"
        />
        <path
          d={linePath}
          stroke="url(#lineStrokeSignUp)"
          strokeWidth="2"
          strokeLinecap="round"
          strokeLinejoin="round"
        />

        {CHART_POINTS.map(([cx, cy], i) => (
          <g key={`${cx}-${cy}`}>
            <circle
              cx={cx}
              cy={cy}
              r="15"
              fill="#6C7CFF"
              opacity="0.4"
              filter="url(#pointGlowSignUp)"
              className={PULSE_INDICES.has(i) ? 'animate-pulse' : undefined}
            />
            <circle cx={cx} cy={cy} r="7.5" fill="#FFFFFF" />
            <circle cx={cx} cy={cy} r="5.5" fill={i === CHART_POINTS.length - 1 ? '#F5B942' : '#5768F0'} />
          </g>
        ))}

        {/* live-market pulse on the newest point */}
        <circle cx={lastX} cy={lastY} r="6" fill="#F5B942" opacity="0.6" className="animate-ping" />
      </svg>

      {/* callouts scattered around the chart's bends */}
      <TrendCallout value="+24.6%" direction="up" className="-top-1 right-[4%]" />
      <TrendCallout value="+9.8%" direction="up" className="top-[24%] left-0" />
      <TrendCallout value="-3.4%" direction="down" className="top-[68%] left-[28%]" />
    </div>
  )
}

function TrendCallout({
  value,
  direction,
  className,
}: {
  value: string
  direction: 'up' | 'down'
  className: string
}) {
  const color = direction === 'up' ? '#16A34A' : '#DC2626'
  return (
    <div className={`absolute z-20 px-3 py-2 ${className}`}>
      <div className="flex items-center gap-1.5">
        <svg width="14" height="14" viewBox="0 0 14 14" fill="none" xmlns="http://www.w3.org/2000/svg">
          {direction === 'up' ? (
            <>
              <path
                d="M2 12L5.5 7.5L8 10L12 4"
                stroke={color}
                strokeWidth="1.6"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
              <path d="M8.6 4H12V7.4" stroke={color} strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" />
            </>
          ) : (
            <>
              <path
                d="M2 2L5.5 6.5L8 4L12 10"
                stroke={color}
                strokeWidth="1.6"
                strokeLinecap="round"
                strokeLinejoin="round"
              />
              <path d="M8.6 10H12V6.6" stroke={color} strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" />
            </>
          )}
        </svg>
        <span className="text-sm font-semibold" style={{ color }}>
          {value}
        </span>
      </div>
    </div>
  )
}
