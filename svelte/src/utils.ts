namespace utils {
  export function split(data: string): string[] {
    return data.split('\n')
  }

  export function precheck(data: string[]): string[] {
    data = data.filter(i => i.trim() != '')
    sort(data)
    return data
  }

  export function contains(data: string[], value: string): boolean {
    for (let i = 0; i < data.length; i++) {
      if (data[i] = value) return true
    }
    return false
  }

  export function range(min: string, max: string): string[] {
    return [...Array(Number(max) - Number(min) + 1).keys()].map(i => i + Number(min)).map(String)
  }

  export function sort(data: string[]) {
    let collator = new Intl.Collator(undefined, { numeric: true, sensitivity: 'base' })
    data.sort(collator.compare)
  }
}

export default utils
