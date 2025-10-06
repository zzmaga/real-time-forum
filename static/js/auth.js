import { updateNav, handleAuthResponse, userToken, connectWebSocket } from './utils.js';
import { navigate} from './script.js'
import { displayPosts } from './posts.js';

const mainContent = document.getElementById('main-content');

export function renderLoginPage() {
    mainContent.innerHTML = `
        <h2>Login</h2>
        <form id="login-form">
            <input type="text" id="login" placeholder="Nickname or Email" required>
            <input type="password" id="password" placeholder="Password" required>
            <span id="errmess"></span>
            <button type="submit">Login</button>
        </form>
    `;
    document.getElementById('login-form').addEventListener('submit', handleLogin);
}

export function renderRegisterPage() {
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
            <input type="password" id="password" minlength="8" placeholder="Password" required>
            <label for="age">Date of Birth:</label>
            <input type="date" id="age" min="1900-01-01" max="${maxDate}" required>
            <select id="gender" required>
                <option value="">Select Gender</option>
                <option value="male">Male</option>
                <option value="female">Female</option>
                <option value="other">Other</option>
            </select>
            <span id="errmess"></span>
            <button type="submit">Register</button>
        </form>
    `;
    document.getElementById('register-form').addEventListener('submit', handleRegister);
}

async function handleLogin(event) {
    event.preventDefault();
    const loginId = document.getElementById('login').value;
    const password = document.getElementById('password').value;
    try {
        const response = await fetch('/api/signin', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ loginId, password })
        });
        
        const data = await response.json();
        
        if (data.success) {
            localStorage.setItem('userToken', data.token);
            const { setUserToken } = await import('./utils.js');
            setUserToken(data.token);
            
            updateNav();
            connectWebSocket();
            navigate('#/posts');
        } else {
            document.getElementById("errmess").innerHTML = data.error;
        }
    } catch (error) {
        console.error('Login error:', error);
        alert('Login failed. Please try again.');
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
        document.getElementById("errmess").innerHTML = data.error;
    }
}