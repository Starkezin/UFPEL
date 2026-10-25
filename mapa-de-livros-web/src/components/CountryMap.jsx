import React, { useState, useEffect } from 'react';
import { MapContainer, TileLayer, Marker, Popup, useMap } from 'react-leaflet';
import L from 'leaflet';
import { Sun, Moon, Map } from 'lucide-react';
import 'leaflet/dist/leaflet.css';

// Componente interno para controlar o foco e zoom do mapa dinamicamente
function MapController({ countries }) {
  const map = useMap();

  useEffect(() => {
    if (!countries || countries.length === 0) return;

    // Filtra coordenadas válidas de latitude/longitude
    const validCoords = countries
      .filter((c) => c.latlng && c.latlng.length === 2 && !isNaN(c.latlng[0]) && !isNaN(c.latlng[1]))
      .map((c) => [c.latlng[0], c.latlng[1]]);

    if (validCoords.length === 0) return;

    if (validCoords.length === 1) {
      // Se apenas um país, foca nele com zoom moderado
      map.setView(validCoords[0], 4, { animate: true, duration: 1 });
    } else {
      // Se múltiplos países, ajusta os limites para conter todos
      const bounds = L.latLngBounds(validCoords);
      map.fitBounds(bounds, {
        padding: [60, 60],
        maxZoom: 6,
        animate: true,
        duration: 1.2,
      });
    }
  }, [countries, map]);

  return null;
}

// Criador de marcadores personalizados usando a bandeira do país
const createFlagIcon = (flagUrl, countryName) => {
  return L.divIcon({
    html: `
      <div class="custom-flag-marker">
        <img src="${flagUrl}" alt="Bandeira de ${countryName}" class="flag-marker-img" />
        <div class="flag-marker-pulse"></div>
      </div>
    `,
    className: 'custom-flag-icon-container',
    iconSize: [38, 38],
    iconAnchor: [19, 19],
    popupAnchor: [0, -19],
  });
};

export default function CountryMap({ countries, selectedLanguageName }) {
  const [mapTheme, setMapTheme] = useState('dark'); // 'dark' ou 'light'

  // Map Tile URLs
  const tileUrls = {
    dark: 'https://{s}.basemaps.cartocdn.com/dark_all/{z}/{x}/{y}{r}.png',
    light: 'https://{s}.basemaps.cartocdn.com/rastertiles/voyager/{z}/{x}/{y}{r}.png',
  };

  const attribution =
    '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> &copy; <a href="https://carto.com/attributions">CARTO</a>';

  // Formata a população para exibição amigável
  const formatPopulation = (num) => {
    if (num >= 1000000000) return `${(num / 1000000000).toFixed(1)}B`;
    if (num >= 1000000) return `${(num / 1000000).toFixed(1)}M`;
    if (num >= 1000) return `${(num / 1000).toFixed(0)}k`;
    return num.toString();
  };

  return (
    <div className="map-wrapper glass-card">
      <div className="map-header">
        <div className="map-title-container">
          <Map size={18} className="text-primary" />
          <h3>Mapa de Idiomas</h3>
        </div>
        
        {selectedLanguageName && (
          <span className="map-badge">
            Exibindo: {selectedLanguageName}
          </span>
        )}

        {/* Botão flutuante para alternar tema do mapa */}
        <button
          onClick={() => setMapTheme(mapTheme === 'dark' ? 'light' : 'dark')}
          className="map-theme-toggle"
          title={`Mudar para mapa ${mapTheme === 'dark' ? 'claro' : 'escuro'}`}
        >
          {mapTheme === 'dark' ? <Sun size={16} /> : <Moon size={16} />}
        </button>
      </div>

      <div className="map-container-relative">
        <MapContainer
          center={[0, 0]}
          zoom={2}
          scrollWheelZoom={true}
          style={{ height: '100%', width: '100%', borderRadius: '12px' }}
        >
          <TileLayer url={tileUrls[mapTheme]} attribution={attribution} />
          
          <MapController countries={countries} />

          {countries &&
            countries
              .filter((c) => c.latlng && c.latlng.length === 2)
              .map((country) => {
                const nameCommon = country.name?.common || 'Desconhecido';
                const nameOfficial = country.name?.official || nameCommon;
                const capital = country.capital ? country.capital[0] : 'Não informada';
                const population = country.population ? formatPopulation(country.population) : 'n/d';
                const region = country.region || 'n/d';
                const flagUrl = country.flags?.svg || country.flags?.png;
                const latlng = country.latlng;

                // Lista de todos os idiomas oficiais falados no país
                const spokeLangs = country.languages 
                  ? Object.values(country.languages).join(', ') 
                  : 'n/d';

                return (
                  <Marker
                    key={country.cca3 || country.cca2 + Math.random()}
                    position={[latlng[0], latlng[1]]}
                    icon={createFlagIcon(flagUrl, nameCommon)}
                  >
                    <Popup className="custom-map-popup">
                      <div className="popup-card">
                        <div className="popup-header">
                          <img src={flagUrl} alt={`Bandeira de ${nameCommon}`} className="popup-flag" />
                          <div>
                            <h4 className="popup-country-name">{nameCommon}</h4>
                            <span className="popup-country-official">{nameOfficial}</span>
                          </div>
                        </div>
                        
                        <div className="popup-details">
                          <div className="popup-detail-row">
                            <span className="popup-label">Capital:</span>
                            <span className="popup-val">{capital}</span>
                          </div>
                          <div className="popup-detail-row">
                            <span className="popup-label">Região:</span>
                            <span className="popup-val">{region}</span>
                          </div>
                          <div className="popup-detail-row">
                            <span className="popup-label">População:</span>
                            <span className="popup-val">{population} habitantes</span>
                          </div>
                          <div className="popup-detail-row">
                            <span className="popup-label">Idiomas Oficiais:</span>
                            <span className="popup-val popup-langs" title={spokeLangs}>{spokeLangs}</span>
                          </div>
                        </div>
                      </div>
                    </Popup>
                  </Marker>
                );
              })}
        </MapContainer>

        {/* Overlay caso não haja países a exibir ou mapa esteja vazio */}
        {(!countries || countries.length === 0) && (
          <div className="map-overlay-message">
            <div className="overlay-glass-card">
              <p>Selecione um livro com idioma mapeável para visualizar a distribuição geográfica.</p>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
