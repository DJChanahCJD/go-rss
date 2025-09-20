// 应用配置
const API_BASE_URL = 'http://localhost:8080'; // 后端API地址

// 工具函数
function showNotification(message, isError = false) {
    console.log(message)
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
    $('#login-form, #register-form, #home-content, #feeds-content, #square-content').hide();
    // 显示指定内容
    $(`#${id}`).show();
}

function setAuthState(isAuthenticated, userInfo = null) {
    if (isAuthenticated && userInfo) {
        $('#login-link, #register-link').hide();
        $('#user-info, #logout-link').show();
        $('#welcome-message').text(`欢迎, ${userInfo.Username}`);
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
            options.headers['Authorization'] = `${userInfo.ApiKey}`;
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
            
            if (data && data.length > 0) {
                data.forEach(post => {
                    const articleItem = $('<div class="article-item"></div>');
                    const title = $(`<div class="article-title"><a href="${post.Url}" target="_blank">${post.Title}</a></div>`);
                    const meta = $(`<div class="article-meta">来源: ${post.FeedName} | 发布时间: ${new Date(post.PublishedAt).toLocaleString()}</div>`);
                    const content = $(`<div class="article-content article-content-expanded" style="color: grey; font-size: 12px;">${post.Description.String || '暂无内容'}</div>`);
                    const divider = $('<hr class="article-divider">');
                    
                    articleItem.append(title).append(meta).append(content).append(divider);
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
    apiCall('GET', '/v1/feed_follows')
        .then(data => {
            const container = $('#feeds-container');
            container.empty();
            
            if (data && data.length > 0) {
                data.forEach(feed => {
                    const feedItem = $('<div class="feed-item ui segment"></div>');
                    const title = $(`<h4 class="ui header">${feed.FeedName}</h4>`);
                    const url = $(`<p><a href="${feed.FeedUrl}" target="_blank">${feed.FeedUrl}</a></p>`);
                    
                    // 取消订阅按钮
                    const unfollowBtn = $('<button class="ui button negative mini">取消订阅</button>');
                    unfollowBtn.click(function() {
                        unfollowFeed(feed.FeedID);
                    });
                    
                    feedItem.append(title).append(url).append(unfollowBtn);
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

// 加载广场订阅源（按订阅数排序）
function loadSquareFeeds() {
    apiCall('GET', '/v1/feeds', null, false)
        .then(data => {
            const container = $('#square-feeds-container');
            container.empty();
            
            if (data && data.length > 0) {
                data.forEach(feed => {
                    const feedItem = $('<div class="square-feed-item ui segment"></div>');
                    const title = $(`<h4 class="ui header">${feed.Name}</h4>`);
                    const url = $(`<p><a href="${feed.Url}" target="_blank">${feed.Url}</a></p>`);
                    const meta = $(`<p class="meta">订阅数: ${feed.FollowsCount.Int64 || 0} | 更新时间: ${new Date(feed.LastFetchedAt.Time).toLocaleDateString()}</p>`);
                    
                    // 订阅按钮
                    const followBtn = $('<button class="ui button primary mini">订阅</button>');
                    followBtn.click(function() {
                        followFeed(feed.ID);
                    });
                    
                    feedItem.append(title).append(url).append(meta).append(followBtn);
                    container.append(feedItem);
                });
            } else {
                container.append('<p>暂无订阅源。</p>');
            }
        })
        .catch(error => {
            showNotification('加载广场订阅源失败: ' + error.message, true);
            console.error('加载广场订阅源失败:', error);
        });
}

// 关注订阅源
function followFeed(feedId) {
    apiCall('POST', '/v1/feed_follows', { feed_id: feedId })
        .then(data => {
            showNotification('订阅成功！');
            loadFeeds(); // 刷新我的订阅列表
        })
        .catch(error => {
            showNotification('订阅失败: ' + error.message, true);
            console.error('订阅失败:', error);
        });
}

// 取消订阅源
function unfollowFeed(feedId) {
    if (confirm('确定要取消订阅吗？')) {
        apiCall('DELETE', `/v1/feed_follows/${feedId}`)
            .then(data => {
                showNotification('取消订阅成功！');
                loadFeeds(); // 刷新我的订阅列表
            })
            .catch(error => {
                showNotification('取消订阅失败: ' + error.message, true);
                console.error('取消订阅失败:', error);
            });
    }
}

// 添加订阅源
function addFeed(name, url) {
    apiCall('POST', '/v1/feeds', { name: name, url: url })
        .then(data => {
            showNotification('添加订阅源成功！');
            loadArticles();
            loadFeeds();
            $('#feed-name').val('')
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
    
    $('#square-link').click(function(e) {
        e.preventDefault();
        showContent('square-content');
        loadSquareFeeds();
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
                if (data) {
                    setAuthState(true, data);
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
                if (data) {
                    setAuthState(true, data);
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
        const name = $('#feed-name').val();
        const url = $('#feed-url').val();
        if (!url) {
            showNotification('请输入RSS订阅源URL', true);
            return;
        }
        addFeed(name, url);
    });
    
    // 检查认证状态
    checkAuthStatus();
});