const mainContent = document.getElementById('main-content');
const loginLink = document.getElementById('login-link');
const registerLink = document.getElementById('register-link');
const postsLink = document.getElementById('posts-link');
const chatsLink = document.getElementById('chats-link');
const logoutLink = document.getElementById('logout-link');

let userToken = localStorage.getItem('userToken');
let ws;

function isLoggedIn() {
    return !!userToken;
}

function updateNav() {
    if (isLoggedIn()) {
        loginLink.style.display = 'none';
        registerLink.style.display = 'none';
        postsLink.style.display = 'inline';
        chatsLink.style.display = 'inline';
        logoutLink.style.display = 'inline';
    } else {
        loginLink.style.display = 'inline';
        registerLink.style.display = 'inline';
        postsLink.style.display = 'none';
        chatsLink.style.display = 'none';
        logoutLink.style.display = 'none';
    }
}

function navigate(path) {
    history.pushState(null, '', path);
    renderPage(path);
}

function handleAuthResponse(response) {
    if (response.status === 401) {
        localStorage.removeItem('userToken');
        userToken = null;
        updateNav();
        navigate('#/login');
        return;
    }
    return response.json();
}

// --- WEBSOCKETS ---
function connectWebSocket() {
    if (ws && ws.readyState === WebSocket.OPEN) return;

    ws = new WebSocket(`ws://localhost:8080/ws?token=${userToken}`);

    ws.onopen = () => {
        console.log('WebSocket connection established.');
    };

    ws.onmessage = (event) => {
    const message = JSON.parse(event.data);
    console.log('Message from server:', message);
    if (message.type === 'new_post') {
        if (window.location.hash === '#/posts') {
            displayPosts();
        }
    } else if (message.type === 'private_message') {
        // ...
    }
};

    ws.onclose = () => {
        console.log('WebSocket connection closed. Attempting to reconnect...');
        setTimeout(connectWebSocket, 1000);
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

function renderLoginPage() {
    mainContent.innerHTML = `
        <h2>Login</h2>
        <form id="login-form">
            <input type="text" id="login" placeholder="Nickname or Email" required>
            <input type="password" id="password" placeholder="Password" required>
            <button type="submit">Login</button>
        </form>
    `;
    document.getElementById('login-form').addEventListener('submit', handleLogin);
}

function renderRegisterPage() {
    const today = new Date();
    const year = today.getFullYear();
    const month = String(today.getMonth() + 1).padStart(2, '0');
    const day = String(today.getDate()).padStart(2, '0');
    const maxDate = `${year}-${month}-${day}`;
    mainContent.innerHTML = `
        <h2>Register</h2>
        <form id="register-form">
            <input type="text" id="nickname" placeholder="Nickname" required>
            <input type="text" id="firstName" placeholder="First Name" required>
            <input type="text" id="lastName" placeholder="Last Name" required>
            <input type="email" id="email" placeholder="E-mail" required>
            <input type="password" id="password" placeholder="Password" required>
            <label for="age">Date of Birth:</label>
            <input type="date" id="age" min="1900-01-01" max="${maxDate}" required>
            <select id="gender" required>
                <option value="">Select Gender</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
                <option value="other">Other</option>
            </select>
            <button type="submit">Register</button>
        </form>
    `;
    document.getElementById('register-form').addEventListener('submit', handleRegister);
}

function renderPostsPage() {
    mainContent.innerHTML = `
        <div class="posts-wrapper">
            <h2>Forum Posts</h2>
            <button id="create-post-btn">Create a post</button>
            <div id="posts-container"></div>
        </div>
    `;
    document.getElementById('create-post-btn').addEventListener('click', showCreatePostModal);
    displayPosts();
}

function renderCreatePostModal() {
    const modalHTML = `
        <div id="create-post-modal" class="modal">
            <div class="modal-content">
                <span class="close-btn">&times;</span>
                <h3>Create a New Post</h3>
                <form id="create-post-form">
                    <input type="text" id="title" placeholder="Title of post" required>
                    <textarea id="post-content" placeholder="What's on your mind?" required></textarea>
                    <h4>Categories:</h4>
                    <div id="category-checkboxes">
                        <label><input type="checkbox" name="category" value="Tech"> Tech</label>
                        <label><input type="checkbox" name="category" value="News"> News</label>
                        <label><input type="checkbox" name="category" value="Sports"> Sports</label>
                    </div>
                    <button type="submit">Create Post</button>
                </form>
            </div>
        </div>
    `;
    document.body.insertAdjacentHTML('beforeend', modalHTML);
    document.getElementById('create-post-form').addEventListener('submit', handleCreatePost);
    document.querySelector('.close-btn').addEventListener('click', hideCreatePostModal);
    window.addEventListener('click', (event) => {
        if (event.target === document.getElementById('create-post-modal')) {
            hideCreatePostModal();
        }
    });
}

function showCreatePostModal() {
    document.getElementById('create-post-modal').style.display = 'flex';
}

function hideCreatePostModal() {
    document.getElementById('create-post-modal').style.display = 'none';
    document.getElementById('create-post-form').reset();
}

async function displayPosts() {
    const postsContainer = document.getElementById('posts-container');
    postsContainer.innerHTML = '';
    let posts = null;
    console.log(userToken);
    try {
        const response = await fetch('/api/posts', {
            method: 'GET',
            headers: { 'Authorization': `${userToken}` }
        });
        posts = await handleAuthResponse(response);
    } catch (error) {
        console.error('Failed to fetch posts:', error);
        posts = []; // Treat the error as an empty list of posts
    }

    if (!posts || posts.length === 0) {
        postsContainer.innerHTML = '<p>No posts yet. Be the first to post!</p>';
        return;
    }
    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.className = 'post-card';
        const categories = Array.isArray(post.category) ? post.category.join(', ') : post.category;
        
        postElement.innerHTML = `
            <h3 class="post-title" style="cursor: pointer;">${post.Title}</h3>
            <p>${post.Content}</p>
            <small>By: ${post.author} | Categories: ${categories}</small>
        `;

        // Add the click event listener to the title
        postElement.querySelector('.post-title').addEventListener('click', () => {
            showPostAndComments(post.id);
        });

        postsContainer.appendChild(postElement);
    });
}

