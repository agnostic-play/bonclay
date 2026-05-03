<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Loader2 } from 'lucide-vue-next'
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogDescription,
  DialogFooter,
} from '@/components/ui/dialog'
import {
  Select,
  SelectContent,
  SelectGroup,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'
import { Button } from '@/components/ui/button'
import { useSquadSession } from '@/composables/useSquadSession'
import { useSquadStore } from '@/composables/useSquadStore'

const router = useRouter()
const { activeSquad, setActiveSquad } = useSquadSession()
const { squads, loading, fetchSquads } = useSquadStore()

const open = ref(false)
const selectedSlug = ref<string>('')

// Show the modal whenever there is no active squad in session
watch(
  () => activeSquad.value,
  (squad) => { open.value = !squad },
  { immediate: true }
)

onMounted(() => void fetchSquads())

const handleConfirm = () => {
  const squad = squads.value.find(s => (s.slug || s.id) === selectedSlug.value)
  if (!squad) return

  setActiveSquad({ id: squad.id, slug: squad.slug, name: squad.name, plan: squad.plan })
  open.value = false
  router.push({ name: 'MockApiTools-Index', params: { slug: squad.slug } })
}
</script>

<template>
  <Dialog :open="open">
    <DialogContent
      :show-close-button="false"
      class="sm:max-w-sm"
      @escape-key-down="(e: Event) => e.preventDefault()"
      @interact-outside="(e: Event) => e.preventDefault()"
    >
      <DialogHeader>
        <DialogTitle>Choose your squad</DialogTitle>
        <DialogDescription>
          Select the squad you want to work with. You can switch squads at any time from the sidebar.
        </DialogDescription>
      </DialogHeader>

      <div class="py-2">
        <div v-if="loading" class="flex items-center gap-2 text-sm text-muted-foreground py-2">
          <Loader2 class="size-4 animate-spin" />
          Loading squads…
        </div>
        <Select v-else v-model="selectedSlug">
          <SelectTrigger class="w-full">
            <SelectValue placeholder="Select a squad…" />
          </SelectTrigger>
          <SelectContent>
            <SelectGroup>
              <SelectItem
                v-for="squad in squads"
                :key="squad.id"
                :value="squad.slug || squad.id"
              >
                {{ squad.name }}
              </SelectItem>
            </SelectGroup>
          </SelectContent>
        </Select>
      </div>

      <DialogFooter>
        <Button class="w-full" :disabled="!selectedSlug || loading" @click="handleConfirm">
          Continue
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>
