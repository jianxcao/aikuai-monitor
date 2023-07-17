const underlineToCamelCaseSingle = (key: string) =>
  typeof key === 'string' ? key.replace(/_(\w)/g, (_, $1) => $1.toUpperCase()) : `${key}`

const camelCaseToUnderlineSingleString = (key: string) =>
  typeof key === 'string' ? key.replace(/([A-Z]\W*)/g, (_, $1) => '_' + $1.toLowerCase()) : `${key}`

const underlineToCamelCase = (data: any): any => {
  if (Object.prototype.toString.call(data) === '[object Object]') {
    const newObj: Record<string, any> = {}
    Object.entries(data).forEach(([key, value]) => {
      const _value = underlineToCamelCase(value)
      newObj[underlineToCamelCaseSingle(key)] = _value
    })
    return newObj
  } else if (Array.isArray(data)) {
    return data.map((item) => underlineToCamelCase(item))
  } else return data
}

const camelCaseToUnderline = (data: any): any => {
  if (Object.prototype.toString.call(data) === '[object Object]') {
    const newObj: Record<string, any> = {}
    Object.entries(data).forEach(([key, value]) => {
      const _value = camelCaseToUnderline(value)
      newObj[camelCaseToUnderlineSingleString(key)] = _value
    })
    return newObj
  } else if (Array.isArray(data)) {
    return data.map((item) => camelCaseToUnderline(item))
  } else return data
}

// 规范化对象的 key
export const normalizeKey = {
  underlineToCamelCase,
  camelCaseToUnderline,
  underlineToCamelCaseSingle,
  camelCaseToUnderlineSingleString
}
