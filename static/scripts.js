document.addEventListener("DOMContentLoaded", () => {
    const createDeckBtn = document.getElementById("create-deck");
    const saveDeckBtn = document.getElementById("save-deck");
    const drawCardBtn = document.getElementById("draw-card");
    const discardCardBtn = document.getElementById("discard-card");

    const tableList = document.getElementById("table-list");
    const handList = document.getElementById("hand-list");

    const sessionId = Date.now();
    console.log("Session ID:", sessionId);

    let socket;

    function initWebSocket() {
        const protocol = window.location.protocol === "https:" ? "wss" : "ws";
        const socketUrl = `${protocol}://${window.location.host}/api/ws/${sessionId}`;
        console.log("Подключаемся к WebSocket:", socketUrl);

        socket = new WebSocket(socketUrl);

        socket.onopen = () => {
            console.log("WebSocket соединение установлено для сессии:", sessionId);
        };

        socket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.queue) updateTableUI(data.queue);
                if (data.deck) updateHandUI(data.deck);
            } catch (error) {
                console.error("Ошибка парсинга JSON:", error);
            }
        };

        socket.onclose = () => {
            console.log("WebSocket соединение закрыто");
        };

        socket.onerror = (error) => {
            console.error("WebSocket ошибка:", error);
        };
    }

    initWebSocket();

    function updateTableUI(cards) {
        tableList.innerHTML = "";
        cards.forEach((card) => {
            let rank, suit;
            if (typeof card === "string") {
                const parts = card.split(" of ");
                rank = parts[0] ? parts[0].trim() : "";
                suit = parts[1] ? parts[1].trim() : "";
            } else {
                rank = card.Rank;
                suit = card.Suit;
            }
            const li = document.createElement("li");
            li.textContent = rank && suit ? `${rank} of ${suit}` : card;
            tableList.appendChild(li);
        });
    }

    function updateHandUI(cards) {
        handList.innerHTML = "";
        cards.forEach((card, index) => {
            let rank, suit;
            if (typeof card === "string") {
                const parts = card.split(" of ");
                rank = parts[0] ? parts[0].trim() : "";
                suit = parts[1] ? parts[1].trim() : "";
            } else {
                rank = card.Rank;
                suit = card.Suit;
            }

            if (!rank || !suit) return;

            const li = document.createElement("li");
            li.dataset.index = index;

            const fileName = `${rank} of ${suit}.png`;
            const img = document.createElement("img");
            img.src = `./static/images/${fileName}`;
            img.alt = `${rank} of ${suit}`;
            img.classList.add("card-image");

            li.appendChild(img);

            li.addEventListener("click", function () {
                this.classList.toggle("selected");
            });

            handList.appendChild(li);
        });
    }

    createDeckBtn.addEventListener("click", async () => {
        try {
            const response = await fetch("/api/deck/create", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-Session-ID": sessionId,
                },
            });
            if (!response.ok) {
                throw new Error(`Ошибка сервера: ${response.status}`);
            }
            const data = await response.json();
            updateHandUI(data.deck || []);
            updateTableUI(data.queue || []);
            console.log("Колода создана");
        } catch (error) {
            console.error("Ошибка при создании колоды:", error);
        }
    });

    saveDeckBtn.addEventListener("click", () => {
        let content = "Стол:\n";
        tableList.querySelectorAll("li").forEach((li) => {
            content += li.textContent + "\n";
        });
        content += "\nВаша Рука:\n";
        handList.querySelectorAll("li").forEach((li) => {
            const img = li.querySelector("img");
            content += img ? img.alt : li.textContent;
            content += "\n";
        });

        const blob = new Blob([content], { type: "text/plain;charset=utf-8" });
        const url = URL.createObjectURL(blob);
        const a = document.createElement("a");
        a.href = url;
        a.download = "deck.txt";
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    });

    drawCardBtn.addEventListener("click", async () => {
        const handCards = Array.from(handList.querySelectorAll("li")).map(li => {
            const img = li.querySelector("img");
            return img ? img.alt : li.textContent;
        });

        const tableCards = Array.from(tableList.querySelectorAll("li")).map(li => li.textContent);

        const selectedCards = Array.from(handList.querySelectorAll("li.selected")).map(li => {
            const img = li.querySelector("img");
            return img ? img.alt : li.textContent;
        });

        if (selectedCards.length === 0) {
            alert("Выберите карту(ы) для выбрасывания.");
            return;
        }

        try {
            const response = await fetch("/api/deck/leave", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-Session-ID": sessionId,
                },
                body: JSON.stringify({
                    hand: handCards,
                    table: tableCards,
                    selected: selectedCards
                })
            });

            if (!response.ok) {
                throw new Error(`Ошибка сервера: ${response.status}`);
            }

            const data = await response.json();
            console.log("Получены данные от сервера:", data);

            updateHandUI(data.deck || []);
            updateTableUI(data.queue || []);
            console.log("Карта выброшена, состояние обновлено");
        } catch (error) {
            console.error("Ошибка при выбрасывании карты:", error);
        }
    });

    discardCardBtn.addEventListener("click", async () => {
        const handCards = Array.from(handList.querySelectorAll("li")).map(li => {
            const img = li.querySelector("img");
            return img ? img.alt : li.textContent;
        });

        const tableCards = Array.from(tableList.querySelectorAll("li")).map(li => li.textContent);

        try {
            const response = await fetch("/api/deck/take", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                    "X-Session-ID": sessionId,
                },
                body: JSON.stringify({
                    hand: handCards,
                    table: tableCards
                })
            });
            if (!response.ok) {
                throw new Error(`Ошибка сервера: ${response.status}`);
            }
            const data = await response.json();
            updateHandUI(data.deck || []);
            updateTableUI(data.queue || []);
            console.log("Карта взята, состояние обновлено");
        } catch (error) {
            console.error("Ошибка при взятии карты:", error);
        }
    });
});
