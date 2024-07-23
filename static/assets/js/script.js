async function lookupHosts() {
  const hostsInput = document.getElementById('hostsInput').value;
  const hostsArray = hostsInput.split('\n').map(line => line.trim()).filter(line => line.length > 0);
  const hosts = hostsArray.join(',');
  const response = await fetch(`/api/lookup?hosts=${hosts}`);
  const data = await response.json();
  displayResults(data);
  updateShareableUrl(hosts);
}

function displayResults(data) {
  const table = document.createElement('table');
  const thead = document.createElement('thead');
  const tbody = document.createElement('tbody');

  // Create header row
  const headerRow = document.createElement('tr');
  ['Domain', 'Host', 'Type', 'Data'].forEach(headerText => {
      const th = document.createElement('th');
      th.textContent = headerText;
      headerRow.appendChild(th);
  });
  thead.appendChild(headerRow);

  // Create rows for each host
  Object.keys(data).forEach(hostKey => {
      const hostRecords = data[hostKey];

      // Merge cells for the host key
      const hostRowSpan = hostRecords.length;
      hostRecords.forEach((record, index) => {
          const tr = document.createElement('tr');

          // Add host key column only for the first row
          if (index === 0) {
              const hostKeyTd = document.createElement('td');
              hostKeyTd.textContent = hostKey;
              hostKeyTd.rowSpan = hostRowSpan;
              tr.appendChild(hostKeyTd);
          }

          const hostTd = document.createElement('td');
          hostTd.textContent = record.host;
          tr.appendChild(hostTd);

          const typeTd = document.createElement('td');
          typeTd.textContent = record.type;
          tr.appendChild(typeTd);

          const dataTd = document.createElement('td');
          dataTd.textContent = record.data;
          tr.appendChild(dataTd);

          tbody.appendChild(tr);
      });
  });

  table.appendChild(thead);
  table.appendChild(tbody);

  const resultsDiv = document.getElementById('results');
  resultsDiv.innerHTML = '';
  resultsDiv.appendChild(table);
}

function updateShareableUrl(hosts) {
  const url = new URL(window.location.href);
  url.searchParams.set('hosts', hosts);
  const shareableUrl = url.toString();
  
  const shareableUrlInput = document.getElementById('shareableUrl');
  shareableUrlInput.value = shareableUrl;
}

function copyShareableUrl() {
  const shareableUrlInput = document.getElementById('shareableUrl');
  shareableUrlInput.select();
  document.execCommand('copy');
  
  const copyButton = document.getElementById('copyButton');
  copyButton.textContent = 'Copied!';
  setTimeout(() => {
    copyButton.textContent = 'Copy URL';
  }, 2000);
}

function loadFromUrl() {
  const urlParams = new URLSearchParams(window.location.search);
  const hosts = urlParams.get('hosts');
  if (hosts) {
    document.getElementById('hostsInput').value = hosts.split(',').join('\n');
    lookupHosts();
  }
}

window.onload = loadFromUrl;
