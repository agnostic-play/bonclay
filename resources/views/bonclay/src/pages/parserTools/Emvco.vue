<template>
  <div class="min-h-screen bg-white p-5 pt-0">
    <div class="mx-auto ">
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent class="p-0">
          <!-- Header -->
          <div class="px-5 pb-4 border-b bg-white/70 backdrop-blur">
            <h1 class="text-2xl font-bold">QR EMV / QRIS Parser</h1>
            <p class="text-sm text-slate-600 mt-1">
              Paste an EMVCo or QRIS TLV string to parse it into structured JSON.
            </p>
          </div>

          <!-- Input Section -->
          <section class="px-5 pt-4 pb-2 space-y-3">
            <div class="flex items-center justify-between">
              <Label for="qr-raw" class="text-base font-semibold">QR String</Label>
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

            <Textarea
                id="qr-raw"
                v-model="qrInput"
                rows="7"
                class="font-mono text-xs"
                placeholder="00020101021126... (EMV TLV string)"
            />

            <div class="flex items-center justify-between">
              <p v-if="qrError" class="text-xs text-red-600">{{ qrError }}</p>
              <div class="flex gap-2">
                <Button size="sm" variant="outline" @click="loadSample">Sample</Button>
                <Button size="sm" variant="outline" @click="clearQr">Clear</Button>
                <Button size="sm" variant="outline" @click="copyOutput">Copy JSON</Button>
                <Button size="sm" @click="parseQr">Parse</Button>
              </div>
            </div>
          </section>

          <div class="mx-5 my-4 h-px bg-slate-100" />

          <!-- JSON Output -->
          <section class="px-5 pb-4 space-y-2">
            <div class="flex items-center justify-between">
              <h2 class="text-sm font-semibold">Parsed Output (JSON)</h2>
              <div class="flex gap-2">
                <Button size="sm" variant="outline" @click="copyOutput">Copy JSON</Button>
                <Button size="sm" variant="outline" @click="downloadJson">Download</Button>
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

          <div class="mx-5 my-4 h-px bg-slate-100" />

          <!-- Parsed Fields -->
          <section class="px-5 pb-6">
            <div class="flex items-center justify-between mb-2">
              <h2 class="text-sm font-semibold">Parsed Fields</h2>
              <div class="flex gap-2">
                <Button size="sm" variant="outline" @click="copyTable">Copy as text</Button>
                <Button size="sm" variant="outline" @click="downloadTable">Download CSV</Button>
              </div>
            </div>

            <div v-if="Object.keys(qrParsed).length" class="overflow-auto border rounded-md">
              <table class="w-full text-sm">
                <thead class="bg-slate-50">
                <tr>
                  <th class="text-left px-3 py-2 border-b w-36">Tag</th>
                  <th class="text-left px-3 py-2 border-b w-56">Field Name</th>
                  <th class="text-left px-3 py-2 border-b">Value</th>
                </tr>
                </thead>
                <tbody>
                <tr
                    v-for="[k, v] in kvRows"
                    :key="k"
                    class="odd:bg-white even:bg-slate-50/40"
                >
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
                  <td
                      class="px-3 py-2 border-b break-all [overflow-wrap:anywhere] font-mono text-[12px]"
                  >
                    {{ v }}
                  </td>
                </tr>
                </tbody>
              </table>
            </div>

            <div
                v-else
                class="text-xs text-slate-500 border rounded-md px-3 py-6 grid place-items-center"
            >
              No parsed fields yet — paste a TLV string above and hit
              <span class="mx-1 font-semibold">Parse</span>.
            </div>
          </section>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { toast } from 'vue-sonner'

import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'

import { parseTag, emvTags } from '@/lib/parser/parser.ts'

/* ---------------- State ---------------- */
const qrInput = ref<string>('')
const qrParsed = ref<Record<string, string>>({})
const qrError = ref<string | null>(null)

/* Derived */
const qrOutput = computed(() => {
  try { return JSON.stringify(qrParsed.value, null, 2) } catch { return '{}' }
})
const kvRows = computed<[string, string][]>(() =>
    Object.entries(qrParsed.value).sort((a, b) => a[0].localeCompare(b[0]))
)
const qrLength = computed(() => (qrInput.value || '').length)
const tagCount = computed(() => {
  try {
    let i = 0, n = 0, s = (qrInput.value || '').trim()
    while (i + 4 <= s.length) {
      const len = Number.parseInt(s.slice(i + 2, i + 4), 10)
      if (Number.isNaN(len)) break
      const end = i + 4 + len
      if (end > s.length) break
      n++; i = end
    }
    return n
  } catch { return 0 }
})

function tagId(key: string) {
  // Extract base tag (first two digits)
  return key.match(/^data(\d{2})/)?.[1] || ''
}

/* ---------------- Actions ---------------- */
function parseQr() {
  const raw = (qrInput.value || '').trim()
  if (!raw) {
    qrParsed.value = {}
    qrError.value = null
    return
  }
  try {
    qrParsed.value = parseTag(raw)
    qrError.value = null
  } catch (e: any) {
    qrParsed.value = {}
    qrError.value = e?.message || 'Invalid EMV string.'
  }
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
  qrError.value = null
}

function loadSample() {
  qrInput.value =
      '00020101021226250006ID.CO.QRIS.WWW011893600915123456789ABC0208COMMERCE520400005303360540412345802ID5910TOKO ABC6010JAKARTA61051234062070702AB6304ABCD'
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
    else {
      qrParsed.value = {}
      qrError.value = null
    }
  }, 250)
})
</script>

<style scoped>
/* Minimal styling handled via Tailwind classes */
</style>
