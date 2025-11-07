<template>
  <div class="min-h-screen bg-white p-5 pt-0">
    <div class="mx-auto ">
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent class="p-0">
          <!-- Header -->
          <div class="px-5 pb-4 border-b bg-white/70 backdrop-blur">
            <h1 class="text-2xl font-bold">ISO 8583 Parser</h1>
            <p class="text-sm text-slate-600 mt-1">
              Paste an ISO 8583:1987 ASCII message (MTI + HEX bitmap + data) and parse using FIXED / LL / LLL.
            </p>
          </div>

          <!-- Input & Controls -->
          <section class="px-5 pt-4 pb-2 space-y-3">
            <div class="flex items-center justify-between gap-3">
              <div class="flex items-center gap-2">
                <Label for="iso-raw" class="text-base font-semibold">Message Specification</Label>
              </div>
            </div>

            <div class="flex-1">
              <select v-model="selectedSpecKey" class="w-full border rounded px-2 py-3 font-mono text-xs">
                <option value="iso87">Spec 87</option>
                <option value="rintis">Spec Custom (Rintis)</option>
                <option value="free_ll">Free (LL)</option>
                <option value="free_lll">Free (LLL)</option>
              </select>
            </div>

            <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-2 mt-5">
              <div class="flex items-center gap-2">
                <Label for="iso-raw" class="text-base font-semibold">Message</Label>
              </div>

              <div class="flex gap-2">
                <Button size="sm" variant="outline" @click="loadSample">Sample ISO</Button>
                <Button size="sm" variant="outline" @click="clearAll"><Trash2/></Button>
              </div>
            </div>

            <!-- MONACO EDITOR -->
            <div class="border rounded-md">
              <div class="h-[150px] overflow-hidden">
                <MonacoEditor
                    :height="'100%'" :width="'100%'"
                    v-model:value="raw"
                    language="plaintext"
                    theme="vs"
                    :options="monacoOpts"
                    @mount="onMonacoMount"
                />
              </div>

              <!-- Status bar -->
              <div class="flex flex-wrap items-center justify-between gap-2 text-[11px] px-2 py-1 bg-slate-50 border-t text-slate-600">
                <div class="font-mono">
                  {{ msgLength }} chars
                </div>
                <div class="font-mono">
                  Ln {{ cursorLine }}, Col {{ cursorCol }} ({{ selLength }} chars selected) · Index {{ cursorIndex }}
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between">
              <p v-if="err" class="text-xs text-red-600">{{ err }}</p>
            </div>
          </section>

          <div class="mx-5 my-4 h-px bg-slate-100"/>

          <!-- Parsed Fields Table -->
          <section class="px-5 pb-6">
            <div class="flex items-center justify-between mb-2">
              <h2 class="text-sm font-semibold">Parsed Fields</h2>
              <div class="flex items-center gap-2 text-xs text-slate-600">
                <span class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5">
                  <span class="font-medium">Length</span>
                  <span class="font-mono">{{ msgLength }}</span>
                </span>
                <span class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5">
                  <span class="font-medium">Fields</span>
                  <span class="font-mono">{{ fieldCount }}</span>
                </span>
                <span class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5">
                  <span class="font-medium">Bitmap Bits</span>
                  <span class="font-mono">{{ presentBits }}</span>
                </span>
                <span class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5" v-if="parsed">
                  <span class="font-medium">Missing DEs</span>
                  <span class="font-mono">{{ missingBitsList.length }}</span>
                </span>
              </div>
            </div>

            <div v-if="parsed && parsed.fields.length" class="overflow-auto border rounded-md">
              <table class="w-full text-sm">
                <thead class="bg-slate-50">
                <tr>
                  <th class="text-center px-3 py-2 border-b w-28">Index</th>
                  <th class="text-center px-3 py-2 border-b w-16">DE</th>
                  <th class="text-left px-3 py-2 border-b w-56">Field Name</th>
                  <th class="text-left px-3 py-2 border-b w-24">Len</th>
                  <th class="text-left px-3 py-2 border-b w-24">Type</th>
                  <th class="text-left px-3 py-2 border-b">Value</th>
                </tr>
                </thead>
                <tbody>
                <tr
                    v-for="pf in parsed.fields"
                    :key="pf.id"
                    class="odd:bg-white even:bg-slate-50/40 hover:bg-sky-50 cursor-pointer"
                    @click="revealField(pf.start, pf.end)"
                    :title="`range ${pf.start}–${pf.end} (len ${pf.end - pf.start})`"
                >

                  <td class="px-3 py-2 border-b text-center font-mono text-[12px]">{{ pf.start }}</td>

                  <td class="px-3 py-2 border-b text-center font-mono text-[12px]">{{ pf.id }}</td>

                  <!-- NEW: Index column (0-based start of VALUE) -->

                  <td class="px-3 py-2 border-b text-[12px]">
                    <span class="font-medium text-slate-800">{{ fieldName(pf.id) }}</span>
                  </td>
                  <td class="px-3 py-2 border-b font-mono text-[12px]">
                    <span v-if="pf.lenType==='FIXED'">{{ pf.actualLength }}</span>
                    <span v-else>{{ pf.declaredLength }} / {{ pf.actualLength }}</span>
                  </td>
                  <td class="px-3 py-2 border-b text-[12px]">{{ pf.lenType }}</td>
                  <td class="px-3 py-2 border-b break-all [overflow-wrap:anywhere] font-mono text-[12px]">{{ pf.value }}</td>
                </tr>
                </tbody>
              </table>
            </div>

            <div v-else class="text-xs text-slate-500 border rounded-md px-3 py-6 grid place-items-center">
              No parsed fields yet — paste an ISO 8583 message above.
            </div>
          </section>

          <div class="mx-5 my-4 h-px bg-slate-100"/>

          <!-- JSON Output -->
          <section class="px-5 pb-4 space-y-2">
            <div class="flex items-center justify-between">
              <h2 class="text-sm font-semibold">Parsed Output (JSON)</h2>
              <div class="flex gap-2">
                <Button size="sm" variant="outline" @click="copyOutput"><ClipboardCopy/></Button>
              </div>
            </div>
            <Textarea
                id="iso-json"
                :model-value="jsonOut"
                readonly
                rows="8"
                class="font-mono text-xs"
            />
          </section>
        </CardContent>
      </Card>

      <!-- Meta summary -->
      <div v-if="parsed" class="mt-4 grid grid-cols-1 md:grid-cols-3 gap-3">
        <div class="border rounded-md p-3">
          <div class="text-xs text-slate-500">MTI</div>
          <div class="font-mono">{{ parsed.mti }}</div>
        </div>
        <div class="border rounded-md p-3">
          <div class="text-xs text-slate-500">Bitmap (HEX)</div>
          <div class="font-mono break-all [overflow-wrap:anywhere]">{{ parsed.bitmapHex }}</div>
        </div>
        <div class="border rounded-md p-3">
          <div class="text-xs text-slate-500">Bits Set</div>
          <div class="font-mono">{{ presentBits }}</div>
        </div>
      </div>
    </div>

    <!-- Missing / Empty summary -->
    <div v-if="parsed" class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-3">
      <div class="border rounded-md p-3">
        <div class="text-xs text-slate-500">Missing DEs (bit = 0 in bitmap)</div>
        <div class="font-mono text-xs break-words">{{ missingBitsList.length ? missingBitsList.join(', ') : '—' }}</div>
      </div>
      <div class="border rounded-md p-3">
        <div class="text-xs text-slate-500">Empty DEs (present but no value)</div>
        <div class="font-mono text-xs break-words">{{ emptyFieldsList.length ? emptyFieldsList.join(', ') : '—' }}</div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch, onBeforeUnmount, nextTick } from 'vue'
