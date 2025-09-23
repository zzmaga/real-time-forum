import { handleAuthResponse, userToken, ws } from './utils.js';
import { navigate } from './script.js';

const mainContent = document.getElementById('main-content');

export function renderPostsPage() {
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

export function renderCreatePostModal() {
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

export async function displayPosts() {
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
    console.log(posts)
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
        smallElement.textContent = `By: ${post.WUser.Nickname} | Categories: ${categories}`;
        const actionsElement = document.createElement('div');
        actionsElement.className = 'post-actions';
        actionsElement.innerHTML = `
            <button class="vote-btn ${post.WUserVote == 1 ? 'active-vote' : ''}" data-post-id="${post.Id}" data-vote-type="1">👍 <span class="like-count">${post.WVoteUp}</span></button>
            <button class="vote-btn ${post.WUserVote == -1 ? 'active-vote' : ''}" data-post-id="${post.Id}" data-vote-type="-1">👎 <span class="dislike-count">${post.WVoteDown}</span></button>
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


export async function showPostAndComments(postId) {
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

function escapeHTML(str) {
    if (typeof str !== 'string') return str;
    return str
        .replace(/&/g, "&amp;")
        .replace(/</g, "&lt;")
        .replace(/>/g, "&gt;")
        .replace(/"/g, "&quot;")
        .replace(/'/g, "&#039;");
}