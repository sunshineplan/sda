namespace utils {
  export function contains(data: string[], value: string): boolean {
    for (let i = 0; i < data.length; i++) {
      if (data[i] == value) return true
    }
    return false
  }

  export function range(min: string, max: string): string[] {
    return [...Array(Number(max) - Number(min) + 1).keys()].map(i => i + Number(min)).map(String)
  }

  export function sort(data: string[]): string[] {
    let collator = new Intl.Collator(undefined, { numeric: true, sensitivity: 'base' })
    data.sort(collator.compare)
    return data
  }
}

export function preprocess(data: string, sort = true): string[] {
  if (!sort) return data.split('\n').map(i => i.trim()).filter(i => i != '')
  return utils.sort(data.split('\n').map(i => i.trim()).filter(i => i != ''))
}

export function format(length: number, content: string[]): string {
  return `\nTotal ${length} ${length == 1 ? "record" : "records"}
\nresult:\n${content.join("\n")}`
}

export async function copy(data: string) {
  if (data.trim() !== "") {
    await navigator.clipboard.writeText(data.trim());
    alert("Text has been copied to clipboard.");
  }
}

export default utils
