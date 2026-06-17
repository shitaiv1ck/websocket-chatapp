const API_BASE = 'http://localhost:8080';
const WS_URL = 'ws://localhost:8080/ws';

const state = {
    currentUser: null,
    activeTab: 'chats',
    chats: [],
    friends: [],
    friendRequests: [],
    users: [],
    activeChat: null,
    messages: [],
    ws: null,
    reconnectAttempts: 0,
    reconnectTimer: null,
    shouldConnect: false
};

function getCsrfToken() {
    const match = document.cookie.match(/(^|;)\s*csrf_token\s*=\s*([^;]+)/);
    return match ? match.pop() : '';
}

function getInitials(username) {
    if (!username) return '?';
    return username.substring(0, 2).toUpperCase();
}

function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatMessageTime(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
}

function formatDate(dateString) {
    if (!dateString) return '';
    const date = new Date(dateString);
    const now = new Date();
    const diff = now - date;
    if (diff < 60000) return 'только что';
    if (diff < 3600000) return Math.floor(diff / 60000) + ' мин назад';
    if (diff < 86400000) return date.toLocaleTimeString('ru-RU', { hour: '2-digit', minute: '2-digit' });
    return date.toLocaleDateString('ru-RU', { day: '2-digit', month: '2-digit' });
}

async function apiRequest(endpoint, method = 'GET', body = null) {
    const options = {
        method: method,
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' }
    };

    const requiresCSRF = ['POST', 'PATCH', 'DELETE'].indexOf(method) !== -1;
    const isExempt = (endpoint === '/api/sessions' && method === 'POST') ||
        (endpoint === '/api/users' && method === 'POST');

    if (requiresCSRF && !isExempt) {
        options.headers['X-CSRF-Token'] = getCsrfToken();
    }

    if (body) {
        options.body = JSON.stringify(body);
    }

    try {
        const response = await fetch(API_BASE + endpoint, options);

        if (response.status === 204) {
            return null;
        }

        if (response.status === 201 && endpoint === '/api/sessions' && method === 'POST') {
            return null;
        }

        let data = {};
        try {
            data = await response.json();
        } catch (e) {
            data = {};
        }

        if (!response.ok) {
            const errorMsg = data.error || data.message || 'HTTP ' + response.status;

            if (response.status === 401) {
                disconnectWebSocket();
                state.currentUser = null;
                showAuthView();
            }

            throw new Error(errorMsg);
        }

        return data;
    } catch (err) {
        if (err.message && err.message.startsWith('HTTP')) {
            throw err;
        }
        throw err;
    }
}

function connectWebSocket() {
    if (!state.shouldConnect) return;
    if (state.ws && (state.ws.readyState === 1 || state.ws.readyState === 0)) return;

    if (state.reconnectTimer) {
        clearTimeout(state.reconnectTimer);
        state.reconnectTimer = null;
    }
    if (state.ws) {
        try { state.ws.close(); } catch (e) { }
        state.ws = null;
    }

    try {
        state.ws = new WebSocket(WS_URL);

        state.ws.onopen = function () {
            state.reconnectAttempts = 0;
        };

        state.ws.onmessage = function (event) {
            try {
                const msg = JSON.parse(event.data);
                handleWebSocketMessage(msg);
            } catch (err) {
                // ignore
            }
        };

        state.ws.onclose = function () {
            state.ws = null;
            if (state.shouldConnect) {
                const delay = Math.min(3000 * Math.pow(1.5, state.reconnectAttempts), 10000);
                state.reconnectAttempts++;
                state.reconnectTimer = setTimeout(connectWebSocket, delay);
            }
        };

        state.ws.onerror = function () {

        };
    } catch (err) {
        state.ws = null;
    }
}

function disconnectWebSocket() {
    state.shouldConnect = false;
    if (state.reconnectTimer) {
        clearTimeout(state.reconnectTimer);
        state.reconnectTimer = null;
    }
    if (state.ws) {
        try { state.ws.close(1000); } catch (e) { }
        state.ws = null;
    }
}

