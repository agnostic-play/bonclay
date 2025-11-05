<script setup lang="ts">
import { ref } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import EndpointAccordionItem from './EndpointAccordionItem.vue'
import type { ApiCategory, Scenario } from '../types/api.types'

interface Props {
  categories: ApiCategory[]
  scenarios: Scenario[]
}

defineProps<Props>()

const emit = defineEmits<{
  scenariosUpdated: [scenarios: Scenario[]]
}>()

const openAccordionItem = ref<string | null>(null)

const isAccordionOpen = (itemValue: string): boolean => {
  return openAccordionItem.value === itemValue
}

const handleAccordionChange = (value: string | null) => {
  openAccordionItem.value = value
}

const handleScenariosUpdate = (updatedScenarios: Scenario[]) => {
  emit('scenariosUpdated', updatedScenarios)
}
</script>

<template>
  <div class="p-8 pt-0 px-5 pb-0 space-y-7 mt-7">
    <div v-for="category in categories" :key="category.id" class="space-y-1">
      <div class="flex items-center gap-3 mb-2 mt-4">
        <h2 class="text-2xl font-semibold text-slate-800">
          {{ category.name }}
        </h2>
        <Badge variant="outline" class="text-xs">
          {{ category.endpoints.length }} endpoints
        </Badge>
      </div>

      <Accordion
          type="single"
          class="w-full space-y-2"
          collapsible
          :model-value="openAccordionItem"
          @update:model-value="handleAccordionChange"
      >
        <EndpointAccordionItem
            v-for="endpoint in category.endpoints"
            :key="`${category.id}-${endpoint.path}`"
            :endpoint="endpoint"
            :category-id="category.id"
            :scenarios="scenarios"
            :is-open="isAccordionOpen(`${category.id}-${endpoint.path}`)"
            @scenarios-updated="handleScenariosUpdate"
        />
      </Accordion>
    </div>
  </div>
</template>