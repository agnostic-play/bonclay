<script setup lang="ts">
import { Plus } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'

interface Props {
  totalEndpoints: number
  activeScenarioCount: number
  categoryCount: number
}

defineProps<Props>()

const emit = defineEmits<{
  editCollection: []
  createNewEndpoint: []
}>()
</script>

<template>
  <div>
    <!-- Title Section -->
    <div class="flex items-start justify-between">
      <div>
        <h1 class="text-4xl font-bold mb-2 text-slate-900">
          Rintis Collection API
        </h1>
        <p class="leading-relaxed max-w-4xl text-slate-600">
          The Rintis Collection API provides a comprehensive mock service for testing payment processing,
          transaction management, and authentication flows. Configure different response scenarios to simulate
          various conditions including success states, errors, and edge cases for robust application testing.
        </p>
      </div>
      <div class="flex gap-2">
        <Button
            variant="outline"
            size="sm"
            @click="emit('editCollection')"
            class="h-9 px-4 font-normal"
        >
          Edit Collection
        </Button>
        <Button
            size="sm"
            @click="emit('createNewEndpoint')"
            class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal"
        >
          <Plus class="w-4 h-4 mr-2" />
          New Endpoint
        </Button>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="flex flex-wrap gap-6 pt-4">
      <div class="flex items-center gap-2">
        <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
        <span class="text-sm text-slate-700">{{ totalEndpoints }} Endpoints</span>
      </div>
      <div class="flex items-center gap-2">
        <div class="w-2 h-2 bg-blue-500 rounded-full"></div>
        <span class="text-sm text-slate-700">{{ activeScenarioCount }} Active Scenarios</span>
      </div>
      <div class="flex items-center gap-2">
        <div class="w-2 h-2 bg-orange-500 rounded-full"></div>
        <span class="text-sm text-slate-700">{{ categoryCount }} Categories</span>
      </div>
    </div>
  </div>
</template>