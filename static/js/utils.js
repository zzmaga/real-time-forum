import { navigate, showAlert } from './script.js';
import { handleIncomingPrivateMessage } from './chats.js';
import { displayPosts } from './posts.js';

export let userToken = localStorage.getItem('userToken');
export let ws;

const loginLink = document.getElementById('login-link');
const registerLink = document.getElementById('register-link');
const postsLink = document.getElementById('posts-link');
const chatsLink = document.getElementById('chats-link');
const logoutLink = document.getElementById('logout-link');

export function isLoggedIn() {
    return !!userToken;
}

export function updateNav() {
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

export function setUserToken(token) {
    userToken = token;
}

export function handleAuthResponse(response) {
    if (response.status === 401) {
        localStorage.removeItem('userToken');
        userToken = null;
        updateNav();
        navigate('#/login');
        return null;
    }
    return response.json();
}

export function connectWebSocket() {
    if (ws && ws.readyState === WebSocket.OPEN) return;
    if (ws) {
        ws.onmessage = null;
        ws.onclose = null;
        ws.onerror = null;
        // Optionally close the old one, just in case
        if (ws.readyState !== WebSocket.CLOSED) {
            ws.close();
        }
    }
    ws = new WebSocket(`ws://localhost:8080/ws?token=${userToken}`);

    ws.onopen = () => {
        console.log('WebSocket connection established.');
    };

    ws.onmessage = (event) => {
        const message = JSON.parse(event.data);
        console.log('Message from server:', message);
        if (message.type === 'new_post') {
            if (window.location.hash === '#/posts') {
                //displayPosts();
                showAlert('infoAlert', 'NEW POST!', 'Update the page now!');
                //showAlert('successAlert');
                console.log("New post! I will do the notification later");
            }
        } else if (message.type === 'private_message') {
            handleIncomingPrivateMessage(message);
        }
    };

    ws.onclose = () => {
        console.log('WebSocket connection closed. Attempting to reconnect...');
        ws = null;
        setTimeout(connectWebSocket, 1000);
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

export function handleLogout() {
    localStorage.removeItem('userToken');
    userToken = null;
    if (ws) {
        ws.close();
    }
    updateNav();
    navigate('#/login');
}