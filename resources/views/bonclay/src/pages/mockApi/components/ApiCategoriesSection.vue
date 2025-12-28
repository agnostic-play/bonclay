<script setup lang="ts">
import { ref, computed, toRef } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from '@/components/ui/accordion'
import EndpointAccordionItem from './EndpointAccordionItem.vue'
import type { EndpointByCategory, ApiCategory, Scenario } from '../types/api.types'

interface Props {
  categories: EndpointByCategory | null
  scenarios: Scenario[]
}

const props = defineProps<Props>()

const safeCategories = computed(() => props.categories ?? ({} as EndpointByCategory))
const scenarios = toRef(props, 'scenarios')

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
    <div v-for="(category, index) in safeCategories" :key="index" class="space-y-1">
      <div class="flex items-center gap-3 mb-2 mt-4">
        <h2 class="text-2xl font-semibold text-slate-800">
          {{ index }}
        </h2>
        <Badge variant="outline" class="text-xs">
          {{ category.length }} endpoints
        </Badge>
      </div>

      <Accordion type="single" class="w-full space-y-2" collapsible :model-value="openAccordionItem"
        @update:model-value="handleAccordionChange">
        <EndpointAccordionItem v-for="endpoint in category" :key="`${index}-${endpoint.path}`" :endpoint="endpoint"
          :category-id="index" :scenarios="scenarios" :is-open="isAccordionOpen(`${index}-${endpoint.path}`)"
          @scenarios-updated="handleScenariosUpdate" />
      </Accordion>
    </div>
  </div>
</template>