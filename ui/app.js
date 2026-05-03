// Zsistant Web UI — Full-featured Assistant Control Center
(function() {
  'use strict';

  const state = {
    view: 'home',
    assistants: [],
    currentAssistantId: null,
    devMode: false,
    theme: 'dark',
    models: [],        // flat model list from /api/models
    providers: [],     // provider list from /api/providers
    activity: [],
  };

  // ---------- Helpers ----------
  const $ = id => document.getElementById(id);
  const esc = t => {
    const d = document.createElement('div');
    d.textContent = t;
    return d.innerHTML;
  };
  const fmt = d => d ? new Date(d).toLocaleString() : '';

  async function api(path, opts = {}) {
    const url = path.startsWith('/') ? path : '/api/' + path;
    try {
      const res = await fetch(url, opts);
      if (res.status === 204) return null;
      if (!res.ok) {
        const body = await res.text();
        return { error: body || res.statusText, status: res.status };
      }
      return await res.json();
    } catch (e) {
      return { error: e.message };
    }
  }
  async function apiGET(path) { return api(path, null); }
  async function apiPOST(path, body) { return api(path, { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) }); }
  async function apiPUT(path, body) { return api(path, { method: 'PUT', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(body) }); }
  async function apiDEL(path) { return api(path, { method: 'DELETE' }); }

  // ---------- Navigation ----------
  function switchView(view) {
    state.view = view;
    document.querySelectorAll('.view').forEach(v => v.classList.toggle('active', v.id === 'view-' + view));
    document.querySelectorAll('.nav-item').forEach(n => n.classList.toggle('active', n.dataset.view === view));
    if (view === 'home') refreshDashboard();
    if (view === 'assistants') loadAssistants();
    if (view === 'jobs') loadGlobalJobs();
    if (view === 'channels') loadChannels();
    if (view === 'tools') loadTools();
    if (view === 'settings') loadSettings();
  }

  function showDetail(id, tab) {
    state.currentAssistantId = id;
    const a = state.assistants.find(x => x.id === id);
    if (!a) return;
    switchView('');
    $('view-home').classList.remove('active');
    $('view-assistants').classList.remove('active');
    $('view-assistant-detail').classList.add('active');
    $('detailTitle').textContent = a.name || a.id;
    $('detailStatus').textContent = a.status || 'unknown';
    $('detailStatus').className = 'status-badge ' + (a.status || 'unknown');
    renderTabs();
    switchTab(tab || 'overview');
  }

  function setBack() {
    $('view-assistant-detail').classList.remove('active');
    $('view-assistants').classList.add('active');
    state.view = 'assistants';
    loadAssistants();
  }

  // ---------- Tabs (dynamically built) ----------
  const TAB_DEFS = [
    { id: 'overview',   label: 'Overview'   },
    { id: 'runtime',    label: 'Runtime'    },
    { id: 'approvals',  label: 'Approvals'  },
    { id: 'chat',       label: 'Chat'       },
    { id: 'channels',   label: 'Channels'   },
    { id: 'tools',      label: 'Tools'      },
    { id: 'knowledge',  label: 'Knowledge'  },
    { id: 'memory',     label: 'Memory'     },
    { id: 'jobs',       label: 'Jobs'       },
    { id: 'browser',    label: 'Browser'    },
    { id: 'logs',       label: 'Logs'       },
    { id: 'settings',   label: 'Settings'   },
    { id: 'developer',  label: 'Developer',  dev: true },
  ];

  function renderTabs() {
    const nav = $('detailTabs');
    nav.innerHTML = TAB_DEFS.map(t => state.devMode || !t.dev
      ? `<button class="tab-btn" data-tab="${t.id}" data-dev="${!!t.dev}">${t.label}</button>`
      : '').join('');
    document.querySelectorAll('.tab-btn').forEach(b => {
      b.addEventListener('click', () => switchTab(b.dataset.tab));
    });

    // Build panels
    const panelContainer = $('detailPanels');
    panelContainer.innerHTML = TAB_DEFS.map(t => `
      <div id="panel-${t.id}" class="tab-panel">
        <div class="panel-content" id="panelContent-${t.id}"></div>
      </div>`).join('');

    if (state.devMode) document.body.classList.add('dev-active');
    else document.body.classList.remove('dev-active');
  }

  function switchTab(tab) {
    document.querySelectorAll('.tab-btn').forEach(b => b.classList.toggle('active', b.dataset.tab === tab));
    document.querySelectorAll('.panel-content').forEach(p => p.closest('.tab-panel').classList.remove('active'));
    const panel = $('panel-' + tab);
    if (panel) panel.classList.add('active');

    const a = state.assistants.find(x => x.id === state.currentAssistantId);
    if (!a) return;
    const pc = $('panelContent-' + tab);
    switch (tab) {
      case 'overview': renderOverview(pc, a); break;
      case 'chat': renderChat(pc, a); break;
      case 'channels': {
        pc.innerHTML = '<div class="animate-in">Loading channels...</div>';
        apiGET('/api/assistants/' + a.id + '/channels').then(data => {
          if (data.error) { pc.innerHTML = `<div class="empty-state error">${esc(data.error)}</div>`; return; }
          renderChannels(pc, data);
        });
        break;
      }
      case 'tools': renderToolsPanel(pc, a); break;
      case 'knowledge': {
        pc.innerHTML = '<div class="animate-in">Loading knowledge...</div>';
        apiGET('/api/assistants/' + a.id + '/knowledge').then(data => {
          if (data.error) { pc.innerHTML = `<div class="empty-state error">${esc(data.error)}</div>`; return; }
          renderKnowledge(pc, data, a);
        });
        break;
      }
      case 'memory': {
        pc.innerHTML = '<div class="animate-in">Loading memories...</div>';
        apiGET('/api/assistants/' + a.id + '/memory').then(data => {
          if (data.error) { pc.innerHTML = `<div class="empty-state error">${esc(data.error)}</div>`; return; }
          renderMemory(pc, data, a);
        });
        break;
      }
      case 'jobs': {
        pc.innerHTML = '<div class="animate-in">Loading jobs...</div>';
        apiGET('/api/assistants/' + a.id + '/jobs').then(data => {
          if (data.error) { pc.innerHTML = `<div class="empty-state error">${esc(data.error)}</div>`; return; }
          renderJobs(pc, data, a);
        });
        break;
      }
      case 'browser': renderBrowser(pc, a); break;
      case 'logs': {
        pc.innerHTML = '<div class="animate-in">Loading logs...</div>';
        apiGET('/api/assistants/' + a.id + '/logs').then(data => {
          if (data.error) { pc.innerHTML = `<div class="empty-state error">${esc(data.error)}</div>`; return; }
          renderLogs(pc, data);
        });
        break;
      }
      case 'settings': renderAssistantSettings(pc, a); break;
      case 'runtime': renderRuntime(pc, a); break;
      case 'approvals': renderApprovals(pc, a); break;
      case 'developer': renderDeveloper(pc, a); break;
    }
  }

  // ---------- Render Assistants list ----------
  async function loadAssistants() {
    const data = await apiGET('/api/assistants');
    if (data.error) { console.error(data.error); return; }
    state.assistants = Array.isArray(data) ? data : [];
    $('navAssCount').textContent = state.assistants.length;
    renderAssistantsList();
    refreshDashboard();
  }

  function renderAssistantsList() {
    const grid = $('assistantGrid');
    if (!state.assistants.length) {
      grid.innerHTML = `<div class="empty-state hero"><h3>No assistants yet.</h3><p>Create your first assistant to get started.</p></div>`;
      return;
    }
    grid.innerHTML = state.assistants.map(a => {
      const chCount = (a.channels || []).filter(c => c.status === 'connected').length;
      return `
        <div class="assistant-card animate-in" data-id="${a.id}">
          <div class="card-header">
            <h3>${esc(a.name)}</h3>
            <span class="status-badge ${a.status}">${a.status}</span>
          </div>
          <p class="card-desc">${esc(a.description || a.purpose || 'No description')}</p>
          <div class="card-meta">
            <span>Model: ${esc(a.default_model || 'default')}</span>
            <span>Channels: ${chCount}</span>
          </div>
          <div class="card-actions">
            <button class="open-btn" onclick="window._app.openDetail('${a.id}')">Open</button>
            <button class="delete-btn" onclick="window._app.delAss('${a.id}')">Delete</button>
          </div>
        </div>`;
    }).join('');
  }

  // ---------- Dashboard ----------
  async function refreshDashboard() {
    const assCard = $('dashAssistantsCard');
    const chCard = $('dashChannelsCard');
    const jobsCard = $('dashJobsCard');
    const actCard = $('dashActivityCard');

    const active = state.assistants.filter(a => a.status === 'active' || a.status === 'created');
    assCard.innerHTML = `<h3>Active Assistants</h3>` + (active.length
      ? active.map(a => `<div class="dash-item">${esc(a.name)} <span class="status-badge ${a.status}">${a.status}</span></div>`).join('')
      : '<div class="empty-state">No assistants yet. Create one in the Assistants tab.</div>');

    const connected = state.assistants.flatMap(a => a.channels || []).filter(c => c.status === 'connected').length;
    chCard.innerHTML = `<h3>Connected Channels</h3>` + (connected
      ? `<div class="dash-stat">${connected}</div><p style="color:var(--text-muted);margin-top:0.5rem;">channel(s) connected</p>`
      : '<div class="empty-state">No channels connected. Configure in an assistant\'s Channels tab.</div>');

    jobsCard.innerHTML = `<h3>Upcoming Jobs</h3><div class="empty-state">Scheduled jobs view coming soon.</div>`;

    try {
      const act = await apiGET('/api/activity');
      if (Array.isArray(act) && act.length) {
        actCard.innerHTML = `<h3>Recent Activity</h3>` + act.slice(0, 5).map(e =>
          `<div class="dash-item">${esc(e.event_type)} <span style="color:var(--text-muted)">${fmt(e.created_at)}</span></div>`
        ).join('');
      } else {
        actCard.innerHTML = `<h3>Recent Activity</h3><div class="empty-state">No recent activity.</div>`;
      }
    } catch (e) { actCard.innerHTML = `<h3>Recent Activity</h3><div class="empty-state">No recent activity.</div>`; }
  }

  // ---------- Overview ----------
  function renderOverview(panel, a) {
    const p = a.persona || {};
    panel.innerHTML = `
      <div class="overview-grid animate-in">
        <div class="ov-card"><h4>Identity</h4>
          <p><strong>Name:</strong> ${esc(a.name)}</p>
          <p><strong>ID:</strong> <code>${esc(a.id)}</code></p>
          <p><strong>Status:</strong> <span class="status-badge ${a.status}">${a.status}</span></p>
          <p><strong>Purpose:</strong> ${esc(a.description || a.purpose || 'Not set')}</p>
        </div>
        <div class="ov-card"><h4>Persona</h4>
          <p><strong>Tone:</strong> ${esc(p.tone || 'Not set')}</p>
          <p><strong>Style:</strong> ${esc(p.style || 'Not set')}</p>
          <p><strong>Role:</strong> ${esc(p.role_description || 'Not set')}</p>
          <p><strong>Boundaries:</strong> ${esc(p.boundaries || 'Not set')}</p>
        </div>
        <div class="ov-card"><h4>Connections</h4>
          <p class="ov-stat">${(a.channels || []).length}</p><p>Channels</p>
          <p class="ov-stat">${(a.tool_permissions || []).length}</p><p>Tools</p>
          <p class="ov-stat">${(a.knowledge || []).length}</p><p>Knowledge items</p>
        </div>
        <div class="ov-card"><h4>Model</h4>
          <p><strong>Model:</strong> ${esc(a.default_model || 'default')}</p>
          <p><strong>Provider:</strong> ${esc(a.provider_name || 'auto')}</p>
          <p><strong>Memory policy:</strong> ${esc(a.memory_policy.scope || 'global')}</p>
          <p><strong>Jobs:</strong> ${a.jobs_enabled ? 'Enabled' : 'Disabled'}</p>
        </div>
      </div>`;
  }

  // ---------- Chat ----------
  function renderChat(panel, a) {
    panel.innerHTML = `
      <div id="chatHistory-${a.id}" class="messages animate-in"></div>
      <div class="composer">
        <textarea id="chatInput" rows="2" placeholder="Message ${esc(a.name)}..."></textarea>
        <button id="chatSendBtn" class="send-btn">Send</button>
      </div>`;
    const send = () => sendChatMessage(a);
    $('chatSendBtn').addEventListener('click', send);
    $('chatInput').addEventListener('keydown', e => { if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); send(); } });
  }

  async function sendChatMessage(a) {
    const input = $('chatInput');
    const text = input.value.trim();
    if (!text) return;
    const h = $('chatHistory-' + a.id);
    h.innerHTML += `<div class="msg user">${esc(text)}</div>`;
    h.scrollTop = h.scrollHeight;
    input.value = '';
    const streamBtn = $('chatSendBtn');
    const origText = streamBtn.textContent;
    streamBtn.textContent = '...';
    streamBtn.disabled = true;

    const msgEl = document.createElement('div');
    msgEl.className = 'msg assistant streaming';
    h.appendChild(msgEl);

    try {
      const resp = await fetch(`/api/assistants/${a.id}/chat/stream`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ message: text })
      });
      if (!resp.ok || !resp.body) {
        throw new Error('stream request failed');
      }
      const reader = resp.body.getReader();
      const decoder = new TextDecoder();
      let full = '';
      while (true) {
        const { done, value } = await reader.read();
        if (done) break;
        const chunk = decoder.decode(value, { stream: true });
        const lines = chunk.split('\n');
        for (const line of lines) {
          if (line.startsWith('data: ')) {
            const dataText = line.slice(6);
            if (dataText === '[DONE]') {
              msgEl.classList.remove('streaming');
              streamBtn.textContent = origText;
              streamBtn.disabled = false;
              break;
            }
            full += dataText;
            msgEl.innerHTML = esc(full);
            h.scrollTop = h.scrollHeight;
          }
        }
      }
      if (msgEl.classList.contains('streaming')) {
        msgEl.classList.remove('streaming');
        streamBtn.textContent = origText;
        streamBtn.disabled = false;
      }
    } catch (e) {
      console.error('stream failed, falling back to sync', e);
      fallbackSyncChat(a, text, msgEl, streamBtn, origText);
    }
  }

  async function fallbackSyncChat(a, text, msgEl, btn, origText) {
    const res = await apiPOST('/api/assistants/' + a.id + '/chat', { message: text });
    btn.textContent = origText;
    btn.disabled = false;
    msgEl.classList.remove('streaming');
    if (res.error) {
      msgEl.className = 'msg error';
      msgEl.innerHTML = `Error: ${esc(res.error)}<br><small style="color:var(--text-muted)">Model: ${esc(a.default_model || 'default')}</small>`;
    } else {
      msgEl.innerHTML = esc(res.response || res.result || '');
    }
    const h = $('chatHistory-' + a.id);
    if (h) h.scrollTop = h.scrollHeight;
  }

  // ---------- Tabs renderers ----------
  function renderChannels(panel, data) {
    if (!data || !data.length) {
      panel.innerHTML = `<div class="empty-state">
        <p>No channels configured.</p>
        <p class="hint">Add a channel to connect this assistant to Telegram, Discord, WhatsApp, Slack, etc.</p>
      </div>`;
      return;
    }
    panel.innerHTML = data.map(c => `
      <div class="channel-row">
        <span class="channel-type">${esc(c.channel_type)}</span>
        <span class="status-badge ${c.status}">${c.status}</span>
        <span style="color:var(--text-muted);font-size:0.8rem">${esc(c.channel_id || '')}</span>
      </div>`).join('');
  }

  function renderToolsPanel(panel, a) {
    // Fetch tools from global registry
    panel.innerHTML = '<div class="animate-in">Loading tools...</div>';
    apiGET('/api/tools').then(data => {
      if (!Array.isArray(data) || !data.length) {
        panel.innerHTML = `<div class="empty-state">No tools registered.</div>`;
        return;
      }
      panel.innerHTML = data.map(t => `
        <div class="tool-row">
          <span class="tool-name">${esc(t.name)}</span>
          <span class="tool-desc">${esc(t.description || '')}</span>
          <span class="status-badge">${a.tool_permissions?.includes(t.name) ? 'allowed' : 'blocked'}</span>
        </div>`).join('');
    });
  }

  function renderKnowledge(panel, data, a) {
    if (!data || !data.length) {
      panel.innerHTML = `<div class="empty-state action">
        <p>No knowledge attached.</p>
        <p class="hint">Upload documents, links, or notes. They become context for this assistant.</p>
        <button class="add-btn" style="margin-top:0.5rem">+ Add Knowledge</button>
      </div>`;
      return;
    }
    panel.innerHTML = data.map(k => `
      <div class="knowledge-row">
        <span><strong>${esc(k.name)}</strong> <span style="color:var(--text-muted)">(${k.type})</span></span>
        <span class="status-badge ${k.status}">${k.status}</span>
      </div>`).join('');
  }

  function renderMemory(panel, data, a) {
    if (!data || !data.length) {
      panel.innerHTML = `<div class="empty-state">
        <p>No memories yet.</p>
        <p class="hint">Memories are facts the assistant remembers about your preferences, projects, and conversations.</p>
      </div>`;
      return;
    }
    panel.innerHTML = data.map(m => `
      <div class="memory-row">
        <span>[${esc(m.category || 'note')}] ${esc(m.content.substring(0, 100))}${m.content.length > 100 ? '...' : ''}</span>
        <span style="color:var(--text-muted);font-size:0.78rem">${fmt(m.created_at)}</span>
      </div>`).join('');
  }

  function renderJobs(panel, data, a) {
    if (!data || !data.length) {
      panel.innerHTML = `<div class="empty-state action">
        <p>No jobs scheduled yet.</p>
        <p class="hint">Jobs are background tasks an assistant runs on a schedule or on demand.</p>
        <button class="add-btn" onclick="window._app.addJob('${a.id}')" style="margin-top:0.5rem">+ Create Job</button>
      </div>`;
      return;
    }
    panel.innerHTML = data.map(j => `
      <div class="job-row">
        <span class="job-name">${esc(j.name || j.type || 'Job')}</span>
        <span class="status-badge ${j.status}">${j.status}</span>
        <span style="color:var(--text-muted);font-size:0.78rem">${fmt(j.created_at)}</span>
      </div>`).join('');
  }

  function renderBrowser(panel, a) {
    apiGET('/api/assistants/' + a.id + '/browser').then(cfg => {
      if (cfg.error) {
        panel.innerHTML = `<div class="empty-state">Browser / MCP status unavailable.</div>`;
        return;
      }
      if (cfg.level === 'not_available') {
        panel.innerHTML = `<div class="empty-state">
          <p>Browser / MCP is not connected.</p>
          <p class="hint">${esc(cfg.setup_message || '')}</p>
          <p class="hint">Configure a Chrome or Playwright MCP server to enable browser actions.</p>
        </div>`;
      } else {
        panel.innerHTML = `<div class="ov-card" style="margin-top:1rem">
          <h4>Browser MCP</h4>
          <p><strong>Level:</strong> ${esc(cfg.level)}</p>
          <p><strong>Servers:</strong> ${(cfg.mcp_servers || []).length}</p>
          <p><strong>Sessions:</strong> ${(cfg.active_sessions || []).length}</p>
        </div>`;
      }
    });
  }

  function renderLogs(panel, data) {
    if (!data || !data.length) {
      panel.innerHTML = `<div class="empty-state">No logs yet.</div>`;
      return;
    }
    panel.innerHTML = data.map(e => `
      <div class="log-row">
        <span><strong>[${esc(e.event_type)}]</strong> ${esc(e.message || '')}</span>
        <span style="color:var(--text-muted);font-size:0.78rem">${fmt(e.created_at)}</span>
      </div>`).join('');
  }

  function renderAssistantSettings(panel, a) {
    panel.innerHTML = `
      <div class="animate-in" style="max-width:600px">
        <h3 style="margin-bottom:1rem">Assistant Settings</h3>
        <div class="form-group">
          <label>Name</label>
          <input type="text" id="setAssName" value="${esc(a.name || '')}">
        </div>
        <div class="form-group">
          <label>Description</label>
          <textarea id="setAssDesc" rows="2">${esc(a.description || '')}</textarea>
        </div>
        <div class="form-group">
          <label>Purpose</label>
          <textarea id="setAssPurpose" rows="3">${esc(a.purpose || '')}</textarea>
        </div>
        <div class="form-group">
          <label>Default Model</label>
          <select id="setAssModel">
            ${state.models.map(m => `<option value="${esc(m.id)}" ${m.id === (a.default_model||'') ? 'selected' : ''}>${esc(m.id)} (${esc(m.provider)})</option>`).join('')}
            <option value="" ${!a.default_model ? 'selected' : ''}>Default (auto)</option>
          </select>
        </div>
        <div class="form-group">
          <label>Persona Tone</label>
          <input type="text" id="setAssTone" value="${esc((a.persona || {}).tone || '')}" placeholder="e.g. friendly, formal">
        </div>
        <div class="form-group">
          <label>Persona Role Description</label>
          <input type="text" id="setAssRole" value="${esc((a.persona || {}).role_description || '')}" placeholder="e.g. senior developer">
        </div>
        <div class="form-group">
          <label>Memory Scope</label>
          <select id="setAssMemory">
            <option value="assistant-only" ${(a.memory_policy || {}).scope === 'assistant-only' ? 'selected' : ''}>Assistant-only</option>
            <option value="global" ${(a.memory_policy || {}).scope === 'global' ? 'selected' : ''}>Global</option>
            <option value="both" ${(a.memory_policy || {}).scope === 'both' ? 'selected' : ''}>Both</option>
          </select>
        </div>
        <button class="primary-btn" id="saveAssSettingsBtn">Save Settings</button>
      </div>`;
    $('saveAssSettingsBtn').addEventListener('click', async () => {
      a.name = $('setAssName').value.trim();
      a.description = $('setAssDesc').value.trim();
      a.purpose = $('setAssPurpose').value.trim();
      a.default_model = $('setAssModel').value;
      const persona = a.persona || {};
      persona.tone = $('setAssTone').value.trim();
      persona.role_description = $('setAssRole').value.trim();
      a.persona = persona;
      a.memory_policy = { ...(a.memory_policy || {}), scope: $('setAssMemory').value };
      const res = await apiPUT('/api/assistants/' + a.id, a);
      if (res && res.error) alert('Error: ' + res.error);
      else {
        state.assistants = state.assistants.filter(x => x.id !== a.id);
        state.assistants.push(res || a);
        switchTab('overview');
      }
    });
  }

  function renderRuntime(panel, a) {
    panel.innerHTML = '<div class="animate-in">Loading runtime state...</div>';
    Promise.all([
      apiGET('/api/assistants/' + a.id + '/state'),
      apiGET('/api/assistants/' + a.id + '/jobs')
    ]).then(([st, jobs]) => {
      if (st.error) { panel.innerHTML = `<div class="empty-state error">${esc(st.error)}</div>`; return; }
      const s = st || {};
      const budget = s.token_budget_per_day || 10000;
      const used = s.tokens_used_today || 0;
      const actB = s.action_budget_per_day || 100;
      const actU = s.actions_used_today || 0;
      const nextCheck = s.next_check_at ? fmt(s.next_check_at) : '—';
      const lastCheck = s.last_check_at ? fmt(s.last_check_at) : '—';
      const controls = s.enabled !== false
        ? `<button class="add-btn" onclick="window._app.pauseAssistant('${a.id}')">Pause</button>`
        : `<button class="add-btn" onclick="window._app.resumeAssistant('${a.id}')">Resume</button>`;
      const runNow = `<button class="add-btn" onclick="window._app.triggerRun('${a.id}')">Run Now</button>`;
      const approvalsLink = `<button class="add-btn" style="margin-left:0.3rem" onclick="window._app.openApprovals('${a.id}')">Approvals</button>`;
      panel.innerHTML = `
        <div class="runtime-panel animate-in">
          <h4>Status</h4>
          <div class="runtime-row"><span>Enabled</span><span class="status-badge ${s.enabled !== false ? 'active' : 'paused'}">${s.enabled !== false ? 'Active' : 'Paused'}</span></div>
          <div class="runtime-row"><span>Status</span><span class="status-badge ${s.status || 'idle'}">${esc(s.status || 'idle')}</span></div>
          <div class="runtime-row"><span>Autonomy</span><span>${esc(s.autonomy_level || 'manual')}</span></div>
          <div class="runtime-row"><span>Interval</span><span>${s.interval || '—'}s</span></div>
          <div class="runtime-row"><span>Last check</span><span>${esc(lastCheck)}</span></div>
          <div class="runtime-row"><span>Next check</span><span>${esc(nextCheck)}</span></div>
          <h4 style="margin-top:1rem">Budget</h4>
          <div class="runtime-row"><span>Tokens today</span><span>${used} / ${budget}</span></div>
          <div class="runtime-row"><span>Actions today</span><span>${actU} / ${actB}</span></div>
          <div class="runtime-row"><span>Failures</span><span>${s.consecutive_failures || 0}</span></div>
          <h4 style="margin-top:1rem">Controls</h4>
          <div style="display:flex;gap:0.5rem;flex-wrap:wrap">${controls} ${runNow} ${approvalsLink}</div>
        </div>`;
    }).catch(err => {
      panel.innerHTML = `<div class="empty-state error">Load failed: ${esc(err.message || err)}</div>`;
    });
  }

  function renderApprovals(panel, a) {
    panel.innerHTML = '<div class="animate-in">Loading approvals...</div>';
    apiGET('/api/approvals?assistant=' + a.id).then(list => {
      if (list.error) { panel.innerHTML = `<div class="empty-state error">${esc(list.error)}</div>`; return; }
      panel.innerHTML = '';
      if (!list.length) {
        panel.innerHTML = `<div class="empty-state"><p>No approval requests.</p><p class="hint">Approvals are created when an assistant needs permission for a risky action.</p></div>`;
        return;
      }
      panel.innerHTML = list.map(r => {
        const status = r.status || 'pending';
        const actions = status === 'pending'
          ? `<button class="add-btn" onclick="window._app.resolveApproval('${r.id}', 'approved')">Approve</button> <button class="delete-btn" onclick="window._app.resolveApproval('${r.id}', 'denied')">Deny</button>`
          : '';
        return `<div class="approval-row">
          <div class="approval-meta"><span class="job-name">${esc(r.action_summary || 'Unknown action')}</span> <span class="status-badge ${status}">${status}</span></div>
          <div class="approval-meta"><span style="color:var(--text-muted);font-size:0.78rem">${esc(r.risk_level || 'low')} risk · requested ${fmt(r.requested_at)}</span></div>
          <div class="approval-actions" style="margin-top:0.4rem">${actions}</div>
        </div>`;
      }).join('');
    });
  }

  function renderDeveloper(panel, a) {
    if (!state.devMode) {
      panel.innerHTML = `<div class="empty-state">Enable Developer Mode to see raw data.</div>`;
      return;
    }
    panel.innerHTML = `
      <div class="dev-panel animate-in" style="margin-bottom:1.5rem">
        <h4>Assistant JSON</h4>
        <pre>${esc(JSON.stringify(a, null, 2))}</pre>
      </div>
      <div class="dev-panel">
        <h4>System Diagnostics</h4>
        <pre>Go version: ${navigator.userAgent.split('\n')[0]}
Assistants: ${state.assistants.length}
Models loaded: ${state.models.length}
Providers loaded: ${state.providers.length}</pre>
      </div>`;
  }

  // ---------- Global views ----------
  function loadGlobalJobs() {
    const el = $('globalJobsList');
    if (!state.assistants.length) { el.innerHTML = '<div class="empty-state">No assistants. Create one first.</div>'; return; }
    el.innerHTML = '<div class="animate-in">Loading jobs...</div>';
    Promise.all(state.assistants.map(a => apiGET('/api/assistants/' + a.id + '/jobs').then(jobs => ({ a, jobs }))))
      .then(results => {
        const all = [];
        results.forEach(({ a, jobs }) => {
          if (!Array.isArray(jobs)) return;
          jobs.forEach(j => { j._assistant_name = a.name; j._assistant_id = a.id; });
          all.push(...jobs);
        });
        if (!all.length) {
          el.innerHTML = '<div class="empty-state"><p>No jobs across all assistants.</p><p class="hint">Jobs are created inside each assistant.</p></div>';
          return;
        }
        all.sort((a, b) => new Date(b.created_at || 0) - new Date(a.created_at || 0));
        el.innerHTML = all.map(j => `
          <div class="job-row" style="cursor:pointer" onclick="window._app.openDetail('${j._assistant_id}', 'jobs')">
            <span class="job-name">${esc(j.name || 'Job')}</span>
            <span class="status-badge ${j.status}">${j.status}</span>
            <span style="color:var(--text-muted);font-size:0.78rem">${esc(j._assistant_name || '')}</span>
            <span style="color:var(--text-muted);font-size:0.78rem">${fmt(j.created_at)}</span>
          </div>`).join('');
      });
  }
  function loadChannels() {
    const el = $('channelsList');
    if (!state.assistants.length) { el.innerHTML = '<div class="empty-state">No assistants. Create one first.</div>'; return; }
    el.innerHTML = state.assistants.flatMap(a => (a.channels || []).map(c =>
      `<div class="channel-row"><span>${esc(a.name)}</span><span class="channel-type">${esc(c.channel_type)}</span><span class="status-badge ${c.status}">${c.status}</span></div>`
    )).join('') || '<div class="empty-state">No channels configured.</div>';
  }
  function loadTools() {
    const el = $('toolsList');
    if (!state.providers.length) return;
    apiGET('/api/tools').then(data => {
      if (!Array.isArray(data) || !data.length) { el.innerHTML = '<div class="empty-state">No tools registered.</div>'; return; }
      el.innerHTML = data.map(t => `<div class="tool-row"><span class="tool-name">${esc(t.name)}</span><span class="tool-desc">${esc(t.description || '')}</span></div>`).join('');
    });
  }

  async function loadSettings() {
    // Load providers and models
    if (!state.models.length) {
      const models = await apiGET('/api/models');
      if (Array.isArray(models)) state.models = models;
    }
    if (!state.providers.length) {
      const providers = await apiGET('/api/providers');
      if (Array.isArray(providers)) state.providers = providers;
    }
    // Populate default model select
    const sel = $('defaultModelSelect');
    if (sel) {
      sel.innerHTML = '<option value="">Choose default model...</option>' +
        state.models.map(m => `<option value="${esc(m.id)}">${esc(m.id)} — ${esc(m.provider)}</option>`).join('');
    }
    // Build provider key inputs for providers needing keys
    const keysEl = $('providerKeysList');
    if (keysEl) {
      const needingKey = state.providers.filter(p => p.name !== 'ollama-local');
      keysEl.innerHTML = needingKey.map(p => `
        <div class="key-row">
          <span>${esc(p.name)}</span>
          <input type="password" data-provider="${esc(p.name)}" placeholder="API key..." value="${state.providers.find(x => x.name === p.name && x.configured) ? '[configured]' : ''}">
        </div>`).join('');
    }
    // System diagnostics area
    const diagEl = $('systemDiagnostics');
    if (diagEl) {
      const configured = state.providers.filter(p => p.configured).length;
      diagEl.innerHTML = `<div class="ov-card" style="padding:1rem;margin-bottom:0.5rem">
        <p><strong>Providers configured:</strong> ${configured} / ${state.providers.length}</p>
        <p><strong>Models available:</strong> ${state.models.length}</p>
      </div>`;
    }
  }

  // ---------- Dev Mode ----------
  function toggleDevMode() {
    state.devMode = !state.devMode;
    const btn = $('devModeToggle');
    btn.textContent = state.devMode ? 'Developer Mode: ON' : 'Developer Mode: OFF';
    btn.classList.toggle('on', state.devMode);
    $('devModeCheckbox').checked = state.devMode;
    if (state.devMode) document.body.classList.add('dev-active');
    else document.body.classList.remove('dev-active');
    // Re-render if in assistant detail
    if (state.currentAssistantId && $('panel-developer')) {
      renderTabs();
    }
  }

  // ---------- Assistant Creation ----------
  function buildCreateModal() {
    const body = $('createAssistantBody');
    body.innerHTML = `
      <div class="form-group"><label>ID (unique, no spaces)</label><input type="text" id="newAssId" placeholder="my-assistant"></div>
      <div class="form-group"><label>Name</label><input type="text" id="newAssName" placeholder="My Assistant"></div>
      <div class="form-group"><label>Description</label><textarea id="newAssDesc" rows="2" placeholder="Short tagline"></textarea></div>
      <div class="form-group"><label>Purpose</label><textarea id="newAssPurpose" rows="3" placeholder="What is this assistant responsible for?"></textarea></div>
      <div class="form-row">
        <div class="form-group"><label>Persona Tone</label><input type="text" id="newAssTone" placeholder="e.g. friendly, formal"></div>
        <div class="form-group"><label>Persona Style</label><input type="text" id="newAssStyle" placeholder="e.g. concise, detailed"></div>
      </div>
      <div class="form-group"><label>Default Model</label><select id="newAssModel"><option value="">Choose model...</option>
        ${state.models.map(m => `<option value="${esc(m.id)}">${esc(m.id)} (${esc(m.provider)})</option>`).join('')}
      </select></div>
      <div class="form-group"><label>Initial Channels</label>
        <div style="display:flex;gap:1rem;flex-wrap:wrap">
          <label class="toggle-row"><input type="checkbox" class="newAssCh" value="web_ui" checked disabled> Web UI</label>
          <label class="toggle-row"><input type="checkbox" class="newAssCh" value="telegram"> Telegram</label>
          <label class="toggle-row"><input type="checkbox" class="newAssCh" value="discord"> Discord</label>
          <label class="toggle-row"><input type="checkbox" class="newAssCh" value="whatsapp"> WhatsApp</label>
          <label class="toggle-row"><input type="checkbox" class="newAssCh" value="slack"> Slack</label>
          <label class="toggle-row"><input type="checkbox" class="newAssCh" value="email"> Email</label>
        </div>
      </div>`;
  }

  async function saveNewAssistant() {
    const id = $('newAssId').value.trim();
    const name = $('newAssName').value.trim();
    if (!id || !name) { alert('ID and Name are required'); return; }
    const channels = [];
    document.querySelectorAll('.newAssCh:checked').forEach(cb => {
      channels.push({ channel_type: cb.value, status: cb.value === 'web_ui' ? 'connected' : 'needs_setup', created_at: new Date().toISOString() });
    });
    const payload = {
      id, name,
      description: $('newAssDesc').value.trim(),
      purpose: $('newAssPurpose').value.trim(),
      persona: { tone: $('newAssTone').value.trim(), style: $('newAssStyle').value.trim() },
      default_model: $('newAssModel').value,
      channels,
      memory_policy: { enabled: true, scope: 'assistant-only', auto_save: true },
      jobs_enabled: true,
    };
    const res = await apiPOST('/api/assistants', payload);
    if (res && res.error) { alert('Error: ' + res.error); return; }
    $('createAssistantDialog').close();
    await loadAssistants();
    switchView('assistants');
  }

  // ---------- App surface ----------
  window._app = {
    openDetail: (id, tab) => showDetail(id, tab),
    delAss: async id => {
      if (!confirm('Delete assistant ' + id + '?')) return;
      await apiDEL('/api/assistants/' + id);
      await loadAssistants();
    },
    addJob: id => {
      const name = prompt('Job name:');
      if (!name) return;
      apiPOST('/api/assistants/' + id + '/jobs', { name, purpose: 'manual', type: 'manual' }).then(() => {
        switchTab('jobs');
      });
    },
    pauseAssistant: async id => {
      await apiPOST('/api/assistants/' + id + '/pause', {});
      switchTab('runtime');
    },
    resumeAssistant: async id => {
      await apiPOST('/api/assistants/' + id + '/resume', {});
      switchTab('runtime');
    },
    triggerRun: async id => {
      await apiPOST('/api/assistants/' + id + '/run', {});
      switchTab('runtime');
    },
    openApprovals: id => {
      showDetail(id, 'approvals');
    },
    resolveApproval: async (id, status) => {
      await apiPOST('/api/approvals/' + id, { status });
      switchTab('approvals');
    },
  };

  // ---------- Init ----------
  function init() {
    document.querySelectorAll('.nav-item').forEach(btn => {
      btn.addEventListener('click', () => switchView(btn.dataset.view));
    });
    $('backBtn').addEventListener('click', setBack);
    $('createAssistantBtn').addEventListener('click', () => {
      buildCreateModal();
      $('createAssistantDialog').showModal();
    });
    $('closeCreateDialog').addEventListener('click', () => $('createAssistantDialog').close());
    $('cancelAssistantBtn').addEventListener('click', () => $('createAssistantDialog').close());
    $('saveAssistantBtn').addEventListener('click', saveNewAssistant);
    $('devModeToggle').addEventListener('click', toggleDevMode);
    $('devModeCheckbox').addEventListener('change', toggleDevMode);
    $('themeSelect').addEventListener('change', e => {
      document.body.classList.remove('dark', 'light');
      document.body.classList.add(e.target.value);
      state.theme = e.target.value;
    });
    $('saveProviderKeysBtn')?.addEventListener('click', async () => {
      const keys = {};
      document.querySelectorAll('#providerKeysList input').forEach(inp => {
        if (inp.value && inp.value !== '[configured]') {
          keys[inp.dataset.provider] = inp.value;
        }
      });
      const res = await apiPOST('/api/settings', { provider_keys: keys, dev_mode: state.devMode, theme: state.theme });
      if (res && res.error) alert('Error: ' + res.error);
      else alert('Provider keys saved.');
    });
    $('reloadProvidersBtn')?.addEventListener('click', loadSettings);

    loadAssistants();
    loadSettings();
    switchView('home');
  }

  init();
})();
