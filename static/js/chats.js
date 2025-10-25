import { userToken, handleAuthResponse, ws, onlineUserIds } from './utils.js';
import { navigate, showAlert } from './script.js';

const mainContent = document.getElementById('main-content');
let selectedRecipientId = null;
let currentMessageOffset = 0;
const messageLimit = 10;

export function renderChatsPage() {
    mainContent.innerHTML = `
        <div class="chat-container">
            <div class="user-list">
                <h3>Users</h3>
                <div id="users-online-list"></div>
            </div>
            <div class="chat-area">
                <div id="messages-container">
                    <p class="chat-placeholder">Select a user to start chatting.</p>
                </div>
                <form id="chat-form" class="chat-input-area" style="display: none;">
                    <input type="text" id="chat-input" placeholder="Type a message..." required>
                    <button type="submit" id="send-button">></button>
                </form>
            </div>
        </div>
    `;
    fetchUsers();
    const chatForm = document.getElementById('chat-form');
    if (chatForm) {
        chatForm.addEventListener('submit', handleSendMessage);
    }
}

// Export function to update users list when online status changes
window.updateUsersList = function() {
    // Re-render users list with updated online status
    fetchUsers();
};

function throttle(func, limit) {
    let inThrottle;
    return function (...args) {
        if (!inThrottle) {
            func.apply(this, args);
            inThrottle = true;
            setTimeout(() => (inThrottle = false), limit);
        }
    };
}

function debounce(func, delay) {
    let timer;
    return function (...args) {
        clearTimeout(timer);
        timer = setTimeout(() => func.apply(this, args), delay);
    };
}

function createMessageElement(msg, currentUserId) {
    const messageElement = document.createElement('div');
    const isSent = msg.SenderID === currentUserId;
    messageElement.className = `chat-message ${isSent ? 'sent' : 'received'}`;
    const dateField = msg.CreatedAt || msg.created_at;
    const messageDate = new Date(dateField).toLocaleString();
    const senderName = isSent ? 'You' : msg.SenderNickname || msg.sender_name;
    messageElement.innerHTML = `
        <div class="message-header">
            <span class="message-author">${senderName}</span> ‎ | ‎  
            <span class="message-date">${messageDate}</span>
        </div>
        <p>${msg.Content || msg.content}</p>
    `;
    return messageElement;
}

function renderMessages(messages, currentUserId, prepend = false) {
    const messagesContainer = document.getElementById('messages-container');
    if (!messagesContainer) return;
    if (messages.length === 0) {
        if (currentMessageOffset === 0) {
            if (messagesContainer.children.length === 1 && messagesContainer.children[0].tagName === 'H3') {
                messagesContainer.innerHTML += '<p id="no-messages-placeholder">There are no messages yet. Start the chat!</p>';
            } else if (currentMessageOffset > 0) {
                const allLoadedIndicator = document.createElement('p');
                allLoadedIndicator.className = 'all-loaded-indicator';
                allLoadedIndicator.textContent = '— End of messages —';
                messagesContainer.insertBefore(allLoadedIndicator, messagesContainer.children[1]);
            }
        }
        return;
    }
    const fragment = document.createDocumentFragment();
    messages.forEach(msg => {
        fragment.appendChild(createMessageElement(msg, currentUserId));
    });
    if (prepend) {
        const firstMessageElement = messagesContainer.children[1];
        messagesContainer.insertBefore(fragment, firstMessageElement);
    } else {
        messagesContainer.appendChild(fragment);
    }
}

