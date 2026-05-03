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
import { updateSquad, deleteSquad } from '@/api'
import { useSquadSession } from '@/composables/useSquadSession'

interface Props {
  open: boolean
  squadId: string
  initialName: string
  initialDesc: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
  updated: []
  deleted: []
}>()

const { activeSquadId, clearActiveSquad } = useSquadSession()

const submitting = ref(false)
const errorMessage = ref<string | null>(null)
const confirmDelete = ref(false)
const deleting = ref(false)

const form = ref({ name: '', desc: '' })

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      form.value = { name: props.initialName, desc: props.initialDesc }
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
    await deleteSquad(props.squadId)
    if (activeSquadId.value === props.squadId) {
      clearActiveSquad()
    }
    emit('update:open', false)
    emit('deleted')
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to delete squad'
  } finally {
    deleting.value = false
    confirmDelete.value = false
  }
}

const handleSubmit = async () => {
  if (!form.value.name.trim()) {
    errorMessage.value = 'Squad name is required'
    return
  }
  if (form.value.name.trim().length < 3) {
    errorMessage.value = 'Squad name must be at least 3 characters'
    return
  }
  submitting.value = true
  errorMessage.value = null
  try {
    await updateSquad(props.squadId, {
      name: form.value.name.trim(),
      desc: form.value.desc.trim(),
    })
    emit('update:open', false)
    emit('updated')
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to update squad'
  } finally {
    submitting.value = false
  }
}
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Edit Squad</DialogTitle>
        <DialogDescription>Update squad details or permanently delete it.</DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="flex flex-col gap-5 py-2">
        <div class="space-y-2">
          <Label for="edit-squad-name">Name <span class="text-red-500">*</span></Label>
          <Input
            id="edit-squad-name"
            v-model="form.name"
            placeholder="e.g. Backend Team"
            :disabled="submitting"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
          />
        </div>

        <div class="space-y-2">
          <Label for="edit-squad-desc">Description</Label>
          <textarea
            id="edit-squad-desc"
            v-model="form.desc"
            placeholder="Brief description of this squad..."
            rows="3"
            :disabled="submitting"
            class="flex w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-gray-400 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 resize-none"
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
                : 'h-9 px-4 font-normal text-red-600 border border-red-200 bg-transparent hover:bg-red-50 hover:border-red-300'"
              @click="handleDelete"
            >
              <Loader2 v-if="deleting" class="w-4 h-4 mr-2 animate-spin" />
              <Trash2 v-else class="w-4 h-4 mr-2" />
              {{ confirmDelete ? 'Confirm delete?' : 'Delete Squad' }}
            </Button>
            <div class="flex gap-2">
              <Button
                type="button"
                variant="outline"
                size="sm"
                class="h-9 font-normal"
                :disabled="submitting || deleting"
                @click="confirmDelete = false; emit('update:open', false)"
              >
                Cancel
              </Button>
              <Button
                type="submit"
                size="sm"
                class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal"
                :disabled="submitting || deleting"
              >
                <Loader2 v-if="submitting" class="w-4 h-4 mr-2 animate-spin" />
                {{ submitting ? 'Saving...' : 'Save Changes' }}
              </Button>
            </div>
          </div>
          <p v-if="confirmDelete" class="text-xs text-red-500 text-left">
            This will permanently delete the squad and all its collections, endpoints, and scenarios. Click again to confirm.
          </p>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
