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
        handleIncomingPrivateMessage(message);
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
            <input type="text" id="nickname" placeholder="Nickname" maxlength="32" required>
            <input type="text" id="firstName" placeholder="First Name" required>
            <input type="text" id="lastName" placeholder="Last Name" required>
            <input type="email" id="email" placeholder="E-mail" maxlength="320" required>
            <input type="password" id="password" placeholder="Password" minlength="8" required>
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
    try {
        const response = await fetch('/api/posts', {
            method: 'GET',
            headers: { 'Authorization': `${userToken}` }
        });
        posts = await handleAuthResponse(response);
    } catch (error) {
        console.error('Failed to fetch posts:', error);
        posts = [];
    }
    if (!posts || posts.length === 0) {
        postsContainer.innerHTML = '<p>No posts yet. Be the first to post!</p>';
        return;
    }
    posts.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt));
    posts.forEach(post => {
        const postElement = document.createElement('div');
        postElement.className = 'post-card';
        const categories = Array.isArray(post.WCategories) ? post.WCategories.join(', ') : post.WCategories;
        const titleElement = document.createElement('h3');
        titleElement.className = 'post-title';
        titleElement.style.cursor = 'pointer';
        titleElement.textContent = post.Title;
        const contentElement = document.createElement('p');
        contentElement.textContent = post.Content;
        const smallElement = document.createElement('small');
        smallElement.textContent = `By: ${post.WUser} | Categories: ${categories}`;
        const actionsElement = document.createElement('div');
        actionsElement.className = 'post-actions';
        actionsElement.innerHTML = `
            <button class="vote-btn" data-post-id="${post.Id}" data-vote-type="1">👍 <span class="like-count">${post.WVoteUp}</span></button>
            <button class="vote-btn" data-post-id="${post.Id}" data-vote-type="-1">👎 <span class="dislike-count">${post.WVoteDown}</span></button>
        `;
        postElement.appendChild(titleElement);
        postElement.appendChild(contentElement);
        postElement.appendChild(smallElement);
        postElement.appendChild(actionsElement);
        titleElement.addEventListener('click', () => {
            navigate(`#/posts/${post.Id}`);
        });
        actionsElement.querySelectorAll('.vote-btn').forEach(button => {
            button.addEventListener('click', handleVote);
        });
        postsContainer.appendChild(postElement);
    });
}