export function handleIncomingPrivateMessage(message) {
    const messagesContainer = document.getElementById('messages-container');
    if (!messagesContainer || selectedRecipientId !== message.payload.sender_id) return;
    const placeholder = document.getElementById('no-messages-placeholder');
    if (placeholder) {
        placeholder.remove();
    }
    const messageElement = createMessageElement(message.payload, null);
    messageElement.className = 'chat-message received';
    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

export async function fetchUsers() {
    const usersList = document.getElementById('users-online-list');
    usersList.innerHTML = '';
    try {
        const response = await fetch('/api/users', {
            method: 'GET',
            headers: { 'Authorization': `${userToken}` }
        });
        if (!response.ok) {
            throw new Error('Failed to fetch users');
        }
        const users = await response.json();
        const validUsers = users.filter(us => us && us.user.Nickname && typeof us.user.Nickname === 'string');
        validUsers.sort((a, b) => a.user.Nickname.localeCompare(b.user.Nickname));
        validUsers.sort((a, b) => {
            const msgA = a.last_message;
            const msgB = b.last_message;
            if (!msgA && !msgB) return 0;
            if (msgA && !msgB) return -1;
            if (!msgA && msgB) return 1;
            const dateA = new Date(msgA.CreatedAt);
            const dateB = new Date(msgB.CreatedAt);
            return dateB - dateA;
        });
        validUsers.forEach(us => {
            const user = us.user;
            const userElement = document.createElement('div');
            userElement.className = 'user-item';
            userElement.dataset.userId = user.ID;
            const isOnline = checkUserOnlineStatus(user.ID);
            const statusDot = isOnline ? '<span class="status-dot online"></span>' : '<span class="status-dot offline"></span>';
            const lastMsg = us.last_message;
            let lastMsgPreview = 'No messages yet';
            if (lastMsg) {
                const content = lastMsg.Content;
                lastMsgPreview = content.length > 30 ? content.substring(0, 30) + '...' : content;
            }
            userElement.innerHTML = `
                <span class="user-nickname">${statusDot} ${user.Nickname}</span>
                <span class="mes">${lastMsgPreview}</span>
            `;
            userElement.addEventListener('click', () => {
                if (window.location.hash !== '#/chats') {
                    navigate('#/chats');
                }
                selectUserForChat(user.ID, user.Nickname);
            });
            usersList.appendChild(userElement);
        });
    } catch (error) {
        console.error('Failed to fetch users:', error);
        usersList.innerHTML = '<p>Failed to load users</p>';
    }
}

function checkUserOnlineStatus(userId) {
    return onlineUserIds.includes(userId);
}

async function fetchAndDisplayMessages(userId, offset, limit, prepend = false) {
    const messagesContainer = document.getElementById('messages-container');
    if (offset > 0 && messagesContainer.querySelector('.all-loaded-indicator')) {
        return;
    }
    const loadingIndicatorId = 'loading-older-messages';
    if (prepend) {
        const loadingIndicator = document.createElement('p');
        loadingIndicator.id = loadingIndicatorId;
        loadingIndicator.textContent = 'Loading...';
        messagesContainer.insertBefore(loadingIndicator, messagesContainer.children[1]);
    }
    let scrollHeightBeforeLoad = messagesContainer.scrollHeight;
    try {
        const response = await fetch('/api/messages', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `${userToken}`
            },
            body: JSON.stringify({
                recipient_id: userId,
                offset: offset,
                limit: limit
            })
        });
        if (response.ok) {
            const messages = await response.json();
            if (!Array.isArray(messages)) {
                console.error('API response is not an array:', messages);
                if (offset === 0) {
                    messagesContainer.innerHTML += '<p>Failed to load messages due to unexpected data format.</p>';
                }
                return;
            }
            const currentUserId = await getCurrentUserId();
            messages.reverse();
            renderMessages(messages, currentUserId, prepend);
            currentMessageOffset += messages.length;
            if (!prepend) {
                messagesContainer.scrollTop = messagesContainer.scrollHeight;
            } else if (messages.length > 0) {
                const newScrollHeight = messagesContainer.scrollHeight;
                messagesContainer.scrollTop = newScrollHeight - scrollHeightBeforeLoad;
            }
        } else {
            if (offset === 0) {
                messagesContainer.innerHTML += '<p>Failed to fetch messages</p>';
            } else {
                console.error('Failed to fetch older messages');
            }
        }
    } catch (error) {
        console.error('Failed to fetch messages:', error);
        if (offset === 0) {
            messagesContainer.innerHTML += '<p>Failed to load messages</p>';
        }
    } finally {
        const loadingIndicator = document.getElementById(loadingIndicatorId);
        if (loadingIndicator) {
            loadingIndicator.remove();
        }
    }
}

