import { navigate, showAlert } from './script.js';
import { handleIncomingPrivateMessage, getCurrentUserId, fetchUsers } from './chats.js';


export let userToken = localStorage.getItem('userToken');
export let ws;
let currentUserId = null;
export let onlineUserIds = []; // Store online user IDs 

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

export async function connectWebSocket() {
    if (currentUserId === null) {
        currentUserId = await getCurrentUserId(); 
    }
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
                if(currentUserId !== message.payload.sender_id){
                    showAlert('infoAlert', 'NEW POST!', 'Update the page now!');
                } else{
                    showAlert('successAlert', 'Post created!', 'Page is updated');
                }
            }
        } else if (message.type === 'private_message') {
            fetchUsers();
            showAlert('infoAlert', 'NEW MESSAGE!', 'from '+message.payload.sender_name);
            handleIncomingPrivateMessage(message);
        } else if (message.type === 'online_users') {
            // Update online users list
            onlineUserIds = message.payload.online_user_ids || [];
            console.log('Online users updated:', onlineUserIds);
            // Trigger UI update if on chats page
            if (window.updateUsersList) {
                window.updateUsersList();
            }
        } else if(message.type === 'new_comment'){
            const text = window.location.hash;
            const postIdMatch = text.match(/^#\/posts\/(\d+)$/);
            const postId = postIdMatch[1];
            if(postId == message.payload.post_id && currentUserId !== message.payload.sender_id){
                console.log(postId);
                showAlert('infoAlert', 'NEW COMMENT', 'Update the page!');
            } else if(currentUserId === message.payload.sender_id){
                showAlert('successAlert', "COMMENT created!", "Page is updated");
            }
        }else if (message.type === 'message_error') {
            // Handle error message when trying to send to offline user
            showAlert('errorAlert', 'Message Failed', message.payload.error);
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