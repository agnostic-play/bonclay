<template>
  <div class="min-h-screen bg-white p-5 pt-0">
    <div class="mx-auto">
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent class="">
          <!-- Header -->
          <div class="flex items-start justify-between gap-4 mb-3">
            <div class="min-w-0 max-w-3xl">
              <h1 class="text-3xl font-bold truncate">
                {{ isLoadingCollection ? 'Loading…' : (collection?.name || 'Untitled Collection') }}
              </h1>
              <p class="mt-1 text-slate-600 whitespace-pre-line break-all [overflow-wrap:anywhere]">
                {{ collection?.description || '' }}
              </p>
            </div>

            <div class="flex gap-2 flex-none">
              <Combobox
                  :model-value="searchText"
                  @update:modelValue="onComboChange"
                  @openChange="onComboOpenChange"
              >
                <ComboboxAnchor as-child>
                  <ComboboxTrigger as-child>
                    <Button class="justify-between min-w-56" variant="outline">
                      {{
                        diagrams.find(d => String(d.id) === String(selectedId))?.title
                        ?? (isLoadingDiagrams ? 'Loading…' : (diagrams.length ? 'Select Diagram' : 'No diagrams'))
                      }}
                      <ChevronsUpDown class="ml-2 h-4 w-3 shrink-0 opacity-50"/>
                    </Button>
                  </ComboboxTrigger>
                </ComboboxAnchor>

                <ComboboxList>
                  <div class="relative w-full max-w-sm items-center">
                    <!-- controlled input so it never mirrors the UUID -->
                    <ComboboxInput
                        :model-value="comboSearch"
                        @update:modelValue="comboSearch = $event"
                        class="pl-9 border-0 border-b rounded-none h-10"
                        placeholder="Search diagram…"
                    />
                    <span class="absolute left-0 inset-y-0 flex items-center justify-center px-3">
                      <Search class="size-4 text-muted-foreground"/>
                    </span>
                  </div>

                  <ComboboxEmpty>No diagrams found.</ComboboxEmpty>

                  <ComboboxGroup v-if="diagrams.length">
                    <ComboboxItem
                        v-for="d in diagrams"
                        :key="d.id"
                        :value="String(d.id)"
                        :text-value="d.title"
                    >
                      {{ d.title }}
                      <ComboboxItemIndicator>
                        <Check :class="cn('ml-auto h-4 w-4')"/>
                      </ComboboxItemIndicator>
                    </ComboboxItem>
                  </ComboboxGroup>

                  <ComboboxItem value="__create__" class="text-primary">
                    <FilePlus2 class="mr-2 h-4 w-4"/>
                    Create new diagram
                  </ComboboxItem>
                </ComboboxList>
              </Combobox>

              <Button variant="outline" size="sm" class="h-9" @click="openEditCollection">
                <Pencil class="mr-2 h-4 w-4"/>
                Edit Collection
              </Button>
            </div>
          </div>

          <!-- Editor + Preview -->
          <div class="h-[90vh] mt-5 w-full flex gap-4">
            <!-- Editor Panel -->
            <div class="w-2/5 h-full flex flex-col rounded-lg border bg-white">
              <div class="flex items-center justify-between p-3 border-b">
                <h2 class="text-sm font-semibold">Editor</h2>
                <div class="text-xs text-muted-foreground flex items-center gap-2">
                  <span v-if="lastUpdatedText" class="text-xs text-muted-foreground">{{ lastUpdatedText }}</span>
                  <Button size="sm" class="px-2 py-1 border text-xs" @click="saveDiagram">
                    <Save/>
                  </Button>

                  <!-- Share Diagram button -->
                  <Button size="sm" variant="outline" class="px-2 py-1 text-xs" @click="openShare">
                    <Share2 class="mr-1 h-3.5 w-3.5" />
                    Share
                  </Button>

                  <!-- Edit Diagram button -->
                  <Button size="sm" variant="outline" class="px-2 py-1 text-xs" @click="openEditDiagram">
                    <Pencil class="mr-1 h-3.5 w-3.5" />
                    Edit
                  </Button>
                </div>
              </div>

              <div ref="editorEl" class="flex-1 min-h-0"/>

              <div v-if="error" class="m-3 rounded-md border border-red-200 bg-red-50 p-3 text-red-700 text-xs">
                <strong class="font-medium">Syntax error:</strong>
                <span class="ml-1">{{ error }}</span>
              </div>
            </div>

            <!-- Preview Panel -->
            <div class="w-3/5 h-full rounded-lg border bg-white flex flex-col">
              <div class="flex items-center justify-between p-3 border-b">
                <h2 class="text-sm font-semibold">Live Preview</h2>
                <div class="flex items-center gap-2">
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="zoomOut()">−</button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="resetView()">Reset</button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="zoomIn()">＋</button>
                  <div class="w-px h-5 bg-gray-200 mx-1"/>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="downloadSVG">SVG</button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="downloadPNG">PNG</button>
                  <div class="w-px h-5 bg-gray-200 mx-1"/>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="openFullScreen">Full
                    screen
                  </button>
                </div>
              </div>

              <div class="relative flex-1 overflow-hidden">
                <!-- Mermaid: render inline SVG -->
                <div v-if="syntaxType === 'mermaid'" ref="previewRef" :class="{ 'opacity-50': isRendering }" class="absolute inset-0 select-none"
                     v-html="svgMarkup"/>
                <!-- PlantUML: render via server image -->
                <div v-else-if="syntaxType === 'plantuml'" class="absolute inset-0 overflow-hidden" ref="previewRef">
                  <img v-if="plantUmlUrl" :src="plantUmlUrl" :class="{ 'opacity-50': isRendering }" class="select-none" alt="PlantUML Diagram" @load="onPlantUmlImgLoad" @error="onPlantUmlImgError" />
                  <div v-else-if="!error" class="absolute inset-0 flex items-center justify-center text-sm text-gray-400">Rendered diagram will appear here</div>
                </div>
                <div v-if="syntaxType === 'mermaid' && !svgMarkup && !error" class="absolute inset-0 grid place-items-center text-sm text-gray-400">
                  Rendered diagram will appear here
                </div>
                <div v-if="error" class="absolute bottom-0 left-0 right-0 m-3 rounded-md border border-red-200 bg-red-50 p-3 text-red-700 text-xs">
                  <strong class="font-medium">Error:</strong>
                  <span class="ml-1">{{ error }}</span>
                </div>
              </div>
            </div>

            <!-- Full Screen Overlay -->
            <div v-if="isFullScreen" class="fixed inset-0 z-50 bg-white/95 backdrop-blur-sm">
              <div class="flex items-center justify-between p-3 border-b bg-white/70">
                <div class="text-sm font-medium">Full Screen Preview</div>
                <div class="flex items-center gap-2">
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="zoomOut('fs')">−</button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="resetView('fs')">Reset
                  </button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="zoomIn('fs')">＋</button>
                  <div class="w-px h-5 bg-gray-200 mx-1"/>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="downloadSVG">SVG</button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="downloadPNG">PNG</button>
                  <div class="w-px h-5 bg-gray-200 mx-1"/>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="closeFullScreen">Close
                    (Esc)
                  </button>
                </div>
              </div>
              <div class="relative h-[calc(100vh-48px)]">
                <div v-if="syntaxType === 'mermaid'" ref="fsPreviewRef" class="absolute inset-0 select-none" v-html="svgMarkup"/>
                <div v-else-if="syntaxType === 'plantuml'" class="absolute inset-0 overflow-hidden" ref="fsPreviewRef">
                  <img v-if="plantUmlUrl" :src="plantUmlUrl" class="select-none" alt="PlantUML Diagram" @load="onPlantUmlFsImgLoad" />
                </div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>

    <!-- Edit Collection Modal -->
    <Dialog v-model:open="showEditCollection">
      <DialogContent class="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Edit Collection</DialogTitle>
          <DialogDescription>Update the collection details.</DialogDescription>
        </DialogHeader>

        <div class="grid gap-3 py-2">
          <div class="grid gap-1.5">
            <Label for="ec-title">Title</Label>
            <Input id="ec-title" v-model="editForm.name" placeholder="Collection title"/>
          </div>
          <div class="grid gap-1.5">
            <Label for="ec-desc">Description</Label>
            <Textarea id="ec-desc" v-model="editForm.description" rows="3"/>
          </div>
        </div>

        <DialogFooter class="flex items-center justify-between">
          <Button variant="destructive" @click="removeCollection">
            <Trash2 class="mr-2 h-4 w-4"/>
            Remove
          </Button>
          <div class="space-x-2">
            <Button variant="outline" @click="showEditCollection = false">Cancel</Button>
            <Button :disabled="!editForm.name.trim()" @click="submitEditCollection">
              <Save class="mr-2 h-4 w-4"/>
              Update
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Create Diagram Modal -->
    <Dialog v-model:open="showCreate">
      <DialogContent class="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Create Diagram</DialogTitle>
          <DialogDescription>Provide a name, syntax type and (optional) initial content.</DialogDescription>
        </DialogHeader>

        <div class="grid gap-3 py-2">
           <div class="grid gap-1.5">
            <Label for="cd-syntax-type">Syntax Type</Label>
            <Select v-model="createForm.syntaxType" @update:modelValue="onCreateSyntaxTypeChange">
              <SelectTrigger id="cd-syntax-type">
                <SelectValue placeholder="Select type"/>
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="mermaid">Mermaid</SelectItem>
                <SelectItem value="plantuml">PlantUML</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="grid gap-1.5">
            <Label for="cd-title">Title</Label>
            <Input id="cd-title" v-model="createForm.title" placeholder="My Diagram"/>
          </div>
          <div class="grid gap-1.5">
            <Label for="cd-syntax">Initial Syntax</Label>
            <Textarea id="cd-syntax" v-model="createForm.syntax" rows="6" class="font-mono text-xs"/>
          </div>
        </div>

        <DialogFooter>
          <Button variant="outline" @click="showCreate = false">Cancel</Button>
          <Button :disabled="!createForm.title.trim()" @click="submitCreateDiagram">
            <FilePlus2 class="mr-2 h-4 w-4"/>
            Create
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Edit Diagram Modal (with Remove button) -->
    <Dialog v-model:open="showEditDiagram">
      <DialogContent class="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Edit Diagram</DialogTitle>
          <DialogDescription>Update the diagram details.</DialogDescription>
        </DialogHeader>

        <div class="grid gap-3 py-2">
             <div class="grid gap-1.5">
            <Label for="ed-syntax-type">Syntax Type</Label>
            <Select v-model="editDiagramForm.syntaxType">
              <SelectTrigger id="ed-syntax-type">
                <SelectValue placeholder="Select type"/>
              </SelectTrigger>
              <SelectContent>
                <SelectItem value="mermaid">Mermaid</SelectItem>
                <SelectItem value="plantuml">PlantUML</SelectItem>
              </SelectContent>
            </Select>
          </div>
          <div class="grid gap-1.5">
            <Label for="ed-title">Title</Label>
            <Input id="ed-title" v-model="editDiagramForm.title" placeholder="Diagram title"/>
          </div>
          <div class="grid gap-1.5">
            <Label for="ed-desc">Description</Label>
            <Textarea id="ed-desc" v-model="editDiagramForm.description" rows="3"/>
          </div>
        </div>

        <DialogFooter class="flex items-center justify-between">
          <Button variant="destructive" @click="removeDiagram">
            <Trash2 class="mr-2 h-4 w-4"/>
            Remove
          </Button>
          <div class="space-x-2">
            <Button variant="outline" @click="showEditDiagram = false">Cancel</Button>
            <Button :disabled="!editDiagramForm.title.trim()" @click="submitEditDiagram">
              <Save class="mr-2 h-4 w-4"/>
              Update
            </Button>
          </div>
        </DialogFooter>
      </DialogContent>
    </Dialog>

    <!-- Share Diagram Modal -->
    <Dialog v-model:open="showShare">
      <DialogContent class="sm:max-w-lg">
        <DialogHeader>
          <DialogTitle>Share Diagram</DialogTitle>
          <DialogDescription>Anyone with this link can view the diagram preview.</DialogDescription>
        </DialogHeader>

        <div class="grid gap-3 py-2">
          <div class="grid gap-1.5">
            <Label for="share-url">Share URL</Label>
            <div class="flex items-center gap-2">
              <Input id="share-url" :model-value="shareUrl" readonly class="font-mono text-xs" />
              <Button variant="outline" size="sm" @click="copyShareUrl">
                <Copy class="mr-2 h-4 w-4" />
                Copy
              </Button>
            </div>
          </div>
        </div>

        <DialogFooter class="justify-end">
          <Button variant="outline" @click="showShare = false">Close</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  </div>
