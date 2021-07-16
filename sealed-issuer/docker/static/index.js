const endpoint = "/api/v1/seal";

let kvPairs = [{ k: "", v: "" }];

function isYAML() {
  const urlParams = new URLSearchParams(window.location.search);
  return urlParams.get("format") && "YAML" === urlParams.get("format");
}

if (isYAML()) {
  document.getElementById("properties-block").style.display = "none";
  document.getElementById("change-format-yaml").style.display = "none";

  document.getElementById("yaml-block").style.display = "block";
  document.getElementById("change-format-properties").style.display = "inline";
}

function parseInput() {
  if (isYAML()) {
    const yaml = document.getElementById("yaml").value;
    if (!yaml) {
      alert("The input is empty!");
      return;
    }
    sendSealRequest({
      yaml: yaml
    });
  } else {
    const name = document.getElementById("name").value;
    const namespace = document.getElementById("namespace").value;
    if (!name) {
      alert("Secret name is missing!");
      return;
    } else if (!namespace) {
      alert("Namespace is missing!");
      return;
    }

    const nameRegex = /[a-z0-9]([-a-z0-9]*[a-z0-9])?/;
    if (!nameRegex.test(name)) {
      alert("Invalid secret name, must satisfy Kubernetes naming rules!");
      return;
    } else if (!nameRegex.test(namespace)) {
      alert("Invalid namespace, must satisfy Kubernetes naming rules!");
      return;
    }
    sendSealRequest({
      kvPairs: kvPairs,
      name: name,
      namespace: namespace
    });
  }
}

async function sendSealRequest(requestBody) {
  const response = await fetch(endpoint, {
    method: "POST",
    body: JSON.stringify(requestBody),
    headers: {
      "Content-Type": "application/json"
    }
  });
  const responseBody = await response.json();
  if (responseBody.errorMessage) {
    alert("server side error: " + responseBody.errorMessage);
    return;
  }
  displayOutput(responseBody);
}

function displayOutput(response) {
  document.getElementById("output").value = response.fileContent;
  setupDownloadButton(response.fileContent, response.fileName);
}

function copy() {
  const copyTextArea = document.getElementById("output");
  copyTextArea.focus();
  copyTextArea.select();
  try {
    const selection = document.getSelection();
    const range = document.createRange();
    range.selectNode(copyTextArea);
    selection.removeAllRanges();
    selection.addRange(range);

    const successful = document.execCommand("copy");
    if (!successful) throw new Error("Unable to copy");
  } catch (err) {
    alert(err.message);
  } finally {
    window.getSelection().removeAllRanges();
  }
}

function setupDownloadButton(content, filename) {
  const data = "data:text/json;charset=utf-8," + encodeURIComponent(content);
  const a = document.getElementById("download-anchor");
  a.setAttribute("href", data);
  a.setAttribute("download", filename);
}

function download() {
  document.getElementById("download-anchor").click();
}

function clearKVListElement() {
  document.getElementById("kv-list").innerHTML = "";
}

function renderKVPairs() {
  clearKVListElement();
  for (let i = 0; i < kvPairs.length; i++) {
    let container = document.createElement("div");
    container.id = "kv" + (i + 1);
    container.className = "kv-group";
    let span = document.createElement("span");
    span.className = "kv-label";
    span.textContent = "key: " + (i + 1);
    let del = document.createElement("button");
    del.id = "del" + (i + 1);
    del.className = "kv-delete";
    del.onclick = deleteKVPair;
    del.textContent = "delete";
    let key = document.createElement("input");
    key.id = "key" + (i + 1);
    key.type = "text";
    key.placeholder = "key";
    key.value = kvPairs[i].k;
    key.oninput = handleKeyUpdate;
    let val = document.createElement("textarea");
    val.id = "value" + (i + 1);
    val.placeholder = "value";
    val.value = kvPairs[i].v;
    val.oninput = handleValueUpdate;
    container.appendChild(span);
    container.appendChild(del);
    container.appendChild(key);
    container.appendChild(val);
    document.getElementById("kv-list").appendChild(container);
  }
}

function addKVPair() {
  kvPairs.push({ k: "", v: "" });
  renderKVPairs();
}

function deleteKVPair(e) {
  let i = parseInt(this.id.replace("del", "")) - 1;
  kvPairs.splice(i, 1);
  renderKVPairs();
}

function handleKeyUpdate(e) {
  let i = parseInt(this.id.replace("key", "")) - 1;
  kvPairs[i].k = e.target.value;
}

function handleValueUpdate(e) {
  let i = parseInt(this.id.replace("value", "")) - 1;
  kvPairs[i].v = e.target.value;
}

renderKVPairs();
