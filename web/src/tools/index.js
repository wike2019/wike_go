function timeFormat(dateStr) {
    const date = new Date(dateStr);
    return  date.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
        second: '2-digit',
        hour12: false // 使用24小时制
    });
}
export {timeFormat}