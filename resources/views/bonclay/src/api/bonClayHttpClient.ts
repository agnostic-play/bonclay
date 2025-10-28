// ================================
// files: bonclay-http-client.ts
// ================================
import {BaseHttpClient} from '@/lib/httpClient/http.ts'
import type {PaginatedResponse, PaginationParams, RequestOptions} from '@/lib/httpClient/types.ts'
import {useNotification} from '@/composables/useNotification'

export interface NotificationConfig {
    showSuccess?: boolean
    showError?: boolean
    successMessage?: string
    errorMessage?: string
}

export abstract class BonclayHttpClient<T> extends BaseHttpClient {
    protected readonly resource: string;
    protected notification = useNotification();

    constructor(resource: string) {
        super(import.meta.env.VITE_BONCLAY_URL || '', {
            timeout: 10_000,
            headers: {'Content-Type': 'application/json'},
        });
        this.resource = resource.replace(/^\/+|\/+$/g, '');
    }

    // Helper method to handle API notifications
    protected handleNotification(response: any, config?: NotificationConfig) {
        // Check if API sends notification in response
        const apiNotification = response?.notification || response?.message;

        if (apiNotification) {
            if (typeof apiNotification === 'string') {
                this.notification.success('Success', apiNotification);
            } else if (typeof apiNotification === 'object') {
                this.notification.showNotification({
                    type: apiNotification.type || 'success',
                    title: apiNotification.title || 'Success',
                    message: apiNotification.message,
                    duration: apiNotification.duration,
                });
            }
        } else if (config?.showSuccess && config?.successMessage) {
            // Show custom success message if configured
            this.notification.success('Success', config.successMessage);
        }
    }

    // Helper method to handle errors
    protected handleError(error: any, config?: NotificationConfig): never {
        const errorData = error.response?.data;
        const errorMessage = errorData?.message || errorData?.error || error.message;

        // Check for API notification structure
        if (errorData?.notification) {
            this.notification.showNotification({
                type: errorData.notification.type || 'error',
                title: errorData.notification.title || 'Error',
                message: errorData.notification.message || errorMessage,
            });
        } else if (config?.showError !== false) {
            // Show error unless explicitly disabled
            const title = this.getErrorTitle(error.response?.status);
            this.notification.error(
                title,
                config?.errorMessage || errorMessage || 'An error occurred'
            );
        }

        throw error;
    }

    // Get appropriate error title based on status code
    protected getErrorTitle(status?: number): string {
        switch (status) {
            case 400:
                return 'Bad Request';
            case 401:
                return 'Unauthorized';
            case 403:
                return 'Forbidden';
            case 404:
                return 'Not Found';
            case 422:
                return 'Validation Error';
            case 500:
                return 'Server Error';
            default:
                return 'Error';
        }
    }

    // --- CRUD (unwrapped to domain directly) ----------------------------------
    async getList(
        params?: PaginationParams,
        options?: RequestOptions & { notification?: NotificationConfig }
    ): Promise<PaginatedResponse<T>> {
        try {
            const response = await this.get<PaginatedResponse<T>>(
                `/${this.resource}`,
                {...options, params}
            );
            this.handleNotification(response, options?.notification);
            return response;
        } catch (error) {
            return this.handleError(error, options?.notification);
        }
    }

    async getDetail(
        id: string | number,
        options?: RequestOptions & { notification?: NotificationConfig }
    ): Promise<T> {
        try {
            const response = await this.get<T>(`/${this.resource}/${id}`, options);
            this.handleNotification(response, options?.notification);
            return response;
        } catch (error) {
            return this.handleError(error, options?.notification);
        }
    }

    async actCreate(
        data: Partial<T>,
        options?: RequestOptions & { notification?: NotificationConfig }
    ): Promise<T> {
        try {
            const response = await this.post<T>(`/${this.resource}`, data, options);
            this.handleNotification(response, {
                showSuccess: true,
                successMessage: 'Item created successfully',
                ...options?.notification,
            });
            return response;
        } catch (error) {
            return this.handleError(error, {
                errorMessage: 'Failed to create item',
                ...options?.notification,
            });
        }
    }

