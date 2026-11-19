export interface ApiErrorDetail {
  code: string
  message: string
}

export interface ApiErrorResponse {
  error: ApiErrorDetail
}

export class HttpClientError extends Error {
  public statusCode: number
  public code: string

  constructor(message: string, statusCode: number, code: string = 'HTTP_ERROR') {
    super(message)
    this.name = 'HttpClientError'
    this.statusCode = statusCode
    this.code = code
  }
}

export class HttpClient {
  private baseUrl: string
  private defaultTimeoutMs: number

  constructor(
    baseUrl: string = (import.meta.env.VITE_API_BASE_URL as string) || '',
    defaultTimeoutMs: number = 10000
  ) {
    let finalBaseUrl = baseUrl ? baseUrl.trim() : ''

    // Si estamos corriendo en un entorno en la nube (ej. Vercel) y no en localhost,
    // garantizamos que las peticiones apunten al backend desplegado en Render (evita fallo de localhost)
    if (typeof window !== 'undefined' && window.location.hostname !== 'localhost' && window.location.hostname !== '127.0.0.1') {
      if (!finalBaseUrl || finalBaseUrl.includes('localhost') || finalBaseUrl.includes('127.0.0.1')) {
        finalBaseUrl = 'https://inverfit-api.onrender.com/api/v1'
      }
    } else if (!finalBaseUrl) {
      finalBaseUrl = 'http://localhost:8081/api/v1'
    }

    this.baseUrl = finalBaseUrl.replace(/\/+$/, '')
    this.defaultTimeoutMs = defaultTimeoutMs
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const url = `${this.baseUrl}/${endpoint.replace(/^\/+/, '')}`

    const controller = new AbortController()
    const timeoutId = setTimeout(() => controller.abort(), this.defaultTimeoutMs)

    const headers: Record<string, string> = {
      'Accept': 'application/json',
      ...((options.headers as Record<string, string>) || {})
    }

    if (options.body && typeof options.body === 'string' && !headers['Content-Type']) {
      headers['Content-Type'] = 'application/json'
    }

    try {
      const response = await fetch(url, {
        ...options,
        headers,
        signal: controller.signal
      })

      clearTimeout(timeoutId)

      // Manejo de respuestas sin contenido (204)
      if (response.status === 204) {
        return {} as T
      }

      // Parseo seguro del cuerpo JSON
      const text = await response.text()
      let data: any = null
      if (text) {
        try {
          data = JSON.parse(text)
        } catch {
          data = text
        }
      }

      if (!response.ok) {
        let errorMsg = `Error HTTP ${response.status}: ${response.statusText}`
        let errorCode = 'HTTP_ERROR'

        if (data && typeof data === 'object' && data.error) {
          errorMsg = data.error.message || errorMsg
          errorCode = data.error.code || errorCode
        }

        throw new HttpClientError(errorMsg, response.status, errorCode)
      }

      return data as T
    } catch (err: unknown) {
      clearTimeout(timeoutId)

      if (err instanceof HttpClientError) {
        throw err
      }

      if (err instanceof DOMException && err.name === 'AbortError') {
        throw new HttpClientError(
          `La solicitud a '${endpoint}' excedió el tiempo límite de espera (${this.defaultTimeoutMs / 1000}s).`,
          408,
          'TIMEOUT'
        )
      }

      // Error de red, conexión rechazada o CORS
      const netMsg = err instanceof Error ? err.message : 'Error de conexión'
      throw new HttpClientError(
        `No se pudo establecer conexión con el backend en ${this.baseUrl}. Verifique que el servicio Go esté en ejecución (${netMsg}).`,
        0,
        'NETWORK_ERROR'
      )
    }
  }

  public get<T>(endpoint: string, params?: Record<string, string>): Promise<T> {
    let finalEndpoint = endpoint
    if (params && Object.keys(params).length > 0) {
      const searchParams = new URLSearchParams()
      for (const [key, val] of Object.entries(params)) {
        if (val !== undefined && val !== null && val !== '') {
          searchParams.append(key, val)
        }
      }
      const qs = searchParams.toString()
      if (qs) {
        finalEndpoint += `?${qs}`
      }
    }
    return this.request<T>(finalEndpoint, { method: 'GET' })
  }

  public post<T>(endpoint: string, body?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'POST',
      body: body !== undefined ? JSON.stringify(body) : undefined
    })
  }

  public put<T>(endpoint: string, body?: any): Promise<T> {
    return this.request<T>(endpoint, {
      method: 'PUT',
      body: body !== undefined ? JSON.stringify(body) : undefined
    })
  }

  public delete<T>(endpoint: string): Promise<T> {
    return this.request<T>(endpoint, { method: 'DELETE' })
  }

  public getBaseUrl(): string {
    return this.baseUrl
  }
}

export const httpClient = new HttpClient()
