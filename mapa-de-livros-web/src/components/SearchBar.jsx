import React, { useState } from 'react';
import { Search, X } from 'lucide-react';

export default function SearchBar({ onSearch, isLoading }) {
  const [query, setQuery] = useState('');

  const handleSubmit = (e) => {
    e.preventDefault();
    if (query.trim()) {
      onSearch(query.trim());
    }
  };

  const handleClear = () => {
    setQuery('');
  };

  return (
    <form onSubmit={handleSubmit} className="search-bar-container">
      <div className="search-input-wrapper">
        <Search className="search-icon-left text-muted" size={20} />
        <input
          type="text"
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          placeholder="Pesquise por título, autor, saga..."
          className="search-input"
          disabled={isLoading}
        />
        {query && !isLoading && (
          <button
            type="button"
            onClick={handleClear}
            className="clear-button"
            aria-label="Limpar pesquisa"
          >
            <X size={18} />
          </button>
        )}
      </div>
      <button
        type="submit"
        disabled={isLoading || !query.trim()}
        className={`search-submit-btn ${isLoading ? 'loading' : ''}`}
      >
        {isLoading ? (
          <div className="spinner-small"></div>
        ) : (
          <>
            <span>Buscar</span>
            <Search size={16} className="btn-icon-right" />
          </>
        )}
      </button>
    </form>
  );
}
