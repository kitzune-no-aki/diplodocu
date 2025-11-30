import { http, HttpResponse } from 'msw';

const API_BASE_URL = 'https://diplodocu.mpech.dev/api';

// Mock data
export const mockBooks = [
  { id: 1, name: 'Harry Potter', nummer: 1, autor: 'J.K. Rowling', sprache: 'English', genre: 'Fantasy' },
  { id: 2, name: 'The Hobbit', nummer: null, autor: 'J.R.R. Tolkien', sprache: 'English', genre: 'Fantasy' },
];

export const mockMangas = [
  { id: 3, name: 'One Piece', nummer: 100, mangaka: 'Eiichiro Oda', sprache: 'Japanese', genre: 'Shonen' },
  { id: 4, name: 'Naruto', nummer: 72, mangaka: 'Masashi Kishimoto', sprache: 'Japanese', genre: 'Shonen' },
];

export const mockGames = [
  { id: 5, name: 'Zelda: Tears of the Kingdom', nummer: null, konsole: 'Nintendo Switch', genre: 'Adventure' },
  { id: 6, name: 'Elden Ring', nummer: null, konsole: 'PC', genre: 'RPG' },
];

export const mockFilmserien = [
  { id: 7, name: 'Breaking Bad', nummer: 5, art: 'Serie', genre: 'Drama' },
  { id: 8, name: 'Inception', nummer: null, art: 'Film', genre: 'Sci-Fi' },
];

export const mockCollections = [
  { id: 1, name: 'My Favorites', webuser_id: 'test-user-123' },
  { id: 2, name: 'To Watch', webuser_id: 'test-user-123' },
];

// API Handlers
export const handlers = [
  // Books
  http.get(`${API_BASE_URL}/books`, () => {
    return HttpResponse.json(mockBooks);
  }),

  http.post(`${API_BASE_URL}/books`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    const newBook = { id: Date.now(), ...body };
    return HttpResponse.json(newBook, { status: 201 });
  }),

  http.get(`${API_BASE_URL}/books/:id`, ({ params }) => {
    const book = mockBooks.find(b => b.id === Number(params.id));
    if (!book) {
      return HttpResponse.json({ error: 'Book not found' }, { status: 404 });
    }
    return HttpResponse.json(book);
  }),

  http.put(`${API_BASE_URL}/books/:id`, async ({ params, request }) => {
    const body = await request.json() as Record<string, unknown>;
    const book = mockBooks.find(b => b.id === Number(params.id));
    if (!book) {
      return HttpResponse.json({ error: 'Book not found' }, { status: 404 });
    }
    return HttpResponse.json({ ...book, ...body });
  }),

  http.delete(`${API_BASE_URL}/books/:id`, ({ params }) => {
    const book = mockBooks.find(b => b.id === Number(params.id));
    if (!book) {
      return HttpResponse.json({ error: 'Book not found' }, { status: 404 });
    }
    return new HttpResponse(null, { status: 204 });
  }),

  // Mangas
  http.get(`${API_BASE_URL}/mangas`, () => {
    return HttpResponse.json(mockMangas);
  }),

  http.post(`${API_BASE_URL}/mangas`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    const newManga = { id: Date.now(), ...body };
    return HttpResponse.json(newManga, { status: 201 });
  }),

  http.delete(`${API_BASE_URL}/mangas/:id`, () => {
    return new HttpResponse(null, { status: 204 });
  }),

  // Games (Spiel)
  http.get(`${API_BASE_URL}/spiel`, () => {
    return HttpResponse.json(mockGames);
  }),

  http.post(`${API_BASE_URL}/spiel`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    const newGame = { id: Date.now(), ...body };
    return HttpResponse.json(newGame, { status: 201 });
  }),

  http.delete(`${API_BASE_URL}/spiel/:id`, () => {
    return new HttpResponse(null, { status: 204 });
  }),

  // Film/Serie
  http.get(`${API_BASE_URL}/filmserie`, () => {
    return HttpResponse.json(mockFilmserien);
  }),

  http.post(`${API_BASE_URL}/filmserie`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    const newFilm = { id: Date.now(), ...body };
    return HttpResponse.json(newFilm, { status: 201 });
  }),

  http.delete(`${API_BASE_URL}/filmserie/:id`, () => {
    return new HttpResponse(null, { status: 204 });
  }),

  // Collections (Sammlungen)
  http.get(`${API_BASE_URL}/sammlungen`, () => {
    return HttpResponse.json(mockCollections);
  }),

  http.post(`${API_BASE_URL}/sammlungen`, async ({ request }) => {
    const body = await request.json() as Record<string, unknown>;
    const newCollection = { id: Date.now(), webuser_id: 'test-user-123', ...body };
    return HttpResponse.json(newCollection, { status: 201 });
  }),

  // User sync
  http.get(`${API_BASE_URL}/sync-user`, () => {
    return HttpResponse.json({ id: 'test-user-123', name: 'Test User' });
  }),
];

// Error handlers for testing error scenarios
export const errorHandlers = {
  networkError: http.get(`${API_BASE_URL}/books`, () => {
    return HttpResponse.error();
  }),

  serverError: http.get(`${API_BASE_URL}/books`, () => {
    return HttpResponse.json({ error: 'Internal Server Error' }, { status: 500 });
  }),

  unauthorized: http.get(`${API_BASE_URL}/books`, () => {
    return HttpResponse.json({ error: 'Unauthorized' }, { status: 401 });
  }),
};
