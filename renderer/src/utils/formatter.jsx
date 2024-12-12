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

export const calculateDaysFromToday = (date) => {
    if (!date) return "";
    const targetDate = new Date(date);
    const currentDate = new Date();
    const diffTime = targetDate - currentDate;
    if (diffTime <= 0) return "0 days"; // If the date has passed or is today
    return `${Math.ceil(diffTime / (1000 * 60 * 60 * 24))} days`; // Calculate remaining days
};
export const formatNumberWithCommas = (num) => {
    return num.toLocaleString(); // Using toLocaleString to format numbers with commas
};