// Store last response for copy functions
let lastResponse = null;

// ── Theme ──────────────────────────────────────────────

function initTheme() {
  const stored = localStorage.getItem("theme");
  const prefersDark = window.matchMedia("(prefers-color-scheme: dark)").matches;
  const isDark = stored === "dark" || (!stored && prefersDark);
  document.documentElement.classList.toggle("dark", isDark);
  document.documentElement.classList.toggle("light", !isDark);
  updateThemeIcons(isDark);
}

function toggleTheme() {
  const isDark = document.documentElement.classList.toggle("dark");
  document.documentElement.classList.toggle("light", !isDark);
  localStorage.setItem("theme", isDark ? "dark" : "light");
  updateThemeIcons(isDark);
}

function updateThemeIcons(isDark) {
  const sun = document.getElementById("sunIcon");
  const moon = document.getElementById("moonIcon");
  if (sun && moon) {
    sun.classList.toggle("hidden", !isDark);
    moon.classList.toggle("hidden", isDark);
  }
}

// ── Toggle Switch ──────────────────────────────────────

function toggleSwitch(btn) {
  if (btn.disabled) return;
  const isOn = btn.getAttribute("aria-checked") === "true";
  btn.setAttribute("aria-checked", String(!isOn));
  const knob = btn.querySelector("span");
  if (!isOn) {
    btn.classList.remove("bg-zinc-200", "dark:bg-zinc-700");
    btn.classList.add("bg-zinc-900", "dark:bg-zinc-100");
    knob.classList.remove("translate-x-0");
    knob.classList.add("translate-x-4");
  } else {
    btn.classList.remove("bg-zinc-900", "dark:bg-zinc-100");
    btn.classList.add("bg-zinc-200", "dark:bg-zinc-700");
    knob.classList.remove("translate-x-4");
    knob.classList.add("translate-x-0");
  }
}

function setSwitchState(btn, on) {
  const current = btn.getAttribute("aria-checked") === "true";
  if (current !== on) toggleSwitch(btn);
}

function getSwitchState(btn) {
  return btn.getAttribute("aria-checked") === "true";
}

function updateIncludeDefaultState() {
  const nameserverInput = document.getElementById("nameserverInput").value.trim();
  const btn = document.getElementById("includeDefault");
  const label = document.getElementById("includeDefaultLabel");
  const hasNameservers = nameserverInput.length > 0;

  if (hasNameservers) {
    btn.removeAttribute("disabled");
    btn.classList.add("cursor-pointer");
    label.classList.remove("opacity-50", "cursor-not-allowed", "text-zinc-400", "dark:text-zinc-500");
    label.classList.add("cursor-pointer", "text-zinc-600", "dark:text-zinc-400");
  } else {
    btn.setAttribute("disabled", "");
    btn.classList.remove("cursor-pointer");
    label.classList.add("opacity-50", "cursor-not-allowed", "text-zinc-400", "dark:text-zinc-500");
    label.classList.remove("cursor-pointer", "text-zinc-600", "dark:text-zinc-400");
  }
}

// Toggle includeDefault and auto re-lookup if results exist
function onIncludeDefaultToggle(btn) {
  if (btn.disabled) return;
  toggleSwitch(btn);
  if (lastResponse) {
    lookupHosts();
  }
}

// ── App Info ───────────────────────────────────────────

async function loadInfo() {
  try {
    const response = await fetch("/api/info");
    if (response.ok) {
      const data = await response.json();
      const ver = data.version || "dev";
      const footerVer = document.getElementById("footerVersion");
      if (footerVer) footerVer.textContent = `(${ver})`;
    }
  } catch {
    // ignore — footer will just show without version
  }
}

// ── DNS Lookup ─────────────────────────────────────────

