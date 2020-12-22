<script lang="ts">
  import CodeMirror from "codemirror";
  import "codemirror/addon/display/placeholder";
  import { onMount } from "svelte";
  import {
    chkDuplicates,
    rmDuplicates,
    chkConsecutive,
    compare,
    diff,
  } from "./sda";
  import { preprocess, format, copy } from "./utils";

  let data1: CodeMirror.EditorFromTextArea,
    data2: CodeMirror.EditorFromTextArea;
  let result = "";
  let source = "Data1";
  let mode = "comm";
  let ignoreDuplicates = true;
  let loading = false;

  onMount(() => {
    data1 = CodeMirror.fromTextArea(
      document.getElementById("inputA") as HTMLTextAreaElement,
      {
        lineNumbers: true,
        lineWrapping: true,
      }
    );
    data2 = CodeMirror.fromTextArea(
      document.getElementById("inputB") as HTMLTextAreaElement,
      {
        lineNumbers: true,
        lineWrapping: true,
      }
    );
    let data = localStorage.getItem("data1");
    if (data) data1.setValue(data);
    data = localStorage.getItem("data2");
    if (data) data2.setValue(data);
    setTimeout(() => {
      data1.refresh();
      data2.refresh();
    }, 200);
  });

  function analyze(operation: string): void {
    loading = true;
    const process = processing();
    const start = new Date().getTime();
    let output = "";
    let r: string[];
    switch (operation) {
      case "chkDuplicates":
        let d: { [k: string]: number };
        if (source == "Data1")
          d = new chkDuplicates(preprocess(data1.getValue())).run();
        else d = new chkDuplicates(preprocess(data2.getValue())).run();
        if (!Object.keys(d).length)
          output = `${source} has no duplicate value.`;
        else
          output =
            `Duplicate values found in ${source}.\n` +
            format(
              Object.keys(d).length,
              Object.keys(d).map((key) => `${key} appears ${d[key]} times.`)
            );
        break;
      case "rmDuplicates":
        if (source == "Data1")
          r = new rmDuplicates(preprocess(data1.getValue())).run();
        else r = new rmDuplicates(preprocess(data2.getValue())).run();
        output = r.join("\n");
        break;
      case "chkConsecutive":
        if (source == "Data1")
          r = new chkConsecutive(preprocess(data1.getValue())).run();
        else r = new chkConsecutive(preprocess(data2.getValue())).run();
        if (!r.length) output = `${source} contains consecutive numbers.`;
        else if (r.length == 1 && r[0] == "!Error!")
          output = `Error!\n${source} contains non-numeric value. Please check!`;
        else
          output = `${source} is not consecutive.
\nThe following numbers are missing:\n${r.join("\n")}`;
        break;
      case "compare":
        if (mode == "comm") {
          r = new compare(
            preprocess(data1.getValue()),
            preprocess(data2.getValue())
          ).run();
          if (!r.length) output = "Two data contain no common value.";
          else
            output = `Common values found between two data.
${format(r.length, r)}`;
        } else {
          const r1 = new compare(
            preprocess(data1.getValue()),
            preprocess(data2.getValue()),
            mode,
            ignoreDuplicates
          ).run();
          const r2 = new compare(
            preprocess(data2.getValue()),
            preprocess(data1.getValue()),
            mode,
            ignoreDuplicates
          ).run();
          if (r1.length + r2.length == 0) {
            output = "Data1 is same as Data2.";
          } else if (!r1.length) {
            output = `Data2 completely contains Data1.\n\nData2 is more than Data1
${format(r2.length, r2)}`;
          } else if (!r2.length) {
            output = `Data1 completely contains Data2.\n\nData1 is more than Data2\n${format(
              r1.length,
              r1
            )}`;
          } else {
            output = `Two files have inconsistent content.
\nData1 is more than Data2\n${format(r1.length, r1)}
\nData2 is more than Data1\n${format(r2.length, r2)}`;
          }
        }
        break;
      case "diff":
        output = new diff(data1.getValue(), data2.getValue())
          .run()
          .replace(`${"=".repeat(67)}\n`, "");
    }
    clearInterval(process);
    result = output + `\n\nDuration for process: ${Date.now() - start}ms`;
    loading = false;
  }

  function clear() {
    data1.setValue("");
    data2.setValue("");
    result = "";
  }

  function swap() {
    const data = data1.getValue();
    data1.setValue(data2.getValue());
    data2.setValue(data);
  }

  function processing(): number {
    return setInterval(() => {
      let dots = result.split("Processing").length + 1;
      if (dots >= 4) dots -= 4;
      result = "Processing" + ".".repeat(dots);
    }, 200);
  }
