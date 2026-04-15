async function refreshData() {
    try {
        const response = await fetch('/api/refresh', { method: 'POST' });
        if (response.ok) {
            window.location.reload();
        }
    } catch (error) {
        console.error('Failed to refresh:', error);
    }
}

async function submitDispatch(event) {
    event.preventDefault();

    const form = event.target;
    const resultDiv = document.getElementById('dispatchResult');

    const storySelect = document.getElementById('storyId');
    const selectedOption = storySelect.options[storySelect.selectedIndex];

    const data = {
        story_id: form.story_id.value,
        repo: form.repo.value,
        description: form.description.value || selectedOption.dataset.desc || ''
    };

    try {
        const response = await fetch('/api/dispatch', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(data)
        });

        const result = await response.json();

        if (response.ok) {
            resultDiv.textContent = `Dispatch created: ${result.path}`;
            resultDiv.className = 'result-message success';
            form.reset();
            setTimeout(() => window.location.reload(), 1500);
        } else {
            resultDiv.textContent = result.error || 'Failed to create dispatch';
            resultDiv.className = 'result-message error';
        }
    } catch (error) {
        resultDiv.textContent = 'Request failed: ' + error.message;
        resultDiv.className = 'result-message error';
    }
}

document.getElementById('storyId')?.addEventListener('change', function() {
    const selectedOption = this.options[this.selectedIndex];
    const repoInput = document.getElementById('repo');
    const descInput = document.getElementById('description');

    if (selectedOption.dataset.repo) {
        repoInput.value = selectedOption.dataset.repo;
    }
    if (selectedOption.dataset.desc) {
        descInput.value = selectedOption.dataset.desc;
    }
});

async function filterSessionLogs() {
    const status = document.getElementById('sessionStatusFilter').value;

    try {
        const response = await fetch(`/api/session-logs?status=${status}`);
        const logs = await response.json();

        const list = document.getElementById('sessionLogsList');
        if (logs.length === 0) {
            list.innerHTML = '<li class="empty-state">No matching session logs</li>';
            return;
        }

        list.innerHTML = logs.map(log => `
            <li class="status-${log.Status}">
                <span class="filename">${log.FileName}</span>
                <span class="badge badge-${log.Status}">${log.Status}</span>
                <span class="timestamp">${formatDate(log.Timestamp)}</span>
            </li>
        `).join('');
    } catch (error) {
        console.error('Failed to filter session logs:', error);
    }
}

function formatDate(dateStr) {
    const date = new Date(dateStr);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function viewFile(path) {
    alert('File: ' + path);
}
