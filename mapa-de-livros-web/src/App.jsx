import React, { useState } from 'react';
import { BookOpen, AlertCircle, RefreshCw } from 'lucide-react';
import SearchBar from './components/SearchBar';
import BookList from './components/BookList';
import BookDetails from './components/BookDetails';
import CountryMap from './components/CountryMap';
import LoadingState from './components/LoadingState';
import { getLanguageName } from './utils/languageMap';

export default function App() {
  const [books, setBooks] = useState([]);
  const [selectedBook, setSelectedBook] = useState(null);
  const [searchLoading, setSearchLoading] = useState(false);
  const [searchError, setSearchError] = useState(null);
  const [lastQuery, setLastQuery] = useState('');
  
  // Estado para os países do idioma da obra selecionada
  const [countriesState, setCountriesState] = useState({
    loading: false,
    error: null,
    data: [],
  });

  // Função para pesquisar obras na Open Library API
  const handleSearch = async (query) => {
    if (!query) return;
    
    setSearchLoading(true);
    setSearchError(null);
    setSelectedBook(null); // Limpa seleção anterior ao buscar
    setLastQuery(query);
    setBooks([]);
    setCountriesState({ loading: false, error: null, data: [] });

    try {
      const response = await fetch(
        `https://openlibrary.org/search.json?title=${encodeURIComponent(query)}&limit=30`
      );
      
      if (!response.ok) {
        throw new Error('Não foi possível obter resposta da Open Library API.');
      }
      
      const data = await response.json();
      
      // Filtra obras duplicadas ou sem título se necessário, mantendo o retorno da API
      const docs = data.docs || [];
      setBooks(docs);
    } catch (err) {
      console.error(err);
      setSearchError('Erro ao buscar livros. Verifique sua conexão e tente novamente.');
    } finally {
      setSearchLoading(false);
    }
  };

  // Função para selecionar um livro e buscar os países pelo idioma
  const handleSelectBook = async (book) => {
    setSelectedBook(book);
    
    // Obtém o código de idioma principal (primeiro da lista de idiomas)
    const langCode = book.language && book.language.length > 0 ? book.language[0] : null;
    
    if (!langCode) {
      setCountriesState({
        loading: false,
        error: 'Idioma não especificado nesta obra na Open Library.',
        data: [],
      });
      return;
    }

    setCountriesState({ loading: true, error: null, data: [] });

    try {
      // Consulta os países que falam o idioma correspondente
      const response = await fetch(
        `https://restcountries.com/v3.1/lang/${encodeURIComponent(langCode)}?fields=name,cca2,cca3,latlng,flags,capital,population,languages,region`
      );

      // Tratamento especial para erro 404 (nenhum país encontrado para o código de idioma)
      if (response.status === 404) {
        setCountriesState({
          loading: false,
          error: null,
          data: [],
        });
        return;
      }

      if (!response.ok) {
        throw new Error(`Serviço REST Countries indisponível (Erro ${response.status}).`);
      }

      const data = await response.json();
      
      // Garante que o retorno seja um array de países válidos
      const countriesList = Array.isArray(data) ? data : [];
      setCountriesState({
        loading: false,
        error: null,
        data: countriesList,
      });
    } catch (err) {
      console.error(err);
      setCountriesState({
        loading: false,
        error: 'Falha ao conectar com o serviço de mapeamento de países.',
        data: [],
      });
    }
  };

  // Limpa a seleção e volta para a lista de resultados da pesquisa atual
  const handleClearSelection = () => {
    setSelectedBook(null);
    setCountriesState({ loading: false, error: null, data: [] });
  };

  // Retorna o nome amigável do idioma principal do livro selecionado para o mapa
  const getSelectedLanguageName = () => {
    if (!selectedBook || !selectedBook.language || selectedBook.language.length === 0) return '';
    return getLanguageName(selectedBook.language[0]);
  };

  return (
    <div className="app-container">
      {/* Cabeçalho */}
      <header className="app-header">
        <div className="logo-section">
          <BookOpen className="logo-icon" size={28} />
          <div>
            <h1 className="logo-title">LibGeo</h1>
            <span className="logo-subtitle">Mapeamento Geográfico de Obras Literárias</span>
          </div>
        </div>
        <div className="app-meta-info">
          <span>Explore a conexão entre livros e países</span>
        </div>
      </header>

      {/* Grid Principal */}
      <main className="main-layout">
        {/* Painel Esquerdo */}
        <section className="side-panel">
          <SearchBar onSearch={handleSearch} isLoading={searchLoading} />
          
          <div className="glass-card flex-1 min-height-0 flex flex-col p-4 relative" style={{ display: 'flex', flexDirection: 'column' }}>
            {selectedBook ? (
              <BookDetails
                book={selectedBook}
                countriesState={countriesState}
                onClear={handleClearSelection}
              />
            ) : (
              <>
                {searchLoading && <LoadingState count={4} />}
                
                {searchError && !searchLoading && (
                  <div className="error-screen">
                    <AlertCircle size={36} className="text-error mb-2" />
                    <h3>Ocorreu um erro</h3>
                    <p>{searchError}</p>
                    <button
                      onClick={() => handleSearch(lastQuery)}
                      className="btn-primary-small flex-center"
                      style={{ gap: '0.4rem' }}
                    >
                      <RefreshCw size={14} />
                      <span>Tentar Novamente</span>
                    </button>
                  </div>
                )}
                
                {!searchLoading && !searchError && (
                  <BookList
                    books={books}
                    onSelectBook={handleSelectBook}
                    selectedBookKey={selectedBook?.key}
                  />
                )}
              </>
            )}
          </div>
        </section>

        {/* Painel Direito (Mapa) */}
        <section className="map-panel">
          <CountryMap
            countries={countriesState.data}
            selectedLanguageName={getSelectedLanguageName()}
          />
        </section>
      </main>
    </div>
  );
}
