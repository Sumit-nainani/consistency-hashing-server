const ringSize = 100;
const radius = 1.4;
const angles = Array.from(
  { length: ringSize },
  (_, i) => -(2 * Math.PI * i) / ringSize
);

const podRowMap = new Map();
const clientRowMap = new Map();

const socket = new WebSocket("ws://localhost:8888/ws-client");

let pods = [];
let clients = [];
let podItems = [];
let clientItems = [];

function drawRing(pods = [], clients = [], podItems = [], clientItems = []) {
  const layout = {
    xaxis: { visible: false },
    yaxis: { visible: false, scaleanchor: "x", scaleratio: 1 },
    width: 1000,
    height: 1000,
    margin: { l: 10, r: 10, t: 10, b: 10 },
    plot_bgcolor: "#fafafa",
    images: [],
  };

  const ring = {
    x: angles.map((a) => radius * Math.cos(a)),
    y: angles.map((a) => radius * Math.sin(a)),
    mode: "lines",
    line: { color: "red" },
    type: "scatter",
    showlegend: false,
  };

  const labelMarkers = {
    x: angles.map((a) => radius * Math.cos(a)),
    y: angles.map((a) => radius * Math.sin(a)),
    mode: "markers+text",
    marker: { size: 22, color: "#daeff0" },
    text: Array.from({ length: ringSize }, (_, i) => String(i)),
    textfont: { size: 10, color: "black" },
    textposition: "middle center",
    type: "scatter",
    showlegend: false,
  };

  const dockerIconSrc = document.getElementById("docker-icon").src;
  const clientIconSrc = document.getElementById("client-icon").src;

  for (let pos of pods) {
    let angle = angles[pos % ringSize];
    let x = radius * Math.cos(angle);
    let y = radius * Math.sin(angle);
    layout.images.push({
      source: dockerIconSrc,
      xref: "x",
      yref: "y",
      x: x,
      y: y,
      sizex: 0.2,
      sizey: 0.2,
      xanchor: "center",
      yanchor: "middle",
      layer: "above",
    });
  }

  for (let pos of clients) {
    let angle = angles[pos % ringSize];
    let x = radius * Math.cos(angle);
    let y = radius * Math.sin(angle);
    layout.images.push({
      source: clientIconSrc,
      xref: "x",
      yref: "y",
      x: x,
      y: y,
      sizex: 0.12,
      sizey: 0.12,
      xanchor: "center",
      yanchor: "middle",
      layer: "above",
    });
  }

  Plotly.newPlot("ring", [ring, labelMarkers], layout);
}

function updateTables(podInfoList, clientInfoList) {
  document.getElementById("rebalancing-message").style.display = "block";
  setTimeout(() => {
    document.getElementById("rebalancing-message").style.display = "none";
  }, 1500);

  const podTable = document.querySelector("#pod-table tbody");
  const clientTable = document.querySelector("#client-table tbody");

  const currentPodHashes = new Set();
  podInfoList.forEach((pod) => {
    const hash = pod.nodeMetaData.nodeHash - 1;
    currentPodHashes.add(hash);
    if (!podRowMap.has(hash)) {
      const row = document.createElement("tr");
      row.innerHTML = `
              <td>${pod.nodeMetaData.nodeName}</td>
              <td>${pod.nodeMetaData.nodeIp}</td>
              <td>${hash}</td>`;
      podTable.appendChild(row);
      podRowMap.set(hash, row);
    }
  });

  for (const [hash, row] of podRowMap.entries()) {
    if (!currentPodHashes.has(hash)) {
      podTable.removeChild(row);
      podRowMap.delete(hash);
    }
  }

  const currentClientHashes = new Set();
  clientInfoList.forEach((client) => {
    const requestHash = client.requestMetaData.requestHash - 1;
    currentClientHashes.add(requestHash);
    if (clientRowMap.has(requestHash)) {
      const row = clientRowMap.get(requestHash);
      row.cells[0].textContent = requestHash;
      row.cells[1].textContent = client.requestMetaData.assignedNodeName;
      row.cells[2].textContent = client.requestMetaData.assignedNodeIp;
      row.cells[3].textContent = client.requestMetaData.assignedNodeHash - 1;
    } else {
      const newRow = document.createElement("tr");
      newRow.innerHTML = `
              <td>${requestHash}</td>
              <td>${client.requestMetaData.assignedNodeName}</td>
              <td>${client.requestMetaData.assignedNodeIp}</td>
              <td>${client.requestMetaData.assignedNodeHash - 1}</td>`;
      clientTable.appendChild(newRow);
      clientRowMap.set(requestHash, newRow);
    }
  });
}

socket.onmessage = function (event) {
  try {
    const raw = JSON.parse(event.data);
    let items = [];

    if (Array.isArray(raw.item)) {
      items = raw.item;
    } else if (Array.isArray(raw)) {
      items = raw;
    } else if (raw.type && (raw.nodeMetaData || raw.requestMetaData)) {
      items = [raw];
    } else {
      console.warn("Unknown message format:", raw);
    }

    items.forEach((item) => {
      if (item.type === "pod") {
        const hash = item.nodeMetaData?.nodeHash - 1 || 64;
        if (item.action === "add") {
          if (!pods.includes(hash)) {
            pods.push(hash);
            podItems.push(item);
          }
        } else if (item.action === "remove") {
          const index = pods.indexOf(hash);
          if (index !== -1) {
            pods.splice(index, 1);
            podItems.splice(index, 1);
          }
        }
      } else if (item.type === "client") {
        clients.push(item.requestMetaData.requestHash - 1);
        clientItems.push(item);
      }
    });

    drawRing(pods, clients, podItems, clientItems);
    updateTables(podItems, clientItems);
  } catch (err) {
    console.error("WebSocket message handling failed:", err);
  }
};