</template>

<style scoped>
:deep(.monaco-scrollable-element)::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

:deep(.monaco-scrollable-element)::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 6px;
}

:deep(.monaco-scrollable-element)::-webkit-scrollbar-track {
  background: transparent;
}
</style>

<!-- src/pages/DiagramCollectionShow.vue -->
<script lang="ts" setup>
import {nextTick, onBeforeUnmount, onMounted, ref, watch, computed} from 'vue'
import {useRoute, useRouter} from 'vue-router'
import mermaid from 'mermaid'
import panzoom, {type PanZoom} from 'panzoom'
import * as monaco from "monaco-editor"
import {loader} from "@guolao/vue-monaco-editor"
import plantumlEncoder from 'plantuml-encoder'

loader.config({monaco})
import {toast} from 'vue-sonner'

import {Card, CardContent} from '@/components/ui/card'
import {Button} from '@/components/ui/button'
import {Check, ChevronsUpDown, FilePlus2, Save, Search, Pencil, Trash2, Share2, Copy} from 'lucide-vue-next'
import {cn} from '@/lib/utils'
import {
  Combobox, ComboboxAnchor, ComboboxEmpty, ComboboxGroup,
  ComboboxInput, ComboboxItem, ComboboxItemIndicator, ComboboxList, ComboboxTrigger,
} from '@/components/ui/combobox'
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogDescription, DialogFooter} from '@/components/ui/dialog'
import {Input} from '@/components/ui/input'
import {Textarea} from '@/components/ui/textarea'
import {Label} from '@/components/ui/label'
import {Select, SelectContent, SelectItem, SelectTrigger, SelectValue} from '@/components/ui/select'

