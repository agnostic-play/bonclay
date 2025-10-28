<script lang="ts" setup>
import { nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import mermaid from 'mermaid'
import panzoom, { type PanZoom } from 'panzoom'
import loader from '@monaco-editor/loader'
import { Card, CardContent } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Check, ChevronsUpDown, FilePlus2, Share2, Search } from 'lucide-vue-next'
import { cn } from '@/lib/utils'
import { Save } from 'lucide-vue-next'
import {
  Combobox,
  ComboboxAnchor,
  ComboboxEmpty,
  ComboboxGroup,
  ComboboxInput,
  ComboboxItem,
  ComboboxItemIndicator,
  ComboboxList,
  ComboboxTrigger,
} from '@/components/ui/combobox'

// -----------------------------
// 👇 NEW: accept a projectId prop
// -----------------------------
const props = defineProps<{ projectId: string }>()

// -----------------------------
// Types for Project & Diagram
// -----------------------------
type Project = { title: string; desc: string }
type DiagramSummary = { id: string; title: string }
type Diagram = {
  id: string
  title: string
  syntax: string
  projectID: string
  desc?: string
  latestUpdated?: string
}

// -----------------------------
// State
// -----------------------------
const project = ref<Project | null>(null)
const isProjectEditing = ref(false)
const isProjectSaving = ref(false)

const diagrams = ref<DiagramSummary[]>([])
const selectedDiagram = ref<Diagram | null>(null)
const isLoadingDiagrams = ref(false)
const isSavingDiagram = ref(false)

// --- Mermaid setup ---
mermaid.initialize({ startOnLoad: false, securityLevel: 'loose', theme: 'default' })

const source = ref<string>(
    `flowchart TD\n  A[Start] --> B{Is it working?}\n  B -- Yes --> C[Ship it]\n  B -- No --> D[Fix it]\n  D --> B\n`
)

const svgMarkup = ref<string>('')
const error = ref<string | null>(null)
const isRendering = ref(false)
const previewRef = ref<HTMLDivElement | null>(null)

// Fullscreen state/refs
const isFullScreen = ref(false)
const fsPreviewRef = ref<HTMLDivElement | null>(null)
let fsPz: PanZoom | null = null

// Monaco
const editorEl = ref<HTMLDivElement | null>(null)
let monaco: typeof import('monaco-editor') | null = null
let editor: import('monaco-editor').editor.IStandaloneCodeEditor | null = null

let pz: PanZoom | null = null
let debounceTimer: number | undefined

const emit = defineEmits<{ (e: 'editCollection'): void; (e: 'createNewEndpoint'): void }>()

function validateSyntax(text: string) {
  try {
    // @ts-expect-error: upstream type gaps
    mermaid.parse(text)
    error.value = null
    return true
  } catch (e: any) {
    error.value = e?.message || 'Mermaid syntax error.'
    return false
  }
}

