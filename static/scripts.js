document.addEventListener("DOMContentLoaded", () => {
    const createDeckBtn = document.getElementById("create-deck");
    const saveDeckBtn = document.getElementById("save-deck");
    const loadDeckBtn = document.getElementById("load-deck");
    const drawCardBtn = document.getElementById("draw-card");
    const discardCardBtn = document.getElementById("discard-card");
    const toggleAICardsBtn = document.getElementById("toggle-ai-cards");

    const tableList = document.getElementById("table-list");
    const handList = document.getElementById("hand-list");
    const aiHandList = document.getElementById("ai-hand-list");

    const sessionId = Date.now();

    let socket;
    let aiCardsRevealed = false;
    let cachedAIDeck = null;

    function initWebSocket() {
        const protocol = window.location.protocol === "https:" ? "wss" : "ws";
        const socketUrl = `${protocol}://${window.location.host}/api/ws/${sessionId}`;

        socket = new WebSocket(socketUrl);

        socket.onopen = () => {
        };

        socket.onmessage = (event) => {
            try {
                const data = JSON.parse(event.data);
                if (data.queue) updateTableUI(data.queue);
                if (data.deck) updateHandUI(data.deck);
                if (data.deckAI) {
                    cachedAIDeck = data.deckAI;
                    updateAIDeckUI(cachedAIDeck);
                }
            } catch (error) {
                console.error("Error parsing JSON:", error);
            }
        };

        socket.onclose = () => {
        };

        socket.onerror = (error) => {
            console.error("WebSocket error:", error);
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
            // Обновляем руку игрока, руку ИИ и стол
            updateHandUI(data.deck || []);
            updateAIDeckUI(data.deckAI);
            updateTableUI(data.queue || []);
            cachedAIDeck = data.deckAI;
        } catch (error) {
            console.error(error);
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

    loadDeckBtn.addEventListener("change", (event) => {
        const file = event.target.files[0];
        if (!file) return;
        const reader = new FileReader();
        reader.onload = function (e) {
            const content = e.target.result;
            const [tableContent, handContent] = content.split("\n\n");
            const tableCards = tableContent.split("\n").slice(1).map(card => card.trim()).filter(card => card);
            updateTableUI(tableCards);
            const handCards = handContent.split("\n").slice(1).map(card => card.trim()).filter(card => card);
            updateHandUI(handCards);
        };
        reader.readAsText(file);
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

        try {
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

            if (!leaveResponse.ok) {
                throw new Error(`Ошибка выбрасывания карт: ${leaveResponse.status}`);
            }

            const leaveData = await leaveResponse.json();

            updateHandUI(leaveData.deck || []);
            updateTableUI(leaveData.queue || []);
            if (leaveData.deckAI) updateAIDeckUI(leaveData.deckAI);

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

            if (!aiResponse.ok) {
                throw new Error(`Ошибка ИИ: ${aiResponse.status}`);
            }

            const aiData = await aiResponse.json();
            if (aiData.data) {
                if (aiData.data.handAI) {
                    updateAIDeckUI(aiData.data.handAI);
                    cachedAIDeck = aiData.data.handAI;
                }
                if (aiData.data.queue) {
                    updateTableUI(aiData.data.queue);
                }
            }
        } catch (error) {
            console.error("Ошибка при выбрасывании карт или запросе ИИ:", error);
            alert("Произошла ошибка, попробуйте снова.");
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
            if (data.deckAI) updateAIDeckUI(data.deckAI);
        } catch (error) {
            console.error("Ошибка при взятии карт:", error);
            alert("Произошла ошибка при взятии карты, попробуйте снова.");
        }
    });
});