/* 🔌 API service */
import diagramCollectionAPI from '@/api/diagramCollectionServices'
import type {DiagramCollection, DiagramDetail, DiagramSummary} from '@/types/entities'
import {useApiFeedback} from "@/composables/useApiFeedback"
import diagramServices from "@/api/diagramServices";

/* ---------------- route & id ---------------- */
const route = useRoute()
const router = useRouter()
const collectionId = computed(() => String(route.params.id ?? ''))

/* ---------------- state ---------------- */
const collection = ref<DiagramCollection | null>(null)
const isLoadingCollection = ref(false)

const diagrams = ref<DiagramSummary[]>([])
const isLoadingDiagrams = ref(false)

const selectedId = ref<string>('')                 // combobox v-model
const searchText = ref<string>('')                 // combobox v-model
const selectedDiagram = ref<DiagramDetail | null>(null)
const isSavingDiagram = ref(false)

const apiError = ref<string | null>(null)

/* combobox search term (prevents UUID showing in the input) */
const comboSearch = ref('')

function onComboOpenChange(open: boolean) {
  if (open) comboSearch.value = ''
}

/* ---------------- mermaid + monaco ---------------- */
mermaid.initialize({startOnLoad: false, securityLevel: 'loose', theme: 'default'})
const source = ref<string>('') // editor text
const svgMarkup = ref<string>('') // rendered svg (mermaid)
const plantUmlUrl = ref<string>('') // rendered img url (plantuml)
const error = ref<string | null>(null)
const isRendering = ref(false)
const previewRef = ref<HTMLDivElement | null>(null)

