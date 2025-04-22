import Keycloak from "keycloak-js"

const keycloak = new Keycloak({
    url: 'https://auth.spech.dev',
    realm: 'bibrex',
    clientId: 'bibrex-frontend',
})

export default keycloak