const STATUS_URL = "http://localhost:8080/status";

// store previous error state to avoid repeated logs
let backendDown = false;
let lastLogTime = null;


async function fetchStatus() {
  try {
    const res = await fetch(STATUS_URL);
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const status = await res.json();

    backendDown = false; // reset error state

    // --- decision box ---
    document.getElementById("algo").textContent = status.current_algo || "-";
    document.getElementById("reason").textContent =
      status.adaptive_reason || "-";
    document.getElementById("selected").textContent =
      status.selected_backend || "-";

    // --- backend table ---
    const tbody = document.querySelector("#backend-table tbody");
    tbody.innerHTML = "";
    status.backends.forEach((b) => {
      const row = document.createElement("tr");
      row.className = b.Alive ? "alive" : "down";
      row.innerHTML = `
        <td>${b.Address}</td>
        <td>${b.Alive}</td>
        <td>${b.ActiveConns}</td>
        <td>${b.Latency}µs</td>
        <td>${b.ErrorCount}</td>
      `;
      tbody.appendChild(row);
    });

    // --- decision logs ---
    const logBox = document.getElementById("log-box");

    if (status.decision_log && status.decision_log.length > 0) {
      status.decision_log.forEach((entry) => {
        // only add logs newer than last rendered
        if (!lastLogTime || new Date(entry.time) > lastLogTime) {
          const div = document.createElement("div");
          div.className = `log-entry ${entry.algo}`;
          div.textContent = `[${new Date(entry.time).toLocaleTimeString()}] ${
            entry.algo
          } → ${entry.backend} (${entry.reason})`;

          logBox.appendChild(div);
          lastLogTime = new Date(entry.time);
        }
      });

    }

  } catch (err) {
    console.error("Failed to fetch status:", err);
    // only show one error until backend recovers
    if (!backendDown) {
      const logBox = document.getElementById("log-box");
      const div = document.createElement("div");
      div.className = "log-entry error";
      div.textContent = `[${new Date().toLocaleTimeString()}] ERROR: Cannot reach backend`;
      logBox.appendChild(div);
      backendDown = true;
    }

    // clear decision box and backend table
    document.getElementById("algo").textContent = "-";
    document.getElementById("reason").textContent = "-";
    document.getElementById("selected").textContent = "-";
    document.querySelector("#backend-table tbody").innerHTML = "";
  }
}

fetchStatus();
setInterval(fetchStatus, 1000);