/* current diagram syntax type */
const syntaxType = computed<'mermaid' | 'plantuml'>(() => {
  const t = (selectedDiagram.value as any)?.syntax_type
      || (selectedDiagram.value as any)?.syntaxType
      || ''
  if (t === 'plantuml') return 'plantuml'
  if (t === 'mermaid') return 'mermaid'
  // fallback: detect from the editor content (works even if backend omits syntax_type)
  const syn = (source.value ?? '').trimStart()
  if (syn.startsWith('@start')) return 'plantuml'
  return 'mermaid'
})

const PLANTUML_SERVER = 'https://www.plantuml.com/plantuml/svg'

function buildPlantUmlUrl(text: string): string {
  try {
    const encoded = plantumlEncoder.encode(text)
    return `${PLANTUML_SERVER}/${encoded}`
  } catch {
    return ''
  }
}

function onPlantUmlImgError() {
  error.value = 'PlantUML server could not render the diagram. Check your syntax.'
}

let pz: PanZoom | null = null
let debounceTimer: number | undefined

const editorEl = ref<HTMLDivElement | null>(null)
let editor: import('monaco-editor').editor.IStandaloneCodeEditor | null = null

/* fullscreen */
const isFullScreen = ref(false)
const fsPreviewRef = ref<HTMLDivElement | null>(null)
let fsPz: PanZoom | null = null

/* computed helpers */
const lastUpdatedText = computed(() => {
  const iso = (selectedDiagram.value as any)?.latestUpdated
  return iso ? `Last updated ${new Date(iso).toLocaleString()}` : ''
})

/* ---------------- sharing ---------------- */
// sharing refs & constants
const showShare = ref(false)
const SHARE_BASE = new URL(diagramServices.getBaseURL()).origin
const shareUrl = ref<string>('')

// small helper to make base64 URL-safe (base64url)
function toBase64Url(s: string) {
  const b64 = btoa(s)
  return b64.replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/g, '')
}

function openShare(): void {
  if (!selectedDiagram.value) {
    toast.warning('Select a diagram first.')
    return
  }

  const id = String(selectedDiagram.value.id || '')
  if (!id) {
    toast.error('Missing diagram id.')
    return
  }

  try {
    // encode the UUID and make it URL-safe
    const b64url = toBase64Url(id)
    shareUrl.value = `${SHARE_BASE}/share/diagram/${b64url}`
  } catch {
    // fallback: still produce a valid URL
    shareUrl.value = `${SHARE_BASE}/share/diagram/${encodeURIComponent(id)}`
  }

  showShare.value = true
}

async function copyShareUrl(): Promise<void> {
  if (!shareUrl.value) {
    toast.error('No share URL available.')
    return
  }
  try {
    await navigator.clipboard.writeText(shareUrl.value)
    toast.success('Link copied to clipboard')
  } catch {
    // fallback for older browsers
    const ta = document.createElement('textarea')
    ta.value = shareUrl.value
    document.body.appendChild(ta)
    ta.select()
    try {
      document.execCommand('copy')
      toast.success('Link copied to clipboard')
    } catch {
      toast.error('Failed to copy link')
    } finally {
      document.body.removeChild(ta)
    }
  }
}


/* ---------------- helpers ---------------- */
function disposePanzoom(i: PanZoom | null) {
  try {
    (i as any)?.dispose?.()
  } catch {}
}

function attachPanZoom(containerRef: typeof previewRef, fullscreen: boolean) {
  if (fullscreen) {
    disposePanzoom(fsPz)
    fsPz = null
  } else {
    disposePanzoom(pz)
    pz = null
  }
  const container = containerRef.value
  if (!container) return

  if (syntaxType.value === 'plantuml') {
    const imgEl = container.querySelector('img') as HTMLImageElement | null
    if (!imgEl) return
    const instance = panzoom(imgEl, { smoothScroll: true, zoomSpeed: 0.065, maxZoom: 8, minZoom: 0.1, bounds: false })
    if (fullscreen) fsPz = instance; else pz = instance
    const rect = container.getBoundingClientRect()
    const nw = imgEl.naturalWidth || 800
    const nh = imgEl.naturalHeight || 600
    if (rect.width > 0 && nw > 0) {
      const scale = Math.min(rect.width / nw, rect.height / nh) * 0.9
      instance.zoomAbs(0, 0, scale)
      instance.moveTo((rect.width - nw * scale) / 2, (rect.height - nh * scale) / 2)
    }
    return
  }

  const svgEl = container.querySelector('svg') as SVGSVGElement | null
  if (!svgEl) return
  svgEl.removeAttribute('width')
  svgEl.removeAttribute('height')
  svgEl.setAttribute('preserveAspectRatio', 'xMidYMid meet')
  svgEl.style.width = '100%'
  svgEl.style.height = '100%'
  const instance = panzoom(svgEl, { smoothScroll: true, zoomSpeed: 0.065, maxZoom: 8, minZoom: 0.25, bounds: false })
  if (fullscreen) fsPz = instance; else pz = instance
}

