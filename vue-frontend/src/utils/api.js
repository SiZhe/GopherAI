import axios from 'axios'

const api = axios.create({
  baseURL: '/api', // 使用代理路径，开发环境会自动代理到后端
  timeout: 0, //不启用超时机制
  // 使用 httpOnly cookie 做鉴权时，需要显式携带凭证
  withCredentials: true
})

// -------- token 刷新：at 失效后用 rt 换新 at（cookie 更新）--------
const STATUS_INVALID_TOKEN = 2006 // CodeInvalidToken
const STATUS_NOT_LOGIN = 2008 // CodeNotLogin
const STATUS_SUCCESS = 1000 // CodeSuccess

// 后端刷新接口在你的 router 里注册为：
// r.POST("/refreshAT", user.RefreshAccessToken)
// 结合 enterRouter.Group("/user") 和前端 baseURL('/api') rewrite:
// 前端请求路径应为 '/user/refreshAT'
const REFRESH_PATH = 'user/refresh-access-token'
const LOGOUT_PATH = 'device/logout'
const SKIP_AUTH_REFRESH = '__skipAuthRefresh'


let refreshPromise = null

async function refreshAccessToken() {
  if (refreshPromise) return refreshPromise

  refreshPromise = (async () => {
    try {
      const res = await api.post(REFRESH_PATH, {}, { [SKIP_AUTH_REFRESH]: true })
      const statusCode = res?.data?.status_code
      return statusCode === STATUS_SUCCESS
    } catch (e) {
      return false
    }
  })()
    .finally(() => {
      refreshPromise = null
    })

  return refreshPromise
}

function redirectToLogin() {
  window.location.href = '/login'
}

async function handleAuthError(config) {
  // 避免无限重试
  if (config && config.__isRetry) return false
  if (!config) return false

  config.__isRetry = true

  const ok = await refreshAccessToken()
  if (!ok) {
    redirectToLogin()
    return false
  }

  // 刷新成功后重试原请求
  return api.request(config)
}

// 响应拦截器：后端鉴权失败使用 status_code + 200
api.interceptors.response.use(
  async response => {
    if (response?.config?.[SKIP_AUTH_REFRESH]) return response

    const statusCode = response?.data?.status_code
    if (statusCode === STATUS_INVALID_TOKEN || statusCode === STATUS_NOT_LOGIN) {
      // config 可能不存在（极少数情况），兜底跳转
      const retryResult = await handleAuthError(response?.config)
      if (retryResult) return retryResult
      redirectToLogin()
    }
    return response
  },
  async error => {
    // 刷新/退出登录等“跳过鉴权刷新”的请求，失败时不应全局重定向
    if (error?.config?.[SKIP_AUTH_REFRESH]) {
      return Promise.reject(error)
    }

    // 少量情况下后端也可能返回 401；这里同样尝试刷新
    if (error?.response?.status === 401) {
      const retryResult = await handleAuthError(error?.config)
      if (retryResult) return retryResult
    }
    redirectToLogin()
    return Promise.reject(error)
  }
)

// 给界面按钮用的退出登录方法（调用后后端清理 cookie）
export async function logout() {
  const res = await api.post(LOGOUT_PATH, {}, { [SKIP_AUTH_REFRESH]: true })
  const statusCode = res?.data?.status_code
  const ok = statusCode === STATUS_SUCCESS
  console.log('[logout] request url:', res?.config?.url, 'status_code:', statusCode, 'ok:', ok)
  return ok
}

// 获取设备列表
export async function getDeviceList() {
  const res = await api.post('/device/device-list', {})
  const statusCode = res?.data?.status_code
  if (statusCode === STATUS_SUCCESS) {
    return { success: true, devices: res?.data?.sessions || [] }
  }
  return { success: false, devices: [] }
}


// 让指定设备下线
export async function offlineDevice(deviceIp, deviceBrowser) {
  const res = await api.post('/device/offline-device', {
    device_ip: deviceIp,
    device_browser: deviceBrowser
  })
  const statusCode = res?.data?.status_code
  return statusCode === STATUS_SUCCESS
}

export default api