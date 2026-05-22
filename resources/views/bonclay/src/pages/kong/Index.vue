<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import { Loader2, TriangleAlert } from 'lucide-vue-next'
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
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { listKongServices, updateKongService, type KongService } from '@/api'
import { useNotification } from '@/composables/useNotification'

const services = ref<KongService[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const searchQuery = ref('')
const togglingId = ref<string | null>(null)
const { success, error: notifyError } = useNotification()

// Confirmation dialog state
const confirmOpen = ref(false)
const confirmTarget = ref<KongService | null>(null)
const confirmLoading = ref(false)

const confirmIsEnable = computed(() => !(confirmTarget.value?.enabled ?? false))
const confirmAction = computed(() => (confirmIsEnable.value ? 'Enable' : 'Disable'))

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
  return services.value.filter((s) => s.name?.toLowerCase().includes(q))
})

const requestToggle = (svc: KongService) => {
  confirmTarget.value = svc
  confirmOpen.value = true
}

const cancelToggle = () => {
  if (confirmLoading.value) return
  confirmOpen.value = false
  confirmTarget.value = null
}

const confirmToggle = async () => {
  const svc = confirmTarget.value
  if (!svc) return

  const target = !svc.enabled
  confirmLoading.value = true
  togglingId.value = svc.id

  try {
    const updated = await updateKongService(svc.id, { enabled: target })
    const idx = services.value.findIndex((x) => x.id === svc.id)
    if (idx >= 0) {
      services.value[idx] = updated ?? { ...services.value[idx], enabled: target }
    }
    success(
      `Service ${target ? 'enabled' : 'disabled'}`,
      svc.name,
    )
    confirmOpen.value = false
    confirmTarget.value = null
  } catch (err: any) {
    notifyError(
      'Failed to update service',
      err?.message || err?.details?.message || 'Unknown error',
    )
  } finally {
    togglingId.value = null
    confirmLoading.value = false
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
            placeholder="Search by name…"
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

        <!-- Confirmation dialog -->
        <Dialog
          :open="confirmOpen"
          @update:open="(v: boolean) => (v ? (confirmOpen = true) : cancelToggle())"
        >
          <DialogContent class="sm:max-w-sm">
            <DialogHeader>
              <div class="flex items-center gap-3">
                <div
                  class="flex size-10 shrink-0 items-center justify-center rounded-full"
                  :class="confirmIsEnable ? 'bg-green-50' : 'bg-red-50'"
                >
                  <TriangleAlert
                    class="size-5"
                    :class="confirmIsEnable ? 'text-green-600' : 'text-red-600'"
                  />
                </div>
                <div>
                  <DialogTitle>{{ confirmAction }} service?</DialogTitle>
                  <DialogDescription class="mt-0.5">
                    Are you sure you want to {{ confirmAction.toLowerCase() }}
                    <span class="font-medium text-gray-900">{{ confirmTarget?.name }}</span
                    >?
                  </DialogDescription>
                </div>
              </div>
            </DialogHeader>

            <DialogFooter class="gap-2 sm:gap-2">
              <Button
                variant="outline"
                size="sm"
                class="h-9 font-normal"
                :disabled="confirmLoading"
                @click="cancelToggle"
              >
                Cancel
              </Button>
              <Button
                size="sm"
                class="h-9 px-4 font-normal text-white"
                :class="confirmIsEnable ? 'bg-green-600 hover:bg-green-700' : 'bg-red-600 hover:bg-red-700'"
                :disabled="confirmLoading"
                @click="confirmToggle"
              >
                <Loader2 v-if="confirmLoading" class="w-4 h-4 mr-2 animate-spin" />
                {{ confirmAction }}
              </Button>
            </DialogFooter>
          </DialogContent>
        </Dialog>

        <!-- Table -->
        <div class="flex-1 min-h-0 p-5 pt-0">
          <ScrollArea class="h-full rounded-lg border bg-white">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>Status</TableHead>
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
                  <TableCell class="text-right">
                    <Button
                      :variant="s.enabled ? 'destructive' : 'default'"
                      size="sm"
                      :disabled="togglingId === s.id"
                      @click="requestToggle(s)"
                    >
                      <template v-if="togglingId === s.id">…</template>
                      <template v-else>{{ s.enabled ? 'Disable' : 'Enable' }}</template>
                    </Button>
                  </TableCell>
                </TableRow>

                <TableRow v-if="!loading && filteredServices.length === 0">
                  <TableCell colspan="3" class="text-center py-8 text-gray-500">
                    No Kong services found.
                  </TableCell>
                </TableRow>

                <TableRow v-if="loading && services.length === 0">
                  <TableCell colspan="3" class="text-center py-8 text-gray-500">
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
