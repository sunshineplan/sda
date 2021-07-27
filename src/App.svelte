<script lang="ts">
  import CodeMirror from "codemirror";
  import "codemirror/addon/display/placeholder";
  import { onMount } from "svelte";
  import {
    chkDuplicates,
    rmDuplicates,
    chkConsecutive,
    compareDiff,
    compareComm,
    diff,
  } from "./sda";
  import { preprocess, format } from "./utils";

  let data1: CodeMirror.EditorFromTextArea,
    data2: CodeMirror.EditorFromTextArea;
  let result = "";
  let source = "Data1";
  let mode = "comm";
  let detail = true;
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
    }, 600);
  });

  const analyze = (operation: string) => {
    const d1 = preprocess(data1.getValue());
    const d2 = preprocess(data2.getValue());
    switch (operation) {
      case "chkDuplicates":
      case "rmDuplicates":
      case "chkConsecutive":
        if (source == "Data1" && !d1.length) {
          result = "Data1 is empty.\nPlease enter something...";
          return;
        } else if (source == "Data2" && !d2.length) {
          result = "Data2 is empty.\nPlease enter something...";
          return;
        }
        break;
      default:
        if (!d1.length) {
          result = "Data1 is empty.\nPlease enter something...";
          return;
        } else if (!d2.length) {
          result = "Data2 is empty.\nPlease enter something...";
          return;
        }
    }
    loading = true;
    const process = processing();
    const start = new Date().getTime();
    let output = "";
    let r: string[];
    switch (operation) {
      case "chkDuplicates":
        let d: { [k: string]: number };
        if (source == "Data1") d = new chkDuplicates(d1).run();
        else d = new chkDuplicates(d2).run();
        if (!Object.keys(d).length)
          output = `${source} has no duplicate value.`;
        else {
          if (detail)
            output =
              `Duplicate values found in ${source}.\n` +
              format(
                Object.keys(d).length,
                Object.keys(d).map((key) => `${key} appears ${d[key]} times.`)
              );
          else output = Object.keys(d).join("\n");
        }
        break;
      case "rmDuplicates":
        if (source == "Data1") r = new rmDuplicates(d1).run();
        else r = new rmDuplicates(d2).run();
        output = r.join("\n");
        break;
      case "chkConsecutive":
        if (source == "Data1") r = new chkConsecutive(d1).run();
        else r = new chkConsecutive(d2).run();
        if (!r.length) output = `${source} contains consecutive numbers.`;
        else if (r.length == 1 && r[0] == "!Error!")
          output = `Error!\n${source} contains non-numeric value. Please check!`;
        else
          output = `${source} is not consecutive.
\nThe following numbers are missing:\n${r.join("\n")}`;
        break;
      case "compare":
        if (mode == "comm") {
          r = new compareComm(d1, d2).run();
          if (!r.length) output = "Two data contain no common value.";
          else
            output = `Common values found between two data.
${format(r.length, r)}`;
        } else {
          const r1 = new compareDiff(d1, d2, ignoreDuplicates).run();
          const r2 = new compareDiff(d2, d1, ignoreDuplicates).run();
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
        output = new diff(d1.join("\n") + "\n", d2.join("\n") + "\n")
          .run()
          .replace(`${"=".repeat(67)}\n`, "");
    }
    clearInterval(process);
    if (
      operation == "rmDuplicates" ||
      (operation == "chkDuplicates" && !detail)
    )
      result = output;
    else result = output + `\n\nDuration for process: ${Date.now() - start}ms`;
    loading = false;
  };

  const clear = () => {
    data1.setValue("");
    data2.setValue("");
    result = "";
  };

  const swap = () => {
    const data = data1.getValue();
    data1.setValue(data2.getValue());
    data2.setValue(data);
  };

  const copy = async () => {
    if (result.trim() !== "")
      if (navigator.clipboard) {
        await navigator.clipboard.writeText(result.trim());
        alert("Text has been copied to clipboard.");
      } else
        alert("This function requires a secure origin. (HTTPS or localhost)");
  };

  const processing = () => {
    return setInterval(() => {
      const s = result.split("Processing");
      let dots = s.length >= 2 ? s[1].length : 0;
      if (dots < 3) dots++;
      else dots -= 3;
      result = "Processing" + ".".repeat(dots);
    }, 200);
  };
</script>

<svelte:window
  on:beforeunload={() => {
    localStorage.setItem("data1", data1.getValue());
    localStorage.setItem("data2", data2.getValue());
  }}
/>

<header class="navbar navbar-expand navbar-light flex-column flex-md-row">
  <a
    class="navbar-brand text-primary m-0 mr-md-3"
    href="/"
    style="font-size:24px"
  >
    Simple Data Analysis
  </a>
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
        on:click={() => analyze("chkDuplicates")}
        type="button"
        class="btn btn-primary btn-block"
        disabled={loading}
      >
        Check Duplicates
      </button>
      <div class="d-flex justify-content-around">
        <div>
          <input type="checkbox" bind:checked={detail} id="detail" />
          <label class="m-0" for="detail">Show Detail</label>
        </div>
      </div>
      <button
        on:click={() => analyze("rmDuplicates")}
        type="button"
        class="btn btn-primary btn-block"
        disabled={loading}
      >
        Remove Duplicates
      </button>
      <button
        on:click={() => analyze("chkConsecutive")}
        type="button"
        class="btn btn-primary btn-block"
        disabled={loading}
      >
        Check Consecutive
      </button>
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
        on:click={() => analyze("compare")}
        type="button"
        class="btn btn-primary btn-block"
        disabled={loading}
      >
        Cross Compare
      </button>
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
            disabled={mode == "comm"}
            id="ignore_duplicates"
          />
          <label class="m-0" for="ignore_duplicates">Ignore Duplicates</label>
        </div>
      </div>
      <br />
      <button
        on:click={() => analyze("diff")}
        type="button"
        class="btn btn-primary btn-block"
        disabled={loading}
      >
        Diff
      </button>
      <br />
      <br />
      <button
        on:click={copy}
        type="button"
        class="btn btn-primary btn-block"
        disabled={loading}
      >
        Copy Result
      </button>
      <br />
      <br />
      <button
        on:click={swap}
        type="button"
        class="btn btn-primary btn-block"
        disabled={loading}
      >
        {@html "Data1<=>Data2"}
      </button>
      <br />
      <br />
      <button
        on:click={clear}
        type="button"
        class="btn btn-primary btn-block"
        disabled={loading}
      >
        Clear
      </button>
    </div>
    <div class="col-4">
      <label for="result"> Result </label>
      <textarea class="form-control" id="result" bind:value={result} readonly />
    </div>
  </div>
</div>

<style>
  .navbar {
    height: 80px;
  }

  .container-fluid {
    position: fixed;
    height: calc(100% - 80px);
  }

  .row {
    height: 100%;
  }

  .row > div {
    height: 100%;
  }

  input {
    margin: 0 5px;
  }

  textarea {
    resize: none;
    font-family: monospace !important;
    height: calc(100% - 47px) !important;
    border: 1px solid #007bff !important;
  }

  :global(.CodeMirror) {
    height: calc(100% - 47px);
    border: 1px solid #007bff;
    border-radius: 0.25em;
    cursor: text;
  }

  :global(.CodeMirror-focused) {
    box-shadow: 0 0 0 0.2rem rgba(0, 123, 255, 0.25);
    outline: 0;
  }

  :global(.CodeMirror-placeholder) {
    color: #999 !important;
  }
</style>
