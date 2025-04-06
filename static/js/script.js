// 注册表单提交事件
document.getElementById('registrationForm').addEventListener('submit', function (e) {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;

    // 这里需要替换为实际的后端注册接口地址
    const registrationUrl = '/api/v3/user/register';
    fetch(registrationUrl, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            name: username,
            email: email,
            password: password
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.code == 200) {
            alert('注册成功.' + data.msg);
        } else {
            alert('注册失败：' + data.msg);
        }
    })
    .catch(error => {
        alert('网络错误：' + error.message);
    });
});