async function lookupHosts() {
  const hostsInput = document.getElementById("hostsInput").value;
  const hostsArray = hostsInput
    .split("\n")
    .flatMap((line) => line.split(","))
    .map((h) => h.trim())
    .filter((h) => h.length > 0);

  if (hostsArray.length === 0) return;

  const hosts = hostsArray.join(",");
  const nameserverInput = document.getElementById("nameserverInput").value.trim();
  const includeDefault = getSwitchState(document.getElementById("includeDefault"));
  const includeAAAA = getSwitchState(document.getElementById("includeAAAA"));

  // Build query string
  let queryParams = `hosts=${encodeURIComponent(hosts)}`;
  if (nameserverInput) {
    const nameservers = nameserverInput
      .split(",")
      .map((ns) => ns.trim())
      .filter((ns) => ns.length > 0)
      .join(",");
    if (nameservers) {
      queryParams += `&nameservers=${encodeURIComponent(nameservers)}`;
    }
  }
  queryParams += `&includeDefault=${includeDefault}`;
  queryParams += `&includeAAAA=${includeAAAA}`;

  try {
    const response = await fetch(`/api/lookup?${queryParams}`);
    const data = await response.json();
    lastResponse = data;
    displayResults(data);
    updateShareableUrl(hosts, nameserverInput, includeDefault, includeAAAA);
  } catch (err) {
    const resultsDiv = document.getElementById("results");
    resultsDiv.innerHTML = `<div class="rounded-lg border border-red-200 dark:border-red-800 bg-red-50 dark:bg-red-950 p-4 text-sm text-red-600 dark:text-red-400">${escapeHtml(err.message)}</div>`;
    document.getElementById("resultsSection").classList.remove("hidden");
  }
}

// ── Display Results ────────────────────────────────────

