<template>
  <div class="min-h-screen bg-white p-5 pt-0">
    <div class="mx-auto">
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent class="p-0">
          <!-- Header -->
          <div class="px-5 pb-4 border-b bg-white/70 backdrop-blur">
            <h1 class="text-2xl font-bold">QR EMV / QRIS Parser</h1>
            <p class="text-sm text-slate-600 mt-1">
              Paste an EMVCo or QRIS TLV string, or upload a QR image to decode and parse it into structured JSON.
            </p>
          </div>

          <!-- Input Section -->
          <section class="px-5 pt-4 pb-2 space-y-3">
            <!-- Upload + state -->
            <div class="flex items-center justify-between gap-2">
              <div class="flex items-center gap-2">
                <input
                    ref="fileInput"
                    type="file"
                    accept="image/*"
                    class="hidden"
                    @change="onFileChange"
                />
                <span v-if="decoding" class="text-xs text-slate-500">Decoding…</span>
                <span v-if="decodeError" class="text-xs text-red-600">{{ decodeError }}</span>
              </div>
            </div>

            <!-- Decode/Preview section -->
            <div v-if="decodeFromImage" class="border rounded-md p-4 bg-slate-50/60">
              <div v-if="!qrImageUrl" class="flex items-center justify-center">
                <Button size="sm" variant="outline" @click="triggerFilePicker">
                  <ImageUp/>
                </Button>
              </div>

              <div v-else class="flex flex-col items-center gap-3">
                <img :src="qrImageUrl" alt="QR preview" class="max-h-56 object-contain mx-auto" />
                <div class="flex items-center gap-2">
                  <Button size="sm" variant="outline" @click="triggerFilePicker">
                    <SquarePen/>
                  </Button>
                  <Button size="sm" variant="outline" @click="clearImage">
                    <Trash2/>
                  </Button>
                </div>
              </div>
            </div>

            <!-- Country + toggle -->
            <div class="flex items-center justify-between gap-3">
              <div class="flex items-center gap-2">
                <Label for="qr-raw" class="text-base font-semibold">Country Specification</Label>
              </div>
              <div class="flex items-center gap-2">
                <label class="flex items-center gap-2 text-xs text-slate-600 select-none">
                  <input type="checkbox" v-model="decodeFromImage" class="rounded border"/>
                  Decode QR from image
                </label>
              </div>
            </div>

            <div class="flex-1">
              <select
                  v-model="qrMode"
                  class="w-full border rounded px-2 py-3 font-mono text-xs"
              >
                <option value="global">EMVCo (Global)</option>
                <option value="qris">QRIS (Indonesia)</option>
              </select>
            </div>

            <!-- Title + actions -->
            <div class="flex items-center justify-between gap-3 mt-5">
              <div class="flex items-center gap-2">
                <Label for="qr-raw" class="text-base font-semibold">QR String</Label>
              </div>
              <div class="flex gap-2">
                <Button size="sm" @click="generateQrFromInput">Generate QR</Button>
                <Button size="sm" variant="outline" @click="loadSample">Sample QR</Button>
                <Button size="sm" variant="outline" @click="clearQr"><Trash2/></Button>
                <!-- NEW: Generate QR -->
              </div>
            </div>

            <!-- MONACO EDITOR + STATUS -->
            <div class="border rounded-md">
              <div class="h-[100px] overflow-hidden">
                <MonacoEditor
                    :height="'100%'" :width="'100%'"
                    v-model:value="qrInput"
                    language="plaintext"
                    theme="vs"
                    :options="monacoOpts"
                    @mount="onMonacoMount"
                />
              </div>

              <!-- Status bar -->
              <div class="flex flex-wrap items-center justify-between gap-2 text-[11px] px-2 py-1 bg-slate-50 border-t text-slate-600">
                <div class="font-mono">
                  {{ qrLength }} chars
                </div>
                <div class="font-mono">
                  Ln {{ cursorLine }}, Col {{ cursorCol }} ({{ selLength }} chars selected) · Index {{ cursorIndex }}
                </div>
              </div>
            </div>

            <div class="flex items-center justify-between">
              <p v-if="qrError" class="text-xs text-red-600">{{ qrError }}</p>
            </div>
          </section>

          <div class="mx-5 my-4 h-px bg-slate-100"/>

          <!-- Parsed Fields -->
          <section class="px-5 pb-6">
            <div class="flex items-center justify-between mb-2">
              <h2 class="text-sm font-semibold">Parsed Fields</h2>
              <div class="flex items-center gap-2 text-xs text-slate-600">
                <span class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5">
                  <span class="font-medium">Length</span>
                  <span class="font-mono">{{ qrLength }}</span>
                </span>
                <span class="inline-flex items-center gap-1 rounded-full border px-2 py-0.5">
                  <span class="font-medium">Tags</span>
                  <span class="font-mono">{{ tagCount }}</span>
                </span>
              </div>
            </div>

            <div v-if="Object.keys(qrParsed).length" class="overflow-auto border rounded-md">
              <table class="w-full text-sm">
                <thead class="bg-slate-50">
                <tr>
                  <th class="text-center px-3 py-2 border-b w-28">Index</th>
                  <th class="text-left px-3 py-2 border-b w-24">Tag</th>
                  <th class="text-left px-3 py-2 border-b w-56">Field Name</th>
                  <th class="text-left px-3 py-2 border-b">Value</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="[k, v] in kvRows" :key="k" class="odd:bg-white even:bg-slate-50/40">
                  <!-- NEW: index column (0-based start of the VALUE) -->
                  <td class="px-3 py-2 border-b text-center font-mono text-[12px]" :title="posOf(k).title">
                    {{ posOf(k).label }}
                  </td>

                  <td class="px-3 py-2 border-b font-mono text-[12px]">{{ k }}</td>

                  <td class="px-3 py-2 border-b text-[12px]">
                      <span
                          :title="emvTags[tagId(k)]?.name || 'Unknown field'"
                          :class="[
                          'font-medium',
                          emvTags[tagId(k)] ? 'text-slate-800' : 'text-slate-400 italic',
                        ]"
                      >
                        {{ emvTags[tagId(k)]?.name || 'Unknown' }}
                      </span>
                  </td>
                  <td class="px-3 py-2 border-b break-all [overflow-wrap:anywhere] font-mono text-[12px]">
                    {{ v }}
                  </td>
                </tr>
                </tbody>
              </table>
            </div>

            <div v-else class="text-xs text-slate-500 border rounded-md px-3 py-6 grid place-items-center">
              No parsed fields yet — paste a TLV string above.
            </div>
          </section>

          <!-- Missing / Issues panel -->
          <div v-if="missing.length" class="mx-5 mb-3 rounded-md border border-amber-300 bg-amber-50 p-3 text-xs">
            <div class="font-semibold mb-1">Missing / Issues ({{ qrModeLabel }})</div>
            <ul class="list-disc pl-4 space-y-1">
              <li v-for="m in missing" :key="m.key">
                <span class="font-mono">{{ m.tag }}</span>
                <span class="mx-1">—</span>
                <span class="font-medium">{{ m.name }}</span>
                <span class="mx-1 text-slate-600">· {{ m.reason }}</span>
              </li>
            </ul>
          </div>
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
                id="qr-json"
                :model-value="qrOutput"
                readonly
                rows="8"
                class="font-mono text-xs"
            />
          </section>
        </CardContent>
      </Card>
    </div>
  </div>

  <!-- NEW: QR Preview Modal -->
  <div
      v-if="showQrModal"
      class="fixed inset-0 z-50"
      @keydown.esc.prevent="closeQrModal"
  >
    <!-- Backdrop -->
    <div class="absolute inset-0 bg-black/40" @click="closeQrModal"></div>

    <!-- Dialog -->
    <div
        class="absolute inset-0 flex items-center justify-center p-4"
        role="dialog"
        aria-modal="true"
    >
      <div class="w-full max-w-sm rounded-xl bg-white shadow-lg border">
        <div class="px-5 py-3 border-b flex items-center justify-between">
          <h3 class="text-sm font-semibold">Generated QR</h3>
          <button
              class="text-slate-500 hover:text-slate-700 text-sm"
              @click="closeQrModal"
              aria-label="Close"
          >✕</button>
        </div>

        <div class="px-5 py-6">
          <div v-if="qrGeneratedUrl" class="flex flex-col items-center gap-3">
            <img :src="qrGeneratedUrl" alt="Generated QR" class="w-56 h-56 object-contain" />
            <p class="text-[11px] text-slate-600 text-center break-all">
              From input ({{ qrLength }} chars)
            </p>
          </div>
          <p v-else class="text-xs text-red-600 text-center">Failed to generate QR.</p>
        </div>

        <div class="px-5 pb-4 flex items-center justify-end gap-2">
          <Button size="sm" variant="outline" @click="downloadGeneratedQr" :disabled="!qrGeneratedUrl">
            Download
          </Button>
          <Button size="sm" @click="closeQrModal">Close</Button>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { computed, ref, watch, onBeforeUnmount, nextTick } from 'vue'