function handleWebSocketMessage(message) {
    const type = message.type;
    const content = message.content;

    switch (type) {
        case 'message.received':
            handleNewMessage(content);
            break;
        case 'message.sent':
            handleSentMessage(content);
            break;
        case 'friend_request.received':
        case 'friend_request.sent':
            handleFriendRequestAdded(content);
            break;
        case 'friend_request.accepted':
            handleFriendRequestAccepted(content);
            break;
        case 'friendship.added':
            handleFriendshipAdded(content);
            break;
        case 'friend_request.declined':
            handleFriendRequestDeclined(content);
            break;
        case 'friendship.deleted':
            handleFriendshipDeleted(content);
            break;
        case 'user.change_status':
        case 'user.change_username':
        case 'user.created':
            handleUserUpdate(content);
            break;
        case 'chat.deleted':
            handleChatDeleted(content);
            break;
        case 'chat.created':
            handleChatCreated(content);
            break;
        default:

    }
}

function handleFriendRequestAdded(content) {
    const requestId = Number(content.friend_request_id);

    if (isNaN(requestId)) {
        return;
    }

    const exists = state.friendRequests.some(function (r) { return r.id === requestId; });
    if (exists) {
        return;
    }

    var normalizedRequest = {
        id: requestId,
        from_user: content.from_user,
        to_user: content.to_user,
        created_at: content.created_at
    };

    state.friendRequests.push(normalizedRequest);
    renderSidebarContent();
    updateRequestsBadge();
}

function handleFriendRequestAccepted(content) {
    const friendshipId = Number(content.friendship_id);
    const requestId = Number(content.friend_request_id);

    if (isNaN(friendshipId)) {
        return;
    }

    var normalizedFriendship = {
        id: friendshipId,
        first_user: content.first_user,
        second_user: content.second_user
    };

    const exists = state.friends.some(function (f) { return f.id === friendshipId; });
    if (!exists) {
        state.friends.push(normalizedFriendship);
    }

    if (!isNaN(requestId)) {
        state.friendRequests = state.friendRequests.filter(function (r) {
            return r.id !== requestId;
        });
    }

    renderSidebarContent();
    updateRequestsBadge();
}

function handleFriendshipAdded(content) {
    const friendshipId = Number(content.friendship_id);

    if (isNaN(friendshipId)) {
        return;
    }

    var normalizedFriendship = {
        id: friendshipId,
        first_user: content.first_user,
        second_user: content.second_user
    };

    const exists = state.friends.some(function (f) { return f.id === friendshipId; });
    if (!exists) {
        state.friends.push(normalizedFriendship);
    }
    renderSidebarContent();
}

function handleFriendRequestDeclined(content) {
    const requestId = Number(content.friend_request_id);

    if (isNaN(requestId)) {
        return;
    }

    state.friendRequests = state.friendRequests.filter(function (r) { return r.id !== requestId; });
    renderSidebarContent();
    updateRequestsBadge();
}

function handleFriendshipDeleted(content) {
    const friendshipId = Number(content.friendship_id);

    if (isNaN(friendshipId)) {
        return;
    }

    state.friends = state.friends.filter(function (f) { return f.id !== friendshipId; });
    renderSidebarContent();
}

function handleChatDeleted(content) {
    const chatId = Number(content.chat_id);

    if (isNaN(chatId)) {
        return;
    }

    state.chats = state.chats.filter(function (c) { return c.id !== chatId; });

    if (state.activeChat && state.activeChat.id === chatId) {
        state.activeChat = null;
        state.messages = [];
        showEmptyChat();
    }

    renderSidebarContent();
}

