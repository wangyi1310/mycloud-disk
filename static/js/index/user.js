// 获取用户信息展示元素
const usernameDisplay = document.getElementById('usernameDisplay');
const avatarDisplay = document.getElementById('avatarDisplay');
const loginButton = document.getElementById('login-btn');
const registerButton = document.getElementById('register-btn');
const logoutButton = document.getElementById('logout-btn');
let isLoggedIn = false;

document.addEventListener('DOMContentLoaded', function () {
    const loginForm = document.getElementById('loginForm');
    const loginErrorMessage = document.getElementById('loginErrorMessage');

    loginForm.addEventListener('submit', async function (e) {
        e.preventDefault();

        const username = document.getElementById('username').value;
        const password = document.getElementById('password').value;

        try {
            const response = await fetch('/api/v3/user/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    username: username,
                    password: password
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                loginErrorMessage.textContent = errorData.message || '登录失败，请重试';
                return;
            }

            const data = await response.json();
            // 登录成功后的处理逻辑，例如跳转到主页
            if (data.code !== 0) {
                alert(data.msg || '登录失败，请重试');
                return;
            }
            alert('登录成功！');
            loginModal.classList.remove('open');
            window.location.reload();
        } catch (error) {
            console.error('登录请求出错:', error);
            loginErrorMessage.textContent = '网络错误，请稍后重试';
        }
    });
});


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
            loginModal.classList.add('open');
            window.location.reload();

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