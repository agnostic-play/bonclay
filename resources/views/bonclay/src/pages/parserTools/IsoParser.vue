
<template>
  <div class="min-h-screen bg-white p-5 pt-0">
    <div class="mx-auto max-w-4xl">
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
            <div class="flex flex-col md:flex-row md:items-center md:justify-between gap-2">
              <div class="flex items-center gap-2">
                <Label for="iso-raw" class="text-base font-semibold">Message</Label>
                <select v-model="selectedSpecKey" class="text-xs border rounded px-2 py-1">
                  <option value="iso87">Spec 87</option>
                  <option value="rintis">Spec Custom (Rintis)</option>
                  <option value="free_ll">Free (LL)</option>
                  <option value="free_lll">Free (LLL)</option>
                </select>
              </div>
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
              </div>
            </div>

            <Textarea
                id="iso-raw"
                v-model="raw"
                rows="10"
                class="font-mono text-xs w-full resize-y whitespace-pre-wrap break-all overflow-auto"
                placeholder="0200F23C448128E080000000000000000000001612345678901234561234560708123456123456123456..."
            />

            <div class="flex items-center justify-between">
              <p v-if="err" class="text-xs text-red-600">{{ err }}</p>
              <div class="flex gap-2">
                <Button size="sm" variant="outline" @click="loadSample">Sample</Button>
                <Button size="sm" variant="outline" @click="clearAll">Clear</Button>
                <Button size="sm" variant="outline" @click="copyOutput">Copy JSON</Button>
                <Button size="sm" @click="doParse">Parse</Button>
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
                id="iso-json"
                :model-value="jsonOut"
                readonly
                rows="8"
                class="font-mono text-xs"
            />
          </section>

          <div class="mx-5 my-4 h-px bg-slate-100" />

          <!-- Parsed Fields Table -->
          <section class="px-5 pb-6">
            <div class="flex items-center justify-between mb-2">
              <h2 class="text-sm font-semibold">Parsed Fields</h2>
              <div class="flex gap-2">
                <Button size="sm" variant="outline" @click="copyTable">Copy as text</Button>
                <Button size="sm" variant="outline" @click="downloadTable">Download CSV</Button>
              </div>
            </div>

            <div v-if="parsed && parsed.fields.length" class="overflow-auto border rounded-md">
              <table class="w-full text-sm">
                <thead class="bg-slate-50">
                <tr>
                  <th class="text-left px-3 py-2 border-b w-20">DE</th>
                  <th class="text-left px-3 py-2 border-b w-56">Field Name</th>
                  <th class="text-left px-3 py-2 border-b w-24">Len</th>
                  <th class="text-left px-3 py-2 border-b w-24">Type</th>
                  <th class="text-left px-3 py-2 border-b">Value</th>
                </tr>
                </thead>
                <tbody>
                <tr v-for="pf in parsed.fields" :key="pf.id" class="odd:bg-white even:bg-slate-50/40">
                  <td class="px-3 py-2 border-b font-mono text-[12px]">{{ pf.id }}</td>
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
              No parsed fields yet — paste an ISO 8583 message above and hit <span class="mx-1 font-semibold">Parse</span>.
            </div>
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
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, watch } from 'vue'
import { toast } from 'vue-sonner'

import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { Label } from '@/components/ui/label'

import { AVAILABLE_SPECS } from '@/lib/parser/iso8583/spec'
import type { MessageSpec } from '@/lib/parser/iso8583/spec'
import { parseISO8583 } from '@/lib/parser/iso8583/parser'

/* ---------------- State ---------------- */
const raw = ref<string>('')
const err = ref<string | null>(null)
const parsed = ref<ReturnType<typeof parseISO8583> | null>(null)
const selectedSpecKey = ref<'iso87' | 'rintis' | 'free_ll' | 'free_lll'>('iso87')

/* Derived */
const activeSpec = computed<MessageSpec>(() => AVAILABLE_SPECS[selectedSpecKey.value])
const jsonOut = computed(() => {
  try { return JSON.stringify(parsed.value?.map || {}, null, 2) } catch { return '{}' }
})
const msgLength = computed(() => (raw.value || '').length)
const fieldCount = computed(() => parsed.value?.fields.length || 0)
const presentBits = computed(() => parsed.value ? parsed.value.bits.reduce((a, b) => a + (b ? 1 : 0), 0) : 0)

function fieldName(id: number) {
  return activeSpec.value.fields[id]?.name || `DE ${id}`
}

/* ---------------- Actions ---------------- */
function doParse() {
  const s = (raw.value || '').trim()
  if (!s) { parsed.value = null; err.value = null; return }
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
  } catch { toast.error('Failed to copy') }
}

async function copyTable() {
  if (!parsed.value) return
  const text = parsed.value.fields
      .map(pf => `${pf.id.toString().padStart(3,'0')} (${pf.lenType}): ${pf.value}`)
      .join('\n')
  try { await navigator.clipboard.writeText(text); toast.success('Parsed fields copied') } catch { toast.error('Failed to copy') }
}

function clearAll() { raw.value = ''; parsed.value = null; err.value = null }

function loadSample() {
  // Simple synthetic sample: 0200 with DE2(LL=16), DE3(6), DE4(12), DE11(6)
  // Bitmap: bit1=0, fields: 2,3,4,11 => positions 2,3,4,11 => 64-bit bitmap
  // Bitmap bits (1-based): 2,3,4,11 => hex: 7? We'll craft via literal for demo.
  // Precomputed primary bitmap with bits 2,3,4,11 set:  F238000000000000
  const mti = '0200'
  const bitmap = '7020000000000000' // demo bitmap (hex ASCII)
  const de2 = '16' + '4567890123456789' // LL + PAN
  const de3 = '000000'
  const de4 = '000000001000'
  const de11 = '123456'
  raw.value = mti + bitmap + de2 + de3 + de4 + de11
  doParse()
}

function downloadJson() {
  try {
    const blob = new Blob([jsonOut.value], { type: 'application/json' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'iso8583-parse.json'
    a.click()
    URL.revokeObjectURL(url)
  } catch { toast.error('Failed to download JSON') }
}

function downloadTable() {
  try {
    if (!parsed.value) return
    const header = 'de,field_name,len_type,declared_len,actual_len,value'
    const rows = parsed.value.fields.map(pf => {
      const name = fieldName(pf.id)
      const decl = pf.declaredLength ?? ''
      return `${pf.id},"${name}",${pf.lenType},${decl},${pf.actualLength},"${String(pf.value).replace(/"/g,'""')}"`
    })
    const csv = [header, ...rows].join('\n')
    const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
    const url = URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = 'iso8583-parse.csv'
    a.click()
    URL.revokeObjectURL(url)
  } catch { toast.error('Failed to download CSV') }
}

/* Auto-parse (debounced) */
let t: number | undefined
watch(raw, (v) => {
  if (t) window.clearTimeout(t)
  t = window.setTimeout(() => { (v || '').trim() ? doParse() : (parsed.value = null, err.value = null) }, 250)
})
</script>

<style scoped>
/* Tailwind covers styling */
</style>
