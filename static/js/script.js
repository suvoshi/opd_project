function changeContent(type) {
    const contentDiv = document.getElementById('content');
    
    if (type === 'first'){
        contentDiv.innerHTML = `
            <h1>Доска объявлений ГЭТ</h1>
        `;
    } else if (type === 'second') {
        contentDiv.innerHTML = `
            <h1>Андросов Виктор Максимович</h1>
            <p>Номер группы: 5130902/40001</p>
            <p>Средний балл: 4.88</p>
            <p>Размер текущей стипендии: 22500 р.</p>
            <p><em>Обновлено: ${new Date().toLocaleString()}</em></p>
        `;
    } else if (type === 'third') {
        contentDiv.innerHTML = `
            <h1>Расписание</h1>
        `;
    } else if (type === 'fourth') {
        contentDiv.innerHTML = `
            <h1>Дисциплины</h1>
            <table>
                <tr><th>Предмет</th><th>Количество занятий</th><th>Количество часов</th><th>Оценки</th><th>Средний балл</th></tr>
                <tr><td>Математика</td><td></td><td></td><td>2</td></tr>
                <tr><td>Физика</td><td></td><td></td><td>4</td></tr>
                <tr><td>Литература</td><td></td><td></td><td>5</td></tr>
                <tr><td>ОБЖ</td><td></td><td></td><td>5</td></tr>
            </table>
        `;
    }
}

function setupButtons() {
    const buttons = document.querySelectorAll('.text-button-alt');
    
    buttons.forEach(button => {
        button.addEventListener('click', function() {
            const type = this.getAttribute('data-type'); 
            changeContent(type);
        });
    });
}

document.addEventListener('DOMContentLoaded', setupButtons);

function showNotification(message) {
    console.log('Уведомление:', message);
}

function updateTimestamp() {
    const timestampElement = document.getElementById('timestamp');
    if (timestampElement) {
        timestampElement.textContent = `Последнее обновление: ${new Date().toLocaleString()}`;
    }
}

setInterval(updateTimestamp, 60000);