function displayResults(data) {
  const resultsDiv = document.getElementById("results");
  resultsDiv.innerHTML = "";

  const resultsSection = document.getElementById("resultsSection");

  if (!data.results || data.results.length === 0) {
    resultsDiv.innerHTML = '<p class="text-sm text-zinc-500 dark:text-zinc-400">No results</p>';
    resultsSection.classList.remove("hidden");
    return;
  }

  resultsSection.classList.remove("hidden");

  // Attach drag-and-drop listeners on the container
  resultsDiv.addEventListener("dragover", onContainerDragOver);
  resultsDiv.addEventListener("drop", onContainerDrop);

  data.results.forEach((nsResult) => {
    // Draggable wrapper
    const wrapper = document.createElement("div");
    wrapper.className = "ns-draggable";
    wrapper.draggable = true;
    wrapper.addEventListener("dragstart", onDragStart);
    wrapper.addEventListener("dragend", onDragEnd);

    // Accordion using <details> — open by default
    const details = document.createElement("details");
    details.open = true;
    details.className = "rounded-lg border border-zinc-200 dark:border-zinc-800 bg-white dark:bg-zinc-900 shadow-sm overflow-hidden";

    // Summary (accordion header + drag handle)
    const summary = document.createElement("summary");
    summary.className = "px-4 py-2.5 flex items-center gap-2 cursor-pointer select-none list-none [&::-webkit-details-marker]:hidden hover:bg-zinc-50 dark:hover:bg-zinc-800/50 transition-colors";
    summary.innerHTML = `
      <span class="drag-handle cursor-grab active:cursor-grabbing text-zinc-300 dark:text-zinc-600 hover:text-zinc-500 dark:hover:text-zinc-400 mr-1 flex-shrink-0" title="Drag to reorder">
        <svg class="w-4 h-4" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="8" y1="6" x2="8" y2="6.01"/><line x1="16" y1="6" x2="16" y2="6.01"/><line x1="8" y1="12" x2="8" y2="12.01"/><line x1="16" y1="12" x2="16" y2="12.01"/><line x1="8" y1="18" x2="8" y2="18.01"/><line x1="16" y1="18" x2="16" y2="18.01"/></svg>
      </span>
      <svg class="w-3 h-3 text-zinc-400 dark:text-zinc-500 transition-transform flex-shrink-0 accordion-chevron" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="6 9 12 15 18 9"/></svg>
      <svg class="w-3.5 h-3.5 text-zinc-400 dark:text-zinc-500 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="2" y="2" width="20" height="8" rx="2" ry="2"/><rect x="2" y="14" width="20" height="8" rx="2" ry="2"/><line x1="6" y1="6" x2="6.01" y2="6"/><line x1="6" y1="18" x2="6.01" y2="18"/></svg>
      <span class="text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider">Nameserver</span>
      <span class="text-sm font-mono text-zinc-900 dark:text-zinc-100">${escapeHtml(nsResult.nameserver)}</span>`;
    details.appendChild(summary);

    // If nameserver has an error, show error state instead of table
    if (nsResult.error) {
      details.className = "rounded-lg border border-red-200 dark:border-red-800/60 bg-white dark:bg-zinc-900 shadow-sm overflow-hidden";
      const errorDiv = document.createElement("div");
      errorDiv.className = "border-t border-red-200 dark:border-red-800/60 bg-red-50 dark:bg-red-950/30 px-4 py-3 flex items-center gap-2";
      errorDiv.innerHTML = `
        <svg class="w-4 h-4 text-red-500 dark:text-red-400 flex-shrink-0" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"/><line x1="12" y1="8" x2="12" y2="12"/><line x1="12" y1="16" x2="12.01" y2="16"/></svg>
        <span class="text-sm text-red-600 dark:text-red-400">${escapeHtml(nsResult.error)}</span>`;
      details.appendChild(errorDiv);
      wrapper.appendChild(details);
      resultsDiv.appendChild(wrapper);
      return;
    }

    // Table wrapper
    const tableWrapper = document.createElement("div");
    tableWrapper.className = "overflow-x-auto border-t border-zinc-200 dark:border-zinc-800";

    const table = document.createElement("table");
    table.className = "ns-table w-full text-sm";

    const thead = document.createElement("thead");
    const headerRow = document.createElement("tr");
    headerRow.className = "border-b border-zinc-200 dark:border-zinc-800";
    ["Domain", "Host", "Type", "Data"].forEach((text) => {
      const th = document.createElement("th");
      th.className = "px-4 py-2 text-left text-xs font-medium text-zinc-500 dark:text-zinc-400 uppercase tracking-wider";
      th.textContent = text;
      headerRow.appendChild(th);
    });
    thead.appendChild(headerRow);

    const tbody = document.createElement("tbody");

    const domains = Object.keys(nsResult.results || {});
    if (domains.length === 0) {
      const tr = document.createElement("tr");
      const td = document.createElement("td");
      td.colSpan = 4;
      td.className = "px-4 py-3 text-sm text-zinc-400 dark:text-zinc-500 italic";
      td.textContent = "No records found";
      tr.appendChild(td);
      tbody.appendChild(tr);
    } else {
      domains.forEach((domain) => {
        const records = nsResult.results[domain];

        if (!Array.isArray(records) || records.length === 0) {
          const tr = document.createElement("tr");
          tr.className = "border-b border-zinc-100 dark:border-zinc-800 last:border-0";

          const domainTd = document.createElement("td");
          domainTd.className = "px-4 py-2 font-medium text-sm whitespace-nowrap";
          domainTd.textContent = domain;
          tr.appendChild(domainTd);

          tr.appendChild(createCell(""));
          tr.appendChild(createCell(""));

          const dataTd = document.createElement("td");
          dataTd.className = "px-4 py-2 text-sm text-zinc-400 dark:text-zinc-500 italic";
          dataTd.textContent = "No record found";
          tr.appendChild(dataTd);

          tbody.appendChild(tr);
        } else {
          records.forEach((record, index) => {
            const tr = document.createElement("tr");
            tr.className = "border-b border-zinc-100 dark:border-zinc-800 last:border-0";

            if (index === 0) {
              const domainTd = document.createElement("td");
              domainTd.className = "px-4 py-2 font-medium text-sm whitespace-nowrap align-top";
              domainTd.textContent = domain;
              domainTd.rowSpan = records.length;
              tr.appendChild(domainTd);
            }

            const hostTd = document.createElement("td");
            hostTd.className = "px-4 py-2 font-mono text-xs text-zinc-600 dark:text-zinc-300";
            hostTd.textContent = record.host || "";
            tr.appendChild(hostTd);

            const typeTd = document.createElement("td");
            typeTd.className = "px-4 py-2 font-mono text-xs";
            const typeSpan = document.createElement("span");
            typeSpan.className = getTypeBadgeClass(record.type);
            typeSpan.textContent = record.type || "";
            typeTd.appendChild(typeSpan);
            tr.appendChild(typeTd);

            const dataTd = document.createElement("td");
            dataTd.className = "px-4 py-2 font-mono text-xs text-zinc-600 dark:text-zinc-300 break-all";
            dataTd.textContent = record.data || "";
            tr.appendChild(dataTd);

            tbody.appendChild(tr);
          });
        }
      });
    }

    table.appendChild(thead);
    table.appendChild(tbody);
    tableWrapper.appendChild(table);
    details.appendChild(tableWrapper);
    wrapper.appendChild(details);
    resultsDiv.appendChild(wrapper);
  });
}

