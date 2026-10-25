import React from 'react';
import { Book, Calendar, User, Globe } from 'lucide-react';
import { getLanguageName } from '../utils/languageMap';

export default function BookList({ books, onSelectBook, selectedBookKey }) {
  if (books.length === 0) {
    return (
      <div className="empty-results-container">
        <Book size={48} className="empty-icon text-muted" />
        <h3>Nenhum livro encontrado</h3>
        <p className="text-muted text-center">
          Tente buscar com outros termos ou verifique a ortografia do título.
        </p>
      </div>
    );
  }

  return (
    <div className="book-list-container">
      <h2 className="section-title">Resultados ({books.length})</h2>
      <div className="book-list-scroll">
        {books.map((book) => {
          const isSelected = selectedBookKey === book.key;
          const authorName = book.author_name ? book.author_name.join(', ') : 'Autor desconhecido';
          const publishYear = book.first_publish_year || 'Ano n/d';
          const languages = book.language ? book.language.slice(0, 3) : [];
          const coverUrl = book.cover_i 
            ? `https://covers.openlibrary.org/b/id/${book.cover_i}-M.jpg` 
            : null;

          return (
            <div
              key={book.key}
              onClick={() => onSelectBook(book)}
              className={`book-card ${isSelected ? 'active' : ''}`}
            >
              {coverUrl ? (
                <div className="book-card-cover-wrapper">
                  <img
                    src={coverUrl}
                    alt={`Capa do livro ${book.title}`}
                    className="book-card-cover"
                    loading="lazy"
                  />
                </div>
              ) : (
                <div className="book-card-cover-placeholder">
                  <Book size={24} />
                  <span className="placeholder-text">Sem Capa</span>
                </div>
              )}
              
              <div className="book-card-details">
                <h3 className="book-card-title" title={book.title}>
                  {book.title}
                </h3>
                
                <div className="book-card-meta">
                  <div className="meta-item" title={`Autor: ${authorName}`}>
                    <User size={14} className="meta-icon" />
                    <span className="meta-text">{authorName}</span>
                  </div>
                  <div className="meta-item">
                    <Calendar size={14} className="meta-icon" />
                    <span className="meta-text">{publishYear}</span>
                  </div>
                </div>

                {languages.length > 0 && (
                  <div className="book-card-tags">
                    <Globe size={12} className="tag-icon text-muted" />
                    {languages.map((lang, idx) => (
                      <span key={idx} className="language-tag">
                        {getLanguageName(lang)}
                      </span>
                    ))}
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>
    </div>
  );
}
