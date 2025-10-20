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
const alertsWrapper = document.getElementById('alerts-wrapper');

export function navigate(path) {
    history.pushState(null, '', path);
    renderPage(path);
}

function renderErrorPage() {
    mainContent.innerHTML = `
        <h2>404 - Page Not Found<p><a href="#/posts">Go to Posts</a></p></h2>
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
                renderErrorPage();
                break;
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

export function setAlertText(alertId, strongText, bodyText) {
    const alertElement = document.getElementById(alertId);
    if (alertElement) {
        const strongElement = alertElement.querySelector('.font__weight-semibold');
        if (strongElement) {
            strongElement.textContent = strongText;
            if (strongElement.nextSibling && strongElement.nextSibling.nodeType === 3) {
                strongElement.nextSibling.textContent = bodyText;
            } else {
                strongElement.insertAdjacentText('afterend', bodyText);
            }
        }
    } else {
        console.error(`Alert element with ID '${alertId}' not found.`);
    }
}

export function showAlert(alertId, maintxt, sectxt) {
    setAlertText(alertId, maintxt, sectxt);
    const alertElement = document.getElementById(alertId);
    if (alertElement) {
        //alertElement.textContent = "New post added! Update the page!";
        alertsWrapper.prepend(alertElement);
        alertElement.classList.remove('hidden');
        setTimeout(() => {
            hideAlert(alertId);
        }, 5000);
    } else {
        console.error(`Alert element with ID '${alertId}' not found.`);
    }
}

export function hideAlert(alertId) {
    const alertElement = document.getElementById(alertId);
    if (alertElement) {
        alertElement.classList.add('hidden');
    } else {
        console.error(`Alert element with ID '${alertId}' not found.`);
    }
}

function setupNotificationControls() {
    document.getElementById('btn-success')?.addEventListener('click', () => showAlert('successAlert'));
    document.getElementById('btn-info')?.addEventListener('click', () => showAlert('infoAlert'));
    document.getElementById('btn-warning')?.addEventListener('click', () => showAlert('warningAlert'));
    document.getElementById('btn-danger')?.addEventListener('click', () => showAlert('dangerAlert'));
    document.getElementById('btn-primary')?.addEventListener('click', () => showAlert('primaryAlert'));
    document.querySelectorAll('.alert-close').forEach(button => {
        button.addEventListener('click', (e) => {
            const alertId = e.currentTarget.getAttribute('data-alert-id');
            if (alertId) {
                hideAlert(alertId);
            }
        });
    });
}

window.addEventListener('popstate', () => renderPage(location.hash));
document.addEventListener('DOMContentLoaded', () => {
    updateNav();
    if (isLoggedIn()) {
        connectWebSocket();
    }
    renderCreatePostModal();
    setupNotificationControls();
    renderPage(location.hash || '#/posts');
});

loginLink.addEventListener('click', () => navigate('#/login'));
registerLink.addEventListener('click', () => navigate('#/register'));
postsLink.addEventListener('click', () => navigate('#/posts'));
chatsLink.addEventListener('click', () => navigate('#/chats'));
logoutLink.addEventListener('click', handleLogout);