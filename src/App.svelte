<script lang="ts">
  import Handsontable from "handsontable/base";
  import {
    AutoColumnSize,
    ContextMenu,
    CopyPaste,
    UndoRedo,
    registerPlugin,
  } from "handsontable/plugins";
  import "handsontable/styles/handsontable.min.css";
  import "handsontable/styles/ht-theme-main.min.css";
  import { onMount } from "svelte";
  import * as sda from "./sda";
  import { format } from "./utils";

  registerPlugin(AutoColumnSize);
  registerPlugin(ContextMenu);
  registerPlugin(CopyPaste);
  registerPlugin(UndoRedo);

  const rowHeaderWidth = 64;

  let table1: Handsontable;
  let table2: Handsontable;
  let result = $state("");
  let source = $state("Data1");
  let mode = $state("comm");
  let detail = $state(true);
  let ignoreDuplicates = $state(true);
  let loading = $state(false);

  const create_table = (elementId: string, name: string, key: string) => {
    let data: any[][] | undefined;
    try {
      data = JSON.parse(localStorage.getItem(key) || "");
    } catch (e) {
      data = undefined;
    }
    const table = new Handsontable(document.getElementById(elementId)!, {
      data,
      colHeaders: [name],
      colWidths() {
        return Math.floor(window.innerWidth * 0.25 - rowHeaderWidth - 24);
      },
      contextMenu: [
        "row_above",
        "row_below",
        "---------",
        "remove_row",
        "---------",
        "undo",
        "redo",
        "---------",
        "copy",
        "cut",
      ],
      height() {
        return window.innerHeight - 80 - 16;
      },
      maxCols: 1,
      minSpareRows: 1,
      rowHeaders: true,
      rowHeaderWidth,
      startRows: 1,
      tabMoves: { row: 1, col: 0 },
      themeName: "ht-theme-main",
      licenseKey: "non-commercial-and-evaluation",
    });
    table.addHook("beforePaste", () =>
      table.updateSettings({ readOnly: true }),
    );
    table.addHook("afterPaste", (data) => {
      table.updateSettings({ readOnly: false });
      table.deselectCell();
      for (let i = data.length - 1; i >= 0; i--)
        if (data[i][0] !== "") {
          table.loadData(data.slice(0, i + 1));
          return;
        }
      table.loadData([[""]]);
    });
    return table;
  };

  onMount(() => {
    table1 = create_table("Table1", "Data1", "data1");
    table2 = create_table("Table2", "Data2", "data2");
  });

  const getData = (table: Handsontable) => {
    return table.getDataAtCol(0).filter((i) => i !== null);
  };

  const analyze = (operation: string) => {
    const d1 = getData(table1).filter((i) => i != "");
    const d2 = getData(table2).filter((i) => i != "");
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
        if (source == "Data1") d = new sda.chkDuplicates(d1).run();
        else d = new sda.chkDuplicates(d2).run();
        if (!Object.keys(d).length)
          output = `${source} has no duplicate value.`;
        else {
          if (detail)
            output =
              `Duplicate values found in ${source}.\n` +
              format(
                Object.keys(d).length,
                Object.keys(d).map((key) => `${key} appears ${d[key]} times.`),
              );
          else output = Object.keys(d).join("\n");
        }
        break;
      case "rmDuplicates":
        if (source == "Data1") r = new sda.rmDuplicates(d1).run();
        else r = new sda.rmDuplicates(d2).run();
        output = r.join("\n");
        break;
      case "chkConsecutive":
        if (source == "Data1") r = new sda.chkConsecutive(d1).run();
        else r = new sda.chkConsecutive(d2).run();
        if (!r.length) output = `${source} contains consecutive numbers.`;
        else if (r.length == 1 && r[0] == "!Error!")
          output = `Error!\n${source} contains non-numeric value. Please check!`;
        else
          output = `${source} is not consecutive.
\nThe following numbers are missing:\n${r.join("\n")}`;
        break;
      case "compare":
        if (mode == "comm") {
          r = new sda.compareComm(d1, d2).run();
          if (!r.length) output = "Two data contain no common value.";
          else
            output = `Common values found between two data.
${format(r.length, r)}`;
        } else {
          const r1 = new sda.compareDiff(d1, d2, ignoreDuplicates).run();
          const r2 = new sda.compareDiff(d2, d1, ignoreDuplicates).run();
          if (r1.length + r2.length == 0) {
            output = "Data1 is same as Data2.";
          } else if (!r1.length) {
            output = `Data2 completely contains Data1.\n\nData2 is more than Data1
${format(r2.length, r2)}`;
          } else if (!r2.length) {
            output = `Data1 completely contains Data2.\n\nData1 is more than Data2\n${format(
              r1.length,
              r1,
            )}`;
          } else {
            output = `Two files have inconsistent content.
\nData1 is more than Data2\n${format(r1.length, r1)}
\nData2 is more than Data1\n${format(r2.length, r2)}`;
          }
        }
        break;
      case "diff":
        output = new sda.diff(d1.join("\n") + "\n", d2.join("\n") + "\n")
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
    if (!confirm("Clear all data?")) return;
    table1.loadData([[""]]);
    table2.loadData([[""]]);
    result = "";
  };

  const swap = () => {
    const data = table1.getData();
    table1.loadData(table2.getData());
    table2.loadData(data);
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
  onbeforeunload={() => {
    localStorage.setItem("data1", JSON.stringify(table1.getData()));
    localStorage.setItem("data2", JSON.stringify(table2.getData()));
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
    <div id="Table1" class="col-3"></div>
    <div id="Table2" class="col-3 pl-0"></div>
    <div class="col-2 p-0 pt-5">
      <button
        onclick={() => analyze("chkDuplicates")}
        type="button"
        class="btn btn-primary w-100"
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
        onclick={() => analyze("rmDuplicates")}
        type="button"
        class="btn btn-primary w-100"
        disabled={loading}
      >
        Remove Duplicates
      </button>
      <button
        onclick={() => analyze("chkConsecutive")}
        type="button"
        class="btn btn-primary w-100"
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
        onclick={() => analyze("compare")}
        type="button"
        class="btn btn-primary w-100"
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
        onclick={() => analyze("diff")}
        type="button"
        class="btn btn-primary w-100"
        disabled={loading}
      >
        Diff
      </button>
      <br />
      <br />
      <button
        onclick={copy}
        type="button"
        class="btn btn-primary w-100"
        disabled={loading}
      >
        Copy Result
      </button>
      <br />
      <br />
      <button
        onclick={swap}
        type="button"
        class="btn btn-primary w-100"
        disabled={loading}
      >
        {@html "Data1<=>Data2"}
      </button>
      <br />
      <br />
      <button
        onclick={clear}
        type="button"
        class="btn btn-danger w-100"
        disabled={loading}
      >
        Clear
      </button>
    </div>
    <div class="col-4">
      <label for="result">Result</label>
      <textarea class="form-control" id="result" bind:value={result} readonly
      ></textarea>
    </div>
  </div>
</div>

<style>
  .navbar {
    height: 80px;
    padding-left: 1rem;
    padding-right: 1rem;
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

  .btn + .btn {
    margin-top: 0.5rem;
  }

  #Table1,
  #Table2 {
    outline: 1px solid var(--ht-border-color);
    outline-offset: -1px;
  }

  textarea {
    resize: none;
    font-family: monospace !important;
    height: calc(100% - 36px) !important;
    border: 1px solid #007bff !important;
  }
</style>
