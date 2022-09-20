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

export const preprocess = (data: string) => {
  return data.split('\n').map(i => i.trim()).map(i => {
    if ((i.split('"').length - 1) % 2) {
      const c = i.match(/^"+/g)
      if (c) if (c[0].length % 2) return i.slice(1)
      return i.slice(0, -1)
    }
    return i
  }).filter(i => i != '')
}

export const format = (length: number, content: string[]) => {
  return `\nTotal ${length} ${length == 1 ? 'record' : 'records'}
\nresult:\n${content.join('\n')}`
}

export default utils
