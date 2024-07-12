async function lookupHosts() {
  const hostsInput = document.getElementById('hostsInput').value;
  const hostsArray = hostsInput.split('\n').map(line => line.trim()).filter(line => line.length > 0);
  const hosts = hostsArray.join(',');
  const response = await fetch(`/api/lookup?hosts=${hosts}`);
  const data = await response.json();
  displayResults(data);
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