    async actUpdate(
        id: string | number,
        data: Partial<T>,
        options?: RequestOptions & { notification?: NotificationConfig }
    ): Promise<T> {
        try {
            const response = await this.put<T>(`/${this.resource}/${id}`, data, options);
            this.handleNotification(response, {
                showSuccess: true,
                successMessage: 'Item updated successfully',
                ...options?.notification,
            });
            return response;
        } catch (error) {
            return this.handleError(error, {
                errorMessage: 'Failed to update item',
                ...options?.notification,
            });
        }
    }

    async actPatch(
        id: string | number,
        data: Partial<T>,
        options?: RequestOptions & { notification?: NotificationConfig }
    ): Promise<T> {
        try {
            const response = await this.patch<T>(`/${this.resource}/${id}`, data, options);
            this.handleNotification(response, {
                showSuccess: true,
                successMessage: 'Item updated successfully',
                ...options?.notification,
            });
            return response;
        } catch (error) {
            return this.handleError(error, {
                errorMessage: 'Failed to update item',
                ...options?.notification,
            });
        }
    }

    async actRemove(
        id: string | number,
        options?: RequestOptions & { notification?: NotificationConfig }
    ): Promise<void> {
        try {
            await this.delete(`/${this.resource}/${id}`, options);
            // For delete operations, we handle notification manually since there's no response data
            if (options?.notification?.showSuccess !== false) {
                this.notification.success(
                    'Success',
                    options?.notification?.successMessage || 'Item deleted successfully'
                );
            }
        } catch (error) {
            this.handleError(error, {
                errorMessage: 'Failed to delete item',
                ...options?.notification,
            });
        }
    }

    hasNextPage(r: PaginatedResponse<T>): boolean {
        return r.next_page !== null;
    }

    hasPreviousPage(r: PaginatedResponse<T>): boolean {
        return r.previous_page !== null;
    }

    async getNextPage(
        current: PaginatedResponse<T>,
        options?: RequestOptions & { notification?: NotificationConfig }
    ): Promise<PaginatedResponse<T> | null> {
        if (!this.hasNextPage(current)) return null;
        return this.getList({
            page_number: current.next_page!,
            page_size: current.page_size,
            keyword: current.keyword || undefined
        }, options);
    }

    async getPreviousPage(
        current: PaginatedResponse<T>,
        options?: RequestOptions & { notification?: NotificationConfig }
    ): Promise<PaginatedResponse<T> | null> {
        if (!this.hasPreviousPage(current)) return null;
        return this.getList({
            page_number: current.previous_page!,
            page_size: current.page_size,
            keyword: current.keyword || undefined
        }, options);
    }
}

// ================================
// Usage Examples
// ================================

/*
// Example 1: Using default notifications (auto-show for create/update/delete)
class UserService extends BonclayHttpClient<User> {
    constructor() {
        super('users');
    }
}

const userService = new UserService();

// This will automatically show "Item created successfully" on success
await userService.actCreate({ name: 'John' });

// Example 2: Custom notification messages
await userService.actCreate(
    { name: 'John' },
    {
        notification: {
            successMessage: 'User has been created successfully!',
            errorMessage: 'Failed to create user. Please try again.',
        }
    }
);

// Example 3: Disable notifications
await userService.actUpdate(
    userId,
    { name: 'Jane' },
    {
        notification: {
            showSuccess: false,
            showError: false,
        }
    }
);

// Example 4: Delete with custom message
await userService.actRemove(
    userId,
    {
        notification: {
            successMessage: 'User has been permanently deleted',
        }
    }
);

// Example 5: Silent operations
await userService.actRemove(userId, {
    notification: { showSuccess: false, showError: false }
});
*/