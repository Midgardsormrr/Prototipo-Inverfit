/**
 * Formats a CLP monetary amount into Chilean locale format, e.g. "$45.000".
 * Returns "Por definir" if null or undefined.
 */
export function formatCurrencyCLP(amount: number | null | undefined): string {
  if (amount === null || amount === undefined) {
    return 'Por definir'
  }
  return `$${amount.toLocaleString('es-CL')}`
}

/**
 * Sanitizes a string to prevent CSV Formula Injection (DDE injection).
 * If the value starts with '=', '+', '-', '@', '\t', or '\r', it prepends a single quote.
 */
export function sanitizeCsvCell(val: string): string {
  const trimmed = (val || '').trim()
  if (trimmed.length > 0) {
    const first = trimmed[0]
    if (first === '=' || first === '+' || first === '-' || first === '@' || first === '\t' || first === '\r') {
      return `'${val}`
    }
  }
  return val
}

/**
 * Gets today's date formatted as YYYY-MM-DD.
 */
export function getTodayDateString(): string {
  const now = new Date()
  const year = now.getFullYear()
  const month = String(now.getMonth() + 1).padStart(2, '0')
  const day = String(now.getDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}
