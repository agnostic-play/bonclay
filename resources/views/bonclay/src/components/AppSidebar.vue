<script setup lang="ts">
import { type SidebarProps } from '@/components/ui/sidebar'
import { SquareTerminal, GitBranch, LockKeyhole, ScanLine } from "lucide-vue-next"
import NavMain from '@/components/NavMain.vue'
import NavUser from '@/components/NavUser.vue'
import {
  Sidebar,
  SidebarContent,
  SidebarFooter,
  SidebarHeader,
  SidebarRail,
} from '@/components/ui/sidebar'
import { useRoute } from "vue-router"
import { computed } from "vue"
import TeamSwitcher from "@/components/TeamSwitcher.vue"

const props = withDefaults(defineProps<SidebarProps>(), {
  collapsible: "icon",
})

const route = useRoute()

const user = {
  name: "shadcn",
  email: "m@example.com",
  avatar: "/avatars/shadcn.jpg",
}

const navMain = computed(() => [
  {
    title: "Mock API",
    url: "#",
    icon: SquareTerminal,
    isActive: route.path.startsWith("/tools/mock-api"),
    items: [
      {
        title: "API Collection",
        routeName: "MockApiTools-Index",
      },
    ],
  },
  {
    title: "Diagram",
    url: "#",
    icon: GitBranch,
    isActive: route.path.startsWith("/tools/mermaid-diagram"),
    items: [
      {
        title: "Collections",
        routeName: "MermaidDiagramTools-Index",
      },
    ],
  },
  {
    title: "Encryption",
    url: "#",
    icon: LockKeyhole,
    isActive: route.path.startsWith("/tools/encryption-tools"),
    items: [
      {
        title: "Encrypt / Decrypt",
        routeName: "EncryptionTools-Index",
      },
    ],
  },
  {
    title: "Parser",
    url: "#",
    icon: ScanLine,
    isActive: route.path.startsWith("/tools/parser"),
    items: [
      {
        title: "EMVCo / QRIS",
        routeName: "ParserTools-Emvco",
      },
      {
        title: "ISO 8583",
        routeName: "ParserTools-ISO",
      },
    ],
  },
])
</script>

<template>
  <Sidebar v-bind="props">
    <SidebarHeader>
      <!-- TeamSwitcher is self-fetching; no props needed -->
      <TeamSwitcher />
    </SidebarHeader>
    <SidebarContent>
      <NavMain :items="navMain" />
    </SidebarContent>
    <SidebarFooter>
      <NavUser :user="user" />
    </SidebarFooter>
    <SidebarRail />
  </Sidebar>
</template>