import { toast } from 'vue-sonner'
import { ClipboardCopy, Trash2 } from 'lucide-vue-next'

import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'

import { AVAILABLE_SPECS } from '@/lib/parser/iso8583/spec'
import type { MessageSpec } from '@/lib/parser/iso8583/spec'
import { parseISO8583 } from '@/lib/parser/iso8583/parser'

/* ---------- Monaco ---------- */
import MonacoEditor from '@guolao/vue-monaco-editor'
import type * as MonacoNS from 'monaco-editor'

const monacoOpts = {
  fontSize: 12,
  lineNumbers: 'on' as const,
  wordWrap: 'on' as const,
  minimap: { enabled: false },
  automaticLayout: true,
  renderWhitespace: 'none' as const,
  folding: false,
  scrollBeyondLastLine: false,
  fontFamily: 'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", monospace',
  fontLigatures: true,
  tabSize: 2,
}

let editorInst: MonacoNS.editor.IStandaloneCodeEditor | null = null
let disposers: Array<{ dispose: () => void }> = []

const cursorLine = ref(1)
const cursorCol = ref(1)
const cursorIndex = ref(0)
const selStart = ref(0)
const selEnd = ref(0)
const selLength = ref(0)

function onMonacoMount(ed: MonacoNS.editor.IStandaloneCodeEditor) {
  editorInst = ed
  nextTick(() => ed.layout())
  setTimeout(() => ed.layout(), 0)
  updateCursorInfo()
  updateSelectionInfo()

  const d1 = ed.onDidChangeCursorPosition(() => {
    updateCursorInfo()
    updateSelectionInfo()
  })
  const d2 = ed.onDidChangeModelContent(() => {
    updateCursorInfo()
    updateSelectionInfo()
  })
  const d3 = ed.onDidChangeCursorSelection(() => {
    updateSelectionInfo()
  })
  disposers.push(d1, d2, d3)
}

