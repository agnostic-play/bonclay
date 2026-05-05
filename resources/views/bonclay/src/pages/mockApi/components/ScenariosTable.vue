<script setup lang="ts">
import { ref } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogDescription,
} from '@/components/ui/dialog'
import ScenarioDetailDialog from './ScenarioDetailDialog.vue'
import { useScenarioUtils } from '@/pages/mockApi/composables/useScenarioUtils'
import type { ScenarioResponse } from '@/types/api.types'
import { setActiveScenario, removeActiveScenario, updateEndpoint } from '@/api'

const PROXY_ROW_ID = '__proxy__'

interface Props {
  scenarios: ScenarioResponse[]
  activeScenarios: string
  endpointId: string
  isProxyEnable?: boolean
  endpointDelay?: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  scenariosUpdated: []
}>()

const { getStatusColor, formatDelay } = useScenarioUtils()
const isToggling = ref(false)

const selectedScenario = ref<ScenarioResponse | null>(null)
const detailOpen = ref(false)

const proxyDelayOpen = ref(false)
const proxyDelayInput = ref<number | null>(null)
const savingDelay = ref(false)

const openProxyDelay = () => {
  proxyDelayInput.value = props.endpointDelay ?? null
  proxyDelayOpen.value = true
}

const saveProxyDelay = async () => {
  savingDelay.value = true
  try {
    await updateEndpoint(props.endpointId, { delay: proxyDelayInput.value ?? undefined })
    emit('scenariosUpdated')
    proxyDelayOpen.value = false
  } catch (err) {
    console.error('Failed to save delay:', err)
  } finally {
    savingDelay.value = false
  }
}

const viewScenario = (scenario: ScenarioResponse) => {
  selectedScenario.value = scenario
  detailOpen.value = true
}

const toggleScenario = async (scenarioId: string): Promise<void> => {
  if (isToggling.value) return
  isToggling.value = true
  try {
    if (scenarioId === PROXY_ROW_ID) {
      await removeActiveScenario(props.endpointId)
    } else {
      await setActiveScenario({ endpoint_id: props.endpointId, scenario_id: scenarioId })
    }
    emit('scenariosUpdated')
  } catch (err) {
    console.error('Failed to toggle scenario:', err)
  } finally {
    isToggling.value = false
  }
}
</script>

<template>
  <Table>
    <TableHeader>
      <TableRow class="bg-slate-50">
        <TableHead class="w-[10%] text-slate-900 text-center font-semibold">HTTP Status</TableHead>
        <TableHead class="w-[55%] text-slate-900 font-semibold">Scenario</TableHead>
        <TableHead class="w-[15%] text-slate-900 text-center font-semibold">Delay</TableHead>
        <TableHead class="w-[10%] text-slate-900 text-center font-semibold">Active</TableHead>
        <TableHead class="w-[10%] text-slate-900 text-center font-semibold">Actions</TableHead>
      </TableRow>
    </TableHeader>
    <TableBody>
            <!-- Forwarded server row — only when proxy is enabled -->
      <TableRow v-if="isProxyEnable" :class="!activeScenarios ? 'bg-blue-50/50' : ''">
        <TableCell class="text-center">
          <Badge variant="secondary" class="bg-slate-100 text-slate-500">—</Badge>
        </TableCell>
        <TableCell>
          <div class="font-medium text-slate-800 italic">Use response from the forwarded server.</div>
        </TableCell>
        <TableCell class="text-center">
          <span class="text-sm">{{ formatDelay(endpointDelay) }}</span>
        </TableCell>
        <TableCell class="text-center">
          <Checkbox
            :model-value="!activeScenarios"
            :disabled="isToggling"
            @update:model-value="() => toggleScenario(PROXY_ROW_ID)"
          />
        </TableCell>
        <TableCell class="text-center">
          <Button
            variant="ghost"
            size="sm"
            class="h-8 px-3 font-normal text-xs text-gray-600 hover:text-gray-800"
            @click="openProxyDelay"
          >
            View
          </Button>
        </TableCell>
      </TableRow>
      
      <TableRow
        v-for="scenario in scenarios"
        :key="scenario.id"
        :class="scenario.id === activeScenarios ? 'bg-blue-50/50' : ''"
      >
        <TableCell class="text-center">
          <Badge variant="secondary" :class="getStatusColor(scenario.status_header)">
            {{ scenario.status_header }}
          </Badge>
        </TableCell>
        <TableCell>
          <div class="font-medium text-slate-800">{{ scenario.desc }}</div>
        </TableCell>
        <TableCell class="text-center">
          <span class="text-sm">{{ formatDelay(scenario.delay) }}</span>
        </TableCell>
        <TableCell class="text-center">
          <Checkbox
            :model-value="scenario.id === activeScenarios"
            :disabled="isToggling"
            @update:model-value="() => toggleScenario(scenario.id)"
          />
        </TableCell>
        <TableCell class="text-center">
          <Button
            variant="ghost"
            size="sm"
            @click="viewScenario(scenario)"
            class="h-8 px-3 font-normal text-xs text-gray-600 hover:text-gray-800"
          >
            View
          </Button>
        </TableCell>
      </TableRow>



      <TableRow v-if="scenarios.length === 0 && !isProxyEnable">
        <TableCell colspan="5" class="text-center text-slate-500 py-8">
          No scenarios configured. Click "New Scenario" to create one.
        </TableCell>
      </TableRow>
    </TableBody>
  </Table>

  <ScenarioDetailDialog
    v-model:open="detailOpen"
    :scenario="selectedScenario"
    @updated="emit('scenariosUpdated')"
    @deleted="emit('scenariosUpdated')"
  />

  <Dialog v-model:open="proxyDelayOpen">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>Forwarded Server Delay</DialogTitle>
        <DialogDescription>Set a delay (in seconds) applied when this endpoint forwards to the proxy server.</DialogDescription>
      </DialogHeader>
      <div class="py-2 space-y-2">
        <Label for="proxy-delay-input">Delay (seconds)</Label>
        <input
          id="proxy-delay-input"
          v-model.number="proxyDelayInput"
          type="number"
          min="0"
          placeholder="0"
          :disabled="savingDelay"
          class="flex h-10 w-full rounded-lg border border-gray-200 bg-transparent px-3 py-2 text-sm shadow-sm focus:border-gray-300 focus:ring-1 focus:ring-gray-300 focus:outline-none disabled:cursor-not-allowed disabled:opacity-50"
        />
      </div>
      <DialogFooter>
        <Button type="button" variant="outline" size="sm" class="h-9 font-normal" :disabled="savingDelay" @click="proxyDelayOpen = false">
          Cancel
        </Button>
        <Button type="button" size="sm" class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal" :disabled="savingDelay" @click="saveProxyDelay">
          {{ savingDelay ? 'Saving...' : 'Update' }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
