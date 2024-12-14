// utils/dateFormatter.js

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

export const calculateDaysFromToday = (date) => {
    if (!date) return "";
    const targetDate = new Date(date);
    const currentDate = new Date();
    const diffTime = targetDate - currentDate;
    if (diffTime <= 0) return "0 days";
    // Use floor instead of ceil to avoid rounding up partial days
    return `${Math.floor(diffTime / (1000 * 60 * 60 * 24))} days`;
};

export const formatNumberWithCommas = (num) => {
    return num.toLocaleString(); // Using toLocaleString to format numbers with commas
};

export const formatSP = (totalSp) => {
    if (!totalSp || totalSp < 1_000_000) {
        return '<1M SP';
    }

    const millions = Math.round(totalSp / 1_000_000);
    return `${millions}M SP`;
}