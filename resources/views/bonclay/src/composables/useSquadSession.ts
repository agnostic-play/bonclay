import { ref, computed } from 'vue'

export interface ActiveSquad {
  id: string
  slug: string
  name: string
  plan?: string
}

const SESSION_KEY = 'bonclay_active_squad'

// Module-level singleton state — shared across all component instances
const activeSquad = ref<ActiveSquad | null>(null)

// Hydrate from sessionStorage once on module load
try {
  const stored = sessionStorage.getItem(SESSION_KEY)
  if (stored) {
    activeSquad.value = JSON.parse(stored) as ActiveSquad
  }
} catch {
  // ignore parse errors
}

export function useSquadSession() {
  const activeSquadId = computed(() => activeSquad.value?.id ?? null)
  const activeSquadSlug = computed(() => activeSquad.value?.slug ?? null)
  const activeSquadName = computed(() => activeSquad.value?.name ?? null)

  const setActiveSquad = (squad: ActiveSquad) => {
    activeSquad.value = squad
    try {
      sessionStorage.setItem(SESSION_KEY, JSON.stringify(squad))
    } catch {
      // ignore
    }
  }

  const clearActiveSquad = () => {
    activeSquad.value = null
    try {
      sessionStorage.removeItem(SESSION_KEY)
    } catch {
      // ignore
    }
  }

  return {
    activeSquad,
    activeSquadId,
    activeSquadSlug,
    activeSquadName,
    setActiveSquad,
    clearActiveSquad,
  }
}
