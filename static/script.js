const SEAT_PRICE = 12;
const HOLD_TIME = 30;

let selectedByYou = [];
let confirmedByYou = [];
let soldByOther = ["F2", "D2"];
let reservedByOther = ["C2"];

let timerInterval = null;
let timeLeft = HOLD_TIME;

const seatGrid = document.getElementById("seat-grid");
const countEl = document.getElementById("count");
const totalEl = document.getElementById("total");
const timerContainer = document.getElementById("timer-container");
const timerEl = document.getElementById("timer");
const bookBtn = document.getElementById("book-btn");

const seatsData = [
    ["A1", "A2", "A3", "A4", "A5", "A6"],
    ["B1", "B2", "B3", "B4", "B5", "B6"],
    ["C1", "C2", "C3", "C4", "C5", "C6"],
    ["spacer"],
    ["D1", "D2", "D3", "D4", "D5", "D6"],
    ["E1", "E2", "E3", "E4", "E5", "E6", "E7"],
    ["F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8"],
];

function init() {
    renderSeats();
    updateSummary();
    startSimulation();
}

function renderSeats() {
    seatGrid.innerHTML = "";
    seatsData.forEach(rowData => {
        if (rowData[0] === "spacer") {
            const spacer = document.createElement("div");
            spacer.style.height = "20px";
            seatGrid.appendChild(spacer);
            return;
        }

        const row = document.createElement("div");
        row.classList.add("rows");

        rowData.forEach(seatNum => {
            const seat = document.createElement("div");
            seat.classList.add("seat");
            seat.textContent = seatNum;
            seat.dataset.num = seatNum;

            if (confirmedByYou.includes(seatNum)) {
                seat.classList.add("confirmed");
            } else if (soldByOther.includes(seatNum)) {
                seat.classList.add("sold");
            } else if (selectedByYou.includes(seatNum)) {
                seat.classList.add("selected");
                seat.addEventListener("click", () => toggleSeat(seatNum));
            } else if (reservedByOther.includes(seatNum)) {
                seat.classList.add("reserved");
            } else {
                seat.classList.add("normal");
                seat.addEventListener("click", () => toggleSeat(seatNum));
            }

            row.appendChild(seat);
        });
        seatGrid.appendChild(row);
    });
}

function toggleSeat(seatNum) {
    if (selectedByYou.includes(seatNum)) {
        selectedByYou = selectedByYou.filter(s => s !== seatNum);
    } else {
        // Can't select if someone else is holding or it's sold
        if (reservedByOther.includes(seatNum) || soldByOther.includes(seatNum) || confirmedByYou.includes(seatNum)) return;
        selectedByYou.push(seatNum);
    }

    if (selectedByYou.length > 0) {
        startTimer();
    } else {
        stopTimer();
    }

    renderSeats();
    updateSummary();
}

function updateSummary() {
    countEl.textContent = selectedByYou.length;
    totalEl.textContent = `$${selectedByYou.length * SEAT_PRICE}`;
    bookBtn.disabled = selectedByYou.length === 0;
}

function startTimer() {
    if (timerInterval) return;
    
    timeLeft = HOLD_TIME;
    timerContainer.style.display = "block";
    updateTimerDisplay();

    timerInterval = setInterval(() => {
        timeLeft--;
        updateTimerDisplay();

        if (timeLeft <= 0) {
            expireSelection();
        }
    }, 1000);
}

function stopTimer() {
    clearInterval(timerInterval);
    timerInterval = null;
    timerContainer.style.display = "none";
}

function updateTimerDisplay() {
    const minutes = Math.floor(timeLeft / 60);
    const seconds = timeLeft % 60;
    timerEl.textContent = `${String(minutes).padStart(2, '0')}:${String(seconds).padStart(2, '0')}`;
}

const modalOverlay = document.getElementById("modal-overlay");
const modalTitle = document.getElementById("modal-title");
const modalMessage = document.getElementById("modal-message");
const modalClose = document.getElementById("modal-close");

function showModal(title, message) {
    modalTitle.textContent = title;
    modalMessage.textContent = message;
    modalOverlay.classList.add("active");
}

function hideModal() {
    modalOverlay.classList.remove("active");
}

modalClose.addEventListener("click", hideModal);

function expireSelection() {
    stopTimer();
    selectedByYou = [];
    renderSeats();
    updateSummary();
    showModal("Session Expired", "Your seat hold has expired. Please select seats again.");
}

// Simulation Logic
function startSimulation() {
    setInterval(() => {
        const allSeats = seatsData.flat().filter(s => s !== "spacer");
        const availableForOther = allSeats.filter(s => 
            !soldByOther.includes(s) && 
            !confirmedByYou.includes(s) && 
            !selectedByYou.includes(s) && 
            !reservedByOther.includes(s)
        );

        if (Math.random() > 0.7 && availableForOther.length > 0) {
            const randomSeat = availableForOther[Math.floor(Math.random() * availableForOther.length)];
            reservedByOther.push(randomSeat);
            renderSeats();
            
            setTimeout(() => {
                reservedByOther = reservedByOther.filter(s => s !== randomSeat);
                if (Math.random() > 0.8) soldByOther.push(randomSeat);
                renderSeats();
            }, 1000 + Math.random() * 1000);
        }
    }, 1000);
}

bookBtn.addEventListener("click", () => {
    showModal("Booking Confirmed!", `Successfully booked: ${selectedByYou.join(", ")}. Enjoy your movie!`);
    selectedByYou.forEach(s => confirmedByYou.push(s));
    selectedByYou = [];
    stopTimer();
    renderSeats();
    updateSummary();
});

init();