import { toast } from 'vue-sonner'
import { ClipboardCopy, Trash2, SquarePen, ImageUp } from 'lucide-vue-next'

import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'

import { emvTags, parseTagWithPos } from '@/lib/parser/emvcoParser.ts'

/* ---------- Monaco (default import) ---------- */
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

/* Cursor + selection UI state */
const cursorLine = ref(1)
const cursorCol = ref(1)
const cursorIndex = ref(0)

const selStart = ref(0)
const selEnd = ref(0)
const selLength = ref(0)

let editorInst: MonacoNS.editor.IStandaloneCodeEditor | null = null
let disposers: Array<{ dispose: () => void }> = []

function onMonacoMount(ed: MonacoNS.editor.IStandaloneCodeEditor) {
  editorInst = ed

  // initial layout + initial positions
  nextTick(() => ed.layout())
  setTimeout(() => ed.layout(), 0)
  updateCursorInfo()
  updateSelectionInfo()

  // listeners
  const d1 = ed.onDidChangeCursorPosition(() => {
    updateCursorInfo()
    updateSelectionInfo() // if selection collapses to caret
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

/* ---------------- NEW: Image decode state ---------------- */
const fileInput = ref<HTMLInputElement | null>(null)
const qrImageUrl = ref<string | null>(null)
const decodeFromImage = ref<boolean>(false)
const decoding = ref<boolean>(false)
const decodeError = ref<string | null>(null)

function triggerFilePicker() {
  fileInput.value?.click()
}

async function onFileChange(e: Event) {
  const target = e.target as HTMLInputElement
  const file = target.files?.[0]
  decodeError.value = null
  if (!file) return
  if (qrImageUrl.value) URL.revokeObjectURL(qrImageUrl.value)
  qrImageUrl.value = URL.createObjectURL(file)
  decoding.value = true
  try {
    const text = await decodeFromImageFile(file)
    if (!text) throw new Error('No QR code detected in the image.')
    qrInput.value = text.trim()
    toast.success('QR decoded from image')
  } catch (err: any) {
    decodeError.value = err?.message ?? 'Failed to decode QR image'
    toast.error(decodeError.value)
  } finally {
    decoding.value = false
    if (fileInput.value) fileInput.value.value = ''
  }
}

function clearImage() {
  if (qrImageUrl.value) URL.revokeObjectURL(qrImageUrl.value)
  qrImageUrl.value = null
  decodeError.value = null
}

/**
 * Try BarcodeDetector first; fallback to @zxing/library if available.
 * Install fallback with: npm i @zxing/library (optional)
 */
async function decodeFromImageFile(file: File): Promise<string | null> {
  const BD: any = (globalThis as any).BarcodeDetector
  if (BD) {
    try {
      const detector = new BD({ formats: ['qr_code'] })
      const bitmap = await createImageBitmap(file)
      const results = await detector.detect(bitmap as any)
      if (results?.length) {
        const val = (results[0] as any).rawValue ?? (results[0] as any).data
        if (val) return String(val)
      }
    } catch {
      // fall through to ZXing
    }
  }

  try {
    const mod = await import(/* @vite-ignore */ '@zxing/library')
    const reader = new (mod as any).BrowserQRCodeReader()
    const url = URL.createObjectURL(file)
    try {
      const res = await reader.decodeFromImageUrl(url)
      return res?.getText?.() ?? res?.text ?? null
    } finally {
      URL.revokeObjectURL(url)
      reader.reset?.()
    }
  } catch {
    throw new Error('Could not decode. Your browser may not support BarcodeDetector. Install optional fallback: npm i @zxing/library')
  }
}

/* ---------------- State ---------------- */
const qrInput = ref<string>('')
const qrParsed = ref<Record<string, string>>({})
const qrPositions = ref<Record<string, { start: number; end: number }>>({})
const qrError = ref<string | null>(null)
const qrMode = ref<'global' | 'qris'>('qris')

/* Derived */
const qrOutput = computed(() => {
  try {
    return JSON.stringify(qrParsed.value, null, 2)
  } catch {
    return '{}'
  }
})
const kvRows = computed<[string, string][]>(() =>
    Object.entries(qrParsed.value).sort((a, b) => a[0].localeCompare(b[0]))
)
const qrLength = computed(() => (qrInput.value || '').length)
const tagCount = computed(() => {
  try {
    let i = 0, n = 0
    const s = (qrInput.value || '').trim()
    while (i + 4 <= s.length) {
      const len = Number.parseInt(s.slice(i + 2, i + 4), 10)
      if (Number.isNaN(len)) break
      const end = i + 4 + len
      if (end > s.length) break
      n++
      i = end
    }
    return n
  } catch {
    return 0
  }
})

function tagId(key: string) {
  return key.match(/^data(\d{2})/)?.[1] || ''
}

/* ---------------- Validation (Missing/Issues) ---------------- */
type MissingItem = { key: string; tag: string; name: string; reason: string }
const qrModeLabel = computed(() => (qrMode.value === 'qris' ? 'QRIS' : 'EMVCo'))

const missing = computed<MissingItem[]>(() => {
  if (!qrInput.value || !qrInput.value.trim()) return []

  const p = qrParsed.value || {}
  const out: MissingItem[] = []

  const has = (id: string) => {
    const v = p[`data${id}`]
    return v != null && String(v).length > 0
  }
  const nameOf = (id: string) => emvTags[id]?.name ?? 'Unknown'

  const requiredRoot = ['00', '52', '53', '58', '59', '60', '63']
  for (const id of requiredRoot) {
    if (!has(id)) out.push({ key: `req-${id}`, tag: id, name: nameOf(id), reason: 'Required by EMVCo' })
  }

  let hasMAI = false
  for (let i = 2; i <= 51; i++) {
    if (has(String(i).padStart(2, '0'))) { hasMAI = true; break }
  }
  if (!hasMAI) {
    out.push({ key: 'req-mai', tag: '02–51', name: 'Merchant Account Information', reason: 'At least one MAI is required (EMVCo)' })
  }

  if (qrMode.value === 'qris') {
    if (has('53') && p['data53'] !== '360') {
      out.push({ key: 'qris-53', tag: '53', name: nameOf('53'), reason: 'QRIS expects IDR (360)' })
    }
    if (has('58') && p['data58'] !== 'ID') {
      out.push({ key: 'qris-58', tag: '58', name: nameOf('58'), reason: 'QRIS expects country code ID' })
    }
    let foundDomesticGUI = false
    for (let id = 26; id <= 51; id++) {
      const k = `data${String(id).padStart(2, '0')}_00`
      const v = p[k]
      if (v) {
        foundDomesticGUI = true
        if (!/^[A-Za-z0-9.-]+\.[A-Za-z0-9.-]+\.[A-Za-z0-9.-]+$/.test(v)) {
          out.push({ key: `qris-gui-${id}`, tag: `${id}.00`, name: 'Globally Unique Identifier', reason: 'Expected reverse-domain/AID/UUID (common QRIS pattern e.g. ID.CO.QRIS.WWW)' })
        }
      }
    }
    if (!foundDomesticGUI) {
      out.push({ key: 'qris-mai00', tag: '26–51.00', name: 'Globally Unique Identifier', reason: 'QRIS MAI should include sub-tag 00 (GUI)' })
    }
  }

  return out
})

/* ---------------- Actions ---------------- */
function parseQr() {
  const raw = (qrInput.value || '').trim()
  if (!raw) { qrParsed.value = {}; qrPositions.value = {}; qrError.value = null; return }
  try {
    const { tags, positions } = parseTagWithPos(raw)
    qrParsed.value = tags
    qrPositions.value = positions
    qrError.value = null
  } catch (e: any) {
    qrParsed.value = {}
    qrPositions.value = {}
    qrError.value = e?.message || 'Invalid EMV string.'
  }
}

function posOf(key: string) {
  const p = qrPositions.value[key]
  if (!p) return { label: '–', title: '' }
  const len = p.end - p.start
  return { label: String(p.start), title: `range ${p.start}–${p.end} (len ${len})` }
}

async function copyOutput() {
  try {
    await navigator.clipboard.writeText(qrOutput.value)
    toast.success('Parsed JSON copied')
  } catch {
    toast.error('Failed to copy')
  }
}

async function copyTable() {
  const text = kvRows.value.map(([k, v]) => `${k}: ${v}`).join('\n')
  try {
    await navigator.clipboard.writeText(text)
    toast.success('Parsed fields copied')
  } catch {
    toast.error('Failed to copy')
  }
}

function clearQr() {
  qrInput.value = ''
  qrParsed.value = {}
  qrPositions.value = {}
  qrError.value = null
}

function loadSample() {
  qrInput.value =
      '00020101021226690017ID.CO.PRIMAQR.WWW011893600998297703010002150000880800000200303UME520454995303360540813000.005802ID5918AYAM GORENG SEKALI6010TANGGERANG61051511562210717C012345678901234563046CF4'
  parseQr()
}

function downloadJson() {
  try {
    const blob = new Blob([qrOutput.value], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'parser-emv-parse.json'
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    toast.error('Failed to download JSON')
  }
}

function downloadTable() {
  try {
    const header = 'field_name,tag,value'
    const rows = kvRows.value.map(([k, v]) => {
      const name = emvTags[tagId(k)]?.name || 'Unknown'
      return `"${name}","${k.replace(/"/g, '""')}","${String(v).replace(/"/g, '""')}"`
    })
    const csv = [header, ...rows].join('\n')
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'parser-emv-parse.csv'
    a.click()
    URL.revokeObjectURL(url)
  } catch {
    toast.error('Failed to download CSV')
  }
}

/* Auto-parse (debounced) */
let t: number | undefined
watch(qrInput, (v) => {
  if (t) window.clearTimeout(t)
  t = window.setTimeout(() => {
    if ((v || '').trim()) parseQr()
    else { qrParsed.value = {}; qrPositions.value = {}; qrError.value = null }
  }, 250)
})

/* ---------- NEW: Generate QR modal state + logic ---------- */
const showQrModal = ref(false)
const qrGeneratedUrl = ref<string | null>(null)

async function generateQrFromInput() {
  const text = (qrInput.value || '').trim()
  if (!text) {
    toast.error('Enter a QR string first')
    return
  }
  try {
    // Dynamic import keeps your main bundle lighter.
    // Be sure to install: npm i qrcode
    const QR = await import(/* @vite-ignore */ 'qrcode')
    const url = await (QR as any).toDataURL(text, {
      errorCorrectionLevel: 'M',
      margin: 1,
      width: 256,
    })
    qrGeneratedUrl.value = url
    showQrModal.value = true
  } catch (e: any) {
    qrGeneratedUrl.value = null
    toast.error(e?.message ?? 'Failed to generate QR')
  }
}

function closeQrModal() {
  showQrModal.value = false
}

function downloadGeneratedQr() {
  if (!qrGeneratedUrl.value) return
  const a = document.createElement('a')
  a.href = qrGeneratedUrl.value
  a.download = 'qr-from-input.png'
  a.click()
}
</script>

<style scoped>
/* Minimal styling handled via Tailwind classes */
</style>
