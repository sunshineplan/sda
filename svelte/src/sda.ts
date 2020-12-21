import utils from './utils'

interface sda {
  run(): string[] | { [k: string]: number }
}

export class chkDuplicates implements sda {
  data: string[]
  constructor(data: string[]) {
    this.data = data
  }
  run(): { [k: string]: number } {
    const result: { [k: string]: number } = {}
    this.data
      .map(i => {
        if (result[i]) result[i]++
        else result[i] = 1
      })
    Object.keys(result).map(key => {
      if (result[key] == 1) delete result[key]
    })
    return Object.fromEntries(
      Object.entries(result).sort(([, a], [, b]) => b - a)
    );
  }
}

export class rmDuplicates implements sda {
  data: string[]
  constructor(data: string[]) {
    this.data = data
  }
  run(): string[] { return [...new Set(this.data)] }
}

export class compare implements sda {
  data1: string[]
  data2: string[]
  mode: string
  ignoreDuplicates: boolean
  constructor(data1: string[], data2: string[], mode: string = 'comm', ignoreDuplicates: boolean = true) {
    if (ignoreDuplicates) {
      this.data1 = new rmDuplicates(data1).run()
      this.data2 = new rmDuplicates(data2).run()
    } else {
      this.data1 = data1
      this.data2 = data2
    }
    this.mode = mode
    this.ignoreDuplicates = ignoreDuplicates
  }
  run(): string[] {
    let result: string[] = []
    if (this.mode == 'diff') {
      const data = [...this.data1]
      for (let i = 0; i < this.data2.length; i++)
        if (utils.contains(data, this.data2[i]))
          data.splice(data.indexOf(this.data2[i]), 1)
      result = data
    } else {
      let data: string[]
      if (this.ignoreDuplicates) data = this.data1
      else data = [...new Set(this.data1)]
      for (let i = 0; i < data.length; i++) {
        if (utils.contains(this.data2, data[i])) result.push(data[i])
      }
    }
    return result
  }
}

export class chkConsecutive implements sda {
  data: string[]
  constructor(data: string[]) {
    this.data = utils.sort(data)
  }
  run(): string[] {
    for (let i = 0; i < this.data.length; i++) {
      const n = Number(this.data[i])
      if (isNaN(n)) return ['!Error!']
    }
    return new compare(
      utils.range(this.data[0], this.data[this.data.length - 1]),
      this.data,
      'diff'
    ).run()
  }
}

export class diff implements sda {
  data1: string[]
  data2: string[]
  constructor(data1: string[], data2: string[]) {
    this.data1 = data1
    this.data2 = data2
  }
  run(): string[] {
    return this.data1
  }
}

export default sda
