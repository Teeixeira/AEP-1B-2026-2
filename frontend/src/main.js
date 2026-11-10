import './style.css'
import { setupMap } from './crime_map.js'

document.querySelector('#app').innerHTML = `
<section id="map-section">
    <div id="map-actions">
        <button id="use-location" type="button">Usar minha localização</button>
        <button id="add-test-points" type="button" disabled>Adicionar pontos de teste</button>
    </div>
    <div id="map">
    </div>
</section>
`

setupMap(document.querySelector('#map'))
