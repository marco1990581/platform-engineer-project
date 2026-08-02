let authorizationHeader = "";

const endpoints = {
  health: "/health",
  hostname: "/hostname",
  memory: "/memory",
  uptime: "/uptime",
  network: "/network",
  filesystem: "/filesystem",
};

function basicAuthorization(username, password) {
  const credentials = new TextEncoder().encode(`${username}:${password}`);
  let binaryCredentials = "";

  for (const byte of credentials) {
    binaryCredentials += String.fromCharCode(byte);
  }

  return `Basic ${btoa(binaryCredentials)}`;
}

async function authenticate(username, password) {
  const candidateAuthorizationHeader = basicAuthorization(username, password);

  const response = await fetch("/hostname", {
    headers: {
      Authorization: candidateAuthorizationHeader,
    },
  });

  if (!response.ok) {
    throw new Error("Authentication failed");
  }

  authorizationHeader = candidateAuthorizationHeader;
}

function authenticatedFetch(endpoint) {
  if (!authorizationHeader) {
    throw new Error("Authentication is required");
  }

  return fetch(endpoint, {
    headers: {
      Authorization: authorizationHeader,
    },
  });
}

function escapeHtml(value) {
  // Runtime values are rendered into HTML templates below, so escape them first.
  return String(value).replace(/[&<>"']/g, (character) => ({
    "&": "&amp;",
    "<": "&lt;",
    ">": "&gt;",
    '"': "&quot;",
    "'": "&#039;",
  })[character]);
}

function formatBytes(bytes) {
  if (bytes === 0) {
    return "0 B";
  }

  const units = ["B", "KiB", "MiB", "GiB", "TiB"];
  const exponent = Math.min(Math.floor(Math.log(bytes) / Math.log(1024)), units.length - 1);
  const value = bytes / 1024 ** exponent;

  return `${value.toFixed(value >= 10 || exponent === 0 ? 0 : 1)} ${units[exponent]}`;
}

function formatDuration(seconds) {
  const totalSeconds = Math.floor(seconds);
  const days = Math.floor(totalSeconds / 86400);
  const hours = Math.floor((totalSeconds % 86400) / 3600);
  const minutes = Math.floor((totalSeconds % 3600) / 60);

  if (days > 0) {
    return `${days}d ${hours}h`;
  }

  if (hours > 0) {
    return `${hours}h ${minutes}m`;
  }

  return `${minutes}m`;
}

function setStatus(healthy) {
  const status = document.querySelector("#status");
  status.textContent = healthy ? "API healthy" : "API unavailable";
  status.className = `status ${healthy ? "status-healthy" : "status-error"}`;
}

function showLoginError(message) {
  const error = document.querySelector("#login-error");
  error.textContent = message;
  error.hidden = false;
}

function showLogin() {
  authorizationHeader = "";
  document.querySelector("#dashboard").hidden = true;
  document.querySelector("#login-screen").hidden = false;
  document.querySelector("#password").value = "";
  document.querySelector("#username").focus();
}

function renderFilesystems(filesystems) {
  document.querySelector("#filesystem-count").textContent =
    `${filesystems.length} visible mounts`;

  const rows = filesystems.slice(0, 8).map((filesystem) => `
    <tr>
      <td>${escapeHtml(filesystem.mount_point)}</td>
      <td>${escapeHtml(filesystem.filesystem_type)}</td>
      <td>${formatBytes(filesystem.used_bytes)} / ${formatBytes(filesystem.total_bytes)}</td>
      <td>${formatBytes(filesystem.available_bytes)}</td>
      <td>${filesystem.usage_percent.toFixed(1)}%</td>
    </tr>
  `);

  document.querySelector("#filesystems").innerHTML = rows.join("");
}

function renderNetwork(interfaces) {
  const cards = interfaces.map((networkInterface) => {
    const state = networkInterface.up ? "up" : "down";
    const addresses = networkInterface.addresses.length
      ? networkInterface.addresses.map(escapeHtml).join("<br>")
      : "No addresses";

    return `
      <article class="interface-card">
        <div class="interface-heading">
          <strong>${escapeHtml(networkInterface.name)}</strong>
          <span class="interface-state interface-${state}">${state}</span>
        </div>
        <p>${addresses}</p>
        <p class="detail">MTU ${networkInterface.mtu}</p>
      </article>
    `;
  });

  document.querySelector("#network").innerHTML = cards.join("");
}

async function loadDashboard() {
  const dashboardError = document.querySelector("#dashboard-error");
  dashboardError.hidden = true;

  try {
    const responses = await Promise.all(
      Object.values(endpoints).map(async (endpoint) => {
        const response = await authenticatedFetch(endpoint);
        if (!response.ok) {
          throw new Error(`${endpoint} returned ${response.status}`);
        }

        return response.json();
      }),
    );

    const [health, hostname, memory, uptime, network, filesystem] = responses;
    const addressCount = network.interfaces.reduce(
      (count, networkInterface) => count + networkInterface.addresses.length,
      0,
    );

    setStatus(health.status === "ok");
    document.querySelector("#hostname").textContent = hostname.hostname;
    document.querySelector("#memory").textContent = formatBytes(memory.available_kb * 1024);
    document.querySelector("#memory-detail").textContent =
      `${formatBytes(memory.free_kb * 1024)} free of ${formatBytes(memory.total_kb * 1024)}`;
    document.querySelector("#uptime").textContent = formatDuration(uptime.seconds);
    document.querySelector("#interfaces").textContent = network.interfaces.length;
    document.querySelector("#addresses").textContent = `${addressCount} assigned addresses`;

    renderFilesystems(filesystem.filesystems);
    renderNetwork(network.interfaces);
  } catch (error) {
    setStatus(false);
    dashboardError.textContent = "Unable to load runtime information. Please try again.";
    dashboardError.hidden = false;
    console.error("Unable to load dashboard", error);
  }
}

document.querySelector("#login-form").addEventListener("submit", async (event) => {
  event.preventDefault();

  const username = document.querySelector("#username").value;
  const password = document.querySelector("#password").value;
  const loginButton = document.querySelector("#login-button");
  const loginError = document.querySelector("#login-error");

  loginError.hidden = true;
  loginButton.disabled = true;

  try {
    await authenticate(username, password);

    document.querySelector("#login-screen").hidden = true;
    document.querySelector("#dashboard").hidden = false;
    document.querySelector("#dashboard").focus();

    await loadDashboard();
  } catch {
    showLoginError("Invalid username or password.");
  } finally {
    loginButton.disabled = false;
  }
});

document.querySelector("#logout-button").addEventListener("click", showLogin);