function handleChatCreated(content) {
    var chat = {
        id: Number(content.id),
        first_user: content.first_user,
        second_user: content.second_user,
        last_message_content: content.last_message_content || null,
        last_message_at: content.last_message_at || new Date().toISOString()
    };

    var exists = state.chats.some(function(c) {
        return c.id === chat.id;
    });

    if (!exists) {
        state.chats.unshift(chat);
        
        if (state.activeTab === 'chats') {
            renderSidebarContent();
        }
        
        if (state.chats.length === 1 && state.activeTab === 'chats') {
            renderSidebarContent();
        }
    }
}

function handleUserUpdate(userData) {
    const userId = Number(userData.id);

    for (var i = 0; i < state.friends.length; i++) {
        var f = state.friends[i];
        if (f.first_user.id === userId) f.first_user = Object.assign({}, f.first_user, userData);
        if (f.second_user.id === userId) f.second_user = Object.assign({}, f.second_user, userData);
    }

    for (var j = 0; j < state.chats.length; j++) {
        var c = state.chats[j];
        if (c.first_user.id === userId) c.first_user = Object.assign({}, c.first_user, userData);
        if (c.second_user.id === userId) c.second_user = Object.assign({}, c.second_user, userData);
    }

    if (state.currentUser && state.currentUser.id === userId) {
        state.currentUser = Object.assign({}, state.currentUser, userData);
        updateUserDisplay();
    }

    if (state.activeChat) {
        var other = state.activeChat.first_user.id === state.currentUser.id ?
            state.activeChat.second_user : state.activeChat.first_user;
        if (other.id === userId) {
            if (state.activeChat.first_user.id === userId) {
                state.activeChat.first_user = Object.assign({}, state.activeChat.first_user, userData);
            } else {
                state.activeChat.second_user = Object.assign({}, state.activeChat.second_user, userData);
            }
            var statusEl = document.getElementById('chat-header-status');
            if (statusEl) {
                statusEl.textContent = userData.is_online ? 'Онлайн' : 'Офлайн';
                statusEl.className = 'chat-header-status' + (userData.is_online ? '' : ' offline');
            }
        }
    }

    renderSidebarContent();
}

function handleNewMessage(data) {
    const chatId = Number(data.chat_id);

    if (state.activeChat && chatId === state.activeChat.id) {
        const exists = state.messages.some(function (m) { return m.id === data.id; });
        if (!exists) {
            state.messages.push(data);
            renderMessages();
            scrollToBottom();
        }
    }
    updateChatLastMessage(chatId, data.content, data.created_at);
}

function handleSentMessage(data) {
    const chatId = Number(data.chat_id);

    if (state.activeChat && chatId === state.activeChat.id) {
        const exists = state.messages.some(function (m) { return m.id === data.id; });
        if (!exists) {
            state.messages.push(data);
            renderMessages();
            scrollToBottom();
        }
    }
    updateChatLastMessage(chatId, data.content, data.created_at);
}

function showAuthView() {
    document.getElementById('auth-view').classList.remove('hidden');
    document.getElementById('app-view').classList.add('hidden');
    state.currentUser = null;
}

function showAppView() {
    document.getElementById('auth-view').classList.add('hidden');
    document.getElementById('app-view').classList.remove('hidden');
    updateUserDisplay();
    if (state.shouldConnect && (!state.ws || state.ws.readyState === 3)) {
        connectWebSocket();
    }
}

function updateUserDisplay() {
    if (state.currentUser) {
        document.getElementById('user-name').textContent = state.currentUser.username;
        document.getElementById('user-avatar').textContent = getInitials(state.currentUser.username);
        var statusEl = document.getElementById('user-status');
        statusEl.textContent = state.currentUser.is_online ? 'Онлайн' : 'Офлайн';
        statusEl.className = 'user-status' + (state.currentUser.is_online ? '' : ' offline');
    }
}

