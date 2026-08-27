import './style.css'
import { setupMap } from './crime_map.js'

document.querySelector('#app').innerHTML = `
<section id="map-section">
    <div id="map">
      
    </div>
</section>
`

setupMap(document.querySelector('#map'))
