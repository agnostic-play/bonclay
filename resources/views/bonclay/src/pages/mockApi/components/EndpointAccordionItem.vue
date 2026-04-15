<script setup lang="ts">
import { ref, computed } from 'vue'
import { Copy, Check, Plus, Clock } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import ScenariosTable from './ScenariosTable.vue'
import EditEndpointDialog from './EditEndpointDialog.vue'
import CreateScenarioDialog from './CreateScenarioDialog.vue'
import { useEndpointUtils } from '../composables/useEndpointUtils'
import type { CollectionEndpoint, ScenarioResponse } from '@/types/api.types'
import MonacoEditor from 'monaco-editor-vue3'

interface Props {
  endpoint: CollectionEndpoint
  categoryId: string
  isOpen: boolean
  baseUrl: string
  isProxyEnable?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  scenariosUpdated: []
}>()

const { getMethodColor, getMethodBgColor } = useEndpointUtils()

const copiedItems = ref<Record<string, boolean>>({})
const code = ref(props.endpoint.script ?? '')
const editDialogOpen = ref(false)
const scenarioDialogOpen = ref(false)

const endpointScenarios = computed<ScenarioResponse[]>(() => props.endpoint.scenario_response ?? [])
const displayPath = computed(() => props.endpoint.path.replace(/\[\^\/\]\+/g, '{query_params}'))
const fullUrl = computed(() => props.baseUrl + displayPath.value)

const copyToClipboard = async (text: string, itemId: string): Promise<void> => {
  try {
    await navigator.clipboard.writeText(text)
    copiedItems.value[itemId] = true
    setTimeout(() => { delete copiedItems.value[itemId] }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

const handleCopyClick = (event: Event, text: string, itemId: string): void => {
  event.stopPropagation()
  copyToClipboard(text, itemId)
}

const handleScenariosUpdate = () => emit('scenariosUpdated')
</script>

<template>
  <AccordionItem
    :value="endpoint.id"
    :class="[
      'rounded-sm border transition-all duration-200',
      isOpen ? 'bg-white border-gray-200 shadow-sm' : getMethodBgColor(endpoint.method),
    ]"
  >
    <AccordionTrigger class="no-underline hover:no-underline px-6 hover:bg-white/50">
      <div class="flex items-center w-full space-x-3">
        <Badge :class="`text-white font-bold text-center min-w-[80px] justify-center ${getMethodColor(endpoint.method)}`">
          {{ endpoint.method }}
        </Badge>

        <div class="flex items-center gap-2 text-slate-700 font-mono text-sm bg-slate-100 px-3 rounded-lg">
          <code>{{ displayPath }}</code>
          <Tooltip>
            <TooltipTrigger as-child>
              <Button
                variant="ghost"
                size="sm"
                class="h-6 w-6 p-0 hover:bg-slate-200"
                @click="(e: MouseEvent) => handleCopyClick(e, fullUrl, endpoint.id)"
              >
                <Check v-if="copiedItems[endpoint.id]" class="w-3 h-3 text-green-600" />
                <Copy v-else class="w-3 h-3 text-slate-600" />
              </Button>
            </TooltipTrigger>
            <TooltipContent><p>Copy full URL</p></TooltipContent>
          </Tooltip>
        </div>

        <span class="ml-auto font-medium text-gray-800">{{ endpoint.desc }}</span>
      </div>
    </AccordionTrigger>

    <AccordionContent class="px-6 pb-6 bg-white rounded-b-xl">
      <Tabs default-value="endpoint">
        <TabsList class="w-full flex justify-center items-center">
          <TabsTrigger class="px-10" value="endpoint">Endpoint</TabsTrigger>
          <TabsTrigger class="px-10" value="script">Script</TabsTrigger>
        </TabsList>

        <TabsContent value="endpoint" class="p-4">
          <div class="space-y-4">
            <div class="flex justify-between items-center">
              <h4 class="text-sm font-medium text-gray-800 flex items-center gap-3">
                <Clock class="w-4 h-4 text-gray-500" />
                Mock Scenarios
              </h4>
              <div class="flex gap-3">
                <Button variant="outline" size="sm" class="h-8 px-3 font-normal text-xs" @click="editDialogOpen = true">
                  Edit Endpoint
                </Button>
                <Button size="sm" class="h-8 px-3 bg-gray-900 hover:bg-gray-800 font-normal text-xs" @click="scenarioDialogOpen = true">
                  <Plus class="w-3.5 h-3.5 mr-1.5" />
                  New Scenario
                </Button>
              </div>
            </div>

            <div class="bg-gray-50/60 rounded-xl border border-gray-100 overflow-hidden">
              <ScenariosTable
                :scenarios="endpointScenarios"
                :active-scenarios="endpoint.active_scenario"
                :endpoint-id="endpoint.id"
                :is-proxy-enable="isProxyEnable"
                :endpoint-delay="endpoint.delay"
                @scenarios-updated="handleScenariosUpdate"
              />
            </div>
          </div>
        </TabsContent>

        <TabsContent value="script" class="p-4">
          <MonacoEditor
            v-model:value="code"
            language="javascript"
            theme="vs-dark"
            height="300px"
            class="border rounded-xl"
            :options="{ padding: { top: 16, bottom: 16 }, minimap: { enabled: false } }"
          />
          <div class="flex justify-end mt-3">
            <Button
              variant="outline"
              size="sm"
              class="h-9 px-4 font-normal bg-[#fc7305] text-white hover:bg-[#ffa55c] hover:text-white hover:cursor-pointer"
            >
              Save Script
            </Button>
          </div>
        </TabsContent>
      </Tabs>
    </AccordionContent>
  </AccordionItem>

  <EditEndpointDialog
    v-model:open="editDialogOpen"
    :endpoint-id="endpoint.id"
    :initial-path="endpoint.path"
    :initial-method="endpoint.method"
    :initial-category="endpoint.category"
    :initial-desc="endpoint.desc"
    :initial-delay="endpoint.delay"
    @updated="handleScenariosUpdate"
  />

  <CreateScenarioDialog
    v-model:open="scenarioDialogOpen"
    :endpoint-id="endpoint.id"
    @created="handleScenariosUpdate"
  />
</template>

<style scoped>
.accordion-trigger:hover,
.accordion-trigger:focus,
.accordion-trigger:active {
  text-decoration: none !important;
}
</style>
