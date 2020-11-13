const sda = Vue.createApp({
  data() { return { output: '' } }
})

sda.component("sda", {
  template: `
<div class='col-2 p-0 pt-5'>
  <button @click='chk_duplicates' type='button' :class='buttonClass' :disabled='loading'>Check Duplicates</button>
  <button @click='rm_duplicates' type='button' :class='buttonClass' :disabled='loading'>Remove Duplicates</button>
  <button @click='chk_consecutive' type='button' :class='buttonClass' :disabled='loading'>Check Consecutive</button>
  <div class='d-flex justify-content-around'>
    <div>
      <input type='radio' v-model='source' value='Data1' id='data1'>
      <label class='m-0' for='data1'>Data1</label>
    </div>
    <div>
      <input type='radio' v-model='source' value='Data2' id='data2'>
      <label class='m-0' for='data2'>Data2</label>
    </div>
  </div>
  <br>
  <button @click='compare' type='button' :class='buttonClass' :disabled='loading'>Cross Compare</button>
  <div class='d-flex justify-content-around'>
    <div>
      <input type='radio' v-model='mode' value='comm' id='comm'>
      <label class='m-0' for='comm'>Comm</label>
    </div>
    <div>
      <input type='radio' v-model='mode' value='diff' id='diff'>
      <label class='m-0' for='diff'>Diff</label>
    </div>
  </div>
  <div class='d-flex justify-content-around'>
    <div>
      <input type='checkbox' v-model='ignore_duplicates' id='ignore_duplicates'>
      <label class='m-0' for='ignore_duplicates'>Ignore Duplicates</label>
    </div>
  </div>
  <br>
  <button @click='diff' type='button' :class='buttonClass' :disabled='loading'>Diff</button>
  <br>
  <br>
  <button @click='copy' type='button' :class='buttonClass' :disabled='loading'>Copy Result</button>
  <br>
  <br>
  <button @click='swap' type='button' :class='buttonClass' :disabled='loading'>Data1<=>Data2</button>
  <br>
  <br>
  <button @click='clear' type='button' :class='buttonClass' :disabled='loading'>Clear</button>
</div>`,
  data() {
    return {
      source: 'Data1',
      mode: 'comm',
      ignore_duplicates: true,
      loading: false,
      process: '',
      buttonClass: ['btn', 'btn-primary', 'btn-block']
    }
  },
  methods: {
    analyze: function (obj) {
      if (obj.source == undefined)
        obj.source = this.source
      switch (obj.source) {
        case 'Data1':
          obj.data1 = inputA.getValue()
          break
        case 'Data2':
          obj.data2 = inputB.getValue()
          break
        default:
          obj.data1 = inputA.getValue()
          obj.data2 = inputB.getValue()
      }
      this.loading = true
      this.processing()
      this.process = setInterval(() => this.processing(), 1000)
      return fetch('analyze', {
        method: "post",
        headers: {
          'Content-Type': 'application/x-www-form-urlencoded'
        },
        body: new URLSearchParams(obj)
      })
        .then(response => response.text())
        .then(text => {
          setTimeout(() => {
            clearInterval(this.process)
            this.$parent.output = text
            this.loading = false
          }, 800);
        })
        .catch(() => setTimeout(() => {
          this.$parent.output = 'Oops, we encountered a problem! Please contact your administrator for assistance.'
          this.loading = false
        }, 800))
    },
    chk_duplicates: function () {
      this.analyze({ func: 'chk_duplicates' })
    },
    rm_duplicates: function () {
      this.analyze({ func: 'rm_duplicates' })
    },
    chk_consecutive: function () {
      this.analyze({ func: 'chk_consecutive' })
    },
    compare: function () {
      this.analyze({
        func: 'compare',
        mode: this.mode,
        source: 'Data1,Data2',
        ignore_duplicates: this.ignore_duplicates
      })
    },
    diff: function () {
      this.analyze({
        func: 'diff',
        source: 'Data1,Data2'
      })
    },
    copy: function () {
      if (this.$parent.output.trim() !== '')
        navigator.clipboard.writeText(this.$parent.output)
          .then(() => alert('Text has been copied to clipboard.'))
          .catch(() => alert('Unable to copy to clipboard.'))
    },
    clear: function () {
      inputA.setValue('')
      inputB.setValue('')
      this.$parent.output = ''
    },
    swap: () => {
      var data1 = inputA.getValue()
      inputA.setValue(inputB.getValue())
      inputB.setValue(data1)
    },
    processing: function () {
      this.$parent.output = 'Processing'
      setTimeout(() => this.$parent.output = 'Processing.', 250)
      setTimeout(() => this.$parent.output = 'Processing..', 500)
      setTimeout(() => this.$parent.output = 'Processing...', 750)
    }
  }
})

sda.mount('#sda')
