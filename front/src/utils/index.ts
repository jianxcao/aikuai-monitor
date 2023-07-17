export { normalizeKey } from './object'

export const FormatSize = (size: number) => {
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  let unitIndex = 0
  while (size >= 1024 && unitIndex < units.length - 1) {
    size /= 1024
    unitIndex++
  }
  return `${size.toFixed(2)}${units[unitIndex]}`
}

export const FormatNum = (size: number) => {
  if (size < 10000) {
    return size
  }
  const res = size / 10000
  return `${res.toFixed(2)}万`
}