async function showPostAndComments(postId) {
    let postData = null;
    try {
        const response = await fetch(`/api/posts/${postId}`, {
            headers: { 'Authorization': `Bearer ${userToken}` }
        });
        postData = await handleAuthResponse(response);
    } catch (error) {
        console.error('Failed to fetch post:', error);
        mainContent.innerHTML = '<p>An error occurred while fetching the post.</p>';
        return;
    }

    if (!postData || !postData.success) {
        mainContent.innerHTML = `<p>${postData.error || 'Failed to load post.'}</p>`;
        return;
    }
    
    const post = postData.post;
    const comments = postData.comments;

    let commentsHtml = '';
    if (comments.length === 0) {
        commentsHtml = '<p>No comments yet.</p>';
    } else {
        comments.forEach(comment => {
            commentsHtml += `
                <div class="comment">
                    <p><strong>${comment.author}</strong> on ${new Date(comment.created_at).toLocaleString()}</p>
                    <p>${comment.content}</p>
                </div>
            `;
        });
    }

    mainContent.innerHTML = `
        <div class="single-post-container">
            <h2>${post.title}</h2>
            <p>${post.content}</p>
            <small>By: ${post.author} | Categories: ${Array.isArray(post.category) ? post.category.join(', ') : post.category}</small>
            <hr>
            <div id="comments-section">
                <h3>Comments</h3>
                ${commentsHtml}
            </div>
            <form id="add-comment-form">
                <textarea id="comment-content" placeholder="Write a comment..." required></textarea>
                <button type="submit">Send Comment</button>
            </form>
        </div>
    `;
    
    document.getElementById('add-comment-form').addEventListener('submit', (e) => handleAddComment(e, postId));
}

