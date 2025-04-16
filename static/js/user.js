// 获取用户信息展示元素
const usernameDisplay = document.getElementById('usernameDisplay');
const avatarDisplay = document.getElementById('avatarDisplay');

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
        }
    } catch (error) {
        console.error('检查登录状态时出错:', error);
        usernameDisplay.textContent = '未登录';
        avatarDisplay.style.display = 'none';
    }
}

// 页面加载时检查登录状态
window.addEventListener('load', checkLoginStatus);