<script setup lang="ts">
import { Loader2, Trash2, TriangleAlert } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'

interface Props {
  open: boolean
  title?: string
  description?: string
  loading?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  title: 'Are you sure?',
})

const emit = defineEmits<{
  'update:open': [value: boolean]
  confirm: []
}>()
</script>

<template>
  <Dialog :open="open" @update:open="emit('update:open', $event)">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <div class="flex items-center gap-3">
          <div class="flex size-10 shrink-0 items-center justify-center rounded-full bg-red-50">
            <TriangleAlert class="size-5 text-red-600" />
          </div>
          <div>
            <DialogTitle>{{ title }}</DialogTitle>
            <DialogDescription v-if="description" class="mt-0.5">
              {{ description }}
            </DialogDescription>
          </div>
        </div>
      </DialogHeader>

      <DialogFooter class="gap-2 sm:gap-2">
        <Button
          variant="outline"
          size="sm"
          class="h-9 font-normal"
          :disabled="loading"
          @click="emit('update:open', false)"
        >
          Cancel
        </Button>
        <Button
          size="sm"
          class="h-9 px-4 bg-red-600 hover:bg-red-700 text-white font-normal"
          :disabled="loading"
          @click="emit('confirm')"
        >
          <Loader2 v-if="loading" class="w-4 h-4 mr-2 animate-spin" />
          <Trash2 v-else class="w-4 h-4 mr-2" />
          Delete
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