function renderChatsPage() {
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
                    <button type="submit" id="send-button">➤</button>
                </form>
            </div>
        </div>
    `;

    fetchUsers();
}

async function fetchUsers() {
    const usersList = document.getElementById('users-online-list');
    usersList.innerHTML = '';
    
    // для примера. А так будем делать запрос 
    const dummyUsers = [
        { id: 2, nickname: 'Alice', is_online: true, last_message: { content: 'Hey, are you there?' } },
        { id: 3, nickname: 'Bob', is_online: false, last_message: { content: 'See you later!' } },
        { id: 4, nickname: 'Charlie', is_online: true },
        { id: 5, nickname: 'David', is_online: true },
        { id: 6, nickname: 'Eve', is_online: false },
    ];
    dummyUsers.sort((a, b) => {
        if (a.last_message && b.last_message) {
            // нужна сортировка по времени
            return 0; 
        }
        if (a.last_message) return -1;
        if (b.last_message) return 1;
        return a.nickname.localeCompare(b.nickname);
    });

    dummyUsers.forEach(user => {
        const userElement = document.createElement('div');
        userElement.className = 'user-item';
        userElement.dataset.userId = user.id;

        const statusDot = user.is_online ? '<span class="status-dot online"></span>' : '<span class="status-dot offline"></span>';
        const lastMessageHtml = user.last_message ? `<div class="last-message-preview">${user.last_message.content}</div>` : '';

        userElement.innerHTML = `
            ${statusDot}
            <span class="user-nickname">${user.nickname}</span>
            ${lastMessageHtml}
        `;
        
        userElement.addEventListener('click', () => {
            selectUserForChat(user.id, user.nickname);
        });

        usersList.appendChild(userElement);
    });
}

async function selectUserForChat(userId, nickname) {
    const messagesContainer = document.getElementById('messages-container');
    const chatForm = document.getElementById('chat-form');
    messagesContainer.innerHTML = `<h3>Chat with ${nickname}</h3>`;
    chatForm.style.display = 'flex';
    chatForm.dataset.recipientId = userId;
    
    // Placeholder to fetch messages
    // You will need a new API endpoint in your Go backend to fetch chat history.
    const dummyMessages = [
        { senderId: 1, content: 'Hi, how are you?', created_at: '2025-09-05T10:00:00Z' },
        { senderId: userId, content: 'I am doing great, thanks!', created_at: '2025-09-05T10:01:00Z' },
    ];
    
    dummyMessages.forEach(msg => {
        const messageElement = document.createElement('div');
        messageElement.className = `chat-message ${msg.senderId === 1 ? 'sent' : 'received'}`;
        
        const messageDate = new Date(msg.created_at).toLocaleString();

        messageElement.innerHTML = `
            <div class="message-header">
                <span class="message-author">${msg.senderId === 1 ? 'You' : nickname}</span>
                <span class="message-date">${messageDate}</span>
            </div>
            <p>${msg.content}</p>
        `;
        messagesContainer.appendChild(messageElement);
    });
    
    // Add event listener for the chat form
    chatForm.addEventListener('submit', handleSendMessage);
}

function renderErrorPage() {
    mainContent.innerHTML = `
        <h2>404 - Page Not Found</h2>
        <p>The page you're looking for doesn't exist.</p>
        <p><a href="#/posts">Go to Posts</a></p>
    `;
}

function renderPage(path) {
    if (isLoggedIn()) {
        switch (path) {
            case '#/login':
            case '#/register':
                navigate('#/posts');
                break;
            case '#/chats':
                renderChatsPage();
                break;
            case '#/posts':
            case '#/':
            default:
                renderPostsPage();
        }
    } else {
        switch (path) {
            case '#/register':
                renderRegisterPage();
                break;
            case '#/error':
                renderErrorPage();
                break;
            case '#/login':
            case '#/':
                renderLoginPage();
                break;
            default:
                navigate('#/login');
        }
    }
}

async function handleLogin(event) {
    event.preventDefault();
    const loginId = document.getElementById('login').value;
    const password = document.getElementById('password').value;
    const response = await fetch('/api/signin', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ loginId, password })
    });
    const data = await response.json();
    if (data.success) {
        localStorage.setItem('userToken', data.token);
        userToken = data.token;
        updateNav();
        navigate('#/posts');
        connectWebSocket();
    } else {
        alert(data.error);
    }
}

async function handleRegister(event) {
    event.preventDefault();
    const newUser = {
        nickname: document.getElementById('nickname').value,
        age: document.getElementById('age').value,
        gender: document.getElementById('gender').value,
        firstName: document.getElementById('firstName').value,
        lastName: document.getElementById('lastName').value,
        email: document.getElementById('email').value,
        password: document.getElementById('password').value,
    };
    const response = await fetch('/api/signup', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(newUser)
    });
    const data = await response.json();
    if (data.success) {
        alert('Registration successful! Please log in.');
        navigate('#/login');
    } else {
        alert(data.error);
    }
}

async function handleCreatePost(event) {
    event.preventDefault();

    const selectedCategories = Array.from(document.querySelectorAll('#category-checkboxes input[name="category"]:checked')).map(checkbox => checkbox.value);

    const postData = {
        title: document.getElementById('title').value,
        content: document.getElementById('post-content').value,
        category: selectedCategories,
    };

    const response = await fetch('/api/posts', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${userToken}`
        },
        body: JSON.stringify(postData)
    });
    const data = await handleAuthResponse(response);
    if (data.success) {
        document.getElementById('title').value = '';
        document.getElementById('post-content').value = '';
        document.querySelectorAll('#category-checkboxes input[type="checkbox"]').forEach(checkbox => checkbox.checked = false);
        ws.send(JSON.stringify({ type: 'new_post', payload: postData }));
    } else {
        alert(data.error);
    }
}

