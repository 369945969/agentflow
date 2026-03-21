class GobangGame {
    constructor() {
        this.boardSize = 15;
        this.board = [];
        this.currentPlayer = 'black';
        this.gameOver = false;
        this.moveHistory = [];
        this.moveCount = 0;
        
        this.initElements();
        this.initBoard();
        this.bindEvents();
        this.updateGameStatus();
    }

    initElements() {
        this.boardElement = document.getElementById('board');
        this.statusMessage = document.getElementById('statusMessage');
        this.moveCountElement = document.getElementById('moveCount');
        this.winOverlay = document.getElementById('winOverlay');
        this.winnerText = document.getElementById('winnerText');
        this.restartBtn = document.getElementById('restartBtn');
        this.undoBtn = document.getElementById('undoBtn');
        this.hintBtn = document.getElementById('hintBtn');
        this.playAgainBtn = document.getElementById('playAgainBtn');
        this.rulesBtn = document.getElementById('rulesBtn');
        this.rulesModal = document.getElementById('rulesModal');
        this.closeRulesBtn = document.getElementById('closeRulesBtn');
        this.blackPlayerIndicator = document.querySelector('.black-player');
        this.whitePlayerIndicator = document.querySelector('.white-player');
    }

    initBoard() {
        this.boardElement.innerHTML = '';
        this.board = Array(this.boardSize).fill().map(() => Array(this.boardSize).fill(null));
        
        for (let row = 0; row < this.boardSize; row++) {
            for (let col = 0; col < this.boardSize; col++) {
                const cell = document.createElement('div');
                cell.className = 'cell';
                cell.dataset.row = row;
                cell.dataset.col = col;
                
                cell.addEventListener('click', () => this.handleCellClick(row, col));
                this.boardElement.appendChild(cell);
            }
        }
    }

    bindEvents() {
        this.restartBtn.addEventListener('click', () => this.restartGame());
        this.undoBtn.addEventListener('click', () => this.undoMove());
        this.hintBtn.addEventListener('click', () => this.showHint());
        this.playAgainBtn.addEventListener('click', () => this.restartGame());
        this.rulesBtn.addEventListener('click', () => this.showRules());
        this.closeRulesBtn.addEventListener('click', () => this.hideRules());
        
        this.rulesModal.addEventListener('click', (e) => {
            if (e.target === this.rulesModal) {
                this.hideRules();
            }
        });
    }

    handleCellClick(row, col) {
        if (this.gameOver || this.board[row][col]) return;
        
        this.makeMove(row, col);
    }

    makeMove(row, col) {
        this.board[row][col] = this.currentPlayer;
        this.moveHistory.push({ row, col, player: this.currentPlayer });
        this.moveCount++;
        
        const cell = this.boardElement.querySelector(`[data-row="${row}"][data-col="${col}"]`);
        cell.classList.add(this.currentPlayer);
        
        if (this.checkWin(row, col)) {
            this.handleWin();
            return;
        }
        
        if (this.moveCount === this.boardSize * this.boardSize) {
            this.handleDraw();
            return;
        }
        
        this.switchPlayer();
        this.updateGameStatus();
    }

    switchPlayer() {
        this.currentPlayer = this.currentPlayer === 'black' ? 'white' : 'black';
        
        this.blackPlayerIndicator.classList.toggle('active', this.currentPlayer === 'black');
        this.whitePlayerIndicator.classList.toggle('active', this.currentPlayer === 'white');
    }

    checkWin(row, col) {
        const player = this.board[row][col];
        const directions = [
            [0, 1],   // 水平
            [1, 0],   // 垂直
            [1, 1],   // 右下斜线
            [1, -1]   // 左下斜线
        ];

        for (const [dx, dy] of directions) {
            let count = 1;
            let winningCells = [{ row, col }];

            for (let i = 1; i < 5; i++) {
                const newRow = row + dx * i;
                const newCol = col + dy * i;
                
                if (this.isValidPosition(newRow, newCol) && this.board[newRow][newCol] === player) {
                    count++;
                    winningCells.push({ row: newRow, col: newCol });
                } else {
                    break;
                }
            }

            for (let i = 1; i < 5; i++) {
                const newRow = row - dx * i;
                const newCol = col - dy * i;
                
                if (this.isValidPosition(newRow, newCol) && this.board[newRow][newCol] === player) {
                    count++;
                    winningCells.push({ row: newRow, col: newCol });
                } else {
                    break;
                }
            }

            if (count >= 5) {
                this.highlightWinningCells(winningCells);
                return true;
            }
        }

        return false;
    }

    isValidPosition(row, col) {
        return row >= 0 && row < this.boardSize && col >= 0 && col < this.boardSize;
    }

    highlightWinningCells(cells) {
        cells.forEach(({ row, col }) => {
            const cell = this.boardElement.querySelector(`[data-row="${row}"][data-col="${col}"]`);
            cell.classList.add('winning');
        });
    }

    handleWin() {
        this.gameOver = true;
        const winner = this.currentPlayer === 'black' ? '黑方' : '白方';
        this.winnerText.textContent = `${winner} 获胜！`;
        this.winOverlay.classList.add('show');
        this.statusMessage.textContent = `${winner} 获得胜利！`;
    }

    handleDraw() {
        this.gameOver = true;
        this.winnerText.textContent = '平局！';
        this.winOverlay.classList.add('show');
        this.statusMessage.textContent = '棋盘已满，平局！';
    }

    undoMove() {
        if (this.moveHistory.length === 0 || this.gameOver) return;
        
        const lastMove = this.moveHistory.pop();
        this.board[lastMove.row][lastMove.col] = null;
        this.moveCount--;
        
        const cell = this.boardElement.querySelector(
            `[data-row="${lastMove.row}"][data-col="${lastMove.col}"]`
        );
        cell.classList.remove('black', 'white', 'winning');
        
        this.currentPlayer = lastMove.player;
        this.gameOver = false;
        this.winOverlay.classList.remove('show');
        
        this.updatePlayerIndicators();
        this.updateGameStatus();
    }

    showHint() {
        if (this.gameOver) return;
        
        const emptyCells = [];
        for (let row = 0; row < this.boardSize; row++) {
            for (let col = 0; col < this.boardSize; col++) {
                if (!this.board[row][col]) {
                    emptyCells.push({ row, col });
                }
            }
        }
        
        if (emptyCells.length === 0) return;
        
        const randomCell = emptyCells[Math.floor(Math.random() * emptyCells.length)];
        const cell = this.boardElement.querySelector(
            `[data-row="${randomCell.row}"][data-col="${randomCell.col}"]`
        );
        
        cell.style.boxShadow = '0 0 0 3px gold';
        setTimeout(() => {
            cell.style.boxShadow = '';
        }, 1000);
        
        this.statusMessage.textContent = `提示：建议在 (${randomCell.row + 1}, ${randomCell.col + 1}) 落子`;
        setTimeout(() => this.updateGameStatus(), 2000);
    }

    restartGame() {
        this.currentPlayer = 'black';
        this.gameOver = false;
        this.moveHistory = [];
        this.moveCount = 0;
        
        this.initBoard();
        this.winOverlay.classList.remove('show');
        this.updatePlayerIndicators();
        this.updateGameStatus();
    }

    updatePlayerIndicators() {
        this.blackPlayerIndicator.classList.toggle('active', this.currentPlayer === 'black');
        this.whitePlayerIndicator.classList.toggle('active', this.currentPlayer === 'white');
    }

    updateGameStatus() {
        const playerName = this.currentPlayer === 'black' ? '黑方' : '白方';
        this.statusMessage.textContent = this.gameOver ? '游戏结束' : `${playerName} 回合`;
        this.moveCountElement.textContent = this.moveCount;
    }

    showRules() {
        this.rulesModal.classList.add('show');
    }

    hideRules() {
        this.rulesModal.classList.remove('show');
    }
}

document.addEventListener('DOMContentLoaded', () => {
    new GobangGame();
});