</script>

<svelte:window
  on:beforeunload={() => {
    localStorage.setItem('data1', data1.getValue());
    localStorage.setItem('data2', data2.getValue());
  }} />

<main>
  <header class="navbar navbar-expand navbar-light flex-column flex-md-row">
    <a
      class="navbar-brand text-primary m-0 mr-md-3"
      href="/"
      style="font-size:24px">Simple Data Analysis</a>
  </header>
  <div class="container-fluid">
    <div class="row">
      <div class="col-3">
        <label for="inputA">Data1</label>
        <textarea id="inputA" placeholder="Paste content here..." />
      </div>
      <div class="col-3 pl-0">
        <label for="inputB">Data2</label>
        <textarea id="inputB" placeholder="Paste content here..." />
      </div>
      <div class="col-2 p-0 pt-5">
        <button
          on:click={() => analyze('chkDuplicates')}
          type="button"
          class="btn btn-primary btn-block"
          disabled={loading}>Check Duplicates</button>
        <button
          on:click={() => analyze('rmDuplicates')}
          type="button"
          class="btn btn-primary btn-block"
          disabled={loading}>Remove Duplicates</button>
        <button
          on:click={() => analyze('chkConsecutive')}
          type="button"
          class="btn btn-primary btn-block"
          disabled={loading}>Check Consecutive</button>
        <div class="d-flex justify-content-around">
          <div>
            <input type="radio" bind:group={source} value="Data1" id="Data1" />
            <label class="m-0" for="Data1">Data1</label>
          </div>
          <div>
            <input type="radio" bind:group={source} value="Data2" id="Data2" />
            <label class="m-0" for="Data2">Data2</label>
          </div>
        </div>
        <br />
        <button
          on:click={() => analyze('compare')}
          type="button"
          class="btn btn-primary btn-block"
          disabled={loading}>Cross Compare</button>
        <div class="d-flex justify-content-around">
          <div>
            <input type="radio" bind:group={mode} value="comm" id="comm" />
            <label class="m-0" for="comm">Comm</label>
          </div>
          <div>
            <input type="radio" bind:group={mode} value="diff" id="diff" />
            <label class="m-0" for="diff">Diff</label>
          </div>
        </div>
        <div class="d-flex justify-content-around">
          <div>
            <input
              type="checkbox"
              bind:checked={ignoreDuplicates}
              id="ignore_duplicates" />
            <label class="m-0" for="ignore_duplicates">Ignore Duplicates</label>
          </div>
        </div>
        <br />
        <button
          on:click={() => analyze('diff')}
          type="button"
          class="btn btn-primary btn-block"
          disabled={loading}>Diff</button>
        <br />
        <br />
        <button
          on:click={async () => copy(result)}
          type="button"
          class="btn btn-primary btn-block"
          disabled={loading}>Copy Result</button>
        <br />
        <br />
        <button
          on:click={swap}
          type="button"
          class="btn btn-primary btn-block"
          disabled={loading}>{@html 'Data1<=>Data2'}</button>
        <br />
        <br />
        <button
          on:click={clear}
          type="button"
          class="btn btn-primary btn-block"
          disabled={loading}>Clear</button>
      </div>
      <div class="col-4">
        <label for="result">Result</label>
        <pre id="result">{result}</pre>
      </div>
    </div>
  </div>
</main>