async function handleAddComment(e, postId) {
    e.preventDefault();
    const commentContent = document.getElementById('comment-content').value;

    const commentData = {
        post_id: postId,
        content: commentContent
    };

    const response = await fetch('/api/comments', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${userToken}`
        },
        body: JSON.stringify(commentData)
    });

    const data = await handleAuthResponse(response);
    if (data && data.success) {
        console.log('Comment added successfully.');
        document.getElementById('comment-content').value = '';
        // Reload the post to show the new comment
        showPostAndComments(postId);
    } else {
        alert(data.error || 'Failed to add comment.');
    }
}
    
    function handleSendMessage(event) {
    event.preventDefault();
    const recipientId = event.target.dataset.recipientId;
    const chatInput = document.getElementById('chat-input');
    const content = chatInput.value;

    if (!content.trim()) return;

    // Send the message via WebSocket
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
        
        // Temporarily add the message to the chat view
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

function handleLogout() {
    localStorage.removeItem('userToken');
    userToken = null;
    if (ws) {
        ws.close();
    }
    updateNav();
    navigate('#/login');
}

window.addEventListener('popstate', () => renderPage(location.hash));
document.addEventListener('DOMContentLoaded', () => {
    updateNav();
    if (isLoggedIn()) {
        connectWebSocket();
    }
    renderCreatePostModal();
    if (location.pathname !== '/') {
        history.replaceState(null, '', '/#/error');
        renderPage('#/error');
    } else {
        renderPage(location.hash || '#/posts');
    }
});

loginLink.addEventListener('click', () => navigate('#/login'));
registerLink.addEventListener('click', () => navigate('#/register'));
postsLink.addEventListener('click', () => navigate('#/posts'));
chatsLink.addEventListener('click', () => navigate('#/chats'));
logoutLink.addEventListener('click', handleLogout);