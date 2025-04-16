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
        } catch (error) {
            console.error('登录请求出错:', error);
            loginErrorMessage.textContent = '网络错误，请稍后重试';
        }
    });
});
