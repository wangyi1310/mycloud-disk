document.addEventListener('DOMContentLoaded', function () {
    const registrationForm = document.getElementById('registrationForm');
    const regErrorMessage = document.getElementById('regErrorMessage');

    registrationForm.addEventListener('submit', async function (e) {
        e.preventDefault();

        const regUsername = document.getElementById('regUsername').value;
        const regEmail = document.getElementById('regEmail').value;
        const regPassword = document.getElementById('regPassword').value;

        try {
            const response = await fetch('/api/v3/user/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    name: regUsername,
                    email: regEmail,
                    password: regPassword
                })
            });

            if (!response.ok) {
                const errorData = await response.json();
                regErrorMessage.textContent = errorData.message || '注册失败，请重试';
                return;
            }

            const data = await response.json();
            alert('注册成功！' + data.msg);
            // 可添加注册成功后的其他逻辑，如跳转页面
        } catch (error) {
            console.error('注册请求出错:', error);
            regErrorMessage.textContent = '网络错误，请稍后重试';
        }
    });
});
