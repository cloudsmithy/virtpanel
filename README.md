# KVMMM

基于 KVM/libvirt 的轻量虚拟机管理面板。Go 后端 + Vue 3 前端。

## 功能

- 虚拟机管理（创建/启动/关机/重启/暂停/克隆/删除）
- 创建时可选系统预设（Linux / Windows / 兼容模式）
  - 自动配置芯片组（i440FX / Q35）、CPU 模式、时钟、磁盘总线、网卡类型
  - 支持手动覆盖每项参数
- VNC 控制台（浏览器内 noVNC）
- 磁盘热挂载/卸载，ISO 挂载/弹出
- 网卡热添加/移除
- 快照管理（创建/恢复/删除）
- 网络管理（NAT 网络创建/启停/删除）
- 存储池和存储卷管理
- ISO 镜像上传和管理
- 仪表盘（主机信息、资源概览）

## 技术栈

| 组件 | 技术 |
|------|------|
| 后端 | Go + Gin + go-libvirt |
| 前端 | Vue 3 + TypeScript + Arco Design + noVNC |
| 虚拟化 | KVM / QEMU / libvirt |

## 环境要求

- Linux 主机，已安装 libvirt、QEMU-KVM
- Go 1.22+
- Node.js 18+、pnpm

## 快速开始

```bash
# 后端
cd backend
go build -o kvmmm ./cmd/main.go
./kvmmm   # 监听 :8080

# 前端（开发）
cd frontend
pnpm install
pnpm dev   # 监听 :5173，自动代理 /api 和 /ws 到后端

# 前端（生产构建）
pnpm build  # 输出到 dist/，用 nginx 反代即可
```

## 项目结构

```
kvmmm/
├── backend/
│   ├── cmd/main.go              # 入口
│   ├── internal/
│   │   ├── handler/             # HTTP 路由处理
│   │   ├── service/             # libvirt 业务逻辑
│   │   └── model/               # 数据模型
│   ├── go.mod
│   └── go.sum
├── frontend/
│   ├── src/
│   │   ├── api/                 # API 请求封装
│   │   ├── views/               # 页面组件
│   │   ├── layout/              # 布局
│   │   ├── router/              # 路由
│   │   └── styles/              # 全局样式
│   ├── package.json
│   └── vite.config.ts
└── .gitignore
```

## Nginx 生产部署示例

```nginx
server {
    listen 80;

    location / {
        root /path/to/frontend/dist;
        try_files $uri $uri/ /index.html;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
    }

    location /ws/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## API 概览

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/host/info | 主机信息 |
| GET | /api/vms | 虚拟机列表 |
| POST | /api/vms | 创建虚拟机 |
| POST | /api/vms/:name/start | 启动 |
| POST | /api/vms/:name/shutdown | 关机 |
| POST | /api/vms/:name/destroy | 强制关机 |
| DELETE | /api/vms/:name | 删除 |
| GET | /api/vms/:name/detail | 虚拟机详情 |
| POST | /api/vms/:name/iso | 挂载 ISO |
| POST | /api/vms/:name/clone | 克隆 |
| GET | /ws/vnc/:name | VNC WebSocket |

完整 API 见 `backend/cmd/main.go` 路由定义。

## License

MIT
