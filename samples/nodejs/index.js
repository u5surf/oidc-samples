const express = require('express')
const fs = require('fs')
const tunnel = require('tunnel')
const { Issuer, custom } = require('openid-client')

;(async () => {
  const app = express()
  const PROXY_HOST = process.env.PROXY_HOST
  const PROXY_PORT = process.env.PROXY_PORT
  const CA_CERT = process.env.CA_CERT
  const HOST = process.env.HOST || '0.0.0.0'
  const PORT = process.env.PORT || 8000
  const ISSUER = process.env.ISSUER || 'https://api.mpin.io'
  const CLIENT_ID = process.env.CLIENT_ID
  const CLIENT_SECRET = process.env.CLIENT_SECRET
  const REDIRECT_URL = process.env.REDIRECT_URL || 'http://localhost:8000/login'

  if (!CLIENT_ID || !CLIENT_SECRET) {
    return console.log('ERROR: Both CLIENT_ID and CLIENT_SECRET are required')
  }

  app.listen(PORT, HOST, (err) => {
    if (err) {
      return console.log('ERROR: ', err)
    }

    console.log(`listening on port ${PORT}`)
  })

  const issuer = await Issuer.discover(ISSUER)

  const client = new issuer.Client({
    client_id: CLIENT_ID,
    client_secret: CLIENT_SECRET,
    redirect_uris: [REDIRECT_URL],
    response_types: ['code']
  })

  if (PROXY_HOST && PROXY_PORT && CA_CERT) {
    client[custom.http_options] = function (options) {
      var agent = tunnel.httpsOverHttp({
        ca: [fs.readFileSync(CA_CERT)],
        proxy: {
          host: PROXY_HOST,
          port: PROXY_PORT
        }
      })
      options.agent = agent
      return options
    }
  }

  const authSessions = {}

  app.get('/', async (request, response) => {
    try {
      const currentTokenSet = authSessions[request.query.session]
      if (currentTokenSet) {
        const userInfo = await client.userinfo(currentTokenSet)
        response.send(userInfo)
      } else {
        const state = Math.floor(Math.random() * 90000) + 10000
        response.redirect(client.authorizationUrl({
          scope: 'openid email',
          state
        }))
      }
    } catch (error) {
      response.status(500).send({
        message: 'Error fetching userinfo',
        error
      })
    }
  })

  app.get('/login', async (request, response) => {
    try {
      const params = client.callbackParams(request)
      const tokenSet = await client.callback(REDIRECT_URL, params, { state: params.state })

      authSessions[params.state] = tokenSet

      response.redirect('/?session=' + params.state)
    } catch (error) {
      response.status(500).send({
        message: 'Error logging in',
        error
      })
    }
  })
})().catch(e => {
  console.log('unhandled error - ' + e)
})
