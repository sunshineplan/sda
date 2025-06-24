namespace utils {
  export const contains = (data: string[], value: string) => {
    for (let i = 0; i < data.length; i++) {
      if (data[i] == value) return true
    }
    return false
  }

  export const range = (min: string, max: string) => {
    return [...Array(Number(max) - Number(min) + 1).keys()].map(i => i + Number(min)).map(String)
  }

  export const sort = (data: string[]) => {
    let collator = new Intl.Collator(undefined, { numeric: true, sensitivity: 'base' })
    data.sort(collator.compare)
    return data
  }
}

export const format = (length: number, content: string[]) => {
  return `\nTotal ${length} ${length == 1 ? 'record' : 'records'}
\nresult:\n${content.join('\n')}`
}

export default utils