function updateCursorInfo() {
  if (!editorInst) return
  const pos = editorInst.getPosition()
  const model = editorInst.getModel()
  if (!pos || !model) return
  cursorLine.value = pos.lineNumber
  cursorCol.value = pos.column
  cursorIndex.value = model.getOffsetAt(pos)
}

function updateSelectionInfo() {
  if (!editorInst) return
  const sel = editorInst.getSelection()
  const model = editorInst.getModel()
  if (!sel || !model) {
    selStart.value = 0
    selEnd.value = 0
    selLength.value = 0
    return
  }
  const start = model.getOffsetAt(sel.getStartPosition())
  const end = model.getOffsetAt(sel.getEndPosition())
  selStart.value = Math.min(start, end)
  selEnd.value = Math.max(start, end)
  selLength.value = selEnd.value - selStart.value
}

onBeforeUnmount(() => {
  disposers.forEach(d => d.dispose())
  disposers = []
  editorInst = null
})

/* ---------------- State ---------------- */
const raw = ref<string>('')
const err = ref<string | null>(null)
const parsed = ref<ReturnType<typeof parseISO8583> | null>(null)
const selectedSpecKey = ref<'iso87' | 'rintis' | 'free_ll' | 'free_lll'>('iso87')

/* Derived */
const activeSpec = computed<MessageSpec>(() => AVAILABLE_SPECS[selectedSpecKey.value])
const jsonOut = computed(() => {
  try {
    return JSON.stringify(parsed.value?.map || {}, null, 2)
  } catch {
    return '{}'
  }
})
const msgLength = computed(() => (raw.value || '').length)
const fieldCount = computed(() => parsed.value?.fields.length || 0)
const presentBits = computed(() => parsed.value ? parsed.value.bits.reduce((a, b) => a + (b ? 1 : 0), 0) : 0)

// --- Missing/Empty detection ---
const allSpecIds = computed(() => Object.keys(activeSpec.value.fields).map(n => parseInt(n, 10)).filter(n => !Number.isNaN(n)))
const presentIdsFromBitmap = computed(() => {
  if (!parsed.value) return [] as number[]
  const ids: number[] = []
  const bits = parsed.value.bits
  for (let i = 1; i < bits.length; i++) if (bits[i] === 1) ids.push(i + 1)
  return ids
})
const missingBitsList = computed(() => allSpecIds.value.filter(id => id >= 2 && !presentIdsFromBitmap.value.includes(id)))
const emptyFieldsList = computed(() => parsed.value ? parsed.value.fields.filter(f => f.actualLength === 0).map(f => f.id) : [])

function fieldName(id: number) {
  return activeSpec.value.fields[id]?.name || `DE ${id}`
}

/* ---------------- Actions ---------------- */
function doParse() {
  const s = (raw.value || '').trim()
  if (!s) {
    parsed.value = null
    err.value = null
    return
  }
  try {
    parsed.value = parseISO8583(s, activeSpec.value)
    err.value = null
  } catch (e: any) {
    parsed.value = null
    err.value = e?.message || 'Invalid ISO8583 message.'
  }
}

async function copyOutput() {
  try {
    await navigator.clipboard.writeText(jsonOut.value)
    toast.success('Parsed JSON copied')
  } catch {
    toast.error('Failed to copy')
  }
}

function clearAll() {
  raw.value = ''
  parsed.value = null
  err.value = null
}

function loadSample() {
  // Demo: MTI 0200 + primary bitmap (hex) + DE2(LL=16), DE3(6), DE4(12), DE11(6)
  const mti = '0200'
  const bitmap = '7020000000000000' // demo primary bitmap (hex ASCII), bits set for 2,3,4,11
  const de2 = '16' + '4567890123456789' // LL + PAN (16)
  const de3 = '000000'
  const de4 = '000000001000'
  const de11 = '123456'
  raw.value = mti + bitmap + de2 + de3 + de4 + de11
  doParse()
}

/* Auto-parse (debounced) */
let t: number | undefined
watch(raw, (v) => {
  if (t) window.clearTimeout(t)
  t = window.setTimeout(() => {
    (v || '').trim() ? doParse() : (parsed.value = null, err.value = null)
  }, 250)
})

/* Click row -> highlight range in Monaco */
function revealField(start: number, end: number) {
  if (!editorInst) return
  const model = editorInst.getModel()
  if (!model) return
  const startPos = model.getPositionAt(start)
  const endPos = model.getPositionAt(end)
  editorInst.setSelection({ startLineNumber: startPos.lineNumber, startColumn: startPos.column, endLineNumber: endPos.lineNumber, endColumn: endPos.column })
  editorInst.revealRangeInCenter({ startLineNumber: startPos.lineNumber, startColumn: startPos.column, endLineNumber: endPos.lineNumber, endColumn: endPos.column })
}
</script>

<style scoped>
/* Tailwind covers styling */
</style>
