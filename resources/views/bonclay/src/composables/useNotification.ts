import { toast } from "vue-sonner"
export type NotificationType = 'success' | 'error' | 'warning' | 'info'

export interface ApiNotification {
    type: NotificationType
    title: string
    message?: string
    duration?: number
}

export function useNotification() {
    const showNotification = ({
                                  type,
                                  title,
                                  message,
                                  duration = 5000,
                              }: ApiNotification) => {
        const variantMap = {
            success: 'default' as const,
            error: 'destructive' as const,
            warning: 'default' as const,
            info: 'default' as const,
        }

        toast({
            title,
            description: message,
            variant: variantMap[type],
            duration,
        })
    }

    const success = (title: string, message?: string) => {
        showNotification({ type: 'success', title, message })
    }

    const error = (title: string, message?: string) => {
        showNotification({ type: 'error', title, message })
    }

    const warning = (title: string, message?: string) => {
        showNotification({ type: 'warning', title, message })
    }

    const info = (title: string, message?: string) => {
        showNotification({ type: 'info', title, message })
    }

    return {
        showNotification,
        success,
        error,
        warning,
        info,
    }
}
