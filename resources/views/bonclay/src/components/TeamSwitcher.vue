<script setup lang="ts">
import { ChevronsUpDown, Plus, Loader2, Settings2 } from "lucide-vue-next"
import { ref, computed, onMounted, watch } from "vue"
import CreateSquadSheet from '@/components/CreateSquadSheet.vue'
import EditSquadDialog from '@/components/EditSquadDialog.vue'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu"
import { useRouter } from 'vue-router'
import { useSquadSession } from '@/composables/useSquadSession'
import { useSquadStore } from '@/composables/useSquadStore'

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  useSidebar,
} from "@/components/ui/sidebar"

const router = useRouter()
const { isMobile } = useSidebar()
const { activeSquad, setActiveSquad } = useSquadSession()
const { squads, loading, fetchSquads } = useSquadStore()

const createSquadRef = ref<InstanceType<typeof CreateSquadSheet> | null>(null)
const editSquadOpen = ref(false)
const editingTeam = ref<typeof squads.value[number] | null>(null)

// The squad shown in the trigger button. Initialised from session if available.
const activeTeam = ref(activeSquad.value ? squads.value.find(s => s.slug === activeSquad.value!.slug) ?? null : null)

// When the squad list arrives, restore the active team from session.
watch(
  squads,
  (list) => {
    if (!Array.isArray(list) || list.length === 0) return
    if (activeSquad.value?.slug) {
      const match = list.find(s => s.slug === activeSquad.value!.slug)
      if (match) activeTeam.value = match
    }
  },
  { immediate: true }
)

// When active squad changes (e.g., from SelectSquadModal), sync the trigger display.
watch(
  activeSquad,
  (squad) => {
    if (!squad) {
      activeTeam.value = null
      return
    }
    const match = squads.value.find(s => s.slug === squad.slug)
    if (match) activeTeam.value = match
  }
)

// Display label for the trigger button
const triggerName = computed(() => activeTeam.value?.name ?? 'Select squad')
const triggerPlan = computed(() => activeTeam.value?.plan ?? '')

const selectTeam = (team: typeof squads.value[number]) => {
  activeTeam.value = team
  // Persist to session so any page/component can read it
  setActiveSquad({
    id: team.id,
    slug: team.slug,
    name: team.name,
    plan: team.plan,
  })
  // Navigate to squad view
  if (team.slug) {
    router.push({ name: 'MockApiTools-Index', params: { slug: team.slug } })
  }
}

const handleAddTeam = () => {
  createSquadRef.value?.openDialog()
}

const handleSquadCreated = async () => {
  const store = useSquadStore()
  store.invalidate()
  await fetchSquads(true)
}

const handleEditSquad = (team: typeof squads.value[number]) => {
  editingTeam.value = team
  editSquadOpen.value = true
}

const handleSquadUpdated = async () => {
  const store = useSquadStore()
  store.invalidate()
  await fetchSquads(true)
}

const handleSquadDeleted = async () => {
  const store = useSquadStore()
  store.invalidate()
  await fetchSquads(true)
}

// Fetch the squad list when the sidebar mounts (cached on subsequent renders)
onMounted(() => {
  void fetchSquads()
})
</script>

<template>
  <SidebarMenu>
    <SidebarMenuItem>
      <DropdownMenu>
        <DropdownMenuTrigger as-child>
          <SidebarMenuButton
              size="lg"
              class="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground"
          >
            <div class="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
              <Loader2 v-if="loading" class="size-4 animate-spin" />
              <component v-else-if="activeTeam?.logo" :is="activeTeam.logo" class="size-4" />
              <ChevronsUpDown v-else class="size-4 opacity-50" />
            </div>
            <div class="grid flex-1 text-left text-sm leading-tight">
              <span class="truncate font-semibold">{{ triggerName }}</span>
              <span class="truncate text-xs">{{ triggerPlan }}</span>
            </div>
            <ChevronsUpDown class="ml-auto" />
          </SidebarMenuButton>
        </DropdownMenuTrigger>

        <DropdownMenuContent
            class="w-[--reka-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            align="start"
            :side="isMobile ? 'bottom' : 'right'"
            :side-offset="4"
        >
          <DropdownMenuLabel class="text-xs text-muted-foreground">
            Squads
          </DropdownMenuLabel>

          <!-- Loading state -->
          <div v-if="loading" class="flex items-center justify-center py-4 gap-2 text-muted-foreground text-xs">
            <Loader2 class="size-3 animate-spin" />
            Loading squads…
          </div>

          <!-- Squad list -->
          <template v-else>
            <DropdownMenuItem
              v-for="(team, index) in squads"
              :key="team.slug"
              class="gap-2 p-2 pr-1"
              :class="{ 'bg-accent': activeTeam?.slug === team.slug }"
              @click="selectTeam(team)"
            >
              <div class="flex size-6 items-center justify-center rounded-sm border">
                <component :is="team.logo" class="size-4 shrink-0" />
              </div>
              <span class="flex-1 truncate">{{ team.name }}</span>
              <DropdownMenuShortcut>⌘{{ index + 1 }}</DropdownMenuShortcut>
              <button
                class="ml-1 p-1 rounded opacity-0 group-hover:opacity-100 hover:bg-sidebar-accent text-muted-foreground hover:text-foreground transition-opacity"
                @click.stop="handleEditSquad(team)"
                title="Edit squad"
              >
                <Settings2 class="size-3.5" />
              </button>
            </DropdownMenuItem>
          </template>

          <DropdownMenuSeparator />

          <DropdownMenuItem class="gap-2 p-2" @click="handleAddTeam">
            <div class="flex size-6 items-center justify-center rounded-md border bg-background">
              <Plus class="size-4" />
            </div>
            <div class="font-medium text-muted-foreground">Add squad</div>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>

      <CreateSquadSheet ref="createSquadRef" @created="handleSquadCreated" />

      <EditSquadDialog
        v-if="editingTeam"
        :open="editSquadOpen"
        :squad-id="editingTeam.id"
        :initial-name="editingTeam.name"
        :initial-desc="editingTeam.plan"
        @update:open="editSquadOpen = $event"
        @updated="handleSquadUpdated"
        @deleted="handleSquadDeleted"
      />
    </SidebarMenuItem>
  </SidebarMenu>
</template>
