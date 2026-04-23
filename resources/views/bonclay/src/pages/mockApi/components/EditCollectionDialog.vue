<script setup lang="ts">
import { ref, watch } from 'vue'
import { Loader2, Trash2 } from 'lucide-vue-next'
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
import { updateCollection, deleteCollection } from '@/api'
import { useRouter } from 'vue-router'

interface Props {
  open: boolean
  collectionId: string
  initialName: string
  initialDesc: string
  initialDocs: string
  initialForwardProxyUrl: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  updated: []
  deleted: []
}>()

const router = useRouter()

const submitting = ref(false)
const errorMessage = ref<string | null>(null)
const confirmDelete = ref(false)
const deleting = ref(false)

const handleDelete = async () => {
  if (!confirmDelete.value) {
    confirmDelete.value = true
    return
  }
  deleting.value = true
  try {
    await deleteCollection(props.collectionId)
    emit('update:open', false)
    emit('deleted')
    router.push({ name: 'MockApiTools-Index' })
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to delete collection'
  } finally {
    deleting.value = false
    confirmDelete.value = false
  }
}

const form = ref({
  name: '',
  desc: '',
  docs: '',
  forward_proxy_url: '',
})

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      form.value = {
        name: props.initialName,
        desc: props.initialDesc,
        docs: props.initialDocs,
        forward_proxy_url: props.initialForwardProxyUrl,
      }
      errorMessage.value = null
      confirmDelete.value = false
    }
  }
)

const handleSubmit = async () => {
  if (!form.value.name.trim()) {
    errorMessage.value = 'Collection name is required'
    return
  }
  if (form.value.name.trim().length < 3) {
    errorMessage.value = 'Collection name must be at least 3 characters'
    return
  }

  submitting.value = true
  errorMessage.value = null
  try {
    await updateCollection(props.collectionId, {
      name: form.value.name.trim(),
      desc: form.value.desc.trim(),
      docs: form.value.docs.trim(),
      forward_proxy_url: form.value.forward_proxy_url.trim(),
    })
    emit('update:open', false)
    emit('updated')
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to update collection'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Edit Collection</DialogTitle>
        <DialogDescription>Update the collection details.</DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="flex flex-col gap-5 py-2">
        <div class="space-y-2">
          <Label for="edit-name">Name <span class="text-red-500">*</span></Label>
          <Input
            id="edit-name"
            v-model="form.name"
            placeholder="e.g. Payment API"
            :disabled="submitting"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <div class="space-y-2">
          <Label for="edit-desc">Description</Label>
          <textarea
            id="edit-desc"
            v-model="form.desc"
            placeholder="Brief description of this collection..."
            rows="3"
            :disabled="submitting"
            class="flex w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-gray-400 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 resize-none"
          />
        </div>

        <div class="space-y-2">
          <Label for="edit-docs">Documentation URL</Label>
          <Input
            id="edit-docs"
            v-model="form.docs"
            placeholder="https://docs.example.com"
            :disabled="submitting"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <div class="space-y-2">
          <Label for="edit-proxy">Forward Proxy URL</Label>
          <Input
            id="edit-proxy"
            v-model="form.forward_proxy_url"
            placeholder="https://api.example.com"
            :disabled="submitting"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <div v-if="errorMessage" class="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">
          {{ errorMessage }}
        </div>

        <DialogFooter class="pt-2 flex-col gap-3">
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
              {{ confirmDelete ? 'Confirm delete?' : 'Delete Collection' }}
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
          <p v-if="confirmDelete" class="text-xs text-red-500 text-left">
            This will permanently delete the collection and all its endpoints and scenarios. Click again to confirm.
          </p>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