function updateRequestsBadge() {
    var badge = document.getElementById('requests-badge');
    if (!badge) return;
    var count = state.friendRequests.length;
    if (count > 0) {
        badge.textContent = count;
        badge.classList.remove('hidden');
    } else {
        badge.classList.add('hidden');
    }
}

function switchTab(tabName) {
    state.activeTab = tabName;
    var tabs = document.querySelectorAll('.tab-btn');
    for (var i = 0; i < tabs.length; i++) {
        tabs[i].classList.toggle('active', tabs[i].dataset.tab === tabName);
    }
    renderSidebarContent();
}

function renderSidebarContent() {
    var container = document.getElementById('sidebar-content');
    if (!container) return;

    if (state.activeTab === 'chats') renderChatsList(container);
    else if (state.activeTab === 'friends') renderFriendsList(container);
    else if (state.activeTab === 'requests') renderRequestsList(container);
    else if (state.activeTab === 'users') renderUsersList(container);
}

function renderChatsList(container) {
    if (state.chats.length === 0) {
        container.innerHTML = '<div style="padding:40px;text-align:center;color:#808080">Нет чатов</div>';
        return;
    }
    var html = '';
    for (var i = 0; i < state.chats.length; i++) {
        var chat = state.chats[i];
        var other = chat.first_user.id === state.currentUser.id ? chat.second_user : chat.first_user;
        var isActive = state.activeChat && state.activeChat.id === chat.id;
        html += '<div class="list-item' + (isActive ? ' active' : '') + '" data-chat-id="' + chat.id + '">' +
            '<div class="list-item-avatar">' + getInitials(other.username) + '</div>' +
            '<div class="list-item-content"><div class="list-item-header"><div class="list-item-name">' + escapeHtml(other.username) + '</div>' +
            '<div class="list-item-time">' + formatDate(chat.last_message_at) + '</div></div>' +
            '<div class="list-item-preview">' + (chat.last_message_content || 'Нет сообщений') + '</div></div></div>';
    }
    container.innerHTML = html;
}

function renderFriendsList(container) {
    if (state.friends.length === 0) {
        container.innerHTML = '<div style="padding:40px;text-align:center;color:#808080">Нет друзей</div>';
        return;
    }
    var html = '';
    for (var i = 0; i < state.friends.length; i++) {
        var f = state.friends[i];
        var other = f.first_user.id === state.currentUser.id ? f.second_user : f.first_user;
        html += '<div class="list-item" data-friend-id="' + other.id + '">' +
            '<div class="list-item-avatar">' + getInitials(other.username) + '</div>' +
            '<div class="list-item-content"><div class="list-item-header"><div class="list-item-name">' + escapeHtml(other.username) + '</div>' +
            '<div class="user-status' + (other.is_online ? '' : ' offline') + '">' + (other.is_online ? 'Онлайн' : 'Офлайн') + '</div></div></div>' +
            '<button class="btn btn-small btn-delete btn-delete-friend" data-friendship-id="' + f.id + '">Удалить</button></div>';
    }
    container.innerHTML = html;
}

