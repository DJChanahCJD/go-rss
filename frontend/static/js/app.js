// 应用配置
const API_BASE_URL = 'http://localhost:8080'; // 后端API地址

// 工具函数
function showNotification(message, isError = false) {
    const notification = $('#notification');
    const notificationMessage = $('#notification-message');
    
    notificationMessage.text(message);
    notification.removeClass('positive negative').addClass(isError ? 'negative' : 'positive');
    notification.show();
    
    setTimeout(() => {
        notification.hide();
    }, 3000);
}

function showContent(id) {
    // 隐藏所有内容
    $('#login-form, #register-form, #home-content, #feeds-content').hide();
    // 显示指定内容
    $(`#${id}`).show();
}

function setAuthState(isAuthenticated, userInfo = null) {
    if (isAuthenticated && userInfo) {
        $('#login-link, #register-link').hide();
        $('#user-info, #logout-link').show();
        $('#welcome-message').text(`欢迎, ${userInfo.username}`);
        localStorage.setItem('user', JSON.stringify(userInfo));
    } else {
        $('#login-link, #register-link').show();
        $('#user-info, #logout-link').hide();
        localStorage.removeItem('user');
    }
}

// 检查用户认证状态
function checkAuthStatus() {
    const user = localStorage.getItem('user');
    if (user) {
        try {
            const userInfo = JSON.parse(user);
            setAuthState(true, userInfo);
            showContent('home-content');
            loadArticles();
        } catch (e) {
            localStorage.removeItem('user');
            showContent('login-form');
        }
    } else {
        showContent('login-form');
    }
}

// API调用函数
function apiCall(method, endpoint, data = null, requiresAuth = true) {
    const options = {
        method: method,
        headers: {
            'Content-Type': 'application/json'
        }
    };
    
    if (requiresAuth) {
        const user = localStorage.getItem('user');
        if (user) {
            const userInfo = JSON.parse(user);
            options.headers['Authorization'] = `ApiKey ${userInfo.apiKey}`;
        }
    }
    
    if (data) {
        options.body = JSON.stringify(data);
    }
    
    return fetch(`${API_BASE_URL}${endpoint}`, options)
        .then(response => {
            if (!response.ok) {
                throw new Error(`API错误: ${response.status}`);
            }
            return response.json();
        });
}

// 加载文章
function loadArticles() {
    apiCall('GET', '/v1/posts')
        .then(data => {
            const container = $('#articles-container');
            container.empty();
            
            if (data.posts && data.posts.length > 0) {
                data.posts.forEach(post => {
                    const articleItem = $('<div class="article-item"></div>');
                    const title = $(`<div class="article-title"><a href="${post.url}" target="_blank">${post.title}</a></div>`);
                    const meta = $(`<div class="article-meta">来源: ${post.feedName} | 发布时间: ${new Date(post.publishedAt).toLocaleString()}</div>`);
                    const content = $(`<div class="article-content">${post.description || '暂无内容'}</div>`);
                    
                    articleItem.append(title).append(meta).append(content);
                    container.append(articleItem);
                });
            } else {
                container.append('<p>暂无文章，请添加RSS订阅源。</p>');
            }
        })
        .catch(error => {
            showNotification('加载文章失败: ' + error.message, true);
            console.error('加载文章失败:', error);
        });
}

// 加载订阅源
function loadFeeds() {
    apiCall('GET', '/v1/feeds')
        .then(data => {
            const container = $('#feeds-container');
            container.empty();
            
            if (data.feeds && data.feeds.length > 0) {
                data.feeds.forEach(feed => {
                    const feedItem = $('<div class="feed-item"></div>');
                    const title = $(`<span class="feed-title">${feed.name}</span>`);
                    const url = $(`<span class="feed-url">${feed.url}</span>`);
                    
                    feedItem.append(title).append(url);
                    container.append(feedItem);
                });
            } else {
                container.append('<p>暂无订阅源，请添加。</p>');
            }
        })
        .catch(error => {
            showNotification('加载订阅源失败: ' + error.message, true);
            console.error('加载订阅源失败:', error);
        });
}

// 添加订阅源
function addFeed(url) {
    apiCall('POST', '/v1/feeds', { url: url })
        .then(data => {
            showNotification('添加订阅源成功！');
            loadArticles();
            $('#feed-url').val('');
        })
        .catch(error => {
            showNotification('添加订阅源失败: ' + error.message, true);
            console.error('添加订阅源失败:', error);
        });
}

// 事件监听
$(document).ready(function() {
    // 导航链接
    $('#login-link').click(function(e) {
        e.preventDefault();
        showContent('login-form');
    });
    
    $('#register-link').click(function(e) {
        e.preventDefault();
        showContent('register-form');
    });
    
    $('#home-link').click(function(e) {
        e.preventDefault();
        showContent('home-content');
        loadArticles();
    });
    
    $('#feeds-link').click(function(e) {
        e.preventDefault();
        showContent('feeds-content');
        loadFeeds();
    });
    
    // 登录按钮
    $('#login-btn').click(function() {
        const username = $('#login-username').val();
        const password = $('#login-password').val();
        
        if (!username || !password) {
            showNotification('请输入用户名和密码', true);
            return;
        }
        
        apiCall('POST', '/v1/users/login', { username: username, password: password }, false)
            .then(data => {
                if (data.user) {
                    setAuthState(true, data.user);
                    showContent('home-content');
                    loadArticles();
                    showNotification('登录成功！');
                }
            })
            .catch(error => {
                showNotification('登录失败: ' + error.message, true);
                console.error('登录失败:', error);
            });
    });
    
    // 注册按钮
    $('#register-btn').click(function() {
        const username = $('#register-username').val();
        const password = $('#register-password').val();
        
        if (!username || !password) {
            showNotification('请输入用户名和密码', true);
            return;
        }
        
        apiCall('POST', '/v1/users', { username: username, password: password }, false)
            .then(data => {
                if (data.user) {
                    setAuthState(true, data.user);
                    showContent('home-content');
                    loadArticles();
                    showNotification('注册成功！');
                }
            })
            .catch(error => {
                showNotification('注册失败: ' + error.message, true);
                console.error('注册失败:', error);
            });
    });
    
    // 退出按钮
    $('#logout-link').click(function(e) {
        e.preventDefault();
        setAuthState(false);
        showContent('login-form');
        showNotification('已退出登录');
    });
    
    // 添加订阅源按钮
    $('#add-feed-btn').click(function() {
        const url = $('#feed-url').val();
        if (!url) {
            showNotification('请输入RSS订阅源URL', true);
            return;
        }
        addFeed(url);
    });
    
    // 检查认证状态
    checkAuthStatus();
});