async function showPostAndComments(postId) {
    let postData = null;
    try {
        const response = await fetch(`/api/posts/${postId}`, {
            headers: { 'Authorization': `${userToken}` }
        });
        postData = await handleAuthResponse(response);
    } catch (error) {
        console.error('Failed to fetch post:', error);
        mainContent.innerHTML = '<p>An error occurred while fetching the post.</p>';
        return;
    }

    if (!postData || !postData.success) {
        mainContent.innerHTML = `<p>${escapeHTML(postData?.error) || 'Failed to load post.'}</p>`;
        return;
    }
    
    const post = postData.post;
    const comments = postData.comments;
    comments.sort((a, b) => new Date(b.CreatedAt) - new Date(a.CreatedAt));

    let commentsHtml = '';
    if (comments.length === 0) {
        commentsHtml = '<p>No comments yet.</p>';
    } else {
        comments.forEach(comment => {
            commentsHtml += `
                <div class="comment-card">
                    <div class="comment-header">
                        <p class="comment-author"><strong>${escapeHTML(comment.Author)}</strong></p>
                        <p class="comment-date">${new Date(comment.CreatedAt).toLocaleString()}</p>
                    </div>
                    <div class="comment-body">
                        <p>${escapeHTML(comment.Content)}</p>
                    </div>
                </div>
            `;
        });
    }

    mainContent.innerHTML = `
        <div class="centered-container">
            <div class="single-post-container">
                <div class="post-content-section">
                    <h2>${escapeHTML(post.Title)}</h2>
                    <p>${escapeHTML(post.Content)}</p>
                    <small>
                        By: ${escapeHTML(post.Author)} | 
                        Categories: ${
                            Array.isArray(post.Category) 
                                ? post.Category.map(c => escapeHTML(c)).join(', ') 
                                : escapeHTML(post.Category)
                        }
                    </small>
                </div>
                <hr>
                <form id="add-comment-form" class="comment-form-container">
                    <textarea id="comment-content" placeholder="Write a comment..." required></textarea>
                    <button type="submit">Send Comment</button>
                </form>
                <div id="comments-section">
                    <h3>Comments</h3>
                    ${commentsHtml}
                </div>
            </div>
        </div>
    `;
    
    document.getElementById('add-comment-form')
        .addEventListener('submit', (e) => handleAddComment(e, postId));
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

function escapeHTML(str) {
    if (typeof str !== 'string') return str;
    return str
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
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
        
        // Sort users by last message time or alphabetically
        users.sort((a, b) => {
            // For now, just sort alphabetically since we don't have last message info yet
            return a.nickname.localeCompare(b.nickname);
        });

        users.forEach(user => {
            const userElement = document.createElement('div');
            userElement.className = 'user-item';
            userElement.dataset.userId = user.id;

            // Check if user is online (connected via WebSocket)
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
    // This would need to be implemented with a global list of online users
    // For now, return false as we don't have this implemented yet
    return false;
}

async function selectUserForChat(userId, nickname) {
    const messagesContainer = document.getElementById('messages-container');
    const chatForm = document.getElementById('chat-form');
    messagesContainer.innerHTML = `<h3>Chat with ${nickname}</h3>`;
    chatForm.style.display = 'flex';
    chatForm.dataset.recipientId = userId;
    
    // Fetch messages from the API
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
            
            // Reverse messages to show oldest first
            messages.reverse().forEach(msg => {
                const messageElement = document.createElement('div');
                const isSent = msg.sender_id === getCurrentUserId(); // We need to get current user ID
                messageElement.className = `chat-message ${isSent ? 'sent' : 'received'}`;
                
                const messageDate = new Date(msg.created_at).toLocaleString();

    //             messageElement.innerHTML = `
    //         <div class="message-header">
    //             <span class="message-author">${msg.senderId === 1 ? 'You' : nickname}</span>
    //             <span class="message-date">${messageDate}</span>
    //         </div>
    //         <p>${msg.content}</p>
    //     `;
    //     messagesContainer.appendChild(messageElement);
    // });

                messageElement.innerHTML = `
                    <div class="message-header">
                        <span class="message-author">${isSent ? 'You' : msg.sender_nickname}</span>
                        <span class="message-date">${messageDate}</span>
                    </div>
                    <p>${msg.content}</p>
                `;
                messagesContainer.appendChild(messageElement);
            });
            
            // Scroll to bottom
            messagesContainer.scrollTop = messagesContainer.scrollHeight;
        } else {
            messagesContainer.innerHTML += '<p>Failed to load messages</p>';
        }
    } catch (error) {
        console.error('Failed to fetch messages:', error);
        messagesContainer.innerHTML += '<p>Failed to load messages</p>';
    }
    
    // Add event listener for the chat form
    chatForm.addEventListener('submit', handleSendMessage);
}

function getCurrentUserId() {
    // This should be stored when user logs in
    // For now, return a placeholder
    return 1;
}

function handleIncomingPrivateMessage(message) {
    const messagesContainer = document.getElementById('messages-container');
    if (!messagesContainer) return; // Not on chat page
    
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

function renderErrorPage() {
    mainContent.innerHTML = `
        <h2>404 - Page Not Found</h2>
        <p>The page you're looking for doesn't exist.</p>
        <p><a href="#/posts">Go to Posts</a></p>
    `;
}

function renderPage(path) {
    const postIdMatch = path.match(/^#\/posts\/(\d+)$/);
    if (isLoggedIn()) {
        if (postIdMatch) {
            const postId = postIdMatch[1];
            showPostAndComments(postId);
            return;
        }
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
            case '':    
                renderPostsPage();
                break;
            case '#/error':
            default:
                navigate('#/error');
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
                renderLoginPage();
                break;
            case '#/posts':
            case '#/chats':
            case '#/':
            case '':
                navigate('#/login');
                break;
            default:
                console.log(path);
                navigate('#/error');
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
        displayPosts();
        hideCreatePostModal();
        ws.send(JSON.stringify({ type: 'new_post', payload: postData }));
    } else {
        alert(data.error);
    }
}

async function handleAddComment(e, postId) {
    e.preventDefault();
    const commentContent = document.getElementById('comment-content').value;

    const commentData = {
        post_id: parseInt(postId, 10),
        content: commentContent
    };

    const response = await fetch('/api/comments', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': `${userToken}`
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

async function handleVote(event) {
    const button = event.currentTarget;
    const postId = button.dataset.postId;
    const voteType = parseInt(button.dataset.voteType, 10);
    const isCurrentlyVoted = button.classList.contains('active-vote');
    const bodyData = {
        post_id: parseInt(postId, 10),
        vote: parseInt(voteType, 10),
    };

    let method;
    let url = '/api/posts/vote';

    if (isCurrentlyVoted) {
        method = 'DELETE';
    } else {
        method = 'POST';
    }

    try {
        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `${userToken}`
            },
            body: JSON.stringify(bodyData)
        });

        const data = await handleAuthResponse(response);
        if (data.success) {
            console.log('Vote action successful. Reloading posts to update counts.');
            displayPosts();
        } else {
            alert(data.error || 'Failed to handle vote.');
        }
    } catch (error) {
        console.error('Failed to handle vote:', error);
        alert('An error occurred while voting.');
    }
}

function handleLogout() {
    localStorage.removeItem('userToken');
    userToken = null;
    if (ws) {
        ws.close();
    }
    console.log("was");
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
    renderPage(location.hash || '#/posts');
});

loginLink.addEventListener('click', () => navigate('#/login'));
registerLink.addEventListener('click', () => navigate('#/register'));
postsLink.addEventListener('click', () => navigate('#/posts'));
chatsLink.addEventListener('click', () => navigate('#/chats'));
logoutLink.addEventListener('click', handleLogout);