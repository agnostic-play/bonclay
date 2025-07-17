import { ref, onMounted, onUnmounted } from 'vue'

export function useResponsive() {
    const isMobile = ref(false)
    const isMobileMenuOpen = ref(false)

    const checkMobile = () => {
        isMobile.value = window.innerWidth < 768
        if (!isMobile.value) {
            isMobileMenuOpen.value = false
        }
    }

    onMounted(() => {
        checkMobile()
        window.addEventListener('resize', checkMobile)
    })

    onUnmounted(() => {
        window.removeEventListener('resize', checkMobile)
    })

    return {
        isMobile,
        isMobileMenuOpen
    }
}