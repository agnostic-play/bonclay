import { useNotification } from '@/composables/useNotification'

export interface NotificationConfig {
    showSuccess?: boolean
    showError?: boolean
    successMessage?: string
    errorMessage?: string
}

interface ApiNotificationPayload {
    type?: 'success' | 'info' | 'warning' | 'error'
    title?: string
    message?: string
    duration?: number
}

function getErrorTitle(status?: number): string {
    switch (status) {
        case 400:
            return 'Bad Request'
        case 401:
            return 'Unauthorized'
        case 403:
            return 'Forbidden'
        case 404:
            return 'Not Found'
        case 422:
            return 'Validation Error'
        case 500:
            return 'Server Error'
        default:
            return 'Error'
    }
}

export function useApiFeedback() {
    const notification = useNotification()

    function handleSuccess<ResponseT = any>(
        response: ResponseT,
        config?: NotificationConfig,
    ): ResponseT {
        // Try to read API-provided notification shape
        const anyRes = response as any
        const apiNotif: ApiNotificationPayload | string | undefined =
            anyRes?.notification || anyRes?.message

        if (apiNotif) {
            if (typeof apiNotif === 'string') {
                notification.success('Success', apiNotif)
            } else if (typeof apiNotif === 'object') {
                notification.showNotification({
                    type: apiNotif.type || 'success',
                    title: apiNotif.title || 'Success',
                    message: apiNotif.message,
                    duration: apiNotif.duration,
                })
            }
        } else if (config?.showSuccess && config?.successMessage) {
            notification.success('Success', config.successMessage)
        }

        return response
    }

    function handleError(error: any, config?: NotificationConfig): never {
        const errData = error?.response?.data
        const message = errData?.message || errData?.error || error?.message

        if (errData?.notification) {
            const n: ApiNotificationPayload = errData.notification
            notification.showNotification({
                type: n.type || 'error',
                title: n.title || 'Error',
                message: n.message || message,
            })
        } else if (config?.showError !== false) {
            const title = getErrorTitle(error?.response?.status)
            notification.error(title, config?.errorMessage || message || 'An error occurred')
        }

        throw error
    }

    /**
     * Convenience wrapper to bind success/error feedback around an API promise.
     *
     * Example:
     *   await withApiFeedback(userService.actCreate(payload), {
     *     successMessage: 'User created',
     *   })
     */
    async function withApiFeedback<T>(
        promise: Promise<T>,
        config?: NotificationConfig,
    ): Promise<T> {
        try {
            const res = await promise
            return handleSuccess(res, config)
        } catch (err) {
            return handleError(err, config)
        }
    }

    return { handleSuccess, handleError, withApiFeedback }
}