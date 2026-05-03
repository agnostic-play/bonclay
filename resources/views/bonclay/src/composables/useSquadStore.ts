import { ref, readonly } from 'vue'
import { Clover } from 'lucide-vue-next'
import type { Component } from 'vue'
import { getSquads } from '@/api'

export interface Squad {
  id: string
  name: string
  slug: string
  plan: string
  logo: Component
}

// ── Module-level singleton state ─────────────────────────────────────────────
const squads = ref<Squad[]>([])
const loading = ref(false)
const error = ref<string | null>(null)
const fetched = ref(false)

function mapSquad(s: { id: string; name: string; slug: string; desc?: string }): Squad {
  return {
    id: s.id,
    name: s.name,
    slug: s.slug || s.id,
    plan: s.desc ?? '',
    logo: Clover,
  }
}

// ── Public composable ────────────────────────────────────────────────────────
export function useSquadStore() {
  /**
   * Fetch squads from the API.
   * @param force – pass `true` to bypass the cache and re-fetch (e.g. after
   *                creating a new squad).
   */
  const fetchSquads = async (force = false) => {
    if (loading.value) return
    if (fetched.value && !force) return

    loading.value = true
    error.value = null

    try {
      const res = await getSquads()
      squads.value = (res?.list ?? []).map(mapSquad)
      fetched.value = true
    } catch (err) {
      error.value = 'Failed to load squads'
      console.error('[useSquadStore] fetchSquads error:', err)
    } finally {
      loading.value = false
    }
  }

  /** Invalidate cache so the next call to fetchSquads() hits the network. */
  const invalidate = () => {
    fetched.value = false
  }

  return {
    squads: readonly(squads),
    loading: readonly(loading),
    error: readonly(error),
    fetchSquads,
    invalidate,
  }
}
