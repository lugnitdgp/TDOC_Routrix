let backends = [];
let requests = [];
let rrIndex = 0;

// Canvas
const canvas = document.getElementById('simulation');
const ctx = canvas.getContext('2d');

ctx.font = '14px Inter, system-ui';
ctx.textBaseline = 'middle';

const METRICS_URL = 'http://localhost:8080/metrics';
const LOAD_BALANCER_X = 50;
const LOAD_BALANCER_Y = canvas.height / 2;

// ---------------- BACKEND POSITIONING ----------------
function setBackendPositions() {
  const n = backends.length;
  backends.forEach((b, i) => {
    b.x = 600;
    b.y = (canvas.height / (n + 1)) * (i + 1) - 30;
    b.vizX = b.x + 70;
    b.vizY = b.y + 30;
  });
}

// ---------------- UTILS ----------------
function roundedRect(x, y, w, h, r) {
  ctx.beginPath();
  ctx.moveTo(x + r, y);
  ctx.lineTo(x + w - r, y);
  ctx.quadraticCurveTo(x + w, y, x + w, y + r);
  ctx.lineTo(x + w, y + h - r);
  ctx.quadraticCurveTo(x + w, y + h, x + w - r, y + h);
  ctx.lineTo(x + r, y + h);
  ctx.quadraticCurveTo(x, y + h, x, y + h - r);
  ctx.lineTo(x, y + r);
  ctx.quadraticCurveTo(x, y, x + r, y);
  ctx.closePath();
}

// ---------------- DRAW LOOP ----------------
function draw() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);

  // ----- LOAD BALANCER -----
  ctx.shadowColor = '#3b82f6';
  ctx.shadowBlur = 15;
  ctx.fillStyle = '#2563eb';
  roundedRect(20, LOAD_BALANCER_Y - 35, 90, 70, 14);
  ctx.fill();
  ctx.shadowBlur = 0;

  ctx.fillStyle = '#ffffff';
  ctx.font = '600 15px Inter';
  ctx.fillText('LOAD', 45, LOAD_BALANCER_Y - 6);
  ctx.fillText('BALANCER', 30, LOAD_BALANCER_Y + 12);

  // ----- BACKENDS -----
  backends.forEach(b => {
    ctx.shadowColor = b.Alive ? '#22c55e' : '#ef4444';
    ctx.shadowBlur = 10;

    ctx.fillStyle = b.Alive ? '#16a34a' : '#dc2626';
    roundedRect(b.x, b.y, 160, 70, 14);
    ctx.fill();

    ctx.shadowBlur = 0;

    ctx.fillStyle = '#ffffff';
    ctx.font = '500 13px Inter';
    ctx.fillText(b.Address, b.x + 10, b.y + 18);

    ctx.font = '12px Inter';
    ctx.fillText(`Conn: ${b.ActiveConns}`, b.x + 10, b.y + 36);
    ctx.fillText(`Latency: ${b.Latency}µs`, b.x + 10, b.y + 52);
    ctx.fillText(`Errors: ${b.ErrorCount}`, b.x + 95, b.y + 36);
  });

  // ----- REQUESTS -----
  requests.forEach(r => {
    ctx.shadowColor = '#facc15';
    ctx.shadowBlur = 12;

    ctx.fillStyle = '#facc15';
    ctx.beginPath();
    ctx.arc(r.x, r.y, 6, 0, Math.PI * 2);
    ctx.fill();

    ctx.shadowBlur = 0;

    r.x += 7;
    r.y += (r.targetY - r.y) * 0.12;
  });

  // ----- CLEANUP -----
  requests = requests.filter(r => {
    if (r.x >= r.targetX) {
      r.backend.ActiveConns = Math.max(0, r.backend.ActiveConns - 1);
      return false;
    }
    return true;
  });

  requestAnimationFrame(draw);
}

// ---------------- LOAD BALANCER LOGIC ----------------
// ... (Keep existing variables: backends, requests, etc.)

function pickBackend() {
  const alive = backends.filter(b => b.Alive);
  if (!alive.length) return null;

  const algo = document.getElementById('algo').value;

  // NEW: No Load Balancer Logic
  if (algo === 'none') {
    // Force all traffic to the first server, even if it's overloaded
    return alive[0]; 
  }

  if (algo === 'roundrobin') {
    const b = alive[rrIndex % alive.length];
    rrIndex++;
    return b;
  }

  if (algo === 'leastconnections') {
    return alive.reduce((a, b) => a.ActiveConns <= b.ActiveConns ? a : b);
  }

  return alive[Math.floor(Math.random() * alive.length)];
}

