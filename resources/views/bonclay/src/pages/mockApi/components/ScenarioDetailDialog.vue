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
import { updateScenarioResponse, type UpdateScenarioPayload } from '@/api'
import type { ScenarioResponse } from '@/types/api.types'

interface Props {
  open: boolean
  scenario: ScenarioResponse | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  updated: []
}>()

const submitting = ref(false)
const errorMessage = ref<string | null>(null)

const form = ref({
  desc: '',
  status_header: 200,
  delay: 0,
  header: '',
  body: '',
})

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen && props.scenario) {
      form.value = {
        desc: props.scenario.desc,
        status_header: props.scenario.status_header,
        delay: props.scenario.delay ?? 0,
        header: props.scenario.header ?? '',
        body: props.scenario.body ?? '',
      }
      errorMessage.value = null
    }
  }
)

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
  if (!props.scenario) return

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
    const payload: UpdateScenarioPayload = {
      endpoint_id: props.scenario.endpoint_id,
      desc: form.value.desc.trim(),
      status_header: form.value.status_header,
      delay: form.value.delay ?? undefined,
      header: form.value.header.trim() || undefined,
      body: form.value.body || undefined,
    }
    await updateScenarioResponse(props.scenario.id, payload)
    emit('update:open', false)
    emit('updated')
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to update scenario'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-lg">
      <DialogHeader>
        <DialogTitle>Edit Scenario</DialogTitle>
        <DialogDescription>Update the mock response scenario.</DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="flex flex-col gap-5 py-2">
        <div class="space-y-2">
          <Label for="sc-edit-desc">Scenario Name <span class="text-red-500">*</span></Label>
          <Input
            id="sc-edit-desc"
            v-model="form.desc"
            placeholder="e.g. Success response"
            :disabled="submitting"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <div class="flex gap-3">
          <div class="space-y-2 flex-1">
            <Label for="sc-edit-status">HTTP Status <span class="text-red-500">*</span></Label>
            <Input
              id="sc-edit-status"
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
            <Label for="sc-edit-delay">Delay (ms)</Label>
            <Input
              id="sc-edit-delay"
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
            <Label for="sc-edit-header">
              Response Headers <span class="text-xs text-gray-400">(JSON)</span>
            </Label>
            <button
              type="button"
              @click="tryFormatJson('header')"
              class="inline-flex items-center gap-1 text-xs text-gray-500 hover:text-gray-800 transition-colors"
            >
              <Wand2 class="w-3 h-3" />
              Format JSON
            </button>
          </div>
          <textarea
            id="sc-edit-header"
            v-model="form.header"
            placeholder='{"Content-Type": "application/json"}'
            rows="2"
            :disabled="submitting"
            class="flex w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm font-mono shadow-sm placeholder:text-gray-400 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 resize-none"
          />
        </div>

        <!-- Response Body -->
        <div class="space-y-2">
          <div class="flex items-center justify-between">
            <Label for="sc-edit-body">Response Body</Label>
            <button
              type="button"
              @click="tryFormatJson('body')"
              class="inline-flex items-center gap-1 text-xs text-gray-500 hover:text-gray-800 transition-colors"
            >
              <Wand2 class="w-3 h-3" />
              Format JSON
            </button>
          </div>
          <textarea
            id="sc-edit-body"
            v-model="form.body"
            placeholder='{"message": "ok"}'
            rows="5"
            :disabled="submitting"
            class="flex w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm font-mono shadow-sm placeholder:text-gray-400 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 resize-none"
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
            {{ submitting ? 'Saving...' : 'Save Changes' }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
