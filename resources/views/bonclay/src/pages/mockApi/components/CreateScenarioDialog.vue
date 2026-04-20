<script setup lang="ts">
import { ref, watch } from 'vue'
import { Loader2, Wand2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { createScenarioResponse } from '@/api'

interface Props {
  open: boolean
  endpointId: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  created: []
}>()

const submitting = ref(false)
const errorMessage = ref<string | null>(null)
const headerViewMode = ref<'edit' | 'preview'>('edit')
const bodyViewMode = ref<'edit' | 'preview'>('edit')

const form = ref({
  desc: '',
  status_header: 200,
  delay: 0,
  header: '{}',
  body: '{}',
})

watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) {
      form.value = { desc: '', status_header: 200, delay: 0, header: '{}', body: '{}' }
      errorMessage.value = null
      headerViewMode.value = 'edit'
      bodyViewMode.value = 'edit'
    }
  }
)

const highlightJson = (val: string): string => {
  if (!val.trim()) return val
  try {
    const formatted = JSON.stringify(JSON.parse(val), null, 2)
    return formatted
      .replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
      .replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, match => {
        if (/^"/.test(match)) return `<span class="${/:$/.test(match) ? 'text-violet-600' : 'text-emerald-600'}">${match}</span>`
        if (/true|false/.test(match)) return `<span class="text-amber-600">${match}</span>`
        if (/null/.test(match)) return `<span class="text-gray-400">${match}</span>`
        return `<span class="text-blue-600">${match}</span>`
      })
  } catch {
    return val.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;')
  }
}

const tryFormatJson = (field: 'body' | 'header') => {
  const val = form.value[field].trim()
  if (!val) return
  try {
    form.value[field] = JSON.stringify(JSON.parse(val), null, 2)
  } catch {
    // not valid JSON — leave as-is
  }
}

const isValidJson = (val: string) => {
  if (!val.trim()) return true
  try { JSON.parse(val); return true } catch { return false }
}

const handleSubmit = async () => {
  if (!form.value.desc.trim()) {
    errorMessage.value = 'Scenario name is required'
    return
  }
  if (!form.value.status_header || form.value.status_header < 100 || form.value.status_header > 599) {
    errorMessage.value = 'HTTP status must be between 100 and 599'
    return
  }
  if (!isValidJson(form.value.header)) {
    errorMessage.value = 'Headers must be valid JSON'
    return
  }

  submitting.value = true
  errorMessage.value = null
  try {
    await createScenarioResponse({
      endpoint_id: props.endpointId,
      desc: form.value.desc.trim(),
      status_header: form.value.status_header,
      delay: form.value.delay || undefined,
      body: form.value.body || undefined,
      header: form.value.header.trim() || undefined,
    })
    emit('update:open', false)
    emit('created')
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to create scenario'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-3xl">
      <DialogHeader>
        <DialogTitle>New Scenario</DialogTitle>
        <DialogDescription>Define a mock response scenario for this endpoint.</DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="flex flex-col gap-5 py-2">
        <div class="space-y-2">
          <Label for="sc-desc">Scenario Name <span class="text-red-500">*</span></Label>
          <Input
            id="sc-desc"
            v-model="form.desc"
            placeholder="e.g. Success response"
            :disabled="submitting"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <div class="flex gap-3">
          <div class="space-y-2 flex-1">
            <Label for="sc-status">HTTP Status <span class="text-red-500">*</span></Label>
            <Input
              id="sc-status"
              v-model.number="form.status_header"
              type="number"
              min="100"
              max="599"
              placeholder="200"
              :disabled="submitting"
              class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
            />
          </div>

          <div class="space-y-2 flex-1">
            <Label for="sc-delay">Delay (ms)</Label>
            <Input
              id="sc-delay"
              v-model.number="form.delay"
              type="number"
              min="0"
              placeholder="0"
              :disabled="submitting"
              class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
            />
          </div>
        </div>

        <!-- Response Headers -->
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <Label for="sc-header">
              Response Headers <span class="text-xs text-gray-400">(JSON)</span>
            </Label>
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="tryFormatJson('header')"
                class="inline-flex items-center gap-1 text-xs text-gray-500 hover:text-gray-800 transition-colors"
              >
                <Wand2 class="w-3 h-3" />
                Format
              </button>
              <div class="flex border border-gray-200 rounded-md overflow-hidden text-xs">
                <button
                  type="button"
                  @click="headerViewMode = 'edit'"
                  :class="['px-2 py-0.5 transition-colors', headerViewMode === 'edit' ? 'bg-gray-900 text-white' : 'text-gray-500 hover:bg-gray-100']"
                >Edit</button>
                <button
                  type="button"
                  @click="headerViewMode = 'preview'"
                  :class="['px-2 py-0.5 transition-colors', headerViewMode === 'preview' ? 'bg-gray-900 text-white' : 'text-gray-500 hover:bg-gray-100']"
                >Preview</button>
              </div>
            </div>
          </div>
          <pre
            v-if="headerViewMode === 'preview'"
            class="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm font-mono h-24 max-h-40 overflow-y-auto"
            v-html="highlightJson(form.header)"
          />
          <textarea
            v-else
            id="sc-header"
            v-model="form.header"
            placeholder='{"Content-Type": "application/json"}'
            :disabled="submitting"
            class="flex w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm font-mono shadow-sm placeholder:text-gray-400 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 resize-y h-24 max-h-40 overflow-y-auto"
          />
        </div>

        <!-- Response Body -->
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <Label for="sc-body">Response Body</Label>
            <div class="flex items-center gap-2">
              <button
                type="button"
                @click="tryFormatJson('body')"
                class="inline-flex items-center gap-1 text-xs text-gray-500 hover:text-gray-800 transition-colors"
              >
                <Wand2 class="w-3 h-3" />
                Format
              </button>
              <div class="flex border border-gray-200 rounded-md overflow-hidden text-xs">
                <button
                  type="button"
                  @click="bodyViewMode = 'edit'"
                  :class="['px-2 py-0.5 transition-colors', bodyViewMode === 'edit' ? 'bg-gray-900 text-white' : 'text-gray-500 hover:bg-gray-100']"
                >Edit</button>
                <button
                  type="button"
                  @click="bodyViewMode = 'preview'"
                  :class="['px-2 py-0.5 transition-colors', bodyViewMode === 'preview' ? 'bg-gray-900 text-white' : 'text-gray-500 hover:bg-gray-100']"
                >Preview</button>
              </div>
            </div>
          </div>
          <pre
            v-if="bodyViewMode === 'preview'"
            class="w-full rounded-lg border border-gray-200 bg-gray-50 px-3 py-2 text-sm font-mono h-48 max-h-64 overflow-y-auto"
            v-html="highlightJson(form.body)"
          />
          <textarea
            v-else
            id="sc-body"
            v-model="form.body"
            placeholder='{"message": "ok"}'
            :disabled="submitting"
            class="flex w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm font-mono shadow-sm placeholder:text-gray-400 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 resize-y h-48 max-h-64 overflow-y-auto"
          />
        </div>

        <div v-if="errorMessage" class="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">
          {{ errorMessage }}
        </div>

        <DialogFooter class="pt-2">
          <Button type="button" variant="outline" size="sm" class="h-9 font-normal" :disabled="submitting" @click="emit('update:open', false)">
            Cancel
          </Button>
          <Button type="submit" size="sm" class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal" :disabled="submitting">
            <Loader2 v-if="submitting" class="w-4 h-4 mr-2 animate-spin" />
            {{ submitting ? 'Creating...' : 'Create Scenario' }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
