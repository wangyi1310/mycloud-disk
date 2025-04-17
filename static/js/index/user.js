// 获取用户信息展示元素
const usernameDisplay = document.getElementById('usernameDisplay');
const avatarDisplay = document.getElementById('avatarDisplay');
const loginButton = document.getElementById('login-btn');
const registerButton = document.getElementById('register-btn');
const logoutButton = document.getElementById('logout-btn');
let isLoggedIn = false;

// 控制登出按钮的显示和隐藏
function toggleLogoutButton() {
    if (isLoggedIn) {
        if (logoutButton.style.display === 'none') {
            logoutButton.style.display = 'block';
        } else {
            logoutButton.style.display = 'none';
        }
    }
}
// 检查用户登录状态
async function checkLoginStatus() {
    try {
        const response = await fetch('/api/v3/user/info');
        if (response.status === 401) {
            // 未登录
            usernameDisplay.textContent = '未登录';
            avatarDisplay.style.display = 'none';
            return;
        }
        const data = await response.json();
        if (data) {
            isLoggedIn = true;
            // 登录成功，更新用户信息
            usernameDisplay.textContent = data.data.Nick || data.data.Email;
            if (data.data.Avatar) {
                avatarDisplay.src = data.data.Avatar;
                avatarDisplay.style.display = 'block';
            } else {
                // 若没有头像，可使用默认头像
                avatarDisplay.src = 'default-avatar.png'; 
                avatarDisplay.style.display = 'block';
            }
            loginButton.style.display = 'none';
            registerButton.style.display = 'none';

        }
    } catch (error) {
        console.error('检查登录状态时出错:', error);
        usernameDisplay.textContent = '未登录';
        avatarDisplay.style.display = 'none';
    }
}

// 处理登出请求
async function logout() {
    if (!isLoggedIn) {
        return;
    }
    try {
        const response = await fetch('/api/v3/user/logout', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            }
        });

        if (response.ok) {
            alert('登出成功');
            isLoggedIn = false;
            // 登出成功后，更新页面状态
            usernameDisplay.textContent = '未登录';
            avatarDisplay.style.display = 'none';
            logoutButton.style.display = 'none';
        } else {
            const errorData = await response.json();
            alert(`登出失败: ${errorData.message || '未知错误'}`);
        }
    } catch (error) {
        console.error('登出请求出错:', error);
        alert('登出请求出错，请稍后重试');
    }
}

// 页面加载时检查登录状态
window.addEventListener('load', checkLoginStatus);