// ---------------- ENHANCED DRAW LOOP ----------------
function draw() {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  const algo = document.getElementById("algo").value;

  // Draw Background Grid (Aesthetic touch)
  ctx.strokeStyle = "rgba(255, 255, 255, 0.03)";
  for (let i = 0; i < canvas.width; i += 40) {
    ctx.beginPath();
    ctx.moveTo(i, 0);
    ctx.lineTo(i, canvas.height);
    ctx.stroke();
  }

  // ----- LOAD BALANCER -----
  // If "none", we draw it faded or with a "Bypass" label
  const lbOpacity = algo === "none" ? 0.2 : 1;
  ctx.globalAlpha = lbOpacity;

  ctx.shadowColor = "#3b82f6";
  ctx.shadowBlur = lbOpacity * 15;
  ctx.fillStyle = "#2563eb";
  roundedRect(20, LOAD_BALANCER_Y - 35, 90, 70, 14);
  ctx.fill();

  ctx.fillStyle = "#ffffff";
  ctx.font = "bold 13px Inter";
  ctx.fillText(
    algo === "none" ? "BYPASSED" : "BALANCER",
    32,
    LOAD_BALANCER_Y + 5
  );
  ctx.globalAlpha = 1;
  ctx.shadowBlur = 0;

  // ----- BACKENDS -----
  // --- UPDATED DRAW SECTION FOR BACKENDS ---
  backends.forEach((b) => {
    // Use a darker card color so white text stands out
    ctx.fillStyle = b.Alive ? "#1e293b" : "#450a0a";
    ctx.strokeStyle = b.Alive ? "#3b82f6" : "#ef4444";
    ctx.lineWidth = 2;
    roundedRect(b.x, b.y, 180, 85, 12); // Slightly larger card
    ctx.fill();
    ctx.stroke();

    // Address - Bold and White
    ctx.fillStyle = "#ffffff";
    ctx.font = "bold 14px Inter";
    ctx.fillText(b.Address, b.x + 15, b.y + 25);

    // Active Connections - Highlighted with a color
    ctx.fillStyle = "#60a5fa"; // Light blue
    ctx.font = "600 13px Inter";
    ctx.fillText(`Active: ${b.ActiveConns}`, b.x + 15, b.y + 48);

    // Latency - Different color for visual separation
    ctx.fillStyle = "#94a3b8"; // Muted grey
    ctx.font = "500 12px Inter";
    ctx.fillText(`Latency: ${b.Latency}ms`, b.x + 15, b.y + 68);
  });

  // ----- REQUESTS -----
  requests.forEach((r) => {
    // If no LB, the request originates from the far left instead of the LB
    const startX = algo === "none" ? -10 : LOAD_BALANCER_X + 90;

    ctx.shadowColor = "#facc15";
    ctx.shadowBlur = 8;
    ctx.fillStyle = "#facc15";
    ctx.beginPath();
    ctx.arc(r.x, r.y, 5, 0, Math.PI * 2);
    ctx.fill();
    ctx.shadowBlur = 0;

    // Movement
    r.x += 6;
    r.y += (r.targetY - r.y) * 0.08;
  });

  // (Keep Cleanup Logic)
  requests = requests.filter((r) => r.x < r.targetX);
  requestAnimationFrame(draw);
}

// ---------------- REQUEST SIMULATION ----------------
function sendRequest() {
  const backend = pickBackend();
  if (!backend) return;

  backend.ActiveConns++;

  requests.push({
    x: LOAD_BALANCER_X + 90,
    y: LOAD_BALANCER_Y,
    targetX: backend.vizX,
    targetY: backend.vizY,
    backend
  });
}

// ---------------- METRICS FETCH ----------------
async function fetchMetrics() {
  try {
    const res = await fetch(METRICS_URL);
    const data = await res.json();

    if (!backends.length) {
      backends = data.map(b => ({ ...b }));
      setBackendPositions();
      return;
    }

    data.forEach(metricBackend => {
      const local = backends.find(
        b => b.Address === metricBackend.Address
      );
      if (!local) return;

      local.Alive = metricBackend.Alive;
      local.Latency = metricBackend.Latency;
      local.ErrorCount = metricBackend.ErrorCount;
    });

    // Traffic rate
    if (Math.random() > 0.35) sendRequest();

  } catch (err) {
    console.error('Metrics fetch failed:', err);
  }
}

// ---------------- START ----------------
setInterval(fetchMetrics, 1000);
draw();
