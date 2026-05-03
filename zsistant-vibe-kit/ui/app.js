(function() {
  'use strict';

  const state = {
    agentID: new URLSearchParams(location.search).get('agent') || '',
    conversation: [],
    isStreaming: false,
    devMode: false,
    theme: localStorage.getItem('zs_theme') || 'dark',
    accent: localStorage.getItem('zs_accent') || '#58a6ff',
    model: '',
    pendingAbort: null,
    providers: [],
    conversations: [],
    agents: [],
    currentConvId: null,
  };

  const $ = id => document.getElementById(id);
  const msgs = $('messages');
  const input = $('composerInput');
  const sendBtn = $('sendBtn');
  const stopBtn = $('stopBtn');
  const modelPicker = $('modelPicker');
  const modelSearch = $('modelSearch');
  const providerBadge = $('providerBadge');
  const inspector = $('inspector');
  const settingsDialog = $('settingsDialog');
  const devModeToggle = $('devModeToggle');

  function init() {
    applyTheme(state.theme);
    applyAccent(state.accent);
    loadProviders();
    loadModels();
    loadAgents();
    loadConversationsFromServer();
    if (state.agentID) {
      $('chatTitle').textContent = 'Chat with ' + state.agentID;
      loadEvents();
    } else {
      $('chatTitle').textContent = 'New Chat';
    }
    sendBtn.addEventListener('click', sendMessage);
    stopBtn.addEventListener('click', stopStream);
    input.addEventListener('keydown', e => {
      if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); sendMessage(); }
    });
    input.addEventListener('input', () => autoResize(input));
    $('settingsBtn').addEventListener('click', () => settingsDialog.showModal());
    $('closeSettings').addEventListener('click', () => settingsDialog.close());
    $('sidebarToggle').addEventListener('click', () => $('sidebar').classList.toggle('collapsed'));
    $('devModeBtn').addEventListener('click', toggleDevMode);
    $('devModeToggle').addEventListener('change', e => { state.devMode = e.target.checked; updateDevUI(); });
    $('closeInspector').addEventListener('click', () => inspector.classList.remove('visible'));
    $('newChatBtn').addEventListener('click', startNewChat);
    $('clearChatBtn').addEventListener('click', clearChat);
    $('saveKeysBtn').addEventListener('click', saveKeys);
    $('themeSelect').addEventListener('change', e => { state.theme = e.target.value; applyTheme(state.theme); localStorage.setItem('zs_theme', state.theme); });
    document.querySelectorAll('.accent-btn').forEach(btn => {
      btn.addEventListener('click', () => { state.accent = btn.dataset.accent; applyAccent(state.accent); localStorage.setItem('zs_accent', state.accent); updateAccentUI(); });
    });
    document.querySelectorAll('.tab-btn').forEach(btn => {
      btn.addEventListener('click', () => switchTab(btn.dataset.tab));
    });
    document.querySelectorAll('.ins-tab').forEach(btn => {
      btn.addEventListener('click', () => switchInsTab(btn.dataset.ins));
    });
    modelPicker.addEventListener('change', e => { state.model = e.target.value; });
    modelSearch.addEventListener('input', () => filterModels(modelSearch.value));
    $('searchInput').addEventListener('input', e => filterConversations(e.target.value));
    settingsDialog.addEventListener('click', e => {
      if (e.target === settingsDialog) settingsDialog.close();
    });
    $('showCreateAgent').addEventListener('click', () => $('createAgentForm').classList.toggle('hidden'));
    $('createAgentCancel').addEventListener('click', () => $('createAgentForm').classList.add('hidden'));
    $('createAgentSave').addEventListener('click', createAgent);
    loadSavedKeys();
    updateAccentUI();
    devModeToggle.checked = state.devMode;
    $('themeSelect').value = state.theme;
  }

  function applyTheme(t) {
    document.body.classList.remove('light', 'dark');
    if (t === 'system') {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
      document.body.classList.add(prefersDark ? 'dark' : 'light');
    } else {
      document.body.classList.add(t);
    }
  }
  function applyAccent(c) { document.documentElement.style.setProperty('--accent', c); }
  function updateAccentUI() {
    document.querySelectorAll('.accent-btn').forEach(b => b.classList.toggle('active', b.dataset.accent === state.accent));
  }

  async function loadProviders() {
    try {
      const res = await fetch('/api/providers');
      const data = await res.json();
      state.providers = data;
      populateProviderList();
      const healthy = data.filter(p => p.status === 'healthy' && p.name !== 'mock');
      if (healthy.length > 0) {
        providerBadge.textContent = healthy[0].name;
        providerBadge.classList.remove('hidden');
      }
    } catch (e) {
      if (state.devMode) logInspector('events', 'Failed to load providers: ' + e.message);
    }
  }

  async function loadModels() {
    try {
      const res = await fetch('/api/models?grouped=true');
      if (!res.ok) throw new Error('status ' + res.status);
      const groups = await res.json();
      if (Array.isArray(groups) && groups.length > 0 && groups[0].models) {
        modelPicker.innerHTML = '<option value="" disabled selected>Model</option>' +
          groups.map(g => {
            const opts = g.models.map(m => `<option value="${m.id}">${m.name}</option>`).join('');
            return `<optgroup label="${g.name}">${opts}</optgroup>`;
          }).join('');
        return;
      }
    } catch (e) {
      if (state.devMode) logInspector('events', 'Grouped models failed, trying flat: ' + e.message);
    }
    try {
      const res = await fetch('/api/models');
      if (!res.ok) throw new Error('status ' + res.status);
      const data = await res.json();
      if (Array.isArray(data) && data.length > 0) {
        modelPicker.innerHTML = '<option value="" disabled selected>Model</option>' +
          data.map(m => `<option value="${m.id}">${m.name}</option>`).join('');
        return;
      }
    } catch (e) {
      if (state.devMode) logInspector('events', 'Failed to load models: ' + e.message);
    }
    populateModelPickerStatic();
  }

  function filterModels(q) {
    const lower = q.toLowerCase();
    const groups = modelPicker.querySelectorAll('optgroup');
    groups.forEach(g => {
      let hasVisible = false;
      g.querySelectorAll('option').forEach(opt => {
        const match = opt.textContent.toLowerCase().includes(lower) || opt.value.toLowerCase().includes(lower);
        opt.hidden = !match;
        if (match) hasVisible = true;
      });
      g.hidden = !hasVisible;
    });
  }

  function populateModelPickerStatic() {
    const models = [
      {id:'gemma3:4b',name:'Gemma 3 4B (Ollama)'},
      {id:'gpt-4o-mini',name:'GPT-4o Mini (OpenAI)'},
      {id:'GLM-5.1',name:'GLM 5.1 (Z.AI)'},
      {id:'gpt-5.1',name:'GPT 5.1 (OpenCode)'},
      {id:'claude-3-5-sonnet',name:'Claude 3.5 Sonnet (Anthropic)'},
    ];
    modelPicker.innerHTML = '<option value="" disabled selected>Model</option>' +
      models.map(m => `<option value="${m.id}">${m.name}</option>`).join('');
  }

  function populateProviderList() {
    const container = $('providerList');
    if (!container) return;
    container.innerHTML = state.providers.map(p => {
      const cls = p.status === 'healthy' ? 'ok' : 'missing';
      const label = p.status === 'healthy' ? 'Ready' : 'No key';
      return `<div class="provider-row"><span class="p-name">${p.name}</span><span class="p-status ${cls}">${label}</span></div>`;
    }).join('');
  }

  async function loadAgents() {
    try {
      const res = await fetch('/api/agents');
      if (!res.ok) return;
      state.agents = await res.json();
      renderAgentList();
    } catch (e) {
      if (state.devMode) logInspector('events', 'Failed to load agents: ' + e.message);
    }
  }

  function renderAgentList() {
    const list = $('agentList');
    if (!state.agents || state.agents.length === 0) {
      list.innerHTML = '<div class="agent-empty">No agents yet</div>';
      return;
    }
    list.innerHTML = state.agents.map(a => {
      const active = a.id === state.agentID ? 'active' : '';
      return `<div class="agent-item ${active}" data-id="${a.id}"><span class="agent-name">${escapeHtml(a.name)}</span><span class="agent-role">${escapeHtml(a.role || '')}</span></div>`;
    }).join('');
    list.querySelectorAll('.agent-item').forEach(el => {
      el.addEventListener('click', () => selectAgent(el.dataset.id));
    });
  }

  function selectAgent(id) {
    const agent = state.agents.find(a => a.id === id);
    if (!agent) return;
    state.agentID = id;
    $('chatTitle').textContent = 'Chat with ' + agent.name;
    history.pushState({}, '', '/?agent=' + encodeURIComponent(id));
    renderAgentList();
  }

  async function createAgent() {
    const id = $('agentFormId').value.trim();
    const name = $('agentFormName').value.trim();
    const role = $('agentFormRole').value.trim();
    if (!id || !name) return;
    try {
      const res = await fetch('/api/agents', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id, name, role }),
      });
      if (res.ok || res.status === 201) {
        $('createAgentForm').classList.add('hidden');
        $('agentFormId').value = '';
        $('agentFormName').value = '';
        $('agentFormRole').value = '';
        loadAgents();
      }
    } catch (e) {
      if (state.devMode) logInspector('events', 'Failed to create agent: ' + e.message);
    }
  }

  function startNewChat() {
    state.conversation = [];
    state.currentConvId = null;
    state.agentID = '';
    history.pushState({}, '', '/');
    $('chatTitle').textContent = 'New Chat';
    msgs.innerHTML = '<div class="empty-state"><h3>How can I help you today?</h3><p>Select a model and start a conversation.</p></div>';
    providerBadge.classList.add('hidden');
    renderAgentList();
  }
  function clearChat() {
    state.conversation = [];
    msgs.innerHTML = '';
    if (!state.agentID) {
      msgs.innerHTML = '<div class="empty-state"><h3>How can I help you today?</h3><p>Select a model and start a conversation.</p></div>';
    }
  }
  function autoResize(ta) { ta.style.height = 'auto'; ta.style.height = ta.scrollHeight + 'px'; }

  async function sendMessage() {
    const text = input.value.trim();
    if (!text || state.isStreaming) return;
    if (!state.agentID) {
      state.agentID = 'web-' + Date.now();
    }

    const userMsg = { role: 'user', content: text, id: 'u-' + Date.now() };
    state.conversation.push(userMsg);
    appendMessage(userMsg);
    input.value = '';
    autoResize(input);

    const botMsg = { role: 'assistant', content: '', id: 'a-' + Date.now(), streaming: true };
    state.conversation.push(botMsg);
    const botEl = appendMessage(botMsg);
    botEl.classList.add('streaming');

    sendBtn.classList.add('hidden');
    stopBtn.classList.remove('hidden');
    state.isStreaming = true;

    if (state.devMode) logInspector('request', JSON.stringify({agent_id: state.agentID, message: text, model: state.model}, null, 2));

    try {
      await streamWithFetch(text, botEl, botMsg);
    } catch (e) {
      botEl.classList.remove('streaming');
      botEl.innerHTML = escapeHtml('Error: ' + e.message);
      botEl.classList.add('error');
      if (state.devMode) logInspector('events', 'Stream error: ' + e.message);
    } finally {
      state.isStreaming = false;
      sendBtn.classList.remove('hidden');
      stopBtn.classList.add('hidden');
      saveConversation();
      loadEvents();
    }
  }

  async function streamWithFetch(text, botEl, botMsg) {
    const controller = new AbortController();
    state.pendingAbort = controller;

    const res = await fetch('/api/chat/stream', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ agent_id: state.agentID, message: text }),
      signal: controller.signal,
    });

    if (!res.ok) {
      const err = await res.json().catch(() => ({error: res.statusText}));
      throw new Error(err.error || res.statusText);
    }

    const reader = res.body.getReader();
    const decoder = new TextDecoder();
    let buffer = '';
    let fullText = '';

    let currentEventType = '';
    let currentEventData = [];

    function flushEvent() {
      if (!currentEventType && currentEventData.length === 0) return;
      const etype = currentEventType || 'message';
      const data = currentEventData.join('\n');
      currentEventType = '';
      currentEventData = [];
      if (state.devMode) logInspector('events', `event: ${etype}\ndata: ${data}`);
      if (etype === 'chunk') {
        const chunk = JSON.parse(data);
        fullText += chunk;
        botMsg.content = fullText;
        botEl.innerHTML = escapeHtml(fullText) + renderMeta(botMsg);
        botEl.scrollIntoView({ behavior: 'smooth', block: 'end' });
      } else if (etype === 'done') {
        const info = JSON.parse(data);
        botMsg.provider = info.provider;
        botMsg.streaming = false;
        botEl.classList.remove('streaming');
        botEl.innerHTML = escapeHtml(fullText) + renderMeta(botMsg);
        providerBadge.textContent = info.provider || 'model';
        providerBadge.classList.remove('hidden');
      } else if (etype === 'error') {
        const err = JSON.parse(data);
        throw new Error(err.error || 'Stream error');
      } else if (etype === 'provider_error') {
        const err = JSON.parse(data);
        if (state.devMode) logInspector('events', 'Provider error: ' + (err.error || ''));
      }
    }

    while (true) {
      const { done, value } = await reader.read();
      if (done) { flushEvent(); break; }
      buffer += decoder.decode(value, { stream: true });
      const lines = buffer.split('\n');
      buffer = lines.pop();

      for (const line of lines) {
        if (line.startsWith('event:')) {
          flushEvent();
          currentEventType = line.slice(6).trim();
        } else if (line.startsWith('data:')) {
          currentEventData.push(line.slice(5).trim());
        } else if (line.trim() === '') {
          flushEvent();
        }
      }
    }
    flushEvent();
  }

  function stopStream() {
    if (state.pendingAbort) {
      state.pendingAbort.abort();
      state.pendingAbort = null;
    }
    state.isStreaming = false;
    sendBtn.classList.remove('hidden');
    stopBtn.classList.add('hidden');
  }

  function appendMessage(msg) {
    const empty = msgs.querySelector('.empty-state');
    if (empty) empty.remove();

    const div = document.createElement('div');
    div.className = 'message ' + msg.role;
    if (msg.streaming) div.classList.add('streaming');
    div.id = msg.id;
    div.innerHTML = escapeHtml(msg.content) + renderMeta(msg) + renderActions(msg);
    msgs.appendChild(div);
    div.scrollIntoView({ behavior: 'smooth', block: 'end' });
    return div;
  }
  function renderMeta(msg) {
    if (!msg.provider && !msg.latency) return '';
    const parts = [];
    if (msg.provider) parts.push(msg.provider);
    if (msg.model) parts.push(msg.model);
    if (msg.latency) parts.push(msg.latency + 'ms');
    return `<div class="msg-meta">${parts.join(' \u00b7 ')}</div>`;
  }
  function renderActions(msg) {
    return `<div class="msg-actions">
      <button onclick="window.copyMsg('${msg.id}')">Copy</button>
      ${msg.role === 'assistant' ? `<button onclick="window.regenerateMsg('${msg.id}')">Regenerate</button>` : ''}
    </div>`;
  }
  window.copyMsg = id => {
    const msg = state.conversation.find(m => m.id === id);
    if (msg) navigator.clipboard.writeText(msg.content);
  };
  window.regenerateMsg = id => {
    const idx = state.conversation.findIndex(m => m.id === id);
    if (idx > 0 && state.conversation[idx - 1].role === 'user') {
      const userText = state.conversation[idx - 1].content;
      state.conversation.splice(idx, 1);
      const el = document.getElementById(id);
      if (el) el.remove();
      input.value = userText;
      sendMessage();
    }
  };
  function escapeHtml(t) {
    const d = document.createElement('div');
    d.textContent = t;
    return d.innerHTML.replace(/\n/g, '<br>');
  }

  function saveConversation() {
    if (state.conversation.length === 0) return;
    const conv = {
      id: state.currentConvId || 'c-' + Date.now(),
      agent_id: state.agentID,
      title: state.conversation[0].content.slice(0, 40) || 'Chat',
      updatedAt: new Date().toISOString(),
      messages: state.conversation,
    };
    state.currentConvId = conv.id;
    const existing = state.conversations.findIndex(c => c.id === conv.id);
    if (existing >= 0) state.conversations[existing] = conv;
    else state.conversations.unshift(conv);
    localStorage.setItem('zs_conversations', JSON.stringify(state.conversations.slice(0, 50)));
    fetch('/api/conversations', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(conv),
    }).catch(() => {});
    renderConversationList();
  }
  async function loadConversationsFromServer() {
    try {
      const res = await fetch('/api/conversations');
      if (!res.ok) throw new Error('status ' + res.status);
      const data = await res.json();
      if (Array.isArray(data) && data.length > 0) {
        state.conversations = data;
      }
    } catch (e) {
      state.conversations = JSON.parse(localStorage.getItem('zs_conversations') || '[]');
      if (state.devMode) logInspector('events', 'Server convs failed, using cache: ' + e.message);
    }
    renderConversationList();
  }

  function renderConversationList() {
    const list = $('conversationList');
    if (state.conversations.length === 0) {
      list.innerHTML = '<div style="padding:1rem; opacity:0.5; font-size:0.85rem;">No conversations yet</div>';
      return;
    }
    list.innerHTML = state.conversations.map(c => {
      const active = c.id === state.currentConvId ? 'active' : '';
      const time = new Date(c.updatedAt).toLocaleTimeString([], {hour:'2-digit', minute:'2-digit'});
      return `<div class="conv-item ${active}" data-id="${c.id}">${escapeHtml(c.title)}<span class="conv-time">${time}</span><button class="conv-del" onclick="window.deleteConv('${c.id}',event)">&times;</button></div>`;
    }).join('');
    list.querySelectorAll('.conv-item').forEach(el => {
      el.addEventListener('click', () => loadConversation(el.dataset.id));
    });
  }
  function loadConversation(id) {
    const conv = state.conversations.find(c => c.id === id);
    if (!conv) return;
    state.currentConvId = id;
    state.agentID = conv.agent_id || conv.agentID || '';
    state.conversation = conv.messages || [];
    $('chatTitle').textContent = (conv.agent_id || conv.agentID || '').startsWith('web-') ? 'Chat' : 'Chat with ' + (conv.agent_id || conv.agentID || '');
    msgs.innerHTML = '';
    state.conversation.forEach(m => appendMessage(m));
    renderConversationList();
    renderAgentList();
  }
  window.deleteConv = (id, ev) => {
    ev.stopPropagation();
    state.conversations = state.conversations.filter(c => c.id !== id);
    localStorage.setItem('zs_conversations', JSON.stringify(state.conversations.slice(0, 50)));
    if (state.currentConvId === id) startNewChat();
    else renderConversationList();
  };
  function filterConversations(q) {
    const list = $('conversationList');
    const items = list.querySelectorAll('.conv-item');
    items.forEach(el => {
      el.style.display = el.textContent.toLowerCase().includes(q.toLowerCase()) ? '' : 'none';
    });
  }

  async function loadEvents() {
    if (!state.agentID) return;
    try {
      const res = await fetch('/api/jobs/' + encodeURIComponent(state.agentID));
      if (!res.ok) return;
      const data = await res.json();
    } catch (e) { /* silent */ }
  }

  function switchTab(tab) {
    document.querySelectorAll('.tab-btn').forEach(b => b.classList.toggle('active', b.dataset.tab === tab));
    document.querySelectorAll('.tab-panel').forEach(p => p.classList.toggle('active', p.id === 'tab-' + tab));
  }
  function switchInsTab(tab) {
    document.querySelectorAll('.ins-tab').forEach(b => b.classList.toggle('active', b.dataset.ins === tab));
    document.querySelectorAll('.ins-panel').forEach(p => p.classList.toggle('active', p.id === 'ins' + tab.charAt(0).toUpperCase() + tab.slice(1)));
  }
  function toggleDevMode() {
    state.devMode = !state.devMode;
    devModeToggle.checked = state.devMode;
    updateDevUI();
  }
  function updateDevUI() {
    const btn = $('devModeBtn');
    btn.style.color = state.devMode ? 'var(--accent)' : '';
    if (state.devMode) {
      inspector.classList.add('visible');
      loadProviderHealth();
    } else {
      inspector.classList.remove('visible');
    }
  }
  function logInspector(panel, text) {
    const el = $('ins' + panel.charAt(0).toUpperCase() + panel.slice(1));
    if (!el) return;
    el.textContent += '\n[' + new Date().toLocaleTimeString() + '] ' + text + '\n';
    el.scrollTop = el.scrollHeight;
  }
  async function loadProviderHealth() {
    try {
      const res = await fetch('/api/providers');
      const data = await res.json();
      $('insHealth').textContent = JSON.stringify(data, null, 2);
    } catch (e) {
      $('insHealth').textContent = 'Failed to load health: ' + e.message;
    }
  }

  async function loadSavedKeys() {
    try {
      const res = await fetch('/api/settings');
      if (!res.ok) return;
      const keys = await res.json();
      $('keyOllama').value = keys['ollama_api_key'] || '';
      $('keyOpenAI').value = keys['openai_api_key'] || '';
      $('keyZAI').value = keys['zai_api_key'] || '';
      $('keyOpenCode').value = keys['opencode_api_key'] || '';
    } catch (e) {
      if (state.devMode) logInspector('events', 'Failed to load keys: ' + e.message);
    }
  }
  async function saveKeys() {
    const keys = {
      ollama_api_key: $('keyOllama').value,
      openai_api_key: $('keyOpenAI').value,
      zai_api_key: $('keyZAI').value,
      opencode_api_key: $('keyOpenCode').value,
    };
    try {
      const res = await fetch('/api/settings', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ secrets: keys }),
      });
      if (res.ok) {
        alert('Keys saved. Refresh the page to activate providers.');
        loadProviders();
      } else {
        alert('Failed to save keys.');
      }
    } catch (e) {
      alert('Error: ' + e.message);
    }
  }

  init();
})();
