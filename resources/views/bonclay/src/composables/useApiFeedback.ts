import { ref } from 'vue'
import { useNotification } from '@/composables/useNotification'

export interface NotificationConfig {
    showSuccess?: boolean
    showError?: boolean
    successMessage?: string
    errorMessage?: string
    /**
     * If true, skips global toast error notification
     * but still captures field errors.
     */
    silentError?: boolean
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

/**
 * Core composable for handling API feedback, toasts, and validation.
 */
export function useApiFeedback() {
    const notification = useNotification()

    const loading = ref(false)
    const error = ref<any>(null)
    const fieldErrors = ref<Record<string, string>>({})
    const data = ref<any>(null)

    function handleSuccess<ResponseT = any>(
        response: ResponseT,
        config?: NotificationConfig
    ): ResponseT {
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


    function handleError(errorObj: any, config?: NotificationConfig): never {
        const errRes = errorObj?.response
        const errData = errRes?.data
        // Prefer nested data.message; fall back to errorObj.message
        const message =
            errData?.message || errData?.error || errorObj?.message || 'Unknown error'

        fieldErrors.value = {}
        error.value = errorObj

        // Drill into nested "data" object if present
        const validationBag =
            (errData && typeof errData === 'object' && errData.data && typeof errData.data === 'object')
                ? errData.data
                : (errData && typeof errData === 'object' ? errData : null)

        if (validationBag) {
            // Collect only string values (field -> message)
            const entries = Object.entries(validationBag)
                .filter(([key, val]) =>
                    typeof val === 'string' && key !== 'message' && key !== 'error'
                ) as Array<[string, string]>

            if (entries.length > 0) {
                fieldErrors.value = Object.fromEntries(entries)
            }
        }

        if (config?.silentError) throw errorObj

        if (errData?.notification) {
            const n: ApiNotificationPayload = errData.notification
            notification.showNotification({
                type: n.type || 'error',
                title: n.title || getErrorTitle(errRes?.status),
                message: n.message || message,
            })
        } else if (config?.showError !== false) {
            const title = getErrorTitle(errRes?.status)
            // If we have field errors, append them to the message for clarity
            const details = Object.entries(fieldErrors.value)
                .map(([k, v]) => `${k}: ${v}`)
                .join(', ')
            notification.error(title, details ? `${message} — ${details}` : message)
        }

        throw errorObj
    }



    // function handleError(errorObj: any, config?: NotificationConfig): never {
    //     const errRes = errorObj?.response
    //     const errData = errRes?.data
    //     const message =
    //         errData?.message || errData?.error || errorObj?.message || 'Unknown error'
    //
    //     fieldErrors.value = {}
    //     error.value = errorObj
    //
    //     // Detect validation fields
    //     if (errData && typeof errData === 'object') {
    //         const possibleFields = Object.entries(errData).filter(
    //             ([key, val]) => typeof val === 'string' && key !== 'message' && key !== 'error'
    //         )
    //         if (possibleFields.length > 0) {
    //             fieldErrors.value = {
    //                 ...fieldErrors.value,
    //             }
    //         }
    //     }
    //
    //     if (config?.silentError) throw errorObj
    //
    //     if (errData?.notification) {
    //         const n: ApiNotificationPayload = errData.notification
    //         notification.showNotification({
    //             type: n.type || 'error',
    //             title: n.title || 'Error',
    //             message: n.message || message,
    //         })
    //     } else if (config?.showError !== false) {
    //         const title = getErrorTitle(errRes?.status)
    //         notification.error(title, config?.errorMessage || message)
    //     }
    //
    //     throw errorObj
    // }

    /**
     * Wrapper to simplify async API calls with built-in notifications & state.
     * Returns { data, error, fieldErrors } for convenient destructuring.
     *
     * Example:
     *   const { data, fieldErrors } = await withApiFeedback(userApi.create(form), {
     *     successMessage: 'User created',
     *   })
     */
    async function withApiFeedback<T>(
        promise: Promise<T>,
        config?: NotificationConfig
    ): Promise<{
        data: T | null
        error: any
        fieldErrors: Record<string, string>
    }> {
        loading.value = true
        error.value = null
        fieldErrors.value = {}
        data.value = null

        try {
            const res = await promise
            const result = handleSuccess(res, config)
            data.value = result
            return { data: result, error: null, fieldErrors: {} }
        } catch (err: any) {
            handleError(err, config)
            return { data: null, error: err, fieldErrors: fieldErrors.value }
        } finally {
            loading.value = false
        }
    }

    return {
        loading,
        error,
        data,
        fieldErrors,
        handleSuccess,
        handleError,
        withApiFeedback,
    }
}
