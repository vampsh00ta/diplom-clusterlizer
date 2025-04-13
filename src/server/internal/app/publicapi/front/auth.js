

function onTelegramAuth(user) {

    const userJsonString = JSON.stringify(user);
    localStorage.setItem("user",userJsonString)
}

function checkAuthWidget() {
    const user = checkAuthorization();
    const authWidget = document.getElementById('authclient-widget');

    if (user) {
        authWidget.textContent = user.username; // Показываем никнейм пользователя
        authWidget.classList.remove('cursor-pointer');
    } else {
        authWidget.textContent = 'Авторизироваться';
        authWidget.addEventListener('click', () => {
            window.location.href = "https://t.me/JobSearcherHelper_bot?start=asdfasdf"
        });
    }
}

// Выполнение проверки авторизации при загрузке страницы
document.addEventListener('DOMContentLoaded', checkAuthWidget);
// Проверка авторизации по JWT токену
function checkAuthorization() {
    const token = getCookie("access_token");
    if (!token) {
        return null;
    }
    try {
        const payload = JSON.parse(atob(token.split('.')[1])); // Декодирование JWT payload
        return payload;
    } catch (error) {
        console.error('Ошибка декодирования токена:', error);
        alert('Ошибка авторизации');
        return null;
    }
}

// Вспомогательная функция для получения значения cookie
function getCookie(name) {
    const matches = document.cookie.match(new RegExp(
        "(?:^|; )" + name.replace(/([\.$?*|{}\(\)\[\]\\\/\+^])/g, '\\$1') + "=([^;]*)"
    ));
    return matches ? decodeURIComponent(matches[1]) : undefined;
}
