// utils.jsx
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
