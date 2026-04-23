<script setup lang="ts">
import { ref, watch } from 'vue'
import { Loader2, Info, Trash2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { updateEndpoint, deleteEndpoint } from '@/api'

interface Props {
  open: boolean
  endpointId: string
  initialPath: string
  initialMethod: string
  initialCategory: string
  initialDesc: string
  initialDelay?: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  updated: []
  deleted: []
}>()

const HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE']

const submitting = ref(false)
const errorMessage = ref<string | null>(null)
const confirmDelete = ref(false)
const deleting = ref(false)

const form = ref({ path: '', method: 'GET', category: '', desc: '', delay: null as number | null })

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      form.value = {
        path: props.initialPath.split('[^/]+').join('{query_params}'),
        method: props.initialMethod,
        category: props.initialCategory,
        desc: props.initialDesc,
        delay: props.initialDelay ?? null,
      }
      errorMessage.value = null
      confirmDelete.value = false
    }
  }
)

const handleDelete = async () => {
  if (!confirmDelete.value) {
    confirmDelete.value = true
    return
  }
  deleting.value = true
  try {
    await deleteEndpoint(props.endpointId)
    emit('update:open', false)
    emit('deleted')
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to delete endpoint'
  } finally {
    deleting.value = false
    confirmDelete.value = false
  }
}

const handleSubmit = async () => {
  if (!form.value.path.trim()) {
    errorMessage.value = 'Path is required'
    return
  }
  if (!form.value.path.startsWith('/')) {
    errorMessage.value = 'Path must start with /'
    return
  }

  submitting.value = true
  errorMessage.value = null
  try {
    await updateEndpoint(props.endpointId, {
      path: form.value.path.trim(),
      method: form.value.method,
      category: form.value.category.trim(),
      desc: form.value.desc.trim(),
      delay: form.value.delay ?? undefined,
    })
    emit('update:open', false)
    emit('updated')
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to update endpoint'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Edit Endpoint</DialogTitle>
        <DialogDescription>Update the endpoint details.</DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="flex flex-col gap-5 py-2">
        <div class="flex gap-3">
          <div class="space-y-2 w-32 flex-shrink-0">
            <Label for="edit-ep-method">Method <span class="text-red-500">*</span></Label>
            <select
              id="edit-ep-method"
              v-model="form.method"
              :disabled="submitting"
              class="flex h-10 w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm shadow-sm focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option v-for="m in HTTP_METHODS" :key="m" :value="m">{{ m }}</option>
            </select>
          </div>

          <div class="space-y-2 flex-1">
            <div class="flex items-center gap-1.5">
              <Label for="edit-ep-path">Path <span class="text-red-500">*</span></Label>
              <Tooltip>
                <TooltipTrigger as-child>
                  <Info class="w-3.5 h-3.5 text-gray-400 cursor-help" />
                </TooltipTrigger>
                <TooltipContent class="max-w-xs text-xs leading-relaxed">
                  Use <code class="font-mono">{query_params}</code> to match any query string segment.<br />
                  e.g. <code class="font-mono">/get-user/{query_params}</code> or <code class="font-mono">/user/{query_params}/remove</code>
                </TooltipContent>
              </Tooltip>
            </div>
            <Input
              id="edit-ep-path"
              v-model="form.path"
              placeholder="/users/:id"
              :disabled="submitting"
              class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
            />
          </div>
        </div>

        <div class="space-y-2">
          <Label for="edit-ep-category">Category</Label>
          <Input
            id="edit-ep-category"
            v-model="form.category"
            placeholder="e.g. Users"
            :disabled="submitting"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <div class="space-y-2">
          <Label for="edit-ep-delay">Proxy Delay (seconds)</Label>
          <Input
            id="edit-ep-delay"
            :value="form.delay ?? undefined"
            @input="(e: Event) => { const v = (e.target as HTMLInputElement).valueAsNumber; form.delay = isNaN(v) ? null : v }"
            type="number"
            min="0"
            placeholder="0"
            :disabled="submitting"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <div class="space-y-2">
          <Label for="edit-ep-desc">Description</Label>
          <textarea
            id="edit-ep-desc"
            v-model="form.desc"
            placeholder="What does this endpoint do?"
            rows="3"
            :disabled="submitting"
            class="flex w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-gray-400 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 resize-none"
          />
        </div>

        <div v-if="errorMessage" class="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">
          {{ errorMessage }}
        </div>

        <DialogFooter class="pt-2 flex-col gap-2">
          <div class="flex justify-between w-full gap-2">
            <Button
              type="button"
              size="sm"
              :disabled="submitting || deleting"
              :class="confirmDelete
                ? 'h-9 px-4 font-normal bg-red-600 hover:bg-red-700 text-white'
                : 'h-9 px-4 font-normal text-red-600 border border-red-200 hover:bg-red-50'"
              @click="handleDelete"
            >
              <Loader2 v-if="deleting" class="w-4 h-4 mr-2 animate-spin" />
              <Trash2 v-else class="w-4 h-4 mr-2" />
              {{ confirmDelete ? 'Confirm delete?' : 'Delete Endpoint' }}
            </Button>
            <div class="flex gap-2">
              <Button type="button" variant="outline" size="sm" class="h-9 font-normal" :disabled="submitting || deleting" @click="confirmDelete = false; emit('update:open', false)">
                Cancel
              </Button>
              <Button type="submit" size="sm" class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal" :disabled="submitting || deleting">
                <Loader2 v-if="submitting" class="w-4 h-4 mr-2 animate-spin" />
                {{ submitting ? 'Saving...' : 'Save Changes' }}
              </Button>
            </div>
          </div>
          <p v-if="confirmDelete" class="text-xs text-red-500 text-left w-full">
            This will permanently delete the endpoint and all its scenarios. Click again to confirm.
          </p>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
