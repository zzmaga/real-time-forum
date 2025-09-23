import { userToken, handleAuthResponse, ws } from './utils.js';

const mainContent = document.getElementById('main-content');
let selectedRecipientId = null;

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
                    <button type="submit" id="send-button">筐､</button>
                </form>
            </div>
        </div>
    `;
    fetchUsers();
}

export function handleIncomingPrivateMessage(message) {
    const messagesContainer = document.getElementById('messages-container');
    if (!messagesContainer || selectedRecipientId !== message.payload.sender_id) return;
    
    const payload = message.payload;
    const messageElement = document.createElement('div');
    messageElement.className = 'chat-message received';
    
    const messageDate = new Date(payload.created_at).toLocaleString();
    
    messageElement.innerHTML = `
        <div class="message-header">
            <span class="message-author">${payload.sender_name}</span>
            <span class="message-date">${messageDate}</span>
        </div>
        <p>${payload.content}</p>
    `;
    
    messagesContainer.appendChild(messageElement);
    messagesContainer.scrollTop = messagesContainer.scrollHeight;
}

async function fetchUsers() {
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
        console.log(users);
        users.sort((a, b) => a.nickname.localeCompare(b.nickname));
        console.log(users);
        users.forEach(user => {
            const userElement = document.createElement('div');
            userElement.className = 'user-item';
            userElement.dataset.userId = user.id;

            const isOnline = checkUserOnlineStatus(user.id);
            const statusDot = isOnline ? '<span class="status-dot online"></span>' : '<span class="status-dot offline"></span>';

            userElement.innerHTML = `
                ${statusDot}
                <span class="user-nickname">${user.nickname}</span>
            `;
            
            userElement.addEventListener('click', () => {
                selectUserForChat(user.id, user.nickname);
            });

            usersList.appendChild(userElement);
        });
    } catch (error) {
        console.error('Failed to fetch users:', error);
        usersList.innerHTML = '<p>Failed to load users</p>';
    }
}

function checkUserOnlineStatus(userId) {
    // Placeholder logic for online status
    return false;
}

async function selectUserForChat(userId, nickname) {
    selectedRecipientId = userId;
    const messagesContainer = document.getElementById('messages-container');
    const chatForm = document.getElementById('chat-form');
    messagesContainer.innerHTML = `<h3>Chat with ${nickname}</h3>`;
    chatForm.style.display = 'flex';
    chatForm.dataset.recipientId = userId;
    
    try {
        const response = await fetch('/api/messages', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `${userToken}`
            },
            body: JSON.stringify({
                recipient_id: userId,
                offset: 0,
                limit: 10
            })
        });
        
        if (response.ok) {
            const messages = await response.json();
            
            messages.reverse().forEach(msg => {
                const messageElement = document.createElement('div');
                const isSent = msg.sender_id === getCurrentUserId(); 
                messageElement.className = `chat-message ${isSent ? 'sent' : 'received'}`;
                
                const messageDate = new Date(msg.created_at).toLocaleString();

                messageElement.innerHTML = `
                    <div class="message-header">
                        <span class="message-author">${isSent ? 'You' : msg.sender_nickname}</span>
                        <span class="message-date">${messageDate}</span>
                    </div>
                    <p>${msg.content}</p>
                `;
                messagesContainer.appendChild(messageElement);
            });
            
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        } else {
            messagesContainer.innerHTML += '<p>Failed to load messages</p>';
        }
    } catch (error) {
        console.error('Failed to fetch messages:', error);
        messagesContainer.innerHTML += '<p>Failed to load messages</p>';
    }
    
    chatForm.addEventListener('submit', handleSendMessage);
}

function getCurrentUserId() {
    // Placeholder function to get the current user's ID
    return 1;
}

function handleSendMessage(event) {
    event.preventDefault();
    const recipientId = event.target.dataset.recipientId;
    const chatInput = document.getElementById('chat-input');
    const content = chatInput.value;

    if (!content.trim()) return;

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
        const messageElement = document.createElement('div');
        messageElement.className = 'chat-message sent';
        
        const messageDate = new Date().toLocaleString();
        
        messageElement.innerHTML = `
            <div class="message-header">
                <span class="message-author">You</span>
                <span class="message-date">${messageDate}</span>
            </div>
            <p>${content}</p>
        `;
        messagesContainer.appendChild(messageElement);
        messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
}