function renderRequestsList(container) {
    if (!state.currentUser) {
        container.innerHTML = '<div style="padding:40px;text-align:center;color:#808080">Загрузка...</div>';
        return;
    }

    var incoming = [];
    var outgoing = [];

    for (var i = 0; i < state.friendRequests.length; i++) {
        var r = state.friendRequests[i];
        if (r.to_user && r.to_user.id === state.currentUser.id) {
            incoming.push(r);
        } else if (r.from_user && r.from_user.id === state.currentUser.id) {
            outgoing.push(r);
        }
    }

    if (incoming.length === 0 && outgoing.length === 0) {
        container.innerHTML = '<div style="padding:40px;text-align:center;color:#808080">Нет заявок</div>';
        return;
    }

    var html = '';
    if (incoming.length > 0) {
        html += '<div class="request-section"><div class="request-section-title">ВХОДЯЩИЕ ЗАЯВКИ (' + incoming.length + ')</div>';
        for (var j = 0; j < incoming.length; j++) {
            var r = incoming[j];
            html += '<div class="list-item" style="flex-direction:column;align-items:stretch">' +
                '<div style="display:flex;align-items:center;gap:12px"><div class="list-item-avatar">' + getInitials(r.from_user.username) + '</div>' +
                '<div class="list-item-content"><div class="list-item-name">' + escapeHtml(r.from_user.username) + '</div></div></div>' +
                '<div class="friend-request-actions"><button class="btn btn-small btn-accept" data-request-id="' + r.id + '">Принять</button>' +
                '<button class="btn btn-small btn-decline" data-request-id="' + r.id + '">Отклонить</button></div></div>';
        }
        html += '</div>';
    }
    if (outgoing.length > 0) {
        html += '<div class="request-section"><div class="request-section-title">ИСХОДЯЩИЕ ЗАЯВКИ (' + outgoing.length + ')</div>';
        for (var k = 0; k < outgoing.length; k++) {
            var req = outgoing[k];
            html += '<div class="list-item" style="flex-direction:column;align-items:stretch">' +
                '<div style="display:flex;align-items:center;gap:12px"><div class="list-item-avatar">' + getInitials(req.to_user.username) + '</div>' +
                '<div class="list-item-content"><div class="list-item-name">' + escapeHtml(req.to_user.username) + '</div></div></div>' +
                '<div style="padding:8px 0;color:#808080;font-size:0.85em">Ожидает ответа...</div></div>';
        }
        html += '</div>';
    }
    container.innerHTML = html;
}

function renderUsersList(container) {
    container.innerHTML = '<div class="search-box"><input type="text" id="user-search" placeholder="Поиск пользователей..."></div><div id="users-list" style="padding-bottom:20px;"></div>';
    var searchInput = document.getElementById('user-search');
    var timeout;
    searchInput.addEventListener('input', function () {
        clearTimeout(timeout);
        timeout = setTimeout(function () { searchUsers(searchInput.value.trim()); }, 300);
    });
    searchUsers('');
}

async function searchUsers(query) {
    var usersList = document.getElementById('users-list');
    if (!usersList) return;

    try {
        var url = '/api/users' + (query ? '?search=' + encodeURIComponent(query) : '');
        var users = await apiRequest(url);
        state.users = users || [];

        var friendIds = {};
        for (var i = 0; i < state.friends.length; i++) {
            var f = state.friends[i];
            friendIds[f.first_user.id === state.currentUser.id ? f.second_user.id : f.first_user.id] = true;
        }

        var html = '';
        for (var j = 0; j < state.users.length; j++) {
            var u = state.users[j];
            if (u.id === state.currentUser.id) continue;
            var isFriend = friendIds[u.id];
            html += '<div class="list-item"><div class="list-item-avatar">' + getInitials(u.username) + '</div>' +
                '<div class="list-item-content"><div class="list-item-header"><div class="list-item-name">' + escapeHtml(u.username) + '</div>' +
                '<div class="user-status' + (u.is_online ? '' : ' offline') + '">' + (u.is_online ? 'Онлайн' : 'Офлайн') + '</div></div></div>' +
                (isFriend ? '' : '<button class="btn-add-friend" data-user-id="' + u.id + '">Добавить</button>') + '</div>';
        }
        usersList.innerHTML = html || '<div style="padding:20px;text-align:center;color:#808080">Никого не найдено</div>';
    } catch (err) {

    }
}

