<template>
  <div class="dash">
    <!-- Top: 3 ring cards -->
    <div class="rings">
      <div class="ring-card">
        <div class="ring-title">CPU</div>
        <a-progress type="circle" :percent="host?.cpu_usage ?? 0" :stroke-width="8" :width="130" color="#007AFF">
          <template #text="{ percent }">
            <div class="ring-inner">
              <div class="ring-num">{{ percent }}%</div>
            </div>
          </template>
        </a-progress>
        <div class="ring-sub">{{ host?.cpu_count ?? 0 }} 核</div>
      </div>
      <div class="ring-card">
        <div class="ring-title">内存</div>
        <a-progress type="circle" :percent="memPercent" :stroke-width="8" :width="130" color="#FF9500">
          <template #text>
            <div class="ring-inner">
              <div class="ring-num">{{ memPercent }}%</div>
            </div>
          </template>
        </a-progress>
        <div class="ring-sub">{{ memUsed }} / {{ memTotal }} GB</div>
      </div>
      <div class="ring-card">
        <div class="ring-title">虚拟机</div>
        <a-progress type="circle" :percent="vmPercent" :stroke-width="8" :width="130" color="#34C759">
          <template #text>
            <div class="ring-inner">
              <div class="ring-num">{{ host?.vm_running ?? 0 }}</div>
            </div>
          </template>
        </a-progress>
        <div class="ring-sub">运行 {{ host?.vm_running ?? 0 }} / 共 {{ host?.vm_total ?? 0 }} 台</div>
      </div>
    </div>

    <!-- Bottom: host info -->
    <div class="info-card">
      <div class="info-row">
        <span class="info-label">主机名</span>
        <span class="info-value">{{ host?.hostname }}</span>
      </div>
      <div class="info-row">
        <span class="info-label">处理器</span>
        <span class="info-value">{{ host?.cpu_model }}</span>
      </div>
      <div class="info-row">
        <span class="info-label">核心数</span>
        <span class="info-value">{{ host?.cpu_count }} 核</span>
      </div>
      <div class="info-row">
        <span class="info-label">总内存</span>
        <span class="info-value">{{ memTotal }} GB</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onBeforeUnmount } from 'vue'
import { hostApi, type HostInfo } from '../../api/host'

const host = ref<HostInfo | null>(null)

const memTotal = computed(() => host.value ? (host.value.memory_total / 1024).toFixed(1) : '0')
const memUsed = computed(() => host.value ? ((host.value.memory_total - host.value.memory_free) / 1024).toFixed(1) : '0')
const memPercent = computed(() => {
  if (!host.value || !host.value.memory_total) return 0
  return Math.round(((host.value.memory_total - host.value.memory_free) / host.value.memory_total) * 100)
})
const vmPercent = computed(() => {
  if (!host.value || !host.value.vm_total) return 0
  return Math.round((host.value.vm_running / host.value.vm_total) * 100)
})

const load = async () => { try { host.value = await hostApi.info() } catch {} }
let timer: ReturnType<typeof setInterval> | null = null
onMounted(() => { load(); timer = setInterval(load, 3000) })
onBeforeUnmount(() => { if (timer) clearInterval(timer) })
</script>

<style scoped>
.dash { max-width: 960px; }

.rings {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 16px;
}
.ring-card {
  background: #fff;
  border-radius: 12px;
  padding: 28px 20px 22px;
  display: flex;
  flex-direction: column;
  align-items: center;
  box-shadow: 0 0.5px 1px rgba(0,0,0,0.04), 0 1px 4px rgba(0,0,0,0.03);
}
.ring-title {
  font-size: 13px;
  font-weight: 600;
  color: #8e8e93;
  margin-bottom: 16px;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}
.ring-inner { text-align: center; }
.ring-num {
  font-size: 28px;
  font-weight: 700;
  color: #1c1c1e;
  letter-spacing: -1px;
}
.ring-sub {
  font-size: 12px;
  color: #8e8e93;
  margin-top: 12px;
}

.info-card {
  background: #fff;
  border-radius: 12px;
  padding: 4px 0;
  box-shadow: 0 0.5px 1px rgba(0,0,0,0.04), 0 1px 4px rgba(0,0,0,0.03);
}
.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 13px 20px;
  border-bottom: 0.5px solid rgba(0,0,0,0.04);
}
.info-row:last-child { border-bottom: none; }
.info-label {
  font-size: 13px;
  color: #1c1c1e;
  font-weight: 500;
}
.info-value {
  font-size: 13px;
  color: #8e8e93;
}
</style>
