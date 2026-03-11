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
import client from '@/api/bonClayHttpClient'

const emit = defineEmits<{
  created: []
}>()

const open = ref(false)
const submitting = ref(false)
const errorMessage = ref<string | null>(null)

const form = ref({
  name: '',
  desc: '',
})

watch(open, (isOpen) => {
  if (!isOpen) {
    form.value = { name: '', desc: '' }
    errorMessage.value = null
  }
})

const openDialog = () => {
  open.value = true
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
    await client.post('/v3/squads/create', {
      name: form.value.name.trim(),
      desc: form.value.desc.trim(),
    })

    open.value = false
    emit('created')
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to create squad'
    console.error('Failed to create squad:', err)
  } finally {
    submitting.value = false
  }
}

defineExpose({ openDialog })
</script>

<template>
  <Dialog v-model:open="open">
    <DialogContent class="sm:max-w-md">
      <DialogHeader>
        <DialogTitle>Create New Squad</DialogTitle>
        <DialogDescription>
          Add a new squad to organize your API collections.
        </DialogDescription>
      </DialogHeader>

      <form @submit.prevent="handleSubmit" class="flex flex-col gap-5 py-2">
        <!-- Name Field -->
        <div class="space-y-2">
          <Label for="squad-name">Name <span class="text-red-500">*</span></Label>
          <Input
            id="squad-name"
            v-model="form.name"
            placeholder="e.g. Backend Team"
            class="h-10 border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
            :disabled="submitting"
          />
        </div>

        <!-- Description Field -->
        <div class="space-y-2">
          <Label for="squad-desc">Description</Label>
          <textarea
            id="squad-desc"
            v-model="form.desc"
            placeholder="Brief description of this squad..."
            rows="3"
            :disabled="submitting"
            class="flex w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm shadow-sm placeholder:text-gray-400 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50 resize-none"
          />
        </div>

        <!-- Error Message -->
        <div v-if="errorMessage" class="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">
          {{ errorMessage }}
        </div>

        <DialogFooter class="pt-2">
          <Button type="button" variant="outline" size="sm" class="h-9 font-normal" :disabled="submitting" @click="open = false">
            Cancel
          </Button>
          <Button
            type="submit"
            size="sm"
            class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal"
            :disabled="submitting"
          >
            <Loader2 v-if="submitting" class="w-4 h-4 mr-2 animate-spin" />
            {{ submitting ? 'Creating...' : 'Create Squad' }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>
