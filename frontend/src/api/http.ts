import axios from 'axios'

const http = axios.create({ baseURL: '/api', timeout: 300000 })

http.interceptors.response.use(
  (res) => res.data,
  (err) => {
    console.error(err)
    return Promise.reject(err)
  }
)

export default http

export function errMsg(e: any, fallback = '操作失败'): string {
  return e?.response?.data?.error || e?.message || fallback
}
