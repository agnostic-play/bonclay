<script setup lang="ts">
import { ref, computed } from 'vue'
import { Copy, Check, Plus, Clock } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import ScenariosTable from './ScenariosTable.vue'
import { useEndpointUtils } from '../composables/useEndpointUtils'
import type { Endpoint, Scenario } from '../types/api.types'

interface Props {
  endpoint: Endpoint
  categoryId: string
  scenarios: Scenario[]
  isOpen: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  scenariosUpdated: [scenarios: Scenario[]]
}>()

const copiedItems = ref<Record<string, boolean>>({})
const baseUrl = 'https://mock.com/test123123123123213'

const { getMethodColor, getMethodBgColor } = useEndpointUtils()

const endpointScenarios = computed(() =>
    props.scenarios.filter(s => s.endpointPath === props.endpoint.path)
)

const copyToClipboard = async (text: string, itemId: string = 'default'): Promise<void> => {
  try {
    await navigator.clipboard.writeText(text)
    copiedItems.value[itemId] = true
    setTimeout(() => {
      delete copiedItems.value[itemId]
    }, 2000)
  } catch (err) {
    console.error('Failed to copy text: ', err)
  }
}

const handleCopyClick = (event: Event, text: string, itemId: string): void => {
  event.stopPropagation()
  copyToClipboard(text, itemId)
}

const editEndpoint = (endpointPath: string): void => {
  console.log('Editing endpoint:', endpointPath)
}

const createNewScenario = (endpointPath: string): void => {
  console.log('Creating new scenario for:', endpointPath)
}

const handleScenariosUpdate = (updatedScenarios: Scenario[]) => {
  emit('scenariosUpdated', updatedScenarios)
}
</script>

<template>
  <AccordionItem
      :value="`${categoryId}-${endpoint.path}`"
      :class="[
      'rounded-sm border transition-all duration-200',
      isOpen
        ? 'bg-white border-gray-200 shadow-sm'
        : getMethodBgColor(endpoint.method)
    ]"
  >
    <AccordionTrigger class="no-underline hover:no-underline px-6 hover:bg-white/50">
      <div class="flex items-center w-full space-x-3">
        <!-- HTTP Method Badge -->
        <Badge
            :class="`text-white font-bold text-center min-w-[80px] justify-center ${getMethodColor(endpoint.method)}`"
        >
          {{ endpoint.method }}
        </Badge>

        <div class="flex items-center gap-2 text-slate-700 font-mono text-sm bg-slate-100 px-3 rounded-lg">
          <code>{{ endpoint.path }}</code>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                  variant="ghost"
                  size="sm"
                  @click="(event) => handleCopyClick(event, baseUrl + endpoint.path, baseUrl + endpoint.path)"
                  class="h-6 w-6 p-0 hover:bg-slate-200"
              >
                <Check v-if="copiedItems[baseUrl + endpoint.path]" class="w-3 h-3 text-green-600" />
                <Copy v-else class="w-3 h-3 text-slate-600" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Copy full URL</p>
            </TooltipContent>
          </Tooltip>
        </div>

        <!-- Endpoint Description -->
        <span class="ml-auto font-medium text-gray-800">
          {{ endpoint.description }}
        </span>
      </div>
    </AccordionTrigger>

    <AccordionContent class="px-6 pb-6 bg-white rounded-b-xl">
      <!-- Scenarios Section -->
      <div class="space-y-4">
        <div class="flex justify-between items-center">
          <h4 class="text-sm font-medium text-gray-800 flex items-center gap-3">
            <Clock class="w-4 h-4 text-gray-500" />
            Mock Scenarios
          </h4>
          <div class="flex gap-3">
            <Button
                variant="outline"
                size="sm"
                @click="editEndpoint(endpoint.path)"
                class="h-8 px-3 font-normal text-xs"
            >
              Edit Endpoint
            </Button>
            <Button
                size="sm"
                @click="createNewScenario(endpoint.path)"
                class="h-8 px-3 bg-gray-900 hover:bg-gray-800 font-normal text-xs"
            >
              <Plus class="w-3.5 h-3.5 mr-1.5" />
              New Scenario
            </Button>
          </div>
        </div>

        <div class="bg-gray-50/60 rounded-xl border border-gray-100 overflow-hidden">
          <ScenariosTable
              :scenarios="endpointScenarios"
              :endpoint-path="endpoint.path"
              :all-scenarios="scenarios"
              @scenarios-updated="handleScenariosUpdate"
          />
        </div>
      </div>
    </AccordionContent>
  </AccordionItem>
</template>

<style scoped>
.accordion-trigger {
  text-decoration: none !important;
}

.accordion-trigger:hover,
.accordion-trigger:focus,
.accordion-trigger:active {
  text-decoration: none !important;
}
</style>