function initEventDelegation() {
    var container = document.getElementById('sidebar-content');

    container.addEventListener('click', function (e) {
        if (e.target.closest('.list-item[data-chat-id]')) {
            var item = e.target.closest('.list-item[data-chat-id]');
            var chatId = Number(item.dataset.chatId);
            selectChat(chatId);
            return;
        }

        if (e.target.classList.contains('btn-delete-friend')) {
            e.stopPropagation();
            var friendshipId = Number(e.target.dataset.friendshipId);
            if (confirm('Удалить друга?')) {
                deleteFriend(friendshipId);
            }
            return;
        }

        if (e.target.closest('.list-item[data-friend-id]')) {
            var item = e.target.closest('.list-item[data-friend-id]');
            var friendId = Number(item.dataset.friendId);
            startChatWithFriend(friendId);
            return;
        }

        if (e.target.classList.contains('btn-accept')) {
            e.stopPropagation();
            var requestId = e.target.getAttribute('data-request-id');
            acceptFriendRequest(requestId);
            return;
        }

        if (e.target.classList.contains('btn-decline')) {
            e.stopPropagation();
            var requestId = e.target.getAttribute('data-request-id');
            declineFriendRequest(requestId);
            return;
        }

        if (e.target.classList.contains('btn-add-friend')) {
            var userId = Number(e.target.dataset.userId);
            sendFriendRequest(userId);
            e.target.textContent = 'Отправлено';
            e.target.disabled = true;
            return;
        }
    });
}

async function selectChat(chatId) {
    var chat = null;
    for (var i = 0; i < state.chats.length; i++) {
        if (state.chats[i].id === chatId) { chat = state.chats[i]; break; }
    }
    if (!chat) return;

    state.activeChat = chat;

    try {
        var msgs = await apiRequest('/api/chats/' + chatId + '/messages');
        state.messages = msgs || [];
        renderActiveChat();
        renderMessages();
        scrollToBottom();
        renderSidebarContent();
    } catch (err) {

    }
}

