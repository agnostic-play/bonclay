<script setup lang="ts">
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Checkbox } from '@/components/ui/checkbox'
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from '@/components/ui/table'
import { useScenarioUtils } from '@/pages/mockApi/composables/useScenarioUtils'
import type { Scenario } from '@/types/api.types'

interface Props {
  scenarios: Scenario[]
  endpointPath: string
  allScenarios: Scenario[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  scenariosUpdated: [scenarios: Scenario[]]
}>()

const { getStatusColor, formatDelay } = useScenarioUtils()

const toggleScenario = (scenarioId: string): void => {
  const updatedScenarios = props.allScenarios.map(s => ({
    ...s,
    isActive: s.endpointPath === props.endpointPath ? s.id === scenarioId : s.isActive
  }))
  emit('scenariosUpdated', updatedScenarios)
}

const viewScenario = (scenarioId: string): void => {
  console.log('Viewing scenario:', scenarioId)
}
</script>

<template>
  <Table>
    <TableHeader>
      <TableRow class="bg-slate-50">
        <TableHead class="w-[10%] text-slate-900 text-center font-semibold">
          HTTP Status
        </TableHead>
        <TableHead class="w-[55%] text-slate-900 font-semibold">
          Scenario
        </TableHead>
        <TableHead class="w-[15%] text-slate-900 text-center font-semibold">
          Delay
        </TableHead>
        <TableHead class="w-[10%] text-slate-900 text-center font-semibold">
          Active
        </TableHead>
        <TableHead class="w-[10%] text-slate-900 text-center font-semibold">
          Actions
        </TableHead>
      </TableRow>
    </TableHeader>
    <TableBody>
      <TableRow
          v-for="scenario in scenarios"
          :key="scenario.id"
          :class="scenario.isActive ? 'bg-blue-50/50' : ''"
      >
        <TableCell class="text-center">
          <Badge
              variant="secondary"
              :class="getStatusColor(scenario.httpStatus)"
          >
            {{ scenario.httpStatus }}
          </Badge>
        </TableCell>
        <TableCell>
          <div>
            <div class="font-medium text-slate-800">{{ scenario.name }}</div>
            <div class="text-xs text-slate-500 mt-0.5">{{ scenario.description }}</div>
          </div>
        </TableCell>
        <TableCell class="text-center">
          <div class="flex items-center justify-center gap-1">
            <span class="text-sm">{{ formatDelay(scenario.responseDelay) }}</span>
          </div>
        </TableCell>
        <TableCell class="text-center">
          <Checkbox
              :checked="scenario.isActive"
              @update:checked="() => toggleScenario(scenario.id)"
          />
        </TableCell>
        <TableCell class="text-center">
          <Button
              variant="ghost"
              size="sm"
              @click="viewScenario(scenario.id)"
              class="h-8 px-3 font-normal text-xs text-gray-600 hover:text-gray-800"
          >
            View
          </Button>
        </TableCell>
      </TableRow>
      <TableRow v-if="scenarios.length === 0">
        <TableCell colspan="5" class="text-center text-slate-500 py-8">
          No scenarios configured. Click "New Scenario" to create one.
        </TableCell>
      </TableRow>
    </TableBody>
  </Table>
</template>