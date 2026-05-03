const SEAT_PRICE = 12;
const HOLD_TIME = 30;

let selectedByYou = [];
let confirmedByYou = [];
let soldByOther = [];
let reservedByOther = [];

let timerInterval = null;
let timeLeft = HOLD_TIME;

const seatGrid = document.getElementById("seat-grid");
const countEl = document.getElementById("count");
const totalEl = document.getElementById("total");
const timerContainer = document.getElementById("timer-container");
const timerEl = document.getElementById("timer");
const bookBtn = document.getElementById("book-btn");

let selectedMovieId = null;
let currentSeatsData = [];
const USER_ID = Math.random().toString(36).substring(2, 15);

async function fetchBookings(movieId) {
    if (!movieId) return;
    try {
        const res = await fetch(`/api/movies/${movieId}/bookings`);
        if (!res.ok) throw new Error();
        const bookings = await res.json();
        
        soldByOther = [];
        confirmedByYou = [];
        
        if (Array.isArray(bookings)) {
            bookings.forEach(booking => {
                if (booking.user_id === USER_ID) {
                    confirmedByYou.push(booking.seat_id);
                } else {
                    soldByOther.push(booking.seat_id);
                }
            });
        }
    } catch (err) {
        soldByOther = [];
        confirmedByYou = [];
    }
    renderSeats();
}

async function init() {
    await renderMovies();
    updateSummary();
}

async function renderMovies() {
    const cont = document.querySelector(".movie-options");
    const movieHeader = document.querySelector(".movie-info h1");
    const genreEl = document.querySelector(".genre");
    const metaEl = document.querySelector(".meta");
    
    let movies = [];
    
    try {
        const res = await fetch("/api/movies");
        if (!res.ok) throw new Error();
        movies = await res.json();
    } catch (err) {
        return;
    }

    cont.innerHTML = "";

    let firstSelection = null;

    for (let i = 0; i < movies.length; i++) {
        const movie = movies[i];
        const movieTag = document.createElement("div");
        movieTag.className = "movie-container";
        movieTag.style.animation = `fadeIn 0.5s ease-out ${i * 0.1}s backwards`;
        movieTag.innerText = movie.title;

        const selectMovie = async () => {
            document.querySelectorAll(".movie-container").forEach(el => el.classList.remove("selected"));
            movieTag.classList.add("selected");
            
            selectedMovieId = movie.id;
            currentSeatsData = movie.seats || [];

            movieHeader.innerText = movie.title;
            genreEl.innerText = movie.genre;
            metaEl.innerHTML = `
                <span>${movie.duration}</span>
                <span>IMDb ${movie.rating}</span>
            `;

            selectedByYou = [];
            await fetchBookings(selectedMovieId);
            updateSummary();
            stopTimer();
        };

        movieTag.addEventListener("click", selectMovie);
        cont.appendChild(movieTag);

        if (i === 0 || movieHeader.innerText === movie.title) {
            firstSelection = selectMovie;
        }
    }

    if (firstSelection) await firstSelection();
}

function renderSeats() {
    seatGrid.innerHTML = "";
    if (!currentSeatsData || currentSeatsData.length === 0) return;

    currentSeatsData.forEach(rowData => {
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

async function toggleSeat(seatNum) {
    if (selectedByYou.includes(seatNum)) {
        selectedByYou = selectedByYou.filter(s => s !== seatNum);
    } else {
        if (reservedByOther.includes(seatNum) || soldByOther.includes(seatNum) || confirmedByYou.includes(seatNum)) return;
        
        if (selectedMovieId) {
            try {
                const res = await fetch(`/api/movies/${selectedMovieId}/seats/${seatNum}/hold`, {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({ user_id: USER_ID })
                });
                if (!res.ok) return;
            } catch (err) {
                return;
            }
        }
        
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

bookBtn.addEventListener("click", () => {
    showModal("Booking Confirmed!", `Successfully booked: ${selectedByYou.join(", ")}. Enjoy your movie!`);
    selectedByYou.forEach(s => confirmedByYou.push(s));
    selectedByYou = [];
    stopTimer();
    renderSeats();
    updateSummary();
});

init();
