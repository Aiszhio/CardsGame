document.addEventListener("DOMContentLoaded", function() {
    const createDeckBtn = document.getElementById("create-deck");
    const uploadDeckInput = document.getElementById("upload-deck");
    const saveDeckBtn = document.getElementById("save-deck");

    const deckList = document.getElementById("deck-list");
    const handList = document.getElementById("hand-list");
    const deckCount = document.getElementById("deck-count");

    let socket;
    let id;

    // Генерация sessionId
    function generateSessionId() {
        return Date.now(); // Генерация sessionId на основе времени (можно добавить случайные символы, если нужно)
    }

    console.log(generateSessionId());

    // Инициализация WebSocket
    function initWebSocket() {
        const protocol = window.location.protocol === 'https:' ? 'wss' : 'ws';
        const socketUrl = `${protocol}://${window.location.host}/api/ws/${id}`;

        console.log("Попытка подключиться к WebSocket:", socketUrl);

        socket = new WebSocket(socketUrl);

        // Обработчик на открытие соединения
        socket.onopen = function() {
            console.log("WebSocket соединение установлено для сессии:", id);
        };

        // Обработчик на получение сообщения
        socket.onmessage = function(event) {
            try {
                const data = JSON.parse(event.data);
                console.log("Получено сообщение:", data);
                if (data.deck) updateDeckUI(data.deck);
                if (data.hand) updateHandUI(data.hand);
            } catch (e) {
                console.error("Ошибка парсинга JSON:", e, "Полученные данные:", event.data);
            }
        };

        // Обработчик на закрытие соединения
        socket.onclose = function() {
            console.log("WebSocket соединение закрыто");
        };

        // Обработчик на ошибку соединения
        socket.onerror = function(event) {
            console.error("WebSocket ошибка:", event);
            console.log("Ошибка в соединении с WebSocket по адресу:", socket.url);
        };
    }

    createDeckBtn.addEventListener("click", function() {
        // 1. Генерация sessionId
        id = generateSessionId();

        // 2. Отправка запроса на сервер для создания колоды
        fetch(`/api/deck/create?session=${id}`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-Session-ID": id
            }
        })
            .then(response => {
                if (!response.ok) throw new Error(`HTTP error! status: ${response.status}`);
                return response.json();
            })
            .then(data => {
                if (data.deck) {
                    // 3. После успешного ответа инициализируем WebSocket
                    initWebSocket();
                    alert("Игра начата!");
                    updateDeckUI(data.deck);
                }
            })
            .catch(error => {
                console.error("Ошибка:", error);
                alert("Ошибка: " + error.message);
            });
    });

    // Обработчик на загрузку колоды из файла
    uploadDeckInput.addEventListener("change", function() {
        const file = this.files[0];
        if (file) {
            const reader = new FileReader();
            reader.onload = function(e) {
                const content = e.target.result;
                fetch("/upload", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json"
                    },
                    body: JSON.stringify({ deck: content })
                })
                    .then(response => response.json())
                    .then(data => {
                        if (data.deck) {
                            alert("Колода успешно загружена!");
                            updateDeckUI(data.deck); // Обновление UI колоды
                            updateHandUI(data.hand); // Обновление UI руки
                        } else {
                            alert("Ошибка загрузки колоды: " + data.message);
                        }
                    })
                    .catch(error => console.error("Ошибка:", error));
            };
            reader.readAsText(file); // Чтение содержимого файла
        }
    });

    // Обработчик на сохранение колоды
    saveDeckBtn.addEventListener("click", function() {
        fetch("/save", {
            method: "GET"
        })
            .then(response => {
                if (response.ok) {
                    return response.blob(); // Сохранение колоды в виде файла
                } else {
                    throw new Error("Ошибка сохранения колоды");
                }
            })
            .then(blob => {
                const url = URL.createObjectURL(blob);
                const a = document.createElement("a");
                a.href = url;
                a.download = "deck.json"; // Имя файла для сохранения
                document.body.appendChild(a);
                a.click(); // Имитируем клик для загрузки
                a.remove();
                URL.revokeObjectURL(url); // Освобождение ресурсов
            })
            .catch(error => {
                console.error("Ошибка:", error);
                alert("Ошибка сохранения колоды: " + error.message);
            });
    });

    // Функция для обновления UI колоды
    function updateDeckUI(deck) {
        deckList.innerHTML = ""; // Очищаем текущую колоду на странице
        deck.forEach(card => {
            const li = document.createElement("li");
            li.textContent = `${card.Rank} of ${card.Suit}`;
            deckList.appendChild(li);
        });

        // Обновляем количество карт в колоде
        deckCount.textContent = `Оставшиеся карты: ${deck.length}`;
    }

    // Функция для обновления UI руки
    function updateHandUI(hand) {
        handList.innerHTML = ""; // Очищаем текущую руку на странице
        hand.forEach(card => {
            const li = document.createElement("li");
            li.textContent = `${card.Rank} of ${card.Suit}`;
            handList.appendChild(li);
        });
    }
});
