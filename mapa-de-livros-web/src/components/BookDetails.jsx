import React from 'react';
import { Book, Calendar, User, Globe, HelpCircle, Building, MapPin, AlertCircle } from 'lucide-react';
import { getLanguageName } from '../utils/languageMap';

export default function BookDetails({ book, countriesState, onClear }) {
  if (!book) {
    return (
      <div className="no-selection-details">
        <div className="glass-card flex-center text-center p-6 select-prompt">
          <Book size={48} className="pulse text-primary mb-4" />
          <h3>Nenhum livro selecionado</h3>
          <p className="text-muted mt-2">
            Pesquise por uma obra e selecione-a para ver detalhes e visualizar no mapa os países que compartilham seu idioma.
          </p>
        </div>
      </div>
    );
  }

  const authorName = book.author_name ? book.author_name.join(', ') : 'Autor desconhecido';
  const publishYear = book.first_publish_year || 'Ano não informado';
  const publishers = book.publisher ? book.publisher.slice(0, 2).join(', ') : 'Editora não informada';
  const primaryLangCode = book.language && book.language.length > 0 ? book.language[0] : null;
  const primaryLangName = primaryLangCode ? getLanguageName(primaryLangCode) : 'Não informado';
  const coverUrl = book.cover_i 
    ? `https://covers.openlibrary.org/b/id/${book.cover_i}-L.jpg` 
    : null;

  return (
    <div className="book-details-card glass-card">
      <div className="details-header">
        <h3 className="section-title">Detalhes da Obra</h3>
        {onClear && (
          <button onClick={onClear} className="btn-text-close" title="Fechar detalhes">
            Limpar
          </button>
        )}
      </div>

      <div className="details-content">
        <div className="details-main-info">
          {coverUrl ? (
            <div className="details-cover-wrapper">
              <img
                src={coverUrl}
                alt={`Capa de ${book.title}`}
                className="details-cover"
                onError={(e) => {
                  e.target.onerror = null;
                  e.target.parentNode.className = 'book-card-cover-placeholder details-cover-placeholder';
                  e.target.outerHTML = '<svg stroke="currentColor" fill="none" stroke-width="2" viewBox="0 0 24 24" stroke-linecap="round" stroke-linejoin="round" height="48" width="48" xmlns="http://www.w3.org/2000/svg"><path d="M4 19.5A2.5 2.5 0 0 1 6.5 17H20"></path><path d="M6.5 2H20v20H6.5A2.5 2.5 0 0 1 4 19.5v-15A2.5 2.5 0 0 1 6.5 2z"></path></svg><span class="placeholder-text">Sem Capa</span>';
                }}
              />
            </div>
          ) : (
            <div className="details-cover-placeholder">
              <Book size={48} className="text-muted" />
              <span className="placeholder-text">Sem Capa</span>
            </div>
          )}

          <div className="details-meta-fields">
            <h2 className="details-title">{book.title}</h2>
            
            <div className="details-field">
              <User size={16} className="details-field-icon" />
              <div>
                <span className="field-label">Autor(es)</span>
                <p className="field-value">{authorName}</p>
              </div>
            </div>

            <div className="details-field">
              <Calendar size={16} className="details-field-icon" />
              <div>
                <span className="field-label">Publicado em</span>
                <p className="field-value">{publishYear}</p>
              </div>
            </div>

            <div className="details-field">
              <Building size={16} className="details-field-icon" />
              <div>
                <span className="field-label">Editora(s)</span>
                <p className="field-value">{publishers}</p>
              </div>
            </div>

            <div className="details-field">
              <Globe size={16} className="details-field-icon" />
              <div>
                <span className="field-label">Idioma Principal</span>
                <p className="field-value highlight-text">
                  {primaryLangName} {primaryLangCode && `(${primaryLangCode.toUpperCase()})`}
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Status Geográfico do Idioma */}
        <div className="geographic-status">
          <div className="geo-status-title">
            <MapPin size={16} className="text-primary" />
            <h4>Mapeamento Geográfico</h4>
          </div>

          {!primaryLangCode ? (
            <div className="status-box warning">
              <AlertCircle size={18} className="status-box-icon" />
              <p>Esta obra não possui informações de idioma. Impossível mapear os países relacionados.</p>
            </div>
          ) : countriesState.loading ? (
            <div className="status-box loading">
              <div className="spinner-small spinner-primary"></div>
              <p>Consultando países falantes de <strong>{primaryLangName}</strong>...</p>
            </div>
          ) : countriesState.error ? (
            <div className="status-box error">
              <AlertCircle size={18} className="status-box-icon" />
              <div className="status-box-content">
                <p>Erro ao obter países:</p>
                <span className="error-msg-detail">{countriesState.error}</span>
              </div>
            </div>
          ) : countriesState.data && countriesState.data.length > 0 ? (
            <div className="status-box success">
              <Globe size={18} className="status-box-icon text-success" />
              <p>
                Idioma <strong>{primaryLangName}</strong> localizado! Mapeado em{' '}
                <strong>{countriesState.data.length}</strong> país(es).
              </p>
            </div>
          ) : (
            <div className="status-box info">
              <HelpCircle size={18} className="status-box-icon text-info" />
              <p>Nenhum país encontrado na base de dados para o idioma {primaryLangName}.</p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