async function startChatWithFriend(friendId) {
    try {
        var chat = await apiRequest('/api/chats', 'POST', { friend_id: friendId });
        var exists = false;
        for (var i = 0; i < state.chats.length; i++) {
            if (state.chats[i].id === chat.id) { exists = true; break; }
        }
        if (!exists) state.chats.unshift(chat);
        await selectChat(chat.id);
        switchTab('chats');
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}

function renderActiveChat() {
    if (!state.activeChat || !state.currentUser) return;

    var other = state.activeChat.first_user.id === state.currentUser.id ?
        state.activeChat.second_user : state.activeChat.first_user;

    document.getElementById('chat-empty').classList.add('hidden');
    document.getElementById('chat-empty').style.display = 'none';
    document.getElementById('chat-active').classList.remove('hidden');
    document.getElementById('chat-header-name').textContent = other.username;
    document.getElementById('chat-header-avatar').textContent = getInitials(other.username);

    var statusEl = document.getElementById('chat-header-status');
    statusEl.textContent = other.is_online ? 'Онлайн' : 'Офлайн';
    statusEl.className = 'chat-header-status' + (other.is_online ? '' : ' offline');

    document.getElementById('delete-chat-btn').onclick = function () {
        if (confirm('Удалить чат?')) deleteChat(state.activeChat.id);
    };
}

function showEmptyChat() {
    document.getElementById('chat-empty').classList.remove('hidden');
    document.getElementById('chat-empty').style.display = 'flex';
    document.getElementById('chat-active').classList.add('hidden');
}

function renderMessages() {
    var container = document.getElementById('chat-messages');
    if (!container) return;

    if (state.messages.length === 0) {
        container.innerHTML = '<div style="text-align:center;color:#808080;padding:40px">Нет сообщений</div>';
        return;
    }
    var html = '';
    for (var i = 0; i < state.messages.length; i++) {
        var msg = state.messages[i];
        var isSent = msg.sender_id === state.currentUser.id;
        html += '<div class="message ' + (isSent ? 'sent' : 'received') + '"><div class="message-content">' +
            escapeHtml(msg.content) + '</div><div class="message-time">' + formatMessageTime(msg.created_at) +
            '</div><div style="clear:both"></div></div>';
    }
    container.innerHTML = html;
}

function scrollToBottom() {
    var container = document.getElementById('chat-messages');
    if (container) container.scrollTop = container.scrollHeight;
}

async function sendMessage() {
    var input = document.getElementById('message-input');
    var content = input.value.trim();
    if (!content || !state.activeChat) return;

    var receiverId = state.activeChat.first_user.id === state.currentUser.id ?
        state.activeChat.second_user.id : state.activeChat.first_user.id;

    try {
        await apiRequest('/api/chats/' + state.activeChat.id + '/messages', 'POST', {
            receiver_id: receiverId,
            content: content
        });
        input.value = '';
        input.style.height = 'auto';
    } catch (err) {
        alert('Ошибка отправки: ' + err.message);
    }
}

async function sendFriendRequest(userId) {
    try {
        await apiRequest('/api/friend-requests', 'POST', { to_user_id: userId });
        alert('Заявка отправлена');
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}

async function acceptFriendRequest(requestId) {
    var id = Number(requestId);
    if (isNaN(id)) {
        alert('Неверный ID');
        return;
    }

    try {
        await apiRequest('/api/friendships', 'POST', { friend_request_id: id });

        state.friendRequests = state.friendRequests.filter(function (r) { return r.id != requestId; });
        renderSidebarContent();
        updateRequestsBadge();
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}

async function declineFriendRequest(requestId) {
    var id = Number(requestId);
    if (isNaN(id)) {
        alert('Неверный ID');
        return;
    }

    try {
        await apiRequest('/api/friend-requests/' + id, 'DELETE');

        state.friendRequests = state.friendRequests.filter(function (r) { return r.id != requestId; });
        renderSidebarContent();
        updateRequestsBadge();
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}

async function deleteFriend(friendshipId) {
    try {
        await apiRequest('/api/friendships/' + friendshipId, 'DELETE');

        state.friends = state.friends.filter(function (f) { return f.id !== friendshipId; });
        renderSidebarContent();
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}

async function deleteChat(chatId) {
    try {
        await apiRequest('/api/chats/' + chatId, 'DELETE');

        state.chats = state.chats.filter(function (c) { return c.id !== chatId; });
        if (state.activeChat && state.activeChat.id === chatId) {
            state.activeChat = null;
            state.messages = [];
            showEmptyChat();
        }
        renderSidebarContent();
    } catch (err) {
        alert('Ошибка: ' + err.message);
    }
}

function updateChatLastMessage(chatId, content, timestamp) {
    var chat = null;
    for (var i = 0; i < state.chats.length; i++) {
        if (state.chats[i].id === chatId) { chat = state.chats[i]; break; }
    }
    if (chat) {
        chat.last_message_content = content;
        chat.last_message_at = timestamp;
        state.chats = state.chats.filter(function (c) { return c.id !== chatId; });
        state.chats.unshift(chat);
        if (state.activeTab === 'chats') renderSidebarContent();
    }
}

function openEditProfileModal() {
    document.getElementById('edit-username').value = state.currentUser.username;
    document.getElementById('edit-old-password').value = '';
    document.getElementById('edit-new-password').value = '';
    document.getElementById('edit-profile-modal').classList.remove('hidden');
}

function closeEditProfileModal() {
    document.getElementById('edit-profile-modal').classList.add('hidden');
}

async function loadData() {
    try {
        var results = await Promise.all([
            apiRequest('/api/chats'),
            apiRequest('/api/friendships'),
            apiRequest('/api/friend-requests?direction=incoming'),
            apiRequest('/api/friend-requests?direction=outgoing')
        ]);

        state.chats = results[0] || [];
        state.friends = results[1] || [];
        var incoming = results[2] || [];
        var outgoing = results[3] || [];
        state.friendRequests = incoming.concat(outgoing);

        updateRequestsBadge();
        renderSidebarContent();
    } catch (err) {

    }
}

function initEventListeners() {
    document.getElementById('show-register').addEventListener('click', function () {
        document.getElementById('login-container').classList.add('hidden');
        document.getElementById('register-container').classList.remove('hidden');
    });

    document.getElementById('show-login').addEventListener('click', function () {
        document.getElementById('register-container').classList.add('hidden');
        document.getElementById('login-container').classList.remove('hidden');
    });

    document.getElementById('register-form').addEventListener('submit', async function (e) {
        e.preventDefault();
        try {
            await apiRequest('/api/users', 'POST', {
                username: document.getElementById('reg-username').value,
                password: document.getElementById('reg-password').value
            });
            alert('Регистрация успешна! Теперь войдите.');
            document.getElementById('show-login').click();
            e.target.reset();
        } catch (err) {
            alert('Ошибка: ' + err.message);
        }
    });

    document.getElementById('login-form').addEventListener('submit', async function (e) {
        e.preventDefault();
        try {
            await apiRequest('/api/sessions', 'POST', {
                username: document.getElementById('login-username').value,
                password: document.getElementById('login-password').value
            });

            var user = await apiRequest('/api/users/me');
            state.currentUser = user;
            state.shouldConnect = true;

            showAppView();
            await loadData();
            e.target.reset();
        } catch (err) {
            alert('Ошибка входа: ' + err.message);
        }
    });

    document.getElementById('logout-btn').addEventListener('click', async function () {
        try {
            await apiRequest('/api/sessions', 'DELETE');
        } catch (err) {
            // ignore
        } finally {
            disconnectWebSocket();
            state.currentUser = null;
            state.chats = [];
            state.friends = [];
            state.friendRequests = [];
            state.users = [];
            state.activeChat = null;
            state.messages = [];
            showAuthView();
        }
    });

    var tabs = document.querySelectorAll('.tab-btn');
    for (var t = 0; t < tabs.length; t++) {
        (function (tab) {
            tabs[t].addEventListener('click', function () { switchTab(tab); });
        })(tabs[t].dataset.tab);
    }

    document.getElementById('send-btn').addEventListener('click', sendMessage);

    document.getElementById('message-input').addEventListener('keydown', function (e) {
        if (e.key === 'Enter' && !e.shiftKey) {
            e.preventDefault();
            sendMessage();
        }
    });

    document.getElementById('message-input').addEventListener('input', function () {
        this.style.height = 'auto';
        this.style.height = Math.min(this.scrollHeight, 120) + 'px';
    });

    document.getElementById('user-info-clickable').addEventListener('click', openEditProfileModal);
    document.getElementById('cancel-edit-profile').addEventListener('click', closeEditProfileModal);

    document.getElementById('edit-profile-form').addEventListener('submit', async function (e) {
        e.preventDefault();
        var body = {};
        var username = document.getElementById('edit-username').value.trim();
        var oldPass = document.getElementById('edit-old-password').value;
        var newPass = document.getElementById('edit-new-password').value;

        if (username) body.username = username;
        if (oldPass) body.old_password = oldPass;
        if (newPass) body.new_password = newPass;

        if (Object.keys(body).length === 0) {
            alert('Нет изменений');
            return;
        }

        try {
            var user = await apiRequest('/api/users', 'PATCH', body);
            state.currentUser = Object.assign({}, state.currentUser, user);
            updateUserDisplay();
            closeEditProfileModal();
            alert('Профиль обновлён');
        } catch (err) {
            alert('Ошибка: ' + err.message);
        }
    });

    window.addEventListener('beforeunload', function () {
        disconnectWebSocket();
    });
}

async function init() {
    initEventListeners();
    initEventDelegation();

    try {
        var user = await apiRequest('/api/users/me');
        state.currentUser = user;
        state.shouldConnect = true;

        showAppView();
        await loadData();
    } catch (err) {
        showAuthView();
    }
}

document.addEventListener('DOMContentLoaded', init);