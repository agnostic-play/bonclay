<script setup lang="ts">
import { ref, computed } from 'vue'
import { Badge } from '@/components/ui/badge'
import { Accordion } from '@/components/ui/accordion'
import EndpointAccordionItem from './EndpointAccordionItem.vue'
import type { EndpointByCategory } from '@/types/api.types'

interface Props {
  categories: EndpointByCategory | null
  baseUrl: string
  isProxyEnable?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  scenariosUpdated: []
}>()

const safeCategories = computed(() => props.categories ?? ({} as EndpointByCategory))

// Each category tracks its own open item independently
const openItems = ref<Record<string, string | null>>({})

const isAccordionOpen = (categoryId: string, itemValue: string): boolean =>
  openItems.value[categoryId] === itemValue

const handleAccordionChange = (categoryId: string, value: string | null) => {
  openItems.value[categoryId] = value ?? null
}
</script>

<template>
  <div class="p-8 pt-0 px-5 pb-0 space-y-7 mt-7">
    <div v-for="(category, index) in safeCategories" :key="index" class="space-y-1">
      <div class="flex items-center gap-3 mb-2 mt-4">
        <h2 class="text-2xl font-semibold text-slate-800">{{ index }}</h2>
        <Badge variant="outline" class="text-xs">{{ category.length }} endpoints</Badge>
      </div>

      <Accordion
        type="single"
        class="w-full space-y-2"
        collapsible
        :model-value="openItems[index] ?? undefined"
        @update:model-value="(val) => handleAccordionChange(index, (val as string) ?? null)"
      >
        <EndpointAccordionItem
          v-for="endpoint in category"
          :key="endpoint.id"
          :endpoint="endpoint"
          :category-id="index"
          :base-url="baseUrl"
          :is-proxy-enable="isProxyEnable"
          :is-open="isAccordionOpen(index, endpoint.id)"
          @scenarios-updated="emit('scenariosUpdated')"
        />
      </Accordion>
    </div>
  </div>
</template>
