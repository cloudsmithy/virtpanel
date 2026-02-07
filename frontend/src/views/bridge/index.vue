<template>
  <a-space direction="vertical" fill size="medium">
    <a-card title="网桥管理">
      <template #extra>
        <a-button type="primary" size="small" @click="showCreate = true">创建网桥</a-button>
      </template>
      <a-table :data="bridges" :loading="loading" row-key="name" :pagination="false">
        <template #columns>
          <a-table-column title="名称" data-index="name">
            <template #cell="{ record }"><code>{{ record.name }}</code></template>
          </a-table-column>
          <a-table-column title="状态">
            <template #cell="{ record }">
              <a-badge :status="record.up ? 'success' : 'default'" :text="record.up ? 'UP' : 'DOWN'" />
            </template>
          </a-table-column>
          <a-table-column title="IP 地址">
            <template #cell="{ record }">{{ record.ip || '-' }}</template>
          </a-table-column>
          <a-table-column title="从属网卡">
            <template #cell="{ record }">
              <a-tag v-for="s in (record.slaves || [])" :key="s" size="small" style="margin-right:4px">{{ s }}</a-tag>
              <span v-if="!record.slaves?.length" style="color:var(--apple-gray)">无</span>
            </template>
          </a-table-column>
          <a-table-column title="操作" :width="100">
            <template #cell="{ record }">
              <a-popconfirm content="删除网桥会将 IP 迁回物理网卡，确认？" @ok="doDelete(record.name)">
                <a-button size="small" status="danger">删除</a-button>
              </a-popconfirm>
            </template>
          </a-table-column>
        </template>
      </a-table>
      <a-empty v-if="!loading && bridges.length === 0" description="暂无 VirtPanel 管理的网桥" style="padding:40px 0" />
    </a-card>

    <a-modal v-model:visible="showCreate" title="创建网桥" :footer="false" unmount-on-close>
      <a-alert style="margin-bottom:16px">创建时绑定物理网卡会短暂中断网络（毫秒级），IP 将从物理网卡迁移到网桥</a-alert>
      <a-form :model="form" layout="vertical">
        <a-form-item label="名称" required extra="实际名称将加上 vp- 前缀">
          <a-input v-model="form.name" placeholder="br0" />
        </a-form-item>
        <a-form-item label="绑定物理网卡">
          <a-select v-model="form.slave_nic" placeholder="可选，不绑定则创建空网桥" allow-clear>
            <a-option v-for="nic in hostNICs" :key="nic.name" :value="nic.name">{{ nic.name }} ({{ nic.ip || 'no ip' }})</a-option>
          </a-select>
        </a-form-item>
        <a-form-item><a-button type="primary" :loading="creating" @click="onCreate">创建</a-button></a-form-item>
      </a-form>
    </a-modal>
  </a-space>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { bridgeApi, type Bridge } from '../../api/bridge'
import { hostApi } from '../../api/host'
import { Message } from '@arco-design/web-vue'
import { errMsg } from '../../api/http'

const bridges = ref<Bridge[]>([])
const loading = ref(false)
const showCreate = ref(false)
const creating = ref(false)
const form = reactive({ name: '', slave_nic: '' })
const hostNICs = ref<{ name: string; ip: string }[]>([])

const load = async () => { loading.value = true; try { bridges.value = await bridgeApi.list() || [] } catch(e: any) { Message.error(errMsg(e, '加载失败')) } loading.value = false }
const loadNICs = async () => { try { hostNICs.value = await hostApi.nics() } catch {} }

const onCreate = async () => {
  if (!form.name.trim()) { Message.warning('请输入名称'); return }
  creating.value = true
  try { await bridgeApi.create(form); Message.success('创建成功'); showCreate.value = false; Object.assign(form, { name: '', slave_nic: '' }); load() } catch(e: any) { Message.error(errMsg(e, '创建失败')) }
  creating.value = false
}

const doDelete = async (name: string) => {
  const short = name.replace(/^vp-/, '')
  try { await bridgeApi.delete(short); Message.success('已删除'); load() } catch(e: any) { Message.error(errMsg(e, '删除失败')) }
}

onMounted(() => { load(); loadNICs() })
</script>
