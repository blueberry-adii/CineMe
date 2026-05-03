const SEAT_PRICE = 12;

let selectedByYou = [];
let confirmedByYou = [];
let soldByOther = [];
let reservedByOther = [];

let timerInterval = null;
let pollingInterval = null;
let sessionExpiries = {}; 
let sessionIds = {}; // Map seatId to sessionId
let lastActionTime = 0;
let pendingSeats = new Set(); 

const seatGrid = document.getElementById("seat-grid");
const countEl = document.getElementById("count");
const totalEl = document.getElementById("total");
const bookBtn = document.getElementById("book-btn");

let selectedMovieId = null;
let currentSeatsData = [];
const USER_ID = Math.random().toString(36).substring(2, 15);

async function fetchBookings(movieId, isInitial = false) {
    if (!movieId) return;
    try {
        const res = await fetch(`/api/movies/${movieId}/bookings`);
        if (!res.ok) throw new Error();
        const bookings = await res.json();
        
        const newReserved = [];
        const newSold = [];
        const newSelected = [];
        const newExpiries = { ...sessionExpiries };
        const newSessionIds = { ...sessionIds };
        
        if (Array.isArray(bookings)) {
            bookings.forEach(booking => {
                if (booking.user_id === USER_ID) {
                    newSelected.push(booking.seat_id);
                    if (booking.expires_at) {
                        newExpiries[booking.seat_id] = new Date(booking.expires_at).getTime();
                    }
                    if (booking.session_id) {
                        newSessionIds[booking.seat_id] = booking.session_id;
                    }
                } else {
                    if (booking.status === "confirmed") {
                        newSold.push(booking.seat_id);
                    } else {
                        newReserved.push(booking.seat_id);
                    }
                }
            });
        }

        reservedByOther = newReserved;
        soldByOther = newSold;
        
        const timeSinceAction = Date.now() - lastActionTime;
        if (isInitial || timeSinceAction > 5000) {
            const mergedSelected = [...new Set([...newSelected, ...Array.from(pendingSeats)])];
            selectedByYou = mergedSelected;
            sessionExpiries = newExpiries;
            sessionIds = newSessionIds;
        }

        renderSeats();
        updateSummary();
        
        if (selectedByYou.length > 0) {
            startTimer();
        } else {
            stopTimer();
        }
    } catch (err) {
        // Silently fail during polling
    }
}

function startPolling() {
    if (pollingInterval) clearInterval(pollingInterval);
    pollingInterval = setInterval(() => {
        if (selectedMovieId) fetchBookings(selectedMovieId);
    }, 1000);
}

async function init() {
    await renderMovies();
    startPolling();
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
            confirmedByYou = []; // Reset confirmed seats when switching movies
            reservedByOther = [];
            soldByOther = []; // Also reset sold seats
            sessionExpiries = {};
            sessionIds = {};
            lastActionTime = 0; 
            await fetchBookings(selectedMovieId, true);
            updateSummary();
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
                
                const timerSpan = document.createElement("span");
                timerSpan.className = "seat-timer";
                timerSpan.id = `timer-${seatNum}`;
                
                const expiry = sessionExpiries[seatNum];
                if (expiry) {
                    const distance = expiry - Date.now();
                    if (distance > 0) {
                        timerSpan.textContent = `${Math.ceil(distance / 1000)}s`;
                    }
                } else {
                    timerSpan.textContent = "10s";
                }
                seat.appendChild(timerSpan);
            } else if (reservedByOther.includes(seatNum)) {
                seat.classList.add("reserved");
            } else {
                seat.classList.add("normal");
            }

            row.appendChild(seat);
        });
        seatGrid.appendChild(row);
    });
}

seatGrid.addEventListener("click", (e) => {
    const seat = e.target.closest(".seat");
    if (!seat || !seat.dataset.num) return;
    toggleSeat(seat.dataset.num);
});

