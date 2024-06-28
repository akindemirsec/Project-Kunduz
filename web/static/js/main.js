document.addEventListener('DOMContentLoaded', function () {
    document.getElementById('dashboard-link').addEventListener('click', loadDashboard);
    document.getElementById('clusters-link').addEventListener('click', loadClusters);
    document.getElementById('cves-link').addEventListener('click', loadCVEs);
    document.getElementById('settings-link').addEventListener('click', loadSettings);
    loadDashboard(); // Default load dashboard
});

function loadDashboard() {
    fetch('/dashboard')
        .then(response => response.text())
        .then(html => {
            document.getElementById('content').innerHTML = html;
        });
}

function loadClusters() {
    fetch('/clusters')
        .then(response => response.json())
        .then(data => {
            const content = document.getElementById('content');
            content.innerHTML = `<h2>Clusters</h2><button onclick="showCreateClusterForm()">Create Cluster</button><div id="clusters-list"></div>`;
            const clustersList = document.getElementById('clusters-list');
            clustersList.innerHTML = '';
            data.forEach(cluster => {
                const clusterDiv = document.createElement('div');
                clusterDiv.innerHTML = `<h3>${cluster.name}</h3><button onclick="showAddSBOMForm(${cluster.id})">Add SBOM</button><button onclick="showAddAssetForm(${cluster.id})">Add Asset</button>`;
                clustersList.appendChild(clusterDiv);
            });
        });
}

function loadCVEs() {
    fetch('/cves')
        .then(response => response.json())
        .then(data => {
            const content = document.getElementById('content');
            content.innerHTML = `<h2>CVEs</h2><button onclick="updateCVE()">Update CVEs</button><div id="cves-list"></div>`;
            const cvesList = document.getElementById('cves-list');
            cvesList.innerHTML = '';
            data.forEach(cve => {
                const cveDiv = document.createElement('div');
                cveDiv.innerHTML = `<h3>${cve.cve_id}</h3><p>${cve.description}</p>`;
                cvesList.appendChild(cveDiv);
            });
        });
}

function loadSettings() {
    const content = document.getElementById('content');
    content.innerHTML = `<h2>Settings</h2><p>Settings content here.</p>`;
}

function showCreateClusterForm() {
    const content = document.getElementById('content');
    content.innerHTML = `
        <h2>Create Cluster</h2>
        <form id="create-cluster-form">
            <label for="cluster-name">Cluster Name:</label>
            <input type="text" id="cluster-name" name="name" required>
            <button type="submit">Create</button>
        </form>
    `;
    document.getElementById('create-cluster-form').addEventListener('submit', createCluster);
}

function createCluster(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    fetch('/clusters', {
        method: 'POST',
        body: JSON.stringify(Object.fromEntries(formData)),
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        alert('Cluster created successfully!');
        loadClusters();
    })
    .catch(error => console.error('Error:', error));
}

function showAddSBOMForm(clusterId) {
    const content = document.getElementById('content');
    content.innerHTML = `
        <h2>Add SBOM</h2>
        <form id="add-sbom-form">
            <input type="hidden" name="cluster_id" value="${clusterId}">
            <label for="sbom-name">SBOM Name:</label>
            <input type="text" id="sbom-name" name="name" required>
            <label for="sbom-data">SBOM Data:</label>
            <textarea id="sbom-data" name="data" required></textarea>
            <button type="submit">Add SBOM</button>
        </form>
    `;
    document.getElementById('add-sbom-form').addEventListener('submit', addSBOM);
}

function addSBOM(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const clusterId = formData.get('cluster_id');
    fetch(`/clusters/${clusterId}/sboms`, {
        method: 'POST',
        body: JSON.stringify(Object.fromEntries(formData)),
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        alert('SBOM added successfully!');
        loadClusters();
    })
    .catch(error => console.error('Error:', error));
}

function showAddAssetForm(clusterId) {
    const content = document.getElementById('content');
    content.innerHTML = `
        <h2>Add Asset</h2>
        <form id="add-asset-form">
            <input type="hidden" name="cluster_id" value="${clusterId}">
            <label for="asset-name">Asset Name:</label>
            <input type="text" id="asset-name" name="name" required>
            <button type="submit">Add Asset</button>
        </form>
    `;
    document.getElementById('add-asset-form').addEventListener('submit', addAsset);
}

function addAsset(event) {
    event.preventDefault();
    const formData = new FormData(event.target);
    const clusterId = formData.get('cluster_id');
    fetch(`/clusters/${clusterId}/assets`, {
        method: 'POST',
        body: JSON.stringify(Object.fromEntries(formData)),
        headers: {
            'Content-Type': 'application/json'
        }
    })
    .then(response => response.json())
    .then(data => {
        alert('Asset added successfully!');
        loadClusters();
    })
    .catch(error => console.error('Error:', error));
}

function updateCVE() {
    fetch('/cves/update', {
        method: 'POST'
    })
    .then(response => response.json())
    .then(data => {
        alert('CVE database updated successfully!');
        loadCVEs();
    })
    .catch(error => console.error('Error:', error));
}
