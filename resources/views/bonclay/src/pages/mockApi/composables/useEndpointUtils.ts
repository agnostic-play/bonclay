export const useEndpointUtils = () => {
    const getMethodColor = (method: string): string => {
        const colors: Record<string, string> = {
            GET: 'bg-green-500 hover:bg-green-600',
            POST: 'bg-blue-500 hover:bg-blue-600',
            PUT: 'bg-orange-500 hover:bg-orange-600',
            DELETE: 'bg-red-500 hover:bg-red-600',
            PATCH: 'bg-purple-500 hover:bg-purple-600'
        }
        return colors[method] || 'bg-gray-500 hover:bg-gray-600'
    }

    const getMethodBgColor = (method: string): string => {
        const colors: Record<string, string> = {
            GET: 'bg-green-300/10',
            POST: 'bg-blue-300/10',
            PUT: 'bg-orange-300/10',
            DELETE: 'bg-red-300/10',
            PATCH: 'bg-purple-300/10'
        }
        return colors[method] || 'bg-gray-300/10'
    }

    return {
        getMethodColor,
        getMethodBgColor
    }
}