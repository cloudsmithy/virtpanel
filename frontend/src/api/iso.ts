import http from './http'
import axios from 'axios'
import type { CancelTokenSource } from 'axios'

export interface ISOFile {
  name: string
  path: string
  size: number
}

export const isoApi = {
  list: () => http.get<any, ISOFile[]>('/isos'),
  upload: (file: File, onProgress?: (percent: number) => void): { source: CancelTokenSource; promise: Promise<any> } => {
    const source = axios.CancelToken.source()
    const form = new FormData()
    form.append('file', file)
    const promise = axios.post('/api/isos/upload', form, {
      baseURL: '',
      timeout: 0,
      cancelToken: source.token,
      onUploadProgress: (e) => {
        if (onProgress && e.total) onProgress(Math.round((e.loaded / e.total) * 100))
      },
    })
    return { source, promise }
  },
  delete: (name: string) => http.delete(`/isos/${name}`),
}
