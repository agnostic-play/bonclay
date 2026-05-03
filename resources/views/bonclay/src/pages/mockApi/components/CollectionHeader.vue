<script setup lang="ts">
import { ref } from 'vue'
import { Plus, Copy, Check, History, ArrowLeft, Braces } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { updateCollection } from '@/api'
import { useRoute, useRouter, RouterLink } from 'vue-router'

const route = useRoute()
const router = useRouter()

interface Props {
  collectionId: string
  collectionName: string
  collectionDesc: string
  totalEndpoints: number
  activeScenarioCount: number
  categoryCount: number
  forwardProxyUrl?: string
  isProxyEnable?: boolean
}

const props = defineProps<Props>()

const emit = defineEmits<{
  editCollection: []
  createNewEndpoint: []
  openVariables: []
  updated: []
}>()

const copiedForward = ref(false)
const togglingProxy = ref(false)

const copyForwardUrl = async () => {
  if (!props.forwardProxyUrl) return
  try {
    await navigator.clipboard.writeText(props.forwardProxyUrl)
    copiedForward.value = true
    setTimeout(() => { copiedForward.value = false }, 2000)
  } catch (err) {
    console.error('Failed to copy:', err)
  }
}

const toggleProxy = async () => {
  if (togglingProxy.value) return
  togglingProxy.value = true
  try {
    await updateCollection(props.collectionId, { is_proxy_enable: !props.isProxyEnable })
    emit('updated')
  } catch (err) {
    console.error('Failed to toggle proxy:', err)
  } finally {
    togglingProxy.value = false
  }
}
</script>

<template>
  <div>
    <!-- Back Navigation -->
    <div class="mb-4">
      <button
        type="button"
        class="inline-flex items-center gap-1.5 text-sm text-slate-500 hover:text-slate-800 transition-colors"
        @click="router.push({ name: 'MockApiTools-Index', params: { slug: route.params.slug } })"
      >
        <ArrowLeft class="w-4 h-4" />
        Back to Collections
      </button>
    </div>

    <!-- Title Section -->
    <div class="flex items-start justify-between">
      <div>
        <h1 class="text-4xl font-bold mb-2 text-slate-900">
          {{ collectionName || 'Collection' }}
        </h1>
        <p class="leading-relaxed max-w-4xl text-slate-600">
          {{ collectionDesc || 'No description provided.' }}
        </p>
      </div>
      <div class="flex gap-2">
        <Button variant="outline" size="sm" @click="emit('editCollection')" class="h-9 px-4 font-normal">
          Edit Collection
        </Button>
        <Button variant="outline" size="sm" @click="emit('openVariables')" class="h-9 px-4 font-normal">
          <Braces class="w-4 h-4 mr-2" />
          Variables
        </Button>
        <RouterLink :to="{ name: 'MockApiTools-CollectionHistory', params: { slug: route.params.slug, collectionSlug: route.params.collectionSlug } }">
          <Button variant="outline" size="sm" class="h-9 px-4 font-normal">
            <History class="w-4 h-4 mr-2" />
            History
          </Button>
        </RouterLink>
        <Button size="sm" @click="emit('createNewEndpoint')" class="h-9 px-4 bg-gray-900 hover:bg-gray-800 font-normal">
          <Plus class="w-4 h-4 mr-2" />
          New Endpoint
        </Button>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="flex flex-wrap items-center gap-6 pt-4">
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

      <!-- Forward Proxy -->
      <div v-if="props.forwardProxyUrl" class="flex items-center gap-2 ml-2 pl-2 border-l border-slate-200">
        <span class="text-sm text-slate-500">Forward proxy:</span>
        <code class="text-xs font-mono text-slate-700 max-w-[200px] truncate">{{ props.forwardProxyUrl }}</code>
        <Tooltip>
          <TooltipTrigger as-child>
            <Button variant="ghost" size="sm" class="h-6 w-6 p-0 hover:bg-slate-100" @click="copyForwardUrl">
              <Check v-if="copiedForward" class="w-3 h-3 text-green-600" />
              <Copy v-else class="w-3 h-3 text-slate-500" />
            </Button>
          </TooltipTrigger>
          <TooltipContent><p>Copy forward URL</p></TooltipContent>
        </Tooltip>
        <Button
          size="sm"
          :disabled="togglingProxy"
          :class="props.isProxyEnable
            ? 'bg-green-600 hover:bg-green-700 text-white'
            : 'bg-slate-200 hover:bg-slate-300 text-slate-600'"
          class="h-6 px-2 text-[11px] font-semibold transition-colors"
          @click="toggleProxy"
        >
          {{ props.isProxyEnable ? 'ON' : 'OFF' }}
        </Button>
      </div>
    </div>
  </div>
</template>
