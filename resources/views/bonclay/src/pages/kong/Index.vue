<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Input } from '@/components/ui/input'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { listKongServices, updateKongService, type KongService } from '@/api'
import { useNotification } from '@/composables/useNotification'

const services = ref<KongService[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const searchQuery = ref('')
const togglingId = ref<string | null>(null)
const { success, error: notifyError } = useNotification()

const fetchServices = async () => {
  loading.value = true
  error.value = null
  try {
    const res = await listKongServices({ size: 1000 })
    services.value = res?.data ?? []
  } catch (err: any) {
    console.error(err)
    error.value = err?.message || 'Failed to fetch Kong services'
  } finally {
    loading.value = false
  }
}

onMounted(fetchServices)

const filteredServices = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return services.value
  return services.value.filter((s) => {
    const tagStr = (s.tags ?? []).join(',').toLowerCase()
    return (
      s.name?.toLowerCase().includes(q) ||
      s.host?.toLowerCase().includes(q) ||
      String(s.port ?? '').includes(q) ||
      s.protocol?.toLowerCase().includes(q) ||
      tagStr.includes(q)
    )
  })
})

const buildUpstream = (s: KongService) => {
  const path = s.path ?? ''
  return `${s.protocol}://${s.host}:${s.port}${path ?? ''}`
}

const formatDate = (unixSec: number) => {
  if (!unixSec) return '-'
  try {
    return new Date(unixSec * 1000).toLocaleString()
  } catch {
    return '-'
  }
}

const toggleEnabled = async (svc: KongService) => {
  const target = !svc.enabled
  togglingId.value = svc.id

  // optimistic update
  const idx = services.value.findIndex((x) => x.id === svc.id)
  if (idx >= 0) services.value[idx] = { ...services.value[idx], enabled: target }

  try {
    const updated = await updateKongService(svc.id, { enabled: target })
    if (idx >= 0 && updated) services.value[idx] = updated
    success(
      `Service ${target ? 'enabled' : 'disabled'}`,
      svc.name,
    )
  } catch (err: any) {
    // revert
    if (idx >= 0) services.value[idx] = { ...services.value[idx], enabled: !target }
    notifyError(
      'Failed to update service',
      err?.message || err?.details?.message || 'Unknown error',
    )
  } finally {
    togglingId.value = null
  }
}
</script>

<template>
  <div class="h-[90vh] w-full flex justify-center rounded-lg bg-muted/10">
    <div class="h-full w-full max-w-7xl">
      <div class="flex flex-col h-full">
        <!-- Header -->
        <div class="p-5 space-y-4 flex-shrink-0">
          <div class="w-full space-y-3">
            <div class="flex items-center justify-between">
              <h1 class="text-2xl font-semibold text-gray-900">Kong Services</h1>
              <Button variant="outline" :disabled="loading" @click="fetchServices">
                {{ loading ? 'Refreshing…' : 'Refresh' }}
              </Button>
            </div>
            <div class="w-12 h-0.5 bg-primary"></div>
          </div>

          <Input
            v-model="searchQuery"
            placeholder="Search by name, host, port, protocol or tag…"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />

          <div v-if="error" class="text-sm text-red-600">{{ error }}</div>

          <div class="text-sm text-gray-500">
            Total: <span class="font-medium text-gray-700">{{ filteredServices.length }}</span>
            <span v-if="filteredServices.length !== services.length">
              / {{ services.length }}
            </span>
          </div>
        </div>

        <!-- Table -->
        <div class="flex-1 min-h-0 p-5 pt-0">
          <ScrollArea class="h-full rounded-lg border bg-white">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead>Protocol</TableHead>
                  <TableHead>Upstream</TableHead>
                  <TableHead>Tags</TableHead>
                  <TableHead>Retries</TableHead>
                  <TableHead>Updated At</TableHead>
                  <TableHead class="text-right">Action</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                <TableRow v-for="s in filteredServices" :key="s.id">
                  <TableCell class="font-medium">{{ s.name }}</TableCell>
                  <TableCell>
                    <Badge :variant="s.enabled ? 'default' : 'secondary'">
                      {{ s.enabled ? 'enabled' : 'disabled' }}
                    </Badge>
                  </TableCell>
                  <TableCell class="uppercase">{{ s.protocol }}</TableCell>
                  <TableCell class="font-mono text-xs">{{ buildUpstream(s) }}</TableCell>
                  <TableCell>
                    <div class="flex flex-wrap gap-1">
                      <Badge
                        v-for="(tag, i) in (s.tags ?? [])"
                        :key="i"
                        variant="outline"
                        class="text-xs"
                      >
                        {{ tag.trim() }}
                      </Badge>
                      <span v-if="!s.tags || s.tags.length === 0" class="text-xs text-gray-400">—</span>
                    </div>
                  </TableCell>
                  <TableCell>{{ s.retries }}</TableCell>
                  <TableCell class="text-xs text-gray-600">{{ formatDate(s.updated_at) }}</TableCell>
                  <TableCell class="text-right">
                    <Button
                      :variant="s.enabled ? 'destructive' : 'default'"
                      size="sm"
                      :disabled="togglingId === s.id"
                      @click="toggleEnabled(s)"
                    >
                      <template v-if="togglingId === s.id">…</template>
                      <template v-else>{{ s.enabled ? 'Disable' : 'Enable' }}</template>
                    </Button>
                  </TableCell>
                </TableRow>

                <TableRow v-if="!loading && filteredServices.length === 0">
                  <TableCell colspan="8" class="text-center py-8 text-gray-500">
                    No Kong services found.
                  </TableCell>
                </TableRow>

                <TableRow v-if="loading && services.length === 0">
                  <TableCell colspan="8" class="text-center py-8 text-gray-500">
                    Loading…
                  </TableCell>
                </TableRow>
              </TableBody>
            </Table>
          </ScrollArea>
        </div>
      </div>
    </div>
  </div>
</template>
