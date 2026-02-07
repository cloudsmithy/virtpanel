import http from './http'

export interface Bridge {
  name: string
  up: boolean
  slaves: string[] | null
  ip: string
}

export const bridgeApi = {
  list: () => http.get<any, Bridge[]>('/bridges'),
  create: (data: { name: string; slave_nic?: string }) => http.post('/bridges', data),
  delete: (name: string) => http.delete(`/bridges/${name}`),
}
