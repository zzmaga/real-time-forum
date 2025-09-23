import { renderLoginPage, renderRegisterPage } from './auth.js';
import { renderPostsPage, showPostAndComments, renderCreatePostModal, displayPosts } from './posts.js';
import { renderChatsPage, handleIncomingPrivateMessage } from './chats.js';
import { isLoggedIn, updateNav, handleLogout, handleAuthResponse, connectWebSocket, ws, userToken } from './utils.js';

const mainContent = document.getElementById('main-content');
const loginLink = document.getElementById('login-link');
const registerLink = document.getElementById('register-link');
const postsLink = document.getElementById('posts-link');
const chatsLink = document.getElementById('chats-link');
const logoutLink = document.getElementById('logout-link');

export function navigate(path) {
    history.pushState(null, '', path);
    renderPage(path);
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