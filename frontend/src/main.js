import './style.css'
import { setupMap } from './crime_map.js'

document.querySelector('#app').innerHTML = `
<section id="map-section" aria-label="Consulta de ocorrências no mapa">
    <div class="map-container">
        <div id="map"></div>
    </div>
    <aside class="search-panel" aria-label="Busca de ocorrências">
        <h1>Ocorrências</h1>
        <form id="search-form">
            <div class="radius-field">
                <label for="search-radius">Raio de busca</label>
                <div class="radius-input">
                    <input id="search-radius" name="radius" type="number" min="0" step="any" value="20" required inputmode="decimal">
                    <span>km</span>
                </div>
            </div>
            <button id="search-button" type="submit">Buscar</button>
            <button id="use-location" class="location-button" type="button">
                <svg viewBox="0 0 24 24" width="16" height="16" fill="none" stroke="currentColor" stroke-width="1.7" stroke-linejoin="round" aria-hidden="true"><path d="m21 3-6.5 18-4-7.5L3 9.5 21 3Z"/></svg>
                Usar minha localização
            </button>
        </form>
        <label class="heat-toggle">
            <input id="show-heat" type="checkbox" checked>
            Áreas de calor
        </label>
        <div class="status-area">
            <p id="search-status" role="status" aria-live="polite" aria-atomic="true" data-type="info">Use sua localização para começar.</p>
            <button id="retry-search" class="text-button" type="button" hidden>Tentar novamente</button>
        </div>
        <details class="test-tools">
            <summary>Pontos de teste</summary>
            <button id="add-test-points" class="text-button" type="button" disabled>Adicionar três pontos próximos</button>
        </details>
    </aside>
    <div class="map-legend">
        <span class="location-label"><i class="location-dot"></i>Sua localização</span>
        <div id="heat-legend">
            <span>Concentração de ocorrências</span>
            <div class="heat-scale"><span>Menor</span><i></i><span>Maior</span></div>
        </div>
    </div>
</section>
`

setupMap(document.querySelector('#map'))