function createCell(text) {
  const td = document.createElement("td");
  td.className = "px-4 py-2 text-sm";
  td.textContent = text;
  return td;
}

function getTypeBadgeClass(type) {
  switch (type) {
    case "CNAME":
      return "inline-block px-1.5 py-0.5 rounded text-xs font-medium bg-blue-50 dark:bg-blue-950 text-blue-700 dark:text-blue-300 border border-blue-200 dark:border-blue-800";
    case "A":
      return "inline-block px-1.5 py-0.5 rounded text-xs font-medium bg-emerald-50 dark:bg-emerald-950 text-emerald-700 dark:text-emerald-300 border border-emerald-200 dark:border-emerald-800";
    case "AAAA":
      return "inline-block px-1.5 py-0.5 rounded text-xs font-medium bg-purple-50 dark:bg-purple-950 text-purple-700 dark:text-purple-300 border border-purple-200 dark:border-purple-800";
    default:
      return "inline-block px-1.5 py-0.5 rounded text-xs font-medium bg-zinc-100 dark:bg-zinc-800 text-zinc-600 dark:text-zinc-400";
  }
}

// ── Drag and Drop ──────────────────────────────────────

let draggedEl = null;
let dropIndicator = null;

function getDropIndicator() {
  if (!dropIndicator) {
    dropIndicator = document.createElement("div");
    dropIndicator.className = "drop-indicator";
  }
  return dropIndicator;
}

function onDragStart(e) {
  draggedEl = this;
  requestAnimationFrame(() => {
    this.classList.add("dragging");
  });
  e.dataTransfer.effectAllowed = "move";
}

function onContainerDragOver(e) {
  e.preventDefault();
  e.dataTransfer.dropEffect = "move";
  if (!draggedEl) return;

  const target = e.target.closest(".ns-draggable");
  if (!target || target === draggedEl) {
    // If hovering outside any card, hide indicator
    const indicator = getDropIndicator();
    indicator.classList.remove("active");
    if (indicator.parentNode) indicator.remove();
    return;
  }

  const container = document.getElementById("results");
  const indicator = getDropIndicator();
  const rect = target.getBoundingClientRect();
  const midY = rect.top + rect.height / 2;

  if (e.clientY < midY) {
    container.insertBefore(indicator, target);
  } else {
    const next = target.nextElementSibling;
    if (next && next === indicator) return; // already in place
    container.insertBefore(indicator, target.nextSibling);
  }
  indicator.classList.add("active");
}

function onContainerDrop(e) {
  e.preventDefault();
  if (!draggedEl) return;

  const indicator = getDropIndicator();
  if (indicator.parentNode) {
    indicator.parentNode.insertBefore(draggedEl, indicator);
    indicator.remove();
  }
}

function onDragEnd() {
  if (draggedEl) {
    draggedEl.classList.remove("dragging");
  }
  draggedEl = null;
  const indicator = getDropIndicator();
  indicator.classList.remove("active");
  if (indicator.parentNode) {
    indicator.remove();
  }
}

function expandAll() {
  document.querySelectorAll("#results details").forEach((d) => (d.open = true));
}

function collapseAll() {
  document.querySelectorAll("#results details").forEach((d) => (d.open = false));
}

// ── Shareable URL ──────────────────────────────────────

function updateShareableUrl(hosts, nameservers, includeDefault, includeAAAA) {
  const url = new URL(window.location.origin + window.location.pathname);
  url.searchParams.set("hosts", hosts);
  if (nameservers) {
    url.searchParams.set("nameservers", nameservers);
  }
  if (!includeDefault) {
    url.searchParams.set("includeDefault", "false");
  }
  if (includeAAAA) {
    url.searchParams.set("includeAAAA", "true");
  }

  const container = document.getElementById("shareableUrlContainer");
  container.classList.remove("hidden");
  document.getElementById("shareableUrl").value = url.toString();
}

// ── Copy Functions ─────────────────────────────────────

async function copyShareableUrl() {
  const url = document.getElementById("shareableUrl").value;
  await copyToClipboard(url, "copyUrlBtn", "Copy URL");
}