function onPlantUmlImgLoad() {
  void nextTick(() => attachPanZoom(previewRef, false))
}

function onPlantUmlFsImgLoad() {
  void nextTick(() => attachPanZoom(fsPreviewRef, true))
}

async function validateSyntax(text: string) {
  const trimmed = (text ?? '').trim()
  if (!trimmed) {
    error.value = null;
    return false
  }
  try {
    await (mermaid as any).parse(trimmed)
    error.value = null
    return true
  } catch (e: any) {
    error.value = e?.message || 'Mermaid syntax error.'
    return false
  }
}

async function renderDiagram() {
  error.value = null
  const trimmed = (source.value ?? '').trim()

  if (syntaxType.value === 'plantuml') {
    svgMarkup.value = ''
    if (!trimmed) {
      plantUmlUrl.value = ''
      return
    }
    isRendering.value = true
    try {
      plantUmlUrl.value = buildPlantUmlUrl(trimmed)
    } finally {
      isRendering.value = false
    }
    return
  }

  // mermaid path
  plantUmlUrl.value = ''
  const ok = await validateSyntax(source.value)
  if (!ok) {
    svgMarkup.value = '';
    return
  }
  isRendering.value = true
  try {
    const id = 'mmd-' + Math.random().toString(36).slice(2)
    const {svg} = await mermaid.render(id, source.value)
    svgMarkup.value = svg
    await nextTick()
    attachPanZoom(previewRef, false)
    if (isFullScreen.value) attachPanZoom(fsPreviewRef, true)
  } catch (e: any) {
    error.value = e?.message || 'Failed to render diagram.'
  } finally {
    isRendering.value = false
  }
}

function onEditorChange() {
  if (debounceTimer) window.clearTimeout(debounceTimer)
  if (editor) source.value = editor.getValue()
  debounceTimer = window.setTimeout(() => {
    void renderDiagram()
  }, 300)
}

function getEditorLanguage(): string {
  return syntaxType.value === 'plantuml' ? 'plantuml' : 'mermaid'
}

