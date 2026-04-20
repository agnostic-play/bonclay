<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { Card, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Tabs, TabsList, TabsTrigger, TabsContent } from '@/components/ui/tabs'
import { getHistory, clearHistory, type RequestHistory } from '@/api/history'
import { Trash2, RefreshCw, ChevronDown, ChevronRight, Copy, Check } from 'lucide-vue-next'

const route = useRoute()

const entries = ref<RequestHistory[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const expandedId = ref<string | null>(null)
const copiedCurls = ref<Record<string, boolean>>({})

const collectionSlug = computed(() => String(route.params.collectionSlug))

const fetchHistory = async () => {
  if (!collectionSlug.value) return
  loading.value = true
  error.value = null
  try {
    entries.value = await getHistory(collectionSlug.value)
  } catch (err) {
    error.value = 'Failed to load history'
    console.error(err)
  } finally {
    loading.value = false
  }
}

const handleClear = async () => {
  if (!confirm('Clear all history for this collection?')) return
  try {
    await clearHistory(collectionSlug.value)
    entries.value = []
  } catch (err) {
    console.error(err)
  }
}

const toggleExpand = (id: string) => {
  expandedId.value = expandedId.value === id ? null : id
}

watch(collectionSlug, fetchHistory, { immediate: true })

const methodColor = (method: string): string => {
  switch (method.toUpperCase()) {
    case 'GET': return 'bg-blue-100 text-blue-700'
    case 'POST': return 'bg-green-100 text-green-700'
    case 'PUT': return 'bg-yellow-100 text-yellow-700'
    case 'PATCH': return 'bg-orange-100 text-orange-700'
    case 'DELETE': return 'bg-red-100 text-red-700'
    default: return 'bg-gray-100 text-gray-700'
  }
}

const statusColor = (code: number): string => {
  if (code >= 200 && code < 300) return 'bg-green-100 text-green-700'
  if (code >= 300 && code < 400) return 'bg-blue-100 text-blue-700'
  if (code >= 400 && code < 500) return 'bg-yellow-100 text-yellow-700'
  return 'bg-red-100 text-red-700'
}

const formatTime = (ts: string): string => {
  const d = new Date(ts)
  return d.toLocaleString()
}

const displayPath = (path: string): string => path.replace(/\[\^\/\]\+/g, '{query_params}')

const formatBody = (body: string): string => {
  if (!body) return '(empty)'
  try {
    return JSON.stringify(JSON.parse(body), null, 2)
  } catch {
    return body
  }
}

const buildCurl = (entry: RequestHistory): string => {
  const url = `${window.location.origin}/mock/${entry.collection_slug}${entry.path}`
  const parts = [`curl -X ${entry.method} '${url}'`]
  for (const [k, v] of Object.entries(entry.request_headers ?? {})) {
    parts.push(`  -H '${k}: ${v}'`)
  }
  if (entry.request_body) {
    const escaped = entry.request_body.replace(/'/g, "'\\''")
    parts.push(`  -d '${escaped}'`)
  }
  return parts.join(' \\\n')
}

const copyCurl = async (entry: RequestHistory) => {
  try {
    await navigator.clipboard.writeText(buildCurl(entry))
    copiedCurls.value[entry.id] = true
    setTimeout(() => { delete copiedCurls.value[entry.id] }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}
</script>

<template>
  <div class="min-h-screen bg-white p-5 pt-0">
    <div class="mx-auto">
      <Card class="overflow-hidden border shadow-sm bg-white">
        <CardContent class="pt-5">

          <div class="flex items-center justify-between mb-4 px-4">
            <div>
              <h2 class="text-lg font-semibold">Request History</h2>
              <p class="text-sm text-muted-foreground">Last 100 requests for <span class="font-mono">{{ collectionSlug }}</span></p>
            </div>
            <div class="flex gap-2">
              <Button variant="outline" size="sm" :disabled="loading" @click="fetchHistory">
                <RefreshCw class="w-4 h-4 mr-1" :class="{ 'animate-spin': loading }" />
                Refresh
              </Button>
              <Button variant="destructive" size="sm" :disabled="entries.length === 0" @click="handleClear">
                <Trash2 class="w-4 h-4 mr-1" />
                Clear
              </Button>
            </div>
          </div>

          <div v-if="loading" class="p-8 text-center text-muted-foreground">
            Loading history...
          </div>

          <div v-else-if="error" class="p-8 text-center text-red-600">
            {{ error }}
          </div>

          <div v-else-if="entries.length === 0" class="p-8 text-center text-muted-foreground">
            No history yet. Make a request to see it here.
          </div>

          <div v-else class="divide-y border rounded-lg mx-4">
            <div
              v-for="entry in entries"
              :key="entry.id"
              class="cursor-pointer"
              @click="toggleExpand(entry.id)"
            >
              <!-- Summary row -->
              <div class="flex items-center gap-3 px-4 py-3 hover:bg-gray-50 transition-colors">
                <component :is="expandedId === entry.id ? ChevronDown : ChevronRight" class="w-4 h-4 text-muted-foreground flex-shrink-0" />

                <span :class="['text-xs font-bold px-2 py-0.5 rounded font-mono', methodColor(entry.method)]">
                  {{ entry.method }}
                </span>

                <span class="font-mono text-sm flex-1 truncate">{{ displayPath(entry.path) }}</span>

                <span :class="['text-xs font-bold px-2 py-0.5 rounded', statusColor(entry.status_code)]">
                  {{ entry.status_code }}
                </span>

                <span v-if="entry.duration_ms" class="text-xs text-muted-foreground font-mono">
                  {{ entry.duration_ms }}ms
                </span>

                <Badge v-if="entry.proxy_used" variant="outline" class="text-xs">proxy</Badge>

                <span class="text-xs text-muted-foreground whitespace-nowrap">{{ formatTime(entry.timestamp) }}</span>
              </div>

              <!-- Expanded detail -->
              <div v-if="expandedId === entry.id" class="border-t px-4 py-4 bg-gray-50/50" @click.stop>
                <Tabs default-value="request">
                  <TabsList class="w-full flex justify-center items-center">
                    <TabsTrigger class="px-10" value="request">Request</TabsTrigger>
                    <TabsTrigger class="px-10" value="response">
                      Response
                      <span :class="['ml-2 text-xs font-bold px-1.5 py-0.5 rounded', statusColor(entry.status_code)]">
                        {{ entry.status_code }}
                      </span>
                      <span v-if="entry.duration_ms" class="ml-1 text-xs text-muted-foreground font-mono">
                        {{ entry.duration_ms }}ms
                      </span>
                    </TabsTrigger>
                  </TabsList>

                  <!-- REQUEST TAB -->
                  <TabsContent value="request" class="p-4">
                    <div class="flex justify-end mb-3">
                      <Button
                        variant="outline"
                        size="sm"
                        class="h-7 px-3 text-xs gap-1.5"
                        @click.stop="copyCurl(entry)"
                      >
                        <Check v-if="copiedCurls[entry.id]" class="w-3 h-3 text-green-600" />
                        <Copy v-else class="w-3 h-3" />
                        {{ copiedCurls[entry.id] ? 'Copied!' : 'Copy cURL' }}
                      </Button>
                    </div>
                    <Accordion type="multiple" :default-value="['req-headers', 'req-body']" class="border rounded-lg overflow-hidden">
                      <AccordionItem value="req-headers" class="border-b last:border-b-0">
                        <AccordionTrigger class="px-4 py-2.5 text-xs font-medium hover:no-underline [&>svg]:h-3.5 [&>svg]:w-3.5">
                          Headers
                        </AccordionTrigger>
                        <AccordionContent class="px-4 pb-3">
                          <pre class="bg-gray-50 border rounded p-2 text-xs overflow-auto max-h-60 font-mono whitespace-pre-wrap">{{ JSON.stringify(entry.request_headers, null, 2) }}</pre>
                        </AccordionContent>
                      </AccordionItem>
                      <AccordionItem value="req-body" class="border-b-0">
                        <AccordionTrigger class="px-4 py-2.5 text-xs font-medium hover:no-underline [&>svg]:h-3.5 [&>svg]:w-3.5">
                          Body
                        </AccordionTrigger>
                        <AccordionContent class="px-4 pb-3">
                          <pre class="bg-gray-50 border rounded p-2 text-xs overflow-auto max-h-60 font-mono whitespace-pre-wrap">{{ formatBody(entry.request_body) }}</pre>
                        </AccordionContent>
                      </AccordionItem>
                    </Accordion>
                  </TabsContent>

                  <!-- RESPONSE TAB -->
                  <TabsContent value="response" class="p-4">
                    <div class="flex items-center gap-2 mb-3">
                      <Badge v-if="entry.proxy_used" variant="outline" class="text-xs">forwarded via proxy</Badge>
                    </div>
                    <Accordion type="multiple" :default-value="['res-headers', 'res-body']" class="border rounded-lg overflow-hidden">
                      <AccordionItem value="res-headers" class="border-b last:border-b-0">
                        <AccordionTrigger class="px-4 py-2.5 text-xs font-medium hover:no-underline [&>svg]:h-3.5 [&>svg]:w-3.5">
                          Headers
                        </AccordionTrigger>
                        <AccordionContent class="px-4 pb-3">
                          <pre class="bg-gray-50 border rounded p-2 text-xs overflow-auto max-h-60 font-mono whitespace-pre-wrap">{{ JSON.stringify(entry.response_headers ?? {}, null, 2) }}</pre>
                        </AccordionContent>
                      </AccordionItem>
                      <AccordionItem value="res-body" class="border-b-0">
                        <AccordionTrigger class="px-4 py-2.5 text-xs font-medium hover:no-underline [&>svg]:h-3.5 [&>svg]:w-3.5">
                          Body
                        </AccordionTrigger>
                        <AccordionContent class="px-4 pb-3">
                          <pre class="bg-gray-50 border rounded p-2 text-xs overflow-auto max-h-60 font-mono whitespace-pre-wrap">{{ formatBody(entry.response_body) }}</pre>
                        </AccordionContent>
                      </AccordionItem>
                    </Accordion>
                  </TabsContent>
                </Tabs>

                <div v-if="entry.scenario_id" class="text-xs text-muted-foreground px-4 pt-1">
                  Scenario ID: <span class="font-mono">{{ entry.scenario_id }}</span>
                </div>
              </div>

            </div>
          </div>

        </CardContent>
      </Card>
    </div>
  </div>
</template>