async function toggleSeat(seatNum) {
    if (reservedByOther.includes(seatNum) || soldByOther.includes(seatNum) || confirmedByYou.includes(seatNum)) return;

    lastActionTime = Date.now(); 

    if (selectedByYou.includes(seatNum)) {
        const sid = sessionIds[seatNum];
        selectedByYou = selectedByYou.filter(s => s !== seatNum);
        delete sessionExpiries[seatNum];
        delete sessionIds[seatNum];
        pendingSeats.delete(seatNum);

        if (sid) {
            fetch(`/api/sessions/${sid}`, {
                method: 'DELETE',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ user_id: USER_ID })
            });
        }
    } else {
        if (selectedByYou.includes(seatNum)) return;
        
        selectedByYou.push(seatNum);
        pendingSeats.add(seatNum);
        sessionExpiries[seatNum] = Date.now() + 10000;
        
        renderSeats();
        updateSummary();
        startTimer();

        if (selectedMovieId) {
            try {
                const res = await fetch(`/api/movies/${selectedMovieId}/seats/${seatNum}/hold`, {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ user_id: USER_ID })
                });
                
                if (!res.ok) {
                    selectedByYou = selectedByYou.filter(s => s !== seatNum);
                    pendingSeats.delete(seatNum);
                    delete sessionExpiries[seatNum];
                } else {
                    const booking = await res.json();
                    if (booking && booking.expires_at) {
                        sessionExpiries[seatNum] = new Date(booking.expires_at).getTime();
                    }
                    if (booking && booking.session_id) {
                        sessionIds[seatNum] = booking.session_id;
                    }
                    pendingSeats.delete(seatNum);
                }
            } catch (err) {
                selectedByYou = selectedByYou.filter(s => s !== seatNum);
                pendingSeats.delete(seatNum);
                delete sessionExpiries[seatNum];
            }
        }
    }

    renderSeats();
    updateSummary();
    if (selectedByYou.length > 0) startTimer(); else stopTimer();
}

function updateSummary() {
    countEl.textContent = selectedByYou.length;
    totalEl.textContent = `$${selectedByYou.length * SEAT_PRICE}`;
    bookBtn.disabled = selectedByYou.length === 0;
}

function startTimer() {
    if (timerInterval) clearInterval(timerInterval);
    
    const update = () => {
        const now = Date.now();
        let anyExpired = false;

        selectedByYou.forEach(seatId => {
            const expiry = sessionExpiries[seatId];
            const el = document.getElementById(`timer-${seatId}`);
            
            if (expiry) {
                const distance = expiry - now;
                if (distance <= 0) {
                    expireSeat(seatId);
                    anyExpired = true;
                } else {
                    if (el) {
                        const s = Math.ceil(distance / 1000);
                        el.textContent = `${s}s`;
                    }
                }
            }
        });

        if (anyExpired) {
            renderSeats();
            updateSummary();
        }
        
        if (selectedByYou.length === 0) stopTimer();
    };

    update();
    timerInterval = setInterval(update, 500);
}

function expireSeat(seatId) {
    selectedByYou = selectedByYou.filter(s => s !== seatId);
    delete sessionExpiries[seatId];
    delete sessionIds[seatId];
}

function stopTimer() {
    if (timerInterval) clearInterval(timerInterval);
    timerInterval = null;
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

bookBtn.addEventListener("click", async () => {
    const seatsToBook = [...selectedByYou];
    
    for (const seatId of seatsToBook) {
        const sid = sessionIds[seatId];
        if (sid) {
            try {
                await fetch(`/api/sessions/${sid}/confirm`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ user_id: USER_ID })
                });
                confirmedByYou.push(seatId);
            } catch (err) {
                console.error("Failed to confirm seat", seatId, err);
            }
        }
    }
    
    showModal("Booking Confirmed!", `Successfully booked: ${seatsToBook.join(", ")}. Enjoy your movie!`);
    selectedByYou = [];
    sessionExpiries = {};
    sessionIds = {};
    stopTimer();
    renderSeats();
    updateSummary();
});

init();
