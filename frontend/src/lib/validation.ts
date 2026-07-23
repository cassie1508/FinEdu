const EMAIL_PATTERN = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
const FULL_NAME_PATTERN = /^[\p{L} '.-]+$/u

export function isValidEmail(email: string): boolean {
  return EMAIL_PATTERN.test(email)
}

// Classification mirrors the backend's rune-based check (unicode.IsLower/
// IsUpper/IsDigit, else special) using Unicode property escapes instead of
// ASCII-only ranges, so a password isn't judged differently by each side.
export function getPasswordIssues(password: string): string[] {
  const issues: string[] = []
  if (password.length < 8) issues.push('At least 8 characters')
  if (!/\p{Ll}/u.test(password)) issues.push('At least one lowercase letter')
  if (!/\p{Lu}/u.test(password)) issues.push('At least one uppercase letter')
  if (!/\p{Nd}/u.test(password)) issues.push('At least one number')
  if (!/[^\p{Ll}\p{Lu}\p{Nd}]/u.test(password)) issues.push('At least one special character')
  return issues
}

export function isValidFullName(name: string): boolean {
  const trimmed = name.trim()
  return trimmed.length >= 2 && trimmed.length <= 100 && FULL_NAME_PATTERN.test(trimmed)
}
