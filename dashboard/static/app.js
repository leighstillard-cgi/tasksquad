// WebSocket connection management
let ws = null;
let reconnectDelay = 1000;
const maxReconnectDelay = 30000;

function connectWebSocket() {
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/ws`;

    ws = new WebSocket(wsUrl);

    ws.onopen = function() {
        console.log('WebSocket connected');
        reconnectDelay = 1000; // Reset delay on successful connection
        updateConnectionStatus(true);
    };

    ws.onmessage = function(event) {
        try {
            const data = JSON.parse(event.data);
            updateAllPanels(data);
        } catch (error) {
            console.error('Failed to parse WebSocket message:', error);
        }
    };

    ws.onclose = function() {
        console.log('WebSocket disconnected, reconnecting in', reconnectDelay, 'ms');
        updateConnectionStatus(false);
        setTimeout(function() {
            reconnectDelay = Math.min(reconnectDelay * 2, maxReconnectDelay);
            connectWebSocket();
        }, reconnectDelay);
    };

    ws.onerror = function(error) {
        console.error('WebSocket error:', error);
    };
}

function updateConnectionStatus(connected) {
    const refreshSpan = document.querySelector('.last-refresh');
    if (refreshSpan) {
        if (connected) {
            refreshSpan.innerHTML = '<span class="ws-status connected"></span> Live';
        } else {
            refreshSpan.innerHTML = '<span class="ws-status disconnected"></span> Reconnecting...';
        }
    }
}

// Panel update functions
function updateAllPanels(data) {
    updateActiveWorkPanel(data.Dispatches);
    updateCompletionsPanel(data.Completions);
    updateEscalationsPanel(data.Escalations);
    updateBacklogPanel(data.Backlog, data.ReadyStories);
    updateSessionLogsPanel(data.SessionLogs);
    updateReadyStoriesDropdown(data.ReadyStories);
}

function updateActiveWorkPanel(dispatches) {
    const panel = document.querySelector('#active-work .panel-content');
    if (!panel) return;

    if (!dispatches || dispatches.length === 0) {
        panel.innerHTML = '<p class="empty-state">No dispatched stories</p>';
        return;
    }

    panel.innerHTML = `
        <table>
            <thead>
                <tr>
                    <th>Story</th>
                    <th>Repo</th>
                    <th>Dispatched</th>
                    <th>Status</th>
                </tr>
            </thead>
            <tbody>
                ${dispatches.map(d => `
                    <tr class="status-${d.Status}">
                        <td>${escapeHtml(d.StoryID)}</td>
                        <td>${escapeHtml(d.Repo)}</td>
                        <td>${formatDate(d.DispatchedAt)}</td>
                        <td><span class="badge badge-${d.Status}">${escapeHtml(d.Status)}</span></td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
    `;
}

function updateCompletionsPanel(completions) {
    const panel = document.querySelector('#completions .panel-content');
    if (!panel) return;

    if (!completions || completions.length === 0) {
        panel.innerHTML = '<p class="empty-state">No completion reports</p>';
        return;
    }

    panel.innerHTML = `
        <ul class="completion-list">
            ${completions.map(c => {
                const filename = c.FilePath.split('/').pop();
                return `
                    <li>
                        <a href="#" onclick="viewCompletion('${escapeHtml(filename)}'); return false;">
                            <strong>${escapeHtml(c.StoryID)}</strong>
                        </a>
                        <span class="badge badge-${c.Status}">${escapeHtml(c.Status)}</span>
                        <span class="timestamp">${formatDateOnly(c.Created)}</span>
                    </li>
                `;
            }).join('')}
        </ul>
    `;
}

function updateEscalationsPanel(escalations) {
    const panel = document.querySelector('#escalations .panel-content');
    const badge = document.querySelector('#escalations h2 .badge');
    if (!panel) return;

    if (!escalations || escalations.length === 0) {
        panel.innerHTML = '<p class="empty-state success">No escalations</p>';
        if (badge) badge.remove();
        return;
    }

    // Update or add badge
    const h2 = document.querySelector('#escalations h2');
    if (badge) {
        badge.textContent = escalations.length;
    } else if (h2) {
        const newBadge = document.createElement('span');
        newBadge.className = 'badge badge-danger';
        newBadge.textContent = escalations.length;
        h2.appendChild(newBadge);
    }

    panel.innerHTML = `
        <ul class="escalation-list">
            ${escalations.map(e => `
                <li>
                    <strong>${escapeHtml(e.StoryID)}</strong>
                    <span class="reason">${escapeHtml(e.Reason)}</span>
                    <span class="timestamp">${formatDate(e.Timestamp)}</span>
                </li>
            `).join('')}
        </ul>
    `;
}

function updateBacklogPanel(backlog, readyStories) {
    const panel = document.querySelector('#backlog .panel-content');
    if (!panel || !backlog) return;

    let readyHtml = '';
    if (readyStories && readyStories.length > 0) {
        readyHtml = `
            <h3>Ready Stories</h3>
            <ul class="story-list">
                ${readyStories.map(s => `
                    <li>
                        <strong>${escapeHtml(s.StoryID)}</strong>
                        <span class="title">${escapeHtml(s.Title)}</span>
                        <span class="repo">${escapeHtml(s.Repo)}</span>
                    </li>
                `).join('')}
            </ul>
        `;
    }

    panel.innerHTML = `
        <div class="status-summary">
            <div class="status-card">
                <span class="count">${backlog.Done ? backlog.Done.length : 0}</span>
                <span class="label">Done</span>
            </div>
            <div class="status-card">
                <span class="count">${backlog.InProgress ? backlog.InProgress.length : 0}</span>
                <span class="label">In Progress</span>
            </div>
            <div class="status-card">
                <span class="count">${backlog.Ready ? backlog.Ready.length : 0}</span>
                <span class="label">Ready</span>
            </div>
            <div class="status-card">
                <span class="count">${backlog.Blocked ? backlog.Blocked.length : 0}</span>
                <span class="label">Blocked</span>
            </div>
        </div>
        ${readyHtml}
    `;
}

function updateSessionLogsPanel(logs) {
    const list = document.getElementById('sessionLogsList');
    if (!list) return;

    if (!logs || logs.length === 0) {
        list.innerHTML = '<li class="empty-state">No session logs</li>';
        return;
    }

    list.innerHTML = logs.map(log => {
        const filename = log.FilePath.split('/').pop();
        return `
            <li class="status-${log.Status}">
                <a href="#" onclick="viewSessionLog('${escapeHtml(filename)}'); return false;" class="filename">${escapeHtml(log.FileName)}</a>
                <span class="badge badge-${log.Status}">${escapeHtml(log.Status)}</span>
                <span class="timestamp">${formatDate(log.Timestamp)}</span>
            </li>
        `;
    }).join('');
}

function updateReadyStoriesDropdown(readyStories) {
    const select = document.getElementById('storyId');
    if (!select) return;

    const currentValue = select.value;
    select.innerHTML = '<option value="">Select a story...</option>';

    if (readyStories && readyStories.length > 0) {
        readyStories.forEach(s => {
            const option = document.createElement('option');
            option.value = s.StoryID;
            option.dataset.repo = s.Repo;
            option.dataset.desc = s.Description;
            option.textContent = `${s.StoryID} - ${s.Title}`;
            select.appendChild(option);
        });
    }

    // Restore selection if still valid
    if (currentValue) {
        select.value = currentValue;
    }
}

// Modal functions
function openModal(content) {
    const modal = document.getElementById('detailModal');
    const body = document.getElementById('modalBody');
    if (modal && body) {
        body.innerHTML = content;
        modal.classList.remove('hidden');
    }
}

function closeModal() {
    const modal = document.getElementById('detailModal');
    if (modal) {
        modal.classList.add('hidden');
    }
}

async function viewCompletion(filename) {
    try {
        const response = await fetch(`/api/completion/${encodeURIComponent(filename)}`);
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }
        const html = await response.text();
        openModal(html);
    } catch (error) {
        console.error('Failed to load completion:', error);
        openModal(`<p class="error">Failed to load completion report: ${error.message}</p>`);
    }
}

async function viewSessionLog(filename) {
    try {
        const response = await fetch(`/api/session-log/${encodeURIComponent(filename)}`);
        if (!response.ok) {
            throw new Error(`HTTP ${response.status}`);
        }
        const html = await response.text();
        openModal(html);
    } catch (error) {
        console.error('Failed to load session log:', error);
        openModal(`<p class="error">Failed to load session log: ${error.message}</p>`);
    }
}

// Utility functions
function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateStr) {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function formatDateOnly(dateStr) {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString();
}

// Original functions (kept for compatibility)
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
            // No need to reload - WebSocket will push the update
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

        list.innerHTML = logs.map(log => {
            const filename = log.FilePath ? log.FilePath.split('/').pop() : log.FileName;
            return `
                <li class="status-${log.Status}">
                    <a href="#" onclick="viewSessionLog('${escapeHtml(filename)}'); return false;" class="filename">${escapeHtml(log.FileName)}</a>
                    <span class="badge badge-${log.Status}">${escapeHtml(log.Status)}</span>
                    <span class="timestamp">${formatDate(log.Timestamp)}</span>
                </li>
            `;
        }).join('');
    } catch (error) {
        console.error('Failed to filter session logs:', error);
    }
}

// Close modal on escape key
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        closeModal();
    }
});

// Close modal on click outside
document.addEventListener('click', function(event) {
    const modal = document.getElementById('detailModal');
    if (event.target === modal) {
        closeModal();
    }
});

// Initialize WebSocket on page load
document.addEventListener('DOMContentLoaded', function() {
    connectWebSocket();
});
