document.addEventListener('DOMContentLoaded', () => {
    const statusEl = document.getElementById('status');
    const quickGameBtn = document.getElementById('quickGameBtn');
    const offlineGameBtn = document.getElementById('offlineGameBtn');
    const gameBoard = document.getElementById('game-board');
    const cells = document.querySelectorAll('.cell');
    
    // WebSocket connection
    const wsProtocol = window.location.protocol === 'https:' ? 'wss://' : 'ws://';
    const wsUrl = wsProtocol + window.location.host + '/ws';
    const socket = new WebSocket(wsUrl);
    
    socket.onopen = () => {
        statusEl.textContent = 'Connected to server';
        console.log('WebSocket connected');
    };
    
    socket.onmessage = (event) => {
        console.log('Message from server:', event.data);
        statusEl.textContent = 'Server says: ' + event.data;
    };
    
    socket.onclose = () => {
        statusEl.textContent = 'Disconnected from server';
    };
    
    socket.onerror = (error) => {
        console.error('WebSocket error:', error);
        statusEl.textContent = 'Connection error';
    };
    
    // Test buttons
    quickGameBtn.addEventListener('click', () => {
        statusEl.textContent = 'Searching for opponent...';
        gameBoard.style.display = 'grid';
        socket.send('Looking for game');
    });
    
    offlineGameBtn.addEventListener('click', () => {
        statusEl.textContent = 'Starting offline game';
        gameBoard.style.display = 'grid';
    });
    
    // Test game board
    cells.forEach(cell => {
        cell.addEventListener('click', () => {
            if (cell.textContent === '') {
                cell.textContent = 'X';
                socket.send(`Move made at ${cell.dataset.index}`);
            }
        });
    });
});