<script setup lang="ts">
import { ref, watch } from 'vue'
import { Plus, Trash2, Pencil, Check, X, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import {
  Sheet,
  SheetContent,
  SheetHeader,
  SheetTitle,
  SheetDescription,
} from '@/components/ui/sheet'
import {
  listCollectionVariables,
  createCollectionVariable,
  updateCollectionVariable,
  deleteCollectionVariable,
} from '@/api'
import type { CustomVariable } from '@/types/api.types'

interface Props {
  open: boolean
  collectionId: string
  collectionSlug: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:open': [value: boolean]
}>()

// ── state ────────────────────────────────────────────────────────────────────

const variables = ref<CustomVariable[]>([])
const loading = ref(false)
const errorMessage = ref<string | null>(null)

// new variable form
const newKey = ref('')
const newValue = ref('')
const adding = ref(false)

// inline edit state
const editingId = ref<string | null>(null)
const editKey = ref('')
const editValue = ref('')
const saving = ref(false)

// delete confirm state
const confirmDeleteId = ref<string | null>(null)
const deletingId = ref<string | null>(null)

// ── data loading ─────────────────────────────────────────────────────────────

const fetchVariables = async () => {
  loading.value = true
  errorMessage.value = null
  try {
    const res = await listCollectionVariables(props.collectionSlug)
    variables.value = res ?? []
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to load variables'
  } finally {
    loading.value = false
  }
}

watch(
  () => props.open,
  (isOpen) => {
    if (isOpen) {
      fetchVariables()
      newKey.value = ''
      newValue.value = ''
      editingId.value = null
      confirmDeleteId.value = null
      errorMessage.value = null
    }
  },
)

// ── add ───────────────────────────────────────────────────────────────────────

const handleAdd = async () => {
  const key = newKey.value.trim()
  const value = newValue.value.trim()
  if (!key || !value) return
  adding.value = true
  errorMessage.value = null
  try {
    const created = await createCollectionVariable(props.collectionId, props.collectionSlug, { key, value })
    variables.value.push(created)
    newKey.value = ''
    newValue.value = ''
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to add variable'
  } finally {
    adding.value = false
  }
}

// ── edit ──────────────────────────────────────────────────────────────────────

const startEdit = (v: CustomVariable) => {
  editingId.value = v.id
  editKey.value = v.key
  editValue.value = v.value
  confirmDeleteId.value = null
}

const cancelEdit = () => {
  editingId.value = null
}

const handleSave = async (id: string) => {
  const key = editKey.value.trim()
  const value = editValue.value.trim()
  if (!key || !value) return
  saving.value = true
  errorMessage.value = null
  try {
    const updated = await updateCollectionVariable(id, { key, value })
    const idx = variables.value.findIndex((v) => v.id === id)
    if (idx !== -1) variables.value[idx] = updated
    editingId.value = null
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to update variable'
  } finally {
    saving.value = false
  }
}

// ── display helpers ───────────────────────────────────────────────────────────

const varSyntaxExample = '{{key}}'
const wrapVar = (key: string) => `{{${key}}}`

// ── delete ────────────────────────────────────────────────────────────────────

const handleDelete = async (id: string) => {
  if (confirmDeleteId.value !== id) {
    confirmDeleteId.value = id
    return
  }
  deletingId.value = id
  try {
    await deleteCollectionVariable(id)
    variables.value = variables.value.filter((v) => v.id !== id)
    confirmDeleteId.value = null
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to delete variable'
  } finally {
    deletingId.value = null
  }
}
</script>

<template>
  <Sheet :open="open" @update:open="emit('update:open', $event)">
    <SheetContent side="right" class="w-[420px] sm:w-[480px] flex flex-col gap-0 p-0">
      <SheetHeader class="px-6 py-5 border-b">
        <SheetTitle>Environment Variables</SheetTitle>
        <SheetDescription>
          Use <code class="font-mono text-xs bg-slate-100 px-1 py-0.5 rounded">{{ varSyntaxExample }}</code>
          in scenario body or headers to inject values at request time.
        </SheetDescription>
      </SheetHeader>

      <div class="flex-1 overflow-y-auto px-6 py-5 space-y-6">

        <!-- Add form -->
        <div class="space-y-3">
          <p class="text-xs font-medium text-slate-500 uppercase tracking-wide">Add Variable</p>
          <div class="flex gap-2">
            <div class="flex-1 space-y-1">
              <Label class="text-xs text-slate-500">Key</Label>
              <Input
                v-model="newKey"
                placeholder="e.g. data.id"
                class="h-9 font-mono text-sm border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
                :disabled="adding"
                @keydown.enter="handleAdd"
              />
            </div>
            <div class="flex-1 space-y-1">
              <Label class="text-xs text-slate-500">Value</Label>
              <Input
                v-model="newValue"
                placeholder="e.g. 123"
                class="h-9 text-sm border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
                :disabled="adding"
                @keydown.enter="handleAdd"
              />
            </div>
            <div class="pt-5">
              <Button
                size="sm"
                class="h-9 px-3 bg-gray-900 hover:bg-gray-800 font-normal"
                :disabled="adding || !newKey.trim() || !newValue.trim()"
                @click="handleAdd"
              >
                <Loader2 v-if="adding" class="w-4 h-4 animate-spin" />
                <Plus v-else class="w-4 h-4" />
              </Button>
            </div>
          </div>
        </div>

        <!-- Error -->
        <div v-if="errorMessage" class="text-sm text-red-600 bg-red-50 px-3 py-2 rounded-lg">
          {{ errorMessage }}
        </div>

        <!-- Variables list -->
        <div class="space-y-2">
          <p class="text-xs font-medium text-slate-500 uppercase tracking-wide">
            Variables
            <span v-if="!loading" class="ml-1 text-slate-400">({{ variables.length }})</span>
          </p>

          <div v-if="loading" class="flex items-center gap-2 text-sm text-slate-400 py-4">
            <Loader2 class="w-4 h-4 animate-spin" />
            Loading…
          </div>

          <div v-else-if="variables.length === 0" class="text-sm text-slate-400 py-4 text-center border border-dashed rounded-lg">
            No variables yet. Add one above.
          </div>

          <div v-else class="space-y-1.5">
            <!-- Editing row -->
            <template v-for="v in variables" :key="v.id">
              <div
                v-if="editingId === v.id"
                class="flex gap-2 items-start p-2 rounded-lg border border-blue-200 bg-blue-50/40"
              >
                <div class="flex-1 space-y-1">
                  <Input
                    v-model="editKey"
                    class="h-8 font-mono text-sm border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-md"
                    :disabled="saving"
                  />
                </div>
                <div class="flex-1 space-y-1">
                  <Input
                    v-model="editValue"
                    class="h-8 text-sm border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-md"
                    :disabled="saving"
                    @keydown.enter="handleSave(v.id)"
                  />
                </div>
                <div class="flex gap-1 pt-0.5">
                  <Button
                    size="sm"
                    class="h-8 w-8 p-0 bg-gray-900 hover:bg-gray-800"
                    :disabled="saving || !editKey.trim() || !editValue.trim()"
                    @click="handleSave(v.id)"
                  >
                    <Loader2 v-if="saving" class="w-3.5 h-3.5 animate-spin" />
                    <Check v-else class="w-3.5 h-3.5" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    class="h-8 w-8 p-0 text-slate-500 hover:text-slate-800"
                    :disabled="saving"
                    @click="cancelEdit"
                  >
                    <X class="w-3.5 h-3.5" />
                  </Button>
                </div>
              </div>

              <!-- Display row -->
              <div
                v-else
                class="flex items-center gap-2 px-3 py-2 rounded-lg border border-slate-100 hover:border-slate-200 hover:bg-slate-50 group transition-colors"
              >
                <code class="text-xs font-mono text-violet-700 bg-violet-50 px-1.5 py-0.5 rounded shrink-0 max-w-[40%] truncate">
                  {{ wrapVar(v.key) }}
                </code>
                <span class="flex-1 text-sm text-slate-600 truncate">{{ v.value }}</span>
                <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
                  <Button
                    variant="ghost"
                    size="sm"
                    class="h-7 w-7 p-0 text-slate-400 hover:text-slate-700"
                    @click="startEdit(v)"
                  >
                    <Pencil class="w-3.5 h-3.5" />
                  </Button>
                  <Button
                    variant="ghost"
                    size="sm"
                    :disabled="deletingId === v.id"
                    :class="confirmDeleteId === v.id
                      ? 'h-7 w-7 p-0 bg-red-600 hover:bg-red-700 text-white rounded-md'
                      : 'h-7 w-7 p-0 text-slate-400 hover:text-red-600'"
                    @click="handleDelete(v.id)"
                  >
                    <Loader2 v-if="deletingId === v.id" class="w-3.5 h-3.5 animate-spin" />
                    <Trash2 v-else class="w-3.5 h-3.5" />
                  </Button>
                </div>
              </div>
            </template>
          </div>
        </div>
      </div>
    </SheetContent>
  </Sheet>
</template>
