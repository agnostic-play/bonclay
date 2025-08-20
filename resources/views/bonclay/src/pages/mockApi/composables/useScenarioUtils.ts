export const useScenarioUtils = () => {
    const getStatusColor = (status: number): string => {
        if (status >= 200 && status < 300) return 'text-green-600 bg-green-50'
        if (status >= 300 && status < 400) return 'text-blue-600 bg-blue-50'
        if (status >= 400 && status < 500) return 'text-orange-600 bg-orange-50'
        return 'text-red-600 bg-red-50'
    }

    const formatDelay = (seconds: number): string => {
        if (seconds === 0) return '0s'
        if (seconds === 1) return '1s'
        return `${seconds}s`
    }

    return {
        getStatusColor,
        formatDelay
    }
}