<script setup lang="ts">
import { ref, watch } from 'vue'
import { Plus, Trash2, Pencil, Check, X, Loader2 } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import ConfirmDeleteDialog from '@/components/ConfirmDeleteDialog.vue'
import {
  listCollectionVariables,
  createCollectionVariable,
  updateCollectionVariable,
  deleteCollectionVariable,
} from '@/api'
import type { CustomVariable } from '@/types/api.types'

interface Props {
  collectionId: string
  collectionSlug: string
}

const props = defineProps<Props>()

const emit = defineEmits<{
  close: []
}>()

// ── state ─────────────────────────────────────────────────────────────────────

const variables = ref<CustomVariable[]>([])
const loading = ref(false)
const errorMessage = ref<string | null>(null)

const newKey = ref('')
const newValue = ref('')
const adding = ref(false)

const editingId = ref<string | null>(null)
const editKey = ref('')
const editValue = ref('')
const saving = ref(false)

const deleteTargetId = ref<string | null>(null)
const deletingId = ref<string | null>(null)

const varSyntaxExample = '{{key}}'
const wrapVar = (key: string) => `{{${key}}}`

// ── load ──────────────────────────────────────────────────────────────────────

const fetchVariables = async () => {
  loading.value = true
  errorMessage.value = null
  try {
    variables.value = (await listCollectionVariables(props.collectionSlug)) ?? []
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to load variables'
  } finally {
    loading.value = false
  }
}

watch(() => props.collectionSlug, fetchVariables, { immediate: true })

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
  deleteTargetId.value = null
}

const cancelEdit = () => { editingId.value = null }

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

// ── delete ────────────────────────────────────────────────────────────────────

const handleConfirmDelete = async () => {
  if (!deleteTargetId.value) return
  const id = deleteTargetId.value
  deletingId.value = id
  try {
    await deleteCollectionVariable(id)
    variables.value = variables.value.filter((v) => v.id !== id)
  } catch (err: any) {
    errorMessage.value = err?.message || 'Failed to delete variable'
  } finally {
    deletingId.value = null
    deleteTargetId.value = null
  }
}
</script>

<template>
  <div class="flex flex-col h-full">
    <!-- Panel header -->
    <div class="flex items-center justify-between px-5 py-4 border-b shrink-0">
      <div>
        <p class="text-sm font-semibold text-slate-900">Environment Variables</p>
        <p class="text-xs text-slate-500 mt-0.5">
          Use <code class="font-mono bg-slate-100 px-1 py-0.5 rounded text-violet-700">{{ varSyntaxExample }}</code>
          in scenario body &amp; headers.
        </p>
      </div>
      <Button
        variant="ghost"
        size="sm"
        class="h-7 w-7 p-0 text-slate-400 hover:text-slate-700 shrink-0"
        @click="emit('close')"
      >
        <X class="w-4 h-4" />
      </Button>
    </div>

    <!-- Panel body (scrollable) -->
    <div class="flex-1 overflow-y-auto px-5 py-4 space-y-5">

      <!-- Add form -->
      <div class="space-y-2">
        <p class="text-xs font-medium text-slate-500 uppercase tracking-wide">Add Variable</p>
        <div class="flex gap-2">
          <div class="flex-1 space-y-1">
            <Label class="text-xs text-slate-500">Key</Label>
            <Input
              v-model="newKey"
              placeholder="e.g. user.id"
              class="h-8 font-mono text-sm border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
              :disabled="adding"
              @keydown.enter="handleAdd"
            />
          </div>
          <div class="flex-1 space-y-1">
            <Label class="text-xs text-slate-500">Value</Label>
            <Input
              v-model="newValue"
              placeholder="e.g. 123"
              class="h-8 text-sm border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-lg"
              :disabled="adding"
              @keydown.enter="handleAdd"
            />
          </div>
          <div class="pt-5">
            <Button
              size="sm"
              class="h-8 w-8 p-0 bg-gray-900 hover:bg-gray-800"
              :disabled="adding || !newKey.trim() || !newValue.trim()"
              @click="handleAdd"
            >
              <Loader2 v-if="adding" class="w-3.5 h-3.5 animate-spin" />
              <Plus v-else class="w-3.5 h-3.5" />
            </Button>
          </div>
        </div>
      </div>

      <!-- Error -->
      <div v-if="errorMessage" class="text-xs text-red-600 bg-red-50 px-3 py-2 rounded-lg">
        {{ errorMessage }}
      </div>

      <!-- Variable list -->
      <div class="space-y-2">
        <p class="text-xs font-medium text-slate-500 uppercase tracking-wide">
          Variables
          <span v-if="!loading" class="ml-1 text-slate-400 normal-case font-normal">({{ variables.length }})</span>
        </p>

        <div v-if="loading" class="flex items-center gap-2 text-sm text-slate-400 py-4">
          <Loader2 class="w-4 h-4 animate-spin" />
          Loading…
        </div>

        <div v-else-if="variables.length === 0"
          class="text-sm text-slate-400 py-6 text-center border border-dashed rounded-lg">
          No variables yet. Add one above.
        </div>

        <div v-else class="space-y-1.5">
          <template v-for="v in variables" :key="v.id">

            <!-- Editing row -->
            <div v-if="editingId === v.id"
              class="flex gap-1.5 items-center p-2 rounded-lg border border-blue-200 bg-blue-50/40">
              <Input
                v-model="editKey"
                class="h-8 flex-1 font-mono text-sm border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-md"
                :disabled="saving"
              />
              <Input
                v-model="editValue"
                class="h-8 flex-1 text-sm border-gray-200 focus:border-gray-300 focus:ring-1 focus:ring-gray-300 rounded-md"
                :disabled="saving"
                @keydown.enter="handleSave(v.id)"
              />
              <Button
                size="sm"
                class="h-8 w-8 p-0 bg-gray-900 hover:bg-gray-800 shrink-0"
                :disabled="saving || !editKey.trim() || !editValue.trim()"
                @click="handleSave(v.id)"
              >
                <Loader2 v-if="saving" class="w-3.5 h-3.5 animate-spin" />
                <Check v-else class="w-3.5 h-3.5" />
              </Button>
              <Button
                variant="ghost"
                size="sm"
                class="h-8 w-8 p-0 text-slate-400 hover:text-slate-700 shrink-0"
                :disabled="saving"
                @click="cancelEdit"
              >
                <X class="w-3.5 h-3.5" />
              </Button>
            </div>

            <!-- Display row -->
            <div v-else
              class="flex items-center gap-2 px-3 py-2 rounded-lg border border-slate-100 hover:border-slate-200 hover:bg-slate-50 group transition-colors">
              <code class="text-xs font-mono text-violet-700 bg-violet-50 px-1.5 py-0.5 rounded shrink-0 max-w-[45%] truncate">
                {{ wrapVar(v.key) }}
              </code>
              <span class="flex-1 text-sm text-slate-600 truncate min-w-0">{{ v.value }}</span>
              <div class="flex gap-0.5 opacity-0 group-hover:opacity-100 transition-opacity shrink-0">
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
                  class="h-7 w-7 p-0 text-slate-400 hover:text-red-600"
                  :disabled="deletingId === v.id"
                  @click="deleteTargetId = v.id"
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
  </div>

  <!-- Delete confirmation modal -->
  <ConfirmDeleteDialog
    :open="deleteTargetId !== null"
    title="Delete variable?"
    description="This will permanently remove this variable from all scenario responses."
    :loading="deletingId !== null"
    @update:open="(v) => { if (!v) deleteTargetId = null }"
    @confirm="handleConfirmDelete"
  />
</template>