function createEditor() {
  if (!editorEl.value) return
  loader.init().then((m) => {
    try {
      // ── Mermaid language ──────────────────────────────────────────────
      m.languages.register({id: 'mermaid'})
      m.languages.setMonarchTokensProvider('mermaid', {
        tokenizer: {
          root: [
            [/\b(flowchart|graph|sequenceDiagram|classDiagram|stateDiagram-v2|erDiagram|journey|gantt|pie)\b/, 'keyword'],
            [/\b(subgraph|end|click|style|linkStyle|accTitle|accDescr)\b/, 'keyword.control'],
            [/-->|==>|===|\+\+|--|==/, 'operator'],
            [/\[[^\]]*\]|\([^\)]*\)|\{[^}]*\}/, 'string'],
            [/"[^"]*"|'[^']*'/, 'string'],
            [/%%.*/, 'comment'],
          ],
        },
      })

      // ── PlantUML language ─────────────────────────────────────────────
      m.languages.register({id: 'plantuml'})
      m.languages.setMonarchTokensProvider('plantuml', {
        tokenizer: {
          root: [
            // @startuml / @enduml
            [/@start\w+|@end\w+/, 'keyword'],
            // single-line comments (' is the comment char in PlantUML, NOT a string delimiter)
            [/'.*$/, 'comment'],
            // block comment open  /'
            [/\/'\//, 'comment', '@blockComment'],
            // preprocessor directives  !include !define etc.
            [/![a-z_]+/, 'keyword.control'],
            // arrows (longest patterns first)
            [/<\|--|--\|>/, 'operator'],
            [/<\*--|--\*>/, 'operator'],
            [/<o--|--o>|<#--|--#>/, 'operator'],
            [/<\.\.|\.\.>/, 'operator'],
            [/<--|-->|->|<-/, 'operator'],
            [/\.\.|--/, 'operator'],
            // keywords
            [/\b(participant|actor|boundary|control|entity|database|collections|queue|usecase|class|interface|enum|abstract|component|package|node|folder|frame|cloud|rectangle|namespace|object|map|state|hide|show|skinparam|title|note|ref|group|alt|else|opt|loop|break|par|critical|activate|deactivate|destroy|autonumber|newpage|box|together|end)\b/, 'keyword'],
            // skin colour tags  #RoyalBlue  #aabbcc
            [/#[a-zA-Z0-9_]+/, 'tag'],
            // double-quoted strings
            [/"[^"]*"/, 'string'],
            // numbers
            [/\b\d+(\.\d+)?\b/, 'number'],
          ],
          blockComment: [
            [/'\//, 'comment', '@pop'],
            [/[\s\S]/, 'comment'],
          ],
        },
      })
    } catch {}

    editor = m.editor.create(editorEl.value!, {
      value: source.value,
      language: getEditorLanguage(),
      automaticLayout: true,
      fontSize: 12,
      wordWrap: 'on',
      minimap: {enabled: false},
      scrollBeyondLastLine: false,
      theme: 'vs',
    })
    editor?.onDidChangeModelContent(onEditorChange)
  })
}

/* ---------------- API glue ---------------- */
async function loadCollection() {
  if (!collectionId.value) return
  isLoadingCollection.value = true;
  apiError.value = null
  try {
    const res: DiagramCollection = await diagramCollectionAPI.getDetail(collectionId.value);
    collection.value = res
  } catch (e: any) {
    apiError.value = e?.message || 'Failed to load collection.'
  } finally {
    isLoadingCollection.value = false
  }
}

const {withApiFeedback} = useApiFeedback();

/** list diagrams; first selection uses list payload (no extra API call) */
async function loadDiagramList() {
  if (!collectionId.value) return
  isLoadingDiagrams.value = true;
  apiError.value = null
  try {
    const data = await withApiFeedback(
        diagramCollectionAPI.getDiagram(collectionId.value),
        {errorMessage: "Failed to load diagrams."}
    )
    const list: DiagramSummary[] = data.data?.list ?? []
    diagrams.value = list.map(d => ({...d, id: String(d.id)}))
    if (diagrams.value.length) {
      selectedId.value = ""   // watcher will render from list
    } else {
      selectedId.value = ''
      selectedDiagram.value = null
      source.value = ''
      svgMarkup.value = ''
      error.value = null
    }
  } catch (e: any) {
    apiError.value = e?.message || 'Failed to load diagrams.'
  } finally {
    isLoadingDiagrams.value = false
  }
}

/** prefer list payload for detail; fallback to API detail call if needed */
function fromListAsDetail(id: string): DiagramDetail | null {
  const hit = diagrams.value.find(d => String(d.id) === String(id))
  if (!hit) return null
  return {
    id: String(hit.id),
    title: hit.title,
    description: (hit as any).description ?? '',
    syntax: (hit as any).syntax ?? '',
    syntax_type: (hit as any).syntax_type ?? hit.syntax_type ?? '',
    syntaxType: (hit as any).syntaxType ?? hit.syntax_type ?? '',
    latestUpdated: (hit as any).latestUpdated,
  } as DiagramDetail
}

async function fetchDetail(id: string) {
  const svc: any = diagramCollectionAPI as any
  const d: DiagramDetail =
      (svc.getDiagramDetail ? await svc.getDiagramDetail(collectionId.value, id)
          : await svc.getDiagram(collectionId.value, id))
  return d
}

async function saveDiagram() {
  if (!selectedDiagram.value) {
    toast.warning("need to select diagram, before update")
    return
  }

  isSavingDiagram.value = true
  try {
    const payload: any = {
      syntax: source.value,
      title: selectedDiagram.value.title,
      syntax_type: syntaxType.value,
    }

    const svc = await withApiFeedback(
        diagramServices.actUpdate(selectedDiagram.value.id, payload),
        { successMessage: 'Diagram saved' }
    )

    if (svc.error || Object.keys(svc.fieldErrors || {}).length > 0) {
      const fe = svc.fieldErrors
      const details = Object.entries(fe).map(([k,v]) => `${k}: ${v}`).join(', ')
      toast.error(details || 'Failed to save diagram.')
      return
    }

    selectedDiagram.value.syntax = source.value
    const idx = diagrams.value.findIndex(d => String(d.id) === String(selectedDiagram.value!.id))
    if (idx >= 0) (diagrams.value[idx] as any).syntax = source.value
    toast.success("Successfully updated diagram.")
    await renderDiagram()
  } catch (e: any) {
    toast.error(e?.message || 'Failed to save diagram.')
  } finally {
    isSavingDiagram.value = false
  }
}

/** update collection */
async function submitEditCollection() {
  try {
    await (diagramCollectionAPI as any).actUpdate(collectionId.value, {
      name: editForm.value.name.trim(),
      description: editForm.value.description,
    })
    toast.success('Collection updated')
    showEditCollection.value = false
    await loadCollection()
  } catch (e: any) {
    toast.error(e?.message || 'Failed to update collection.')
  }
}

/** remove collection */
async function removeCollection() {
  try {
    await (diagramCollectionAPI as any).remove?.(collectionId.value)
    toast('Collection removed')
    router.replace({name: 'MermaidDiagramTools-ProjectList'})
  } catch (e: any) {
    toast.error(e?.message || 'Failed to remove collection.')
  }
}

/** create new diagram */
const showCreate = ref(false)

const MERMAID_DEFAULT_SYNTAX = 'flowchart TD\n  A[Start] --> B[End]\n'
const PLANTUML_DEFAULT_SYNTAX = '@startuml\nA -> B: Hello\n@enduml\n'

const createForm = ref({
  title: '',
  syntax: MERMAID_DEFAULT_SYNTAX,
  syntaxType: 'mermaid' as 'mermaid' | 'plantuml',
})

function onCreateSyntaxTypeChange(val: unknown) {
  val = String(val ?? '')
  const t = val as 'mermaid' | 'plantuml'
  // only reset syntax if it's still the default for the other type
  if (
    (t === 'plantuml' && (createForm.value.syntax === MERMAID_DEFAULT_SYNTAX || createForm.value.syntax === '')) ||
    (t === 'mermaid' && (createForm.value.syntax === PLANTUML_DEFAULT_SYNTAX || createForm.value.syntax === ''))
  ) {
    createForm.value.syntax = t === 'plantuml' ? PLANTUML_DEFAULT_SYNTAX : MERMAID_DEFAULT_SYNTAX
  }
  createForm.value.syntaxType = t
}

async function submitCreateDiagram() {
  if (!createForm.value.title.trim()) return
  try {
    const created = await diagramServices.actCreate({
      collection_id: collectionId.value,
      title: createForm.value.title.trim(),
      syntax: createForm.value.syntax,
      syntax_type: createForm.value.syntaxType,
    })
    toast.success('Diagram created')
    showCreate.value = false
    createForm.value = {title: '', syntax: MERMAID_DEFAULT_SYNTAX, syntaxType: 'mermaid'}
    await loadDiagramList()
    if (created?.id) selectedId.value = String(created.id)
  } catch (e: any) {
    toast.error('Failed to create diagram.')
  }
}

/* ---------------- edit collection modal ---------------- */
const showEditCollection = ref(false)
const editForm = ref({name: '', description: ''})

function openEditCollection() {
  editForm.value.name = collection.value?.name ?? ''
  editForm.value.description = collection.value?.description ?? ''
  showEditCollection.value = true
}

/* ---------------- edit diagram modal ---------------- */
const showEditDiagram = ref(false)
const editDiagramForm = ref({ title: '', description: '', syntaxType: 'mermaid' as 'mermaid' | 'plantuml' })

function openEditDiagram() {
  if (!selectedDiagram.value) {
    toast.warning('Select a diagram first.')
    return
  }
  editDiagramForm.value.title = selectedDiagram.value.title || ''
  editDiagramForm.value.description = selectedDiagram.value.description || ''
  editDiagramForm.value.syntaxType = syntaxType.value
  showEditDiagram.value = true
}

async function submitEditDiagram() {
  if (!selectedDiagram.value) return
  try {
    const payload: any = {
      title: editDiagramForm.value.title.trim(),
      description: editDiagramForm.value.description,
      syntax_type: editDiagramForm.value.syntaxType,
    }

    const svc = await withApiFeedback(
        diagramServices.actUpdate(selectedDiagram.value.id, payload),
        { successMessage: 'Diagram updated' }
    )

    if (svc?.error || Object.keys(svc?.fieldErrors || {}).length > 0) {
      const fe = svc.fieldErrors || {}
      const details = Object.entries(fe).map(([k, v]) => `${k}: ${v}`).join(', ')
      toast.error(details || 'Failed to update diagram.')
      return
    }

    // sync local state
    selectedDiagram.value.title = payload.title
    selectedDiagram.value.description = payload.description
    ;(selectedDiagram.value as any).syntax_type = payload.syntax_type
    ;(selectedDiagram.value as any).syntaxType = payload.syntax_type

    const idx = diagrams.value.findIndex(d => String(d.id) === String(selectedDiagram.value!.id))
    if (idx >= 0) {
      diagrams.value[idx].title = payload.title
      ;(diagrams.value[idx] as any).description = payload.description
      ;(diagrams.value[idx] as any).syntax_type = payload.syntax_type
      ;(diagrams.value[idx] as any).syntaxType = payload.syntax_type
    }

    showEditDiagram.value = false
    // re-render with potentially changed type
    await renderDiagram()
  } catch (e: any) {
    toast.error(e?.message || 'Failed to update diagram.')
  }
}

/** remove diagram (from Edit Diagram modal) */
async function removeDiagram() {
  if (!selectedDiagram.value) return
  try {
    await diagramServices.actRemove(selectedDiagram.value.id)
    toast.success('Diagram removed')
    showEditDiagram.value = false

    // reload list + clear current selection and preview
    await loadDiagramList()
    selectedId.value = ''
    selectedDiagram.value = null
    source.value = ''
    svgMarkup.value = ''
    error.value = null
  } catch (e: any) {
    toast.error(e?.message || 'Failed to remove diagram.')
  }
}

/* ---------------- downloads & fullscreen ---------------- */
function downloadSVG() {
  if (syntaxType.value === 'plantuml') {
    try {
      const encoded = plantumlEncoder.encode(source.value)
      window.open(`https://www.plantuml.com/plantuml/svg/${encoded}`, '_blank')
    } catch { toast.error('Failed to encode PlantUML diagram.') }
    return
  }
  if (!svgMarkup.value) { toast.error('No SVG to download.'); return }
  const blob = new Blob([svgMarkup.value], { type: 'image/svg+xml' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = (selectedDiagram.value?.title || 'diagram') + '.svg'
  a.click()
  URL.revokeObjectURL(url)
}

function downloadPNG() {
  if (syntaxType.value === 'plantuml') {
    try {
      const encoded = plantumlEncoder.encode(source.value)
      const a = document.createElement('a')
      a.href = `https://www.plantuml.com/plantuml/png/${encoded}`
      a.download = (selectedDiagram.value?.title || 'diagram') + '.png'
      a.click()
    } catch { toast.error('Failed to encode PlantUML diagram.') }
    return
  }
  if (!svgMarkup.value) { toast.error('No diagram to download.'); return }
  const blob = new Blob([svgMarkup.value], { type: 'image/svg+xml' })
  const url = URL.createObjectURL(blob)
  const img = new Image()
  img.onload = () => {
    const probe = document.createElement('div')
    probe.style.cssText = 'position:absolute;left:-9999px'
    probe.innerHTML = svgMarkup.value
    document.body.appendChild(probe)
    const s = probe.querySelector('svg')
    let w = 1600, h = 900
    try {
      const vb = (s?.getAttribute('viewBox') || '').split(/\s+/).map(Number)
      if (vb.length === 4 && vb[2] > 0 && vb[3] > 0) {
        w = Math.min(4000, Math.max(800, vb[2]))
        h = Math.min(4000, Math.max(600, vb[3]))
      }
    } catch {}
    document.body.removeChild(probe)
    const canvasEl = document.createElement('canvas')
    canvasEl.width = w; canvasEl.height = h
    const ctx = canvasEl.getContext('2d')!
    ctx.fillStyle = '#ffffff'; ctx.fillRect(0, 0, w, h)
    ctx.drawImage(img, 0, 0, w, h)
    canvasEl.toBlob((b) => {
      if (!b) return
      const a = document.createElement('a')
      a.href = URL.createObjectURL(b)
      a.download = (selectedDiagram.value?.title || 'diagram') + '.png'
      a.click()
    }, 'image/png', 1)
    URL.revokeObjectURL(url)
  }
  img.onerror = () => URL.revokeObjectURL(url)
  img.src = url
}

function openFullScreen() {
  isFullScreen.value = true
  nextTick(() => attachPanZoom(fsPreviewRef, true))
}

function closeFullScreen() {
  isFullScreen.value = false
  disposePanzoom(fsPz);
  fsPz = null
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && isFullScreen.value) {
    e.stopPropagation();
    closeFullScreen()
  }
}

function zoomIn(target: 'fs' | 'normal' = 'normal') {
  const inst = target === 'fs' ? fsPz : pz;
  inst?.smoothZoom?.(0, 0, 1.25)
}

function zoomOut(target: 'fs' | 'normal' = 'normal') {
  const inst = target === 'fs' ? fsPz : pz;
  inst?.smoothZoom?.(0, 0, 0.8)
}

function resetView(target: 'fs' | 'normal' = 'normal') {
  const inst = target === 'fs' ? fsPz : pz
  if (!inst) return
  if (syntaxType.value === 'plantuml') {
    const containerRef = target === 'fs' ? fsPreviewRef : previewRef
    const container = containerRef.value
    const imgEl = container?.querySelector('img') as HTMLImageElement | null
    if (imgEl && container) {
      const rect = container.getBoundingClientRect()
      const nw = imgEl.naturalWidth || 800
      const nh = imgEl.naturalHeight || 600
      const scale = Math.min(rect.width / nw, rect.height / nh) * 0.9
      inst.zoomAbs(0, 0, scale)
      inst.moveTo((rect.width - nw * scale) / 2, (rect.height - nh * scale) / 2)
    }
    return
  }
  inst.moveTo?.(0, 0)
  inst.zoomAbs?.(0, 0, 1)
}

/* ---------------- lifecycle ---------------- */
onMounted(async () => {
  createEditor()
  await loadCollection()
  await loadDiagramList()
  window.addEventListener('keydown', onKeydown)
})
onBeforeUnmount(() => {
  editor?.dispose();
  disposePanzoom(pz)
})

/* ---------------- Combobox behavior ---------------- */
function onComboChange(val: unknown) {
  const v = String(val ?? '')
  if (v === '__create__') {
    showCreate.value = true;
    return;
  }
  // accept only known IDs
  const exists = diagrams.value.some(d => String(d.id) === v)
  if (!exists) return
  selectedId.value = v
}

/* When selectedId changes: use list payload first; if absent, fetch detail */
watch(selectedId, async (id) => {
  if (!id || id === '__create__') return
  const fromList = fromListAsDetail(id)
  if (fromList) {
    selectedDiagram.value = fromList
    source.value = fromList.syntax ?? ''
    editor?.setValue(source.value)
    // update Monaco language to match syntax type
    await nextTick()
    if (editor) {
      const model = editor.getModel()
      if (model) monaco.editor.setModelLanguage(model, getEditorLanguage())
    }
    await renderDiagram()
  } else {
    try {
      const d = await fetchDetail(id)
      selectedDiagram.value = d
      source.value = d?.syntax ?? ''
      editor?.setValue(source.value)
      // update Monaco language to match syntax type
      await nextTick()
      if (editor) {
        const model = editor.getModel()
        if (model) monaco.editor.setModelLanguage(model, getEditorLanguage())
      }
      await renderDiagram()
    } catch (e: any) {
      toast.error(e?.message || 'Failed to load diagram.')
    }
  }
})
</script>