function handleScrollForOlderMessages(event) {
    const container = event.currentTarget;
    if (container.scrollTop < 20 && selectedRecipientId !== null) {
        container.removeEventListener('scroll', container._scrollHandler);
        fetchAndDisplayMessages(selectedRecipientId, currentMessageOffset, messageLimit, true).then(() => {
            container.addEventListener('scroll', container._scrollHandler);
        });
    }
}

async function selectUserForChat(userId, nickname) {
    selectedRecipientId = userId;
    currentMessageOffset = 0;
    const messagesContainer = document.getElementById('messages-container');
    const chatForm = document.getElementById('chat-form');
    messagesContainer.innerHTML = `<h3>Chat with ${nickname}</h3>`;
    if (messagesContainer._scrollHandler) {
        messagesContainer.removeEventListener('scroll', messagesContainer._scrollHandler);
    }
    const throttledScrollHandler = throttle(handleScrollForOlderMessages, 300);
    messagesContainer._scrollHandler = throttledScrollHandler;
    messagesContainer.addEventListener('scroll', throttledScrollHandler);
    chatForm.style.display = 'flex';
    chatForm.dataset.recipientId = userId;
    const inputContainer = document.getElementById('chat-input');
    const submitCont = document.getElementById('send-button')
    if(checkUserOnlineStatus(userId) === true){
        inputContainer.disabled = false;
        inputContainer.placeholder = 'Type a message';
        submitCont.disabled = false;
        submitCont.textContent = '>';
    } else{
        inputContainer.disabled = true;
        inputContainer.placeholder = 'The user is offline';
        submitCont.disabled = true;
        submitCont.textContent = 'X';
    }
    await fetchAndDisplayMessages(userId, currentMessageOffset, messageLimit, false);
}

export async function getCurrentUserId() {
    try {
        const response = await fetch('/api/users/profile', {
            method: 'GET',
            headers: { 'Authorization': `${userToken}` }
        });
        if (response.ok) {
            const data = await response.json();
            return data.user.ID;
        }
    } catch (error) {
        console.error('Failed to get current user ID:', error);
    }
    return null;
}

function handleSendMessage(event) {
    event.preventDefault();
    const recipientId = event.target.dataset.recipientId;
    const chatInput = document.getElementById('chat-input');
    const content = chatInput.value;
    if (!content.trim()) return;
    
    // Check if recipient is online
    if (!onlineUserIds.includes(parseInt(recipientId))) {
        showAlert('warningAlert', "User if offline", 'You cannot send messages to offline users.')
        return;
    }
    
    if (ws && ws.readyState === WebSocket.OPEN) {
        const message = {
            type: 'private_message',
            payload: {
                recipient_id: parseInt(recipientId),
                content: content,
            }
        };
        ws.send(JSON.stringify(message));
        chatInput.value = '';
        const messagesContainer = document.getElementById('messages-container');
        const placeholder = document.getElementById('no-messages-placeholder');
        if (placeholder) {
            placeholder.remove();
        }
        const messageElement = document.createElement('div');
        messageElement.className = 'chat-message sent';
        const messageDate = new Date().toLocaleString();
        messageElement.innerHTML = `
            <div class="message-header">
                <span class="message-author">You</span> ‎ | ‎
                <span class="message-date">${messageDate}</span>
            </div>
            <p>${content}</p>
        `;
        messagesContainer.appendChild(messageElement);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
}
