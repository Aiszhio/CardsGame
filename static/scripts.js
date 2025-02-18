document.addEventListener("DOMContentLoaded", () => {
    const createDeckBtn = document.getElementById("create-deck");
    const saveDeckBtn = document.getElementById("save-deck");
    const loadDeckBtn = document.getElementById("load-deck");
    const drawCardBtn = document.getElementById("draw-card"); // кнопка для выбрасывания карт
    const toggleAICardsBtn = document.getElementById("toggle-ai-cards");

    const tableList = document.getElementById("table-list");
    const handList = document.getElementById("hand-list");
    const aiHandList = document.getElementById("ai-hand-list");
    const gamingTableList = document.getElementById("gaming-table-list");

    const sessionId = Date.now();

    let socket;
    let aiCardsRevealed = false;
    let cachedAIDeck = null;
    let gamingTableCards = [];

    function initWebSocket() {
        const protocol = window.location.protocol === "https:" ? "wss" : "ws";
        const socketUrl = `${protocol}://${window.location.host}/api/ws/${sessionId}`;

        socket = new WebSocket(socketUrl);

        socket.onopen = () => {
            console.log("WebSocket подключён");
        };

        socket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.queue) updateTableUI(data.queue);
                if (data.deck) updateHandUI(data.deck);
                if (data.deckAI) {
                    console.log("Received AI deck:", data.deckAI);
                    cachedAIDeck = data.deckAI;
                    updateAIDeckUI(cachedAIDeck);
                }
            } catch (error) {
                console.error("Error parsing JSON:", error);
            }
        };

        socket.onclose = () => {
            console.log("WebSocket закрыт");
        };

        socket.onerror = (error) => {
            console.error("WebSocket error:", error);
        };
    }

    initWebSocket();

    function updateGamingTableUI(cards) {
        gamingTableList.innerHTML = "";
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
            if (!rank || !suit) return;
            const li = document.createElement("li");
            li.dataset.card = `${rank} of ${suit}`;
            const fileName = `${rank} of ${suit}.png`;
            const img = document.createElement("img");
            img.src = `./static/images/${fileName}`;
            img.alt = `${rank} of ${suit}`;
            img.classList.add("card-image");
            li.appendChild(img);
            gamingTableList.appendChild(li);
        });
    }

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
            li.dataset.card = `${rank} of ${suit}`;
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

    function updateAIDeckUI(aiDeck) {
        aiHandList.innerHTML = "";
        let cards = [];
        if (Array.isArray(aiDeck)) {
            cards = aiDeck;
        } else if (aiDeck && typeof aiDeck === "object" && aiDeck.Currents && Array.isArray(aiDeck.Currents.Cards)) {
            cards = aiDeck.Currents.Cards;
        }
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
            li.dataset.card = `${rank} of ${suit}`;
            const img = document.createElement("img");
            if (!aiCardsRevealed) {
                img.src = `./static/images/Close.jpg`;
                img.alt = "Закрытая карта";
            } else {
                const fileName = `${rank} of ${suit}.png`;
                img.src = `./static/images/${fileName}`;
                img.alt = `${rank} of ${suit}`;
            }
            img.classList.add("card-image");
            li.appendChild(img);
            aiHandList.appendChild(li);
        });
    }

    toggleAICardsBtn.addEventListener("click", () => {
        aiCardsRevealed = !aiCardsRevealed;
        toggleAICardsBtn.textContent = aiCardsRevealed ? "Скрыть карты ИИ" : "Показать карты ИИ";
        if (cachedAIDeck) {
            updateAIDeckUI(cachedAIDeck);
        }
    });

    createDeckBtn.addEventListener("click", async () => {
        const response = await fetch("/api/deck/create", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-Session-ID": sessionId,
            },
        });
        const data = await response.json();
        if (data.error) {
            alert(data.error);
            return;
        }
        // Очищаем игровое поле выброшенных карт при начале игры
        gamingTableCards = [];
        updateGamingTableUI(gamingTableCards);
        updateHandUI(data.deck || []);
        updateAIDeckUI(data.deckAI);
        updateTableUI(data.queue || []);
        cachedAIDeck = data.deckAI;
    });

    saveDeckBtn.addEventListener("click", async () => {
        const handCards = Array.from(handList.querySelectorAll("li")).map(li => {
            const img = li.querySelector("img");
            return img ? img.alt : li.textContent;
        });
        const tableCards = Array.from(tableList.querySelectorAll("li")).map(li => li.textContent);
        const aiCards = Array.from(document.querySelectorAll("#ai-hand-list li")).map(li =>
            li.dataset.card || (li.querySelector("img") ? li.querySelector("img").alt : li.textContent)
        );

        const payload = {
            handCards: handCards,
            aiCards: aiCards,
            tableCards: tableCards
        };
        console.log("Отправляем данные для сохранения:", payload);

        const response = await fetch("/api/deck/save", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-Session-ID": sessionId,
            },
            body: JSON.stringify(payload)
        });
        const data = await response.json();
        if (data.error) {
            alert(data.error);
            return;
        }
        console.log("Получен URL для скачивания:", data.downloadUrl);

        const a = document.createElement("a");
        a.href = data.downloadUrl;
        a.download = data.downloadUrl.split("/").pop();
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
    });

    loadDeckBtn.addEventListener("change", async (event) => {
        const file = event.target.files[0];
        if (!file) return;
        const response = await fetch(`/api/deck/load?filename=${encodeURIComponent(file.name)}`, {
            method: "GET"
        });
        const data = await response.json();
        if (data.error) {
            alert(data.error);
            return;
        }
        updateTableUI(data.tableCards || []);
        updateHandUI(data.handCards || []);
        updateAIDeckUI(data.aiCards || []);
    });

    function getAICardNames() {
        return Array.from(document.querySelectorAll("#ai-hand-list li")).map(li =>
            li.dataset.card || (li.querySelector("img") ? li.querySelector("img").alt : li.textContent)
        );
    }

    function cardToString(card) {
        if (typeof card === "object" && card.Rank && card.Suit) {
            return `${card.Rank} of ${card.Suit}`;
        }
        return card;
    }

    async function checkStatus() {
        const handCards = Array.from(handList.querySelectorAll("li")).map(li => {
            const img = li.querySelector("img");
            return img ? img.alt : li.textContent;
        });
        const tableCards = Array.from(tableList.querySelectorAll("li")).map(li => li.textContent);
        const aiCards = getAICardNames();

        const statusPayload = {
            hand: handCards,
            deckAI: aiCards,
            table: tableCards
        };
        console.log("Отправляем данные на /api/deck/status:", statusPayload);

        const statusResponse = await fetch("/api/deck/status", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-Session-ID": sessionId,
            },
            body: JSON.stringify(statusPayload)
        });
        const statusData = await statusResponse.json();
        console.log("Получен ответ от /api/deck/status:", statusData);

        if (statusData.error) {
            alert(statusData.error);
            return;
        }
        updateHandUI(statusData.handCards);
        updateAIDeckUI(statusData.aiCards);
        updateTableUI(statusData.tableCards);
    }


    drawCardBtn.addEventListener("click", async () => {
        const handCards = Array.from(handList.querySelectorAll("li")).map(li => {
            const img = li.querySelector("img");
            return img ? img.alt : li.textContent;
        });
        const tableCards = Array.from(tableList.querySelectorAll("li")).map(li => li.textContent);

        const selectedCards = Array.from(handList.querySelectorAll("li.selected")).map(li =>
            li.dataset.card || (li.querySelector("img") ? li.querySelector("img").alt : '')
        ).filter(card => card !== '');

        if (selectedCards.length === 0) {
            alert("Выберите карту(ы) для обмена.");
            return;
        }

        // Проверяем, что все выбранные карты имеют одинаковый ранг.
        const ranks = selectedCards.map(card => {
            const parts = card.split(" of ");
            return parts[0] ? parts[0].trim() : "";
        });
        const firstRank = ranks[0];
        const allSame = ranks.every(rank => rank === firstRank);
        if (!allSame) {
            alert("Выбранные карты должны быть одной величины!");
            return;
        }

        // Если проверка пройдена, добавляем карту(ы) в область "выброшенные карты"
        selectedCards.forEach(card => {
            gamingTableCards.push(card);
        });
        updateGamingTableUI(gamingTableCards);

        console.log("Отправляем данные на /api/deck/leave:", {
            hand: handCards,
            table: tableCards,
            selected: selectedCards
        });

        const leaveResponse = await fetch("/api/deck/leave", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-Session-ID": sessionId,
            },
            body: JSON.stringify({
                hand: handCards,
                table: tableCards,
                selected: selectedCards,
            })
        });
        const leaveData = await leaveResponse.json();
        console.log("Ответ от /api/deck/leave:", leaveData);
        if (leaveData.error) {
            alert(leaveData.error);
            return;
        }
        updateHandUI(leaveData.deck || []);
        updateTableUI(leaveData.queue || []);

        console.log("Отправляем данные на /api/deck/ai:", {
            selected: selectedCards,
            handAI: getAICardNames()
        });

        const aiResponse = await fetch("/api/deck/ai", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "X-Session-ID": sessionId,
            },
            body: JSON.stringify({
                selected: selectedCards,
                handAI: getAICardNames()
            })
        });
        const aiData = await aiResponse.json();
        console.log("Ответ от /api/deck/ai:", aiData);
        if (aiData.error) {
            alert(aiData.error);
            return;
        }
        if (aiData.data) {
            if (aiData.data.handAI) {
                updateAIDeckUI(aiData.data.handAI);
                cachedAIDeck = aiData.data.handAI;
            }
            if (aiData.data.selectedAI && aiData.data.selectedAI.length > 0) {
                aiData.data.selectedAI.forEach(card => {
                    gamingTableCards.push(card);
                });
                updateGamingTableUI(gamingTableCards);
            }
            if (aiData.data.queue) {
                updateTableUI(aiData.data.queue);
            }
            console.log("Selected AI:", aiData.data.selectedAI);
        }
        if (gamingTableCards.length >= 2) {
            setTimeout(() => {
                gamingTableCards = [];
                updateGamingTableUI(gamingTableCards);
            }, 4000);
        }
        await checkStatus();
    });
});
