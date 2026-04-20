<script setup lang="ts">
import { ref, watch } from 'vue'
import { Loader2 } from 'lucide-vue-next'
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
import { createEndpoint } from '@/api'
import type { CollectionEndpoint } from '@/types/api.types'

interface Props {
  open: boolean
  collectionId: string
  existingEndpoints?: CollectionEndpoint[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  created: []
}>()

const HTTP_METHODS = ['GET', 'POST', 'PUT', 'PATCH', 'DELETE']

const submitting = ref(false)
const errorMessage = ref<string | null>(null)

const form = ref({
  path: '',
  method: 'GET',
  category: '',
  desc: '',
})

watch(
  () => props.open,
  (isOpen) => {
    if (!isOpen) {
      form.value = { path: '', method: 'GET', category: '', desc: '' }
      errorMessage.value = null
    }
  }
)

const handleSubmit = async () => {
  if (!form.value.path.trim()) {
    errorMessage.value = 'Endpoint path is required'
    return
  }
  if (!form.value.path.startsWith('/')) {
    errorMessage.value = 'Path must start with /'
    return
  }

  const duplicate = (props.existingEndpoints ?? []).some(
    ep => ep.method.toUpperCase() === form.value.method && ep.path === form.value.path.trim()
  )
  if (duplicate) {
    errorMessage.value = `${form.value.method} ${form.value.path.trim()} already exists in this collection`
    return
  }

  submitting.value = true
  errorMessage.value = null
  try {
    await createEndpoint({
      collection_id: props.collectionId,
      path: form.value.path.trim(),
      method: form.value.method,
      category: form.value.category.trim() || undefined,
      desc: form.value.desc.trim() || undefined,
    })
    emit('update:open', false)
    emit('created')
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to create endpoint'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>New Endpoint</DialogTitle>
        <DialogDescription>Add a new mock endpoint to this collection.</DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="flex flex-col gap-5 py-2">
        <div class="flex gap-3">
          <div class="space-y-2 w-32 flex-shrink-0">
            <Label for="ep-method">Method <span class="text-red-500">*</span></Label>
            <select
              id="ep-method"
              v-model="form.method"
              :disabled="submitting"
              class="flex h-10 w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm shadow-sm focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
            >
              <option v-for="m in HTTP_METHODS" :key="m" :value="m">{{ m }}</option>
            </select>
          </div>

          <div class="space-y-2 flex-1">
            <Label for="ep-path">Path <span class="text-red-500">*</span></Label>
            <Input
              id="ep-path"
              v-model="form.path"
              placeholder="/users/:id"
              :disabled="submitting"
              class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
            />
          </div>
        </div>

        <div class="space-y-2">
          <Label for="ep-category">Category</Label>
          <Input
            id="ep-category"
            v-model="form.category"
            placeholder="e.g. Users"
            :disabled="submitting"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <div class="space-y-2">
          <Label for="ep-desc">Description</Label>
          <textarea
            id="ep-desc"
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

        <DialogFooter class="pt-2">
          <Button type="button" variant="outline" size="sm" class="h-9 font-normal" :disabled="submitting" @click="emit('update:open', false)">
            Cancel
          </Button>
          <Button type="submit" size="sm" class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal" :disabled="submitting">
            <Loader2 v-if="submitting" class="w-4 h-4 mr-2 animate-spin" />
            {{ submitting ? 'Creating...' : 'Create Endpoint' }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
