import { describe, it, expect, vi } from 'vitest';
import { mockBooks, mockMangas, mockGames, mockFilmserien } from '../mocks/handlers';

const API_BASE_URL = 'https://diplodocu.mpech.dev/api';
const mockToken = 'mock-jwt-token';

const fetchWithAuth = async (url: string, options: RequestInit = {}) => {
  return fetch(url, {
    ...options,
    headers: {
      'Authorization': `Bearer ${mockToken}`,
      'Content-Type': 'application/json',
      ...options.headers,
    },
  });
};

describe('API Integration Tests', () => {
  describe('Books API', () => {
    it('fetches all books', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/books`);
      const data = await response.json();

      expect(response.ok).toBe(true);
      expect(data).toHaveLength(mockBooks.length);
      expect(data[0]).toHaveProperty('name', 'Harry Potter');
    });

    it('creates a new book', async () => {
      const newBook = {
        name: 'New Book',
        nummer: 1,
        autor: 'New Author',
        sprache: 'German',
        genre: 'Mystery',
      };

      const response = await fetchWithAuth(`${API_BASE_URL}/books`, {
        method: 'POST',
        body: JSON.stringify(newBook),
      });
      const data = await response.json();

      expect(response.status).toBe(201);
      expect(data).toHaveProperty('name', 'New Book');
      expect(data).toHaveProperty('id');
    });

    it('fetches a single book by ID', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/books/1`);
      const data = await response.json();

      expect(response.ok).toBe(true);
      expect(data).toHaveProperty('name', 'Harry Potter');
      expect(data).toHaveProperty('autor', 'J.K. Rowling');
    });

    it('returns 404 for non-existent book', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/books/999`);

      expect(response.status).toBe(404);
    });

    it('updates an existing book', async () => {
      const updatedData = {
        name: 'Harry Potter Updated',
        nummer: 2,
      };

      const response = await fetchWithAuth(`${API_BASE_URL}/books/1`, {
        method: 'PUT',
        body: JSON.stringify(updatedData),
      });
      const data = await response.json();

      expect(response.ok).toBe(true);
      expect(data).toHaveProperty('name', 'Harry Potter Updated');
    });

    it('deletes a book', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/books/1`, {
        method: 'DELETE',
      });

      expect(response.status).toBe(204);
    });
  });

  describe('Mangas API', () => {
    it('fetches all mangas', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/mangas`);
      const data = await response.json();

      expect(response.ok).toBe(true);
      expect(data).toHaveLength(mockMangas.length);
      expect(data[0]).toHaveProperty('mangaka');
    });

    it('creates a new manga', async () => {
      const newManga = {
        name: 'New Manga',
        nummer: 1,
        mangaka: 'New Mangaka',
        sprache: 'Japanese',
        genre: 'Seinen',
      };

      const response = await fetchWithAuth(`${API_BASE_URL}/mangas`, {
        method: 'POST',
        body: JSON.stringify(newManga),
      });
      const data = await response.json();

      expect(response.status).toBe(201);
      expect(data).toHaveProperty('name', 'New Manga');
    });

    it('deletes a manga', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/mangas/3`, {
        method: 'DELETE',
      });

      expect(response.status).toBe(204);
    });
  });

  describe('Games (Spiel) API', () => {
    it('fetches all games', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/spiel`);
      const data = await response.json();

      expect(response.ok).toBe(true);
      expect(data).toHaveLength(mockGames.length);
      expect(data[0]).toHaveProperty('konsole');
    });

    it('creates a new game', async () => {
      const newGame = {
        name: 'New Game',
        konsole: 'PlayStation 5',
        genre: 'Action',
      };

      const response = await fetchWithAuth(`${API_BASE_URL}/spiel`, {
        method: 'POST',
        body: JSON.stringify(newGame),
      });
      const data = await response.json();

      expect(response.status).toBe(201);
      expect(data).toHaveProperty('name', 'New Game');
    });

    it('deletes a game', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/spiel/5`, {
        method: 'DELETE',
      });

      expect(response.status).toBe(204);
    });
  });

  describe('Film/Serie API', () => {
    it('fetches all film/series', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/filmserie`);
      const data = await response.json();

      expect(response.ok).toBe(true);
      expect(data).toHaveLength(mockFilmserien.length);
      expect(data[0]).toHaveProperty('art');
    });

    it('creates a new film', async () => {
      const newFilm = {
        name: 'New Film',
        art: 'Film',
        genre: 'Action',
      };

      const response = await fetchWithAuth(`${API_BASE_URL}/filmserie`, {
        method: 'POST',
        body: JSON.stringify(newFilm),
      });
      const data = await response.json();

      expect(response.status).toBe(201);
      expect(data).toHaveProperty('name', 'New Film');
    });

    it('creates a new serie', async () => {
      const newSerie = {
        name: 'New Serie',
        nummer: 3,
        art: 'Serie',
        genre: 'Comedy',
      };

      const response = await fetchWithAuth(`${API_BASE_URL}/filmserie`, {
        method: 'POST',
        body: JSON.stringify(newSerie),
      });
      const data = await response.json();

      expect(response.status).toBe(201);
      expect(data).toHaveProperty('art', 'Serie');
    });

    it('deletes a film/serie', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/filmserie/7`, {
        method: 'DELETE',
      });

      expect(response.status).toBe(204);
    });
  });

  describe('Collections API', () => {
    it('fetches user collections', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/sammlungen`);
      const data = await response.json();

      expect(response.ok).toBe(true);
      expect(Array.isArray(data)).toBe(true);
    });

    it('creates a new collection', async () => {
      const newCollection = {
        name: 'My New Collection',
      };

      const response = await fetchWithAuth(`${API_BASE_URL}/sammlungen`, {
        method: 'POST',
        body: JSON.stringify(newCollection),
      });
      const data = await response.json();

      expect(response.status).toBe(201);
      expect(data).toHaveProperty('name', 'My New Collection');
      expect(data).toHaveProperty('webuser_id');
    });
  });

  describe('User Sync API', () => {
    it('syncs user data', async () => {
      const response = await fetchWithAuth(`${API_BASE_URL}/sync-user`);
      const data = await response.json();

      expect(response.ok).toBe(true);
      expect(data).toHaveProperty('id');
      expect(data).toHaveProperty('name');
    });
  });
});
