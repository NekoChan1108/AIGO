// 使用相对路径，触发 Vite Proxy 代理
// 前端请求 /api/v1/... -> Vite 代理 -> http://localhost:9999/api/v1/...
const BASE_URL = '/api/v1'
const TOKEN_KEY = 'access_token'

interface RequestOptions extends RequestInit {
  params?: Record<string, string>;
}

// 基础请求处理函数
async function fetcher<T>(endpoint: string, options: RequestOptions = {}): Promise<T> {
  const { params, ...init } = options
  
  // 处理 URL 参数
  let url = `${BASE_URL}${endpoint}`
  if (params) {
    const searchParams = new URLSearchParams(params)
    url += `?${searchParams.toString()}`
  }

  // 设置默认 Headers
  const headers = new Headers(init.headers)
  if (!headers.has('Content-Type') && !(init.body instanceof FormData)) {
    headers.set('Content-Type', 'application/json')
  }

  const response = await fetch(url, { ...init, headers })

  // 1. 全局处理：检查后端是否下发了新 Token (自动刷新机制)
  const newToken = response.headers.get('X-New-Access-Token')
  if (newToken) {
    const tokenValue = newToken.replace('Bearer ', '')
    localStorage.setItem(TOKEN_KEY, tokenValue)
  }

  const data = await response.json()

  // 2. 错误处理
  if (!response.ok) {
    // 401 未授权：清除 Token 并跳转登录 (仅针对需要鉴权的接口)
    if (response.status === 401) {
      localStorage.removeItem(TOKEN_KEY)
      if (!window.location.pathname.includes('/login')) {
        window.location.href = '/login'
      }
    }
    throw new Error(data.msg || data.message || 'Request failed')
  }

  return data
}

// ==========================================
// 1. Public Client (不需要 Token)
// 用于: 登录, 注册, 忘记密码, 公开数据
// ==========================================
export const publicClient = {
  get: <T>(url: string, params?: Record<string, string>) => 
    fetcher<T>(url, { method: 'GET', params }),
    
  post: <T>(url: string, body: any) => 
    fetcher<T>(url, { method: 'POST', body: body instanceof FormData ? body : JSON.stringify(body) }),
}

// ==========================================
// 2. Auth Client (自动携带 Token)
// 用于: 聊天, 用户信息, 文件上传
// ==========================================
export const authClient = {
  request: <T>(url: string, options: RequestOptions) => {
    const token = localStorage.getItem(TOKEN_KEY)
    const headers = new Headers(options.headers)
    
    if (token) {
      // 这里的 Key 要看你后端 middleware/auth.go 怎么取
      // 通常是 Authorization
      headers.set('Authorization', `Bearer ${token}`) 
    }

    return fetcher<T>(url, { ...options, headers })
  },

  get: <T>(url: string, params?: Record<string, string>) => 
    authClient.request<T>(url, { method: 'GET', params }),

  post: <T>(url: string, body: any) => 
    authClient.request<T>(url, { method: 'POST', body: body instanceof FormData ? body : JSON.stringify(body) }),
    
  delete: <T>(url: string, body?: any) => 
    authClient.request<T>(url, { method: 'DELETE', body: body ? JSON.stringify(body) : undefined }),
}