async function copyAsJSON() {
  if (!lastResponse) return;
  const json = JSON.stringify(lastResponse, null, 2);
  await copyToClipboard(json, "copyJsonBtn", "JSON");
}

async function copyAsText() {
  if (!lastResponse || !lastResponse.results) return;

  let text = "";
  lastResponse.results.forEach((nsResult, i) => {
    if (i > 0) text += "\n";
    text += `Nameserver: ${nsResult.nameserver}\n`;

    let maxDomain = 6, maxHost = 4, maxType = 4, maxData = 4;
    const rows = [];

    Object.keys(nsResult.results || {}).forEach((domain) => {
      const records = nsResult.results[domain];
      if (!records || records.length === 0) {
        rows.push({ domain, host: "", type: "", data: "No record found" });
        maxDomain = Math.max(maxDomain, domain.length);
        maxData = Math.max(maxData, 15);
      } else {
        records.forEach((r) => {
          rows.push({ domain, host: r.host, type: r.type, data: r.data });
          maxDomain = Math.max(maxDomain, domain.length);
          maxHost = Math.max(maxHost, (r.host || "").length);
          maxType = Math.max(maxType, (r.type || "").length);
          maxData = Math.max(maxData, (r.data || "").length);
        });
      }
    });

    const pad = (s, w) => (s || "").padEnd(w);
    const sep = `+-${"-".repeat(maxDomain)}-+-${"-".repeat(maxHost)}-+-${"-".repeat(maxType)}-+-${"-".repeat(maxData)}-+`;

    text += sep + "\n";
    text += `| ${pad("Domain", maxDomain)} | ${pad("Host", maxHost)} | ${pad("Type", maxType)} | ${pad("Data", maxData)} |\n`;
    text += sep + "\n";
    rows.forEach((r) => {
      text += `| ${pad(r.domain, maxDomain)} | ${pad(r.host, maxHost)} | ${pad(r.type, maxType)} | ${pad(r.data, maxData)} |\n`;
    });
    text += sep + "\n";
  });

  await copyToClipboard(text, "copyTextBtn", "Text");
}

async function copyToClipboard(text, buttonId, originalLabel) {
  try {
    await navigator.clipboard.writeText(text);
  } catch {
    const textarea = document.createElement("textarea");
    textarea.value = text;
    textarea.style.position = "fixed";
    textarea.style.opacity = "0";
    document.body.appendChild(textarea);
    textarea.select();
    document.execCommand("copy");
    document.body.removeChild(textarea);
  }
  const btn = document.getElementById(buttonId);
  const originalHTML = btn.innerHTML;
  btn.innerHTML = `<svg class="w-3.5 h-3.5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"/></svg> Copied!`;
  setTimeout(() => {
    btn.innerHTML = originalHTML;
  }, 2000);
}

// ── Load from URL ──────────────────────────────────────

function loadFromUrl() {
  const params = new URLSearchParams(window.location.search);
  const hosts = params.get("hosts");
  if (hosts) {
    document.getElementById("hostsInput").value = hosts.split(",").join("\n");
  }

  const nameservers = params.get("nameservers");
  if (nameservers) {
    document.getElementById("nameserverInput").value = nameservers;
  }

  const includeDefault = params.get("includeDefault");
  if (includeDefault === "false") {
    setSwitchState(document.getElementById("includeDefault"), false);
  }

  const includeAAAA = params.get("includeAAAA");
  if (includeAAAA === "true") {
    setSwitchState(document.getElementById("includeAAAA"), true);
  }

  // Update include-default toggle state based on nameserver input
  updateIncludeDefaultState();

  if (hosts) {
    lookupHosts();
  }
}

// ── Helpers ────────────────────────────────────────────

function escapeHtml(str) {
  const div = document.createElement("div");
  div.textContent = str;
  return div.innerHTML;
}

// ── Init ───────────────────────────────────────────────

initTheme();
window.onload = function () {
  loadInfo();
  loadFromUrl();

  // Watch nameserver input to enable/disable include-default toggle
  const nsInput = document.getElementById("nameserverInput");
  nsInput.addEventListener("input", updateIncludeDefaultState);
  nsInput.addEventListener("change", updateIncludeDefaultState);
  nsInput.addEventListener("paste", function() { setTimeout(updateIncludeDefaultState, 0); });
  updateIncludeDefaultState();
};
