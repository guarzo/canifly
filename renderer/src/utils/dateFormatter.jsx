// crypto/dateFormatter.js

export const formatDate = (isoString) => {
    return new Date(isoString).toLocaleString(undefined, {
        year: undefined, // Remove the year
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit',
        hour12: false, // Use 24-hour time
    });
};

