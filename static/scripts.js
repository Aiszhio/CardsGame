(() => {
    const qs               = s  => document.querySelector(s);
    const createDeckBtn    = qs('#create-deck');
    const saveDeckBtn      = qs('#save-deck');
    const loadDeckInput    = qs('#load-deck');
    const drawCardBtn      = qs('#draw-card');
    const toggleAICardsBtn = qs('#toggle-ai-cards');
    const tableList        = qs('#table-list');
    const handList         = qs('#hand-list');
    const aiHandList       = qs('#ai-hand-list');
    const gamingTableList  = qs('#gaming-table-list');

    const sessionId         = Date.now();
    let socket;
    let aiCardsRevealed    = false;
    let cachedAIDeck       = null;
    let gamingTableCards   = [];
    let inFlight           = false;
    let isLoadingFromSave = false;

    const sleep            = ms => new Promise(r => setTimeout(r, ms));
    const cardToString     = c  => (typeof c === 'object' && c?.Rank) ? `${c.Rank} of ${c.Suit}` : String(c);
    const renderList       = (ul, arr, fn) => {
        ul.innerHTML = '';
        arr.forEach((c, i) => ul.appendChild(fn(c, i)));
    };
    const makeImgLi        = (card, hidden = false, i = null) => {
        const li  = document.createElement('li');
        li.dataset.index = i;
        li.dataset.card  = cardToString(card);
        const img = document.createElement('img');
        img.className    = 'card-image fade-in';
        if (hidden) {
            img.src = './static/images/Close.jpg';
            img.alt = 'Closed';
        } else {
            const [rank, suit] = cardToString(card).split(' of ');
            img.src = `./static/images/${rank} of ${suit}.png`;
            img.alt = cardToString(card);
        }
        li.appendChild(img);
        return li;
    };

    const api = async (url, opt = {}) => {
        const res  = await fetch(url, opt);
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        if (data.error) throw new Error(data.error);
        return data;
    };

    function updateGamingTableUI(cards) {
        renderList(gamingTableList, cards, c => makeImgLi(c, false));
    }

    function updateTableUI(cards) {
        renderList(tableList, cards, c => {
            const li = document.createElement('li');
            li.textContent = cardToString(c);
            li.className   = 'fade-in';
            return li;
        });
    }

    function updateHandUI(cards) {
        if (!Array.isArray(cards)) return;
        renderList(handList, cards, (c, i) => {
            const li = makeImgLi(c, false, i);
            li.addEventListener('click', () => li.classList.toggle('selected'));
            return li;
        });
    }

    function updateAIDeckUI(deck) {
        const cards = Array.isArray(deck)
            ? deck
            : deck?.Currents?.Cards ?? [];
        renderList(aiHandList, cards, (c, i) => makeImgLi(c, !aiCardsRevealed, i));
    }

    const selHandCards = () => Array.from(
        handList.querySelectorAll('li.selected')
    ).map(li => li.dataset.card);

    const allHandCards = () => Array.from(
        handList.querySelectorAll('li')
    ).map(li => li.dataset.card);

    const allTableCards = () => Array.from(
        tableList.querySelectorAll('li')
    ).map(li => li.textContent);

    const allAICards = () => Array.from(
        aiHandList.querySelectorAll('li')
    ).map(li => li.dataset.card);

    async function checkStatus() {
        try {
            const payload = {
                hand:   allHandCards(),
                deckAI: allAICards(),
                table:  allTableCards()
            };
            const data = await api(
                '/api/deck/status',
                {
                    method: 'POST',
                    headers: {
                        'Content-Type':   'application/json',
                        'X-Session-ID':   sessionId
                    },
                    body: JSON.stringify(payload)
                }
            );
            if (data.handCards)  updateHandUI(data.handCards);
            if (data.aiCards) {
                cachedAIDeck = data.aiCards;
                updateAIDeckUI(cachedAIDeck);
            }
            if (data.tableCards) updateTableUI(data.tableCards);
        } catch (e) {
            console.error(e);
        }
    }

    toggleAICardsBtn.addEventListener('click', () => {
        aiCardsRevealed = !aiCardsRevealed;
        toggleAICardsBtn.textContent = aiCardsRevealed
            ? 'Скрыть карты ИИ'
            : 'Показать карты ИИ';
        if (cachedAIDeck) updateAIDeckUI(cachedAIDeck);
    });

    createDeckBtn.addEventListener('click', async () => {
        try {
            const d = await api(
                '/api/deck/create',
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-Session-ID': sessionId
                    }
                }
            );
            gamingTableCards = [];
            updateGamingTableUI(gamingTableCards);
            updateHandUI(d.deck || []);
            updateAIDeckUI(d.deckAI || []);
            cachedAIDeck = d.deckAI;
            updateTableUI(d.queue || []);
        } catch (e) {
            alert(e.message);
        }
    });

    saveDeckBtn.addEventListener('click', async () => {
        const payload = {
            handCards:  allHandCards(),
            aiCards:    allAICards(),
            tableCards: allTableCards()
        };
        try {
            const { downloadUrl } = await api(
                '/api/deck/save',
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-Session-ID': sessionId
                    },
                    body: JSON.stringify(payload)
                }
            );
            const a = document.createElement('a');
            a.href    = downloadUrl;
            a.download = downloadUrl.split('/').pop();
            document.body.appendChild(a);
            a.click();
            a.remove();
        } catch (e) {
            alert(e.message);
        }
    });


    loadDeckInput.addEventListener('change', async e => {
        const file = e.target.files[0];
        if (!file) return;

        isLoadingFromSave = true;

        try {
            const d = await api(
                `/api/deck/load?filename=${encodeURIComponent(file.name)}`
            );

            gamingTableCards = [];
            updateTableUI(d.tableCards  || []);
            updateHandUI (d.handCards    || []);
            updateAIDeckUI(d.aiCards     || []);

            document.querySelectorAll('button')
                .forEach(b => b.disabled = false);

            alert('♻️ Игра загружена! Продолжаем с позиции сохранения.');

            setTimeout(async () => {
                await checkStatus();
                isLoadingFromSave = false;
            }, 200);

        } catch (err) {
            alert(err.message);
            isLoadingFromSave = false;
        }
    });

    drawCardBtn.addEventListener('click', async () => {
        if (inFlight) return;
        const selected = selHandCards();
        if (!selected.length) {
            alert('Выбери карту, шустрик!');
            return;
        }
        const ranks = [...new Set(selected.map(c => c.split(' of ')[0]))];
        if (ranks.length !== 1) {
            alert('Все выброшенные карты должны быть одинакового ранга!');
            return;
        }
        inFlight = true;
        gamingTableCards.push(...selected);
        updateGamingTableUI(gamingTableCards);

        try {
            const leaveData = await api(
                '/api/deck/leave',
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-Session-ID': sessionId
                    },
                    body: JSON.stringify({
                        hand:     allHandCards(),
                        table:    allTableCards(),
                        selected
                    })
                }
            );
            updateHandUI(leaveData.deck || []);
            updateTableUI(leaveData.queue || []);

            const aiData = await api(
                '/api/deck/ai',
                {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                        'X-Session-ID': sessionId
                    },
                    body: JSON.stringify({
                        selected,
                        handAI: allAICards()
                    })
                }
            );
            if (aiData.data) {
                if (aiData.data.handAI) {
                    cachedAIDeck = aiData.data.handAI;
                    updateAIDeckUI(cachedAIDeck);
                }
                if (aiData.data.selectedAI?.length) {
                    gamingTableCards.push(...aiData.data.selectedAI);
                    updateGamingTableUI(gamingTableCards);
                }
                if (aiData.data.queue) updateTableUI(aiData.data.queue);
            }
            if (gamingTableCards.length >= 2) {
                await sleep(4000);
                gamingTableCards = [];
                updateGamingTableUI(gamingTableCards);
            }
            await checkStatus();
        } catch (err) {
            alert(err.message);
        }
        inFlight = false;
    });

    setInterval(checkStatus, 8000);

    function initWebSocket() {
        const proto = location.protocol === 'https:' ? 'wss' : 'ws';
        socket = new WebSocket(`${proto}://${location.host}/api/ws/${sessionId}`);
        socket.addEventListener('message', e => {
            let data;
            try { data = JSON.parse(e.data); }
            catch { return; }

            if (data.type === 'gameOver') {
                if (isLoadingFromSave) return;
                alert(
                    data.winner === 'player'
                        ? '🎉 Вы выиграли!'
                        : '😞 Вы проиграли'
                );
                document.querySelectorAll('button')
                    .forEach(b => b.disabled = true);
                return;
            }

            if (data.queue)   updateTableUI(data.queue);
            if (data.deck)    updateHandUI(data.deck);
            if (data.deckAI) {
                cachedAIDeck = data.deckAI;
                updateAIDeckUI(data.deckAI);
            }
        });

    }

    initWebSocket();
})();