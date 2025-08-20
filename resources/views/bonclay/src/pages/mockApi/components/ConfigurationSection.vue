<script setup lang="ts">
import { ref } from 'vue'
import { Copy, ExternalLink, Check } from 'lucide-vue-next'
import { Button } from '@/components/ui/button'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'

interface Props {
  baseUrl: string
  documentationUrl: string
}

defineProps<Props>()

const copiedItems = ref<Record<string, boolean>>({})

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
</script>

<template>
  <div class="bg-muted rounded-lg p-4">
    <div class="flex flex-col md:flex-row md:items-center gap-4">
      <!-- Base URL -->
      <div class="flex items-center gap-2 flex-1">
        <span class="text-slate-700 font-medium text-sm">Base URL:</span>
        <div class="flex items-center gap-2 bg-background rounded-md px-3 py-1.5 flex-1 border">
          <code class="text-slate-900 font-mono text-xs">
            {{ baseUrl }}
          </code>
          <Tooltip>
            <TooltipTrigger asChild>
              <Button
                  variant="ghost"
                  size="sm"
                  @click="copyToClipboard(baseUrl, 'baseUrl')"
                  class="h-6 w-6 p-0 hover:bg-slate-100"
              >
                <Check v-if="copiedItems['baseUrl']" class="w-3 h-3 text-green-600" />
                <Copy v-else class="w-3 h-3 text-slate-600" />
              </Button>
            </TooltipTrigger>
            <TooltipContent>
              <p>Copy base URL</p>
            </TooltipContent>
          </Tooltip>
        </div>
      </div>

      <!-- Documentation Link -->
      <Button variant="outline" size="sm" asChild>
        <a :href="documentationUrl" target="_blank" class="inline-flex items-center gap-2">
          <ExternalLink class="w-3 h-3" />
          Documentation
        </a>
      </Button>
    </div>
  </div>
</template>