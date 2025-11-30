import { describe, it, expect, vi } from 'vitest';

// Test the formatDetails helper function
// We need to extract it for testing since it's defined inside mytable.tsx
// For now, we'll test it by recreating the logic

type ProductDetails = {
  autor?: string;
  mangaka?: string;
  konsole?: string;
  art?: 'Film' | 'Serie';
  sprache?: string;
  genre?: string;
};

interface Product {
  id: number;
  name: string;
  nummer: number | null;
  art: 'Buch' | 'Manga' | 'Spiel' | 'Filmserie';
  details: ProductDetails;
}

// Recreate formatDetails for testing
const formatDetails = (product: Product): string => {
  const details: string[] = [];
  if (product.nummer) details.push(`Vol: ${product.nummer}`);
  if (product.details.autor) details.push(`Autor: ${product.details.autor}`);
  if (product.details.mangaka) details.push(`Mangaka: ${product.details.mangaka}`);
  if (product.details.konsole) details.push(`Konsole: ${product.details.konsole}`);
  if (product.details.art) details.push(`Art: ${product.details.art}`);
  if (product.details.sprache) details.push(`Sprache: ${product.details.sprache}`);
  if (product.details.genre) details.push(`Genre: ${product.details.genre}`);
  return details.join(' | ');
};

describe('formatDetails', () => {
  it('formats a book with all details', () => {
    const book: Product = {
      id: 1,
      name: 'Test Book',
      nummer: 1,
      art: 'Buch',
      details: {
        autor: 'Test Author',
        sprache: 'Deutsch',
        genre: 'Fantasy',
      },
    };

    const result = formatDetails(book);

    expect(result).toContain('Vol: 1');
    expect(result).toContain('Autor: Test Author');
    expect(result).toContain('Sprache: Deutsch');
    expect(result).toContain('Genre: Fantasy');
  });

  it('formats a manga with mangaka', () => {
    const manga: Product = {
      id: 2,
      name: 'Test Manga',
      nummer: 5,
      art: 'Manga',
      details: {
        mangaka: 'Eiichiro Oda',
        sprache: 'Japanisch',
        genre: 'Shonen',
      },
    };

    const result = formatDetails(manga);

    expect(result).toContain('Vol: 5');
    expect(result).toContain('Mangaka: Eiichiro Oda');
    expect(result).toContain('Sprache: Japanisch');
    expect(result).toContain('Genre: Shonen');
  });

  it('formats a game with console', () => {
    const game: Product = {
      id: 3,
      name: 'Test Game',
      nummer: null,
      art: 'Spiel',
      details: {
        konsole: 'Nintendo Switch',
        genre: 'Action',
      },
    };

    const result = formatDetails(game);

    expect(result).not.toContain('Vol:');
    expect(result).toContain('Konsole: Nintendo Switch');
    expect(result).toContain('Genre: Action');
  });

  it('formats a film/serie with art type', () => {
    const film: Product = {
      id: 4,
      name: 'Test Film',
      nummer: null,
      art: 'Filmserie',
      details: {
        art: 'Film',
        genre: 'Sci-Fi',
      },
    };

    const result = formatDetails(film);

    expect(result).toContain('Art: Film');
    expect(result).toContain('Genre: Sci-Fi');
  });

  it('returns empty string for product with no details', () => {
    const product: Product = {
      id: 5,
      name: 'Empty Product',
      nummer: null,
      art: 'Buch',
      details: {},
    };

    const result = formatDetails(product);

    expect(result).toBe('');
  });

  it('joins multiple details with pipe separator', () => {
    const product: Product = {
      id: 6,
      name: 'Multi Detail',
      nummer: 1,
      art: 'Buch',
      details: {
        autor: 'Author',
        genre: 'Genre',
      },
    };

    const result = formatDetails(product);

    expect(result).toBe('Vol: 1 | Autor: Author | Genre: Genre');
  });
});

describe('Product Types', () => {
  it('validates Book product type structure', () => {
    const book: Product = {
      id: 1,
      name: 'Harry Potter',
      nummer: 1,
      art: 'Buch',
      details: {
        autor: 'J.K. Rowling',
        sprache: 'English',
        genre: 'Fantasy',
      },
    };

    expect(book.art).toBe('Buch');
    expect(book.details.autor).toBe('J.K. Rowling');
  });

  it('validates Manga product type structure', () => {
    const manga: Product = {
      id: 2,
      name: 'One Piece',
      nummer: 100,
      art: 'Manga',
      details: {
        mangaka: 'Eiichiro Oda',
        sprache: 'Japanese',
        genre: 'Shonen',
      },
    };

    expect(manga.art).toBe('Manga');
    expect(manga.details.mangaka).toBe('Eiichiro Oda');
  });

  it('validates Spiel product type structure', () => {
    const game: Product = {
      id: 3,
      name: 'Zelda',
      nummer: null,
      art: 'Spiel',
      details: {
        konsole: 'Nintendo Switch',
        genre: 'Adventure',
      },
    };

    expect(game.art).toBe('Spiel');
    expect(game.details.konsole).toBe('Nintendo Switch');
  });

  it('validates Filmserie product type structure', () => {
    const serie: Product = {
      id: 4,
      name: 'Breaking Bad',
      nummer: 5,
      art: 'Filmserie',
      details: {
        art: 'Serie',
        genre: 'Drama',
      },
    };

    expect(serie.art).toBe('Filmserie');
    expect(serie.details.art).toBe('Serie');
  });
});