async function renderDiagram() {
  if (!validateSyntax(source.value)) return
  isRendering.value = true
  try {
    const id = 'mmd-' + Math.random().toString(36).slice(2)
    const { svg } = await mermaid.render(id, source.value)
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

function disposePanzoom(instance: PanZoom | null) {
  try {
    // @ts-expect-error: optional dispose
    instance?.dispose?.()
  } catch {/* noop */}
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
  const svgEl = container.querySelector('svg') as SVGSVGElement | null
  if (!svgEl) return

  svgEl.removeAttribute('width')
  svgEl.removeAttribute('height')
  svgEl.setAttribute('preserveAspectRatio', 'xMidYMid meet')
  svgEl.style.width = '100%'
  svgEl.style.height = '100%'

  const instance = panzoom(svgEl, { smoothScroll: true, zoomSpeed: 0.065, maxZoom: 8, minZoom: 0.25, bounds: false })
  if (fullscreen) fsPz = instance
  else pz = instance
}

function zoomIn(target: 'fs' | 'normal' = 'normal') {
  const inst = target === 'fs' ? fsPz : pz
  // @ts-expect-error: smoothZoom present at runtime
  inst?.smoothZoom?.(0, 0, 1.25)
}

function zoomOut(target: 'fs' | 'normal' = 'normal') {
  const inst = target === 'fs' ? fsPz : pz
  // @ts-expect-error
  inst?.smoothZoom?.(0, 0, 0.8)
}

function resetView(target: 'fs' | 'normal' = 'normal') {
  const inst = target === 'fs' ? fsPz : pz
  // @ts-expect-error
  if (inst) {
    inst.moveTo?.(0, 0)
    inst.zoomAbs?.(0, 0, 1)
  }
}

function downloadSVG() {
  if (!svgMarkup.value) return
  const blob = new Blob([svgMarkup.value], { type: 'image/svg+xml;charset=utf-8' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = 'diagram.svg'
  document.body.appendChild(a)
  a.click()
  a.remove()
  URL.revokeObjectURL(url)
}

async function downloadPNG() {
  if (!svgMarkup.value) return
  try {
    const svg = new Blob([svgMarkup.value], { type: 'image/svg+xml;charset=utf-8' })
    const url = URL.createObjectURL(svg)
    const img = new Image()
    img.crossOrigin = 'anonymous'

    await new Promise<void>((resolve, reject) => {
      img.onload = () => resolve()
      img.onerror = (e) => reject(e)
      img.src = url
    })

    const scale = 2
    const w = (img.naturalWidth || 1200) * scale
    const h = (img.naturalHeight || 800) * scale
    const canvas = document.createElement('canvas')
    canvas.width = w
    canvas.height = h
    const ctx = canvas.getContext('2d')!
    ctx.fillStyle = '#ffffff'
    ctx.fillRect(0, 0, w, h)
    ctx.drawImage(img, 0, 0, w, h)

    await new Promise<void>((resolve) => {
      canvas.toBlob((blob) => {
        if (!blob) return resolve()
        const pngUrl = URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = pngUrl
        a.download = 'diagram.png'
        document.body.appendChild(a)
        a.click()
        a.remove()
        URL.revokeObjectURL(pngUrl)
        resolve()
      }, 'image/png')
    })

    URL.revokeObjectURL(url)
  } catch (e) {
    console.error(e)
  }
}

function onEditorChange() {
  if (debounceTimer) window.clearTimeout(debounceTimer)
  if (editor) source.value = editor.getValue()
  validateSyntax(source.value)
  debounceTimer = window.setTimeout(renderDiagram, 250)
}

function createEditor() {
  if (!editorEl.value) return
  loader.init().then((m) => {
    monaco = m as any

    // Register a lightweight Mermaid language for syntax highlighting
    try {
      m.languages.register({ id: 'mermaid' })
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
    } catch {/* already registered */}

    editor = m.editor.create(editorEl.value!, {
      value: source.value,
      language: 'mermaid',
      fontFamily:
          'ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace',
      fontSize: 12,
      lineNumbers: 'on',
      automaticLayout: true,
      wordWrap: 'on',
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
      bracketPairColorization: { enabled: true },
      tabSize: 2,
      theme: 'vs',
    })

    editor.onDidChangeModelContent(onEditorChange)
  })
}

function openFullScreen() {
  isFullScreen.value = true
  nextTick(() => attachPanZoom(fsPreviewRef, true))
}

function closeFullScreen() {
  isFullScreen.value = false
  disposePanzoom(fsPz)
  fsPz = null
}

function onKeydown(e: KeyboardEvent) {
  if (e.key === 'Escape' && isFullScreen.value) {
    e.stopPropagation()
    closeFullScreen()
  }
}

// ---------------------------------
// 🔌 API INTEGRATION (NEW)
// ---------------------------------
const isProjectLoading = ref(false)
const apiError = ref<string | null>(null)

async function loadProject() {
  isProjectLoading.value = true
  apiError.value = null
  try {
    // @ts-ignore - API provided at runtime
    const res = await mermaidDiagramToolsProjectAPI.get(props.projectId)
    project.value = { title: res?.title ?? 'Untitled Project', desc: res?.description ?? '' }
  } catch (e: any) {
    apiError.value = e?.message || 'Failed to load project.'
  } finally {
    isProjectLoading.value = false
  }
}

async function loadDiagrams() {
  isLoadingDiagrams.value = true
  apiError.value = null
  try {
    // @ts-ignore - API provided at runtime
    const list = await mermaidDiagramToolsProjectAPI.getAllDiagram(props.projectId)
    diagrams.value = (list || []).map((d: any) => ({ id: String(d.id ?? d.diagramID ?? ''), title: d.title ?? 'Untitled' }))

    // Auto-select first diagram
    if (diagrams.value.length > 0) {
      await selectDiagram(diagrams.value[0].id)
    } else {
      selectedDiagram.value = null
      source.value = ''
      svgMarkup.value = ''
    }
  } catch (e: any) {
    apiError.value = e?.message || 'Failed to load diagrams.'
  } finally {
    isLoadingDiagrams.value = false
  }
}

async function selectDiagram(diagramId: string) {
  if (!diagramId) return
  apiError.value = null
  try {
    // @ts-ignore - API provided at runtime
    const d = await mermaidDiagramToolsProjectDiagramAPI.get(props.projectId, diagramId)
    const normalized: Diagram = {
      id: String(diagramId),
      title: d?.title ?? 'Untitled',
      syntax: d?.syntax ?? '',
      projectID: props.projectId,
      desc: d?.desc,
      latestUpdated: d?.latestUpdated,
    }
    selectedDiagram.value = normalized

    // Put into editor and preview
    source.value = normalized.syntax || ''
    if (editor) editor.setValue(source.value)
    await renderDiagram()
  } catch (e: any) {
    apiError.value = e?.message || 'Failed to load diagram.'
  }
}

async function saveDiagram() {
  if (!selectedDiagram.value) return
  isSavingDiagram.value = true
  apiError.value = null
  try {
    // @ts-ignore - API provided at runtime
    await mermaidDiagramToolsProjectDiagramAPI.save(props.projectId, selectedDiagram.value.id, source.value)
    // Optionally refresh metadata such as latestUpdated
    await selectDiagram(selectedDiagram.value.id)
  } catch (e: any) {
    apiError.value = e?.message || 'Failed to save diagram.'
  } finally {
    isSavingDiagram.value = false
  }
}

async function saveProjectMeta() {
  if (!project.value) return
  isProjectSaving.value = true
  apiError.value = null
  try {
    // @ts-ignore - API provided at runtime
    await mermaidDiagramToolsProjectAPI.update(props.projectId, { title: project.value.title, description: project.value.desc })
    isProjectEditing.value = false
  } catch (e: any) {
    apiError.value = e?.message || 'Failed to save project.'
  } finally {
    isProjectSaving.value = false
  }
}

// ---------------------------------
// Lifecycle
// ---------------------------------
onMounted(async () => {
  createEditor()
  renderDiagram()
  window.addEventListener('keydown', onKeydown)

  // Initial fetches
  await loadProject()
  await loadDiagrams()
})

onBeforeUnmount(() => {
  window.removeEventListener('keydown', onKeydown)
  if (editor) {
    editor.dispose()
    editor = null
  }
  disposePanzoom(pz)
  disposePanzoom(fsPz)
})

// If user changes selected diagram from the combobox, fetch & update editor
const selectedId = ref<string>('')
watch(selectedId, async (id) => {
  if (id) await selectDiagram(id)
})
</script>

<template>
  <div class="min-h-screen bg-white p-5 pt-0">
    <div class="mx-auto">
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent class="pt-5">
          <!-- Title Section -->
          <div class="flex items-start justify-between">
            <div class="max-w-3xl">
              <template v-if="!isProjectEditing">
                <h1 class="text-4xl font-bold mb-2 text-slate-900">
                  {{ project?.title || (isProjectLoading ? 'Loading…' : 'Untitled Project') }}
                </h1>
                <p class="leading-relaxed text-slate-600 whitespace-pre-line">
                  {{ project?.desc || '' }}
                </p>
              </template>
              <template v-else>
                <input
                    v-model="project!.title"
                    class="w-full text-3xl font-bold mb-2 text-slate-900 border-b outline-none"
                    placeholder="Project title"
                />
                <textarea
                    v-model="project!.desc"
                    rows="3"
                    class="w-full text-slate-700 border rounded-md p-2"
                    placeholder="Project description"
                />
              </template>
              <p v-if="apiError" class="text-sm text-red-600 mt-2">{{ apiError }}</p>
            </div>
            <div class="flex gap-2">
              <Button v-if="!isProjectEditing" class="h-9 px-4 font-normal" size="sm" type="button" variant="outline"
                      @click="isProjectEditing = true">
                Edit Project
              </Button>
              <Button v-else class="h-9 px-4 font-normal" size="sm" type="button" :disabled="isProjectSaving"
                      @click="saveProjectMeta">
                <Save class="mr-1" />
                {{ isProjectSaving ? 'Saving…' : 'Save Project' }}
              </Button>
            </div>
          </div>

          <!-- Toolbar Row -->
          <div class="flex items-center gap-2 mt-3 justify-between">
            <div>
              <!-- Combobox of Diagrams (populated from API) -->
              <Combobox v-model="selectedId" by="label">
                <ComboboxAnchor as-child>
                  <ComboboxTrigger as-child>
                    <Button class="justify-between min-w-56" variant="outline">
                      {{
                        (diagrams.find(d => d.id === selectedId)?.title) || (isLoadingDiagrams ? 'Loading…' : 'Select Diagram')
                      }}
                      <ChevronsUpDown class="ml-2 h-4 w-3 shrink-0 opacity-50" />
                    </Button>
                  </ComboboxTrigger>
                </ComboboxAnchor>

                <ComboboxList>
                  <div class="relative w-full max-w-sm items-center">
                    <ComboboxInput class="pl-9 focus-visible:ring-0 border-0 border-b rounded-none h-10" placeholder="Search diagram…" />
                    <span class="absolute start-0 inset-y-0 flex items-center justify-center px-3">
                      <Search class="size-4 text-muted-foreground" />
                    </span>
                  </div>

                  <ComboboxEmpty>
                    No diagrams found.
                  </ComboboxEmpty>

                  <ComboboxGroup>
                    <ComboboxItem v-for="d in diagrams" :key="d.id" :value="d.id">
                      {{ d.title }}
                      <ComboboxItemIndicator>
                        <Check :class="cn('ml-auto h-4 w-4')" />
                      </ComboboxItemIndicator>
                    </ComboboxItem>
                  </ComboboxGroup>
                </ComboboxList>
              </Combobox>
            </div>

            <div class="flex items-center gap-2 mt-3 justify-end">
              <Button class="h-9 px-4 font-normal" size="sm" type="button" variant="outline" @click="emit('editCollection')">
                <Share2 />
                Share Diagram
              </Button>
              <Button size="sm" @click="emit('createNewEndpoint')" class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal">
                <FilePlus2 />
                New Diagram
              </Button>
              <Button class="h-9 px-4 font-normal bg-primary text-white" size="sm" type="button" variant="outline"
                      :disabled="!selectedDiagram || isSavingDiagram" @click="saveDiagram">
                <Save />
                {{ isSavingDiagram ? 'Saving…' : 'Save Diagram' }}
              </Button>
            </div>
          </div>

          <div class="h-[90vh] mt-5 w-full flex gap-4">
            <!-- Editor Panel (Monaco) -->
            <div class="w-2/5 h-full flex flex-col rounded-lg border bg-white">
              <div class="flex items-center justify-between p-3 border-b">
                <h2 class="text-sm font-semibold">Mermaid Editor</h2>
                <div class="text-xs text-muted-foreground">Monaco • Live validation</div>
              </div>

              <div ref="editorEl" class="flex-1 min-h-0" />

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
                  <div class="w-px h-5 bg-gray-200 mx-1" />
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="downloadSVG">SVG</button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="downloadPNG">PNG</button>
                  <div class="w-px h-5 bg-gray-200 mx-1" />
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="openFullScreen">Full
                    screen
                  </button>
                </div>
              </div>

              <div class="relative flex-1 overflow-hidden">
                <div ref="previewRef" :class="{ 'opacity-50': isRendering }" class="absolute inset-0 select-none" v-html="svgMarkup" />
                <div v-if="!svgMarkup && !error" class="absolute inset-0 grid place-items-center text-sm text-gray-400">
                  Rendered diagram will appear here
                </div>
              </div>
            </div>

            <!-- Full Screen Overlay -->
            <div v-if="isFullScreen" class="fixed inset-0 z-50 bg-white/95 backdrop-blur-sm">
              <div class="flex items-center justify-between p-3 border-b bg-white/70">
                <div class="text-sm font-medium">Full Screen Preview</div>
                <div class="flex items-center gap-2">
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="zoomOut('fs')">−</button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="resetView('fs')">Reset</button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="zoomIn('fs')">＋</button>
                  <div class="w-px h-5 bg-gray-200 mx-1" />
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="downloadSVG">SVG</button>
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="downloadPNG">PNG</button>
                  <div class="w-px h-5 bg-gray-200 mx-1" />
                  <button class="px-2 py-1 rounded-md border text-xs" type="button" @click="closeFullScreen">Close (Esc)</button>
                </div>
              </div>
              <div class="relative h-[calc(100vh-48px)]">
                <div ref="fsPreviewRef" class="absolute inset-0 select-none" v-html="svgMarkup" />
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  </div>
</template>

<style scoped>
:deep(.monaco-scrollable-element)::-webkit-scrollbar { width: 8px; height: 8px; }
:deep(.monaco-scrollable-element)::-webkit-scrollbar-thumb { background: #d1d5db; border-radius: 6px; }
:deep(.monaco-scrollable-element)::-webkit-scrollbar-track { background: transparent; }
</style>
