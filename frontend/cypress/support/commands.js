Cypress.Commands.add('login', (username, password) => {
  cy.visit('/login')
  cy.get('[data-test="login-username-input"]').type(username)
  cy.get('[data-test="password-input"]').type(password)
  cy.get('[data-test="login-button"]').click()
})

Cypress.Commands.add('signup', (fullname, password) => {
  const username = `testuser_${Math.floor(Math.random() * 100000)}`

  cy.visit('/signup')
  cy.get('[data-test="username-input"]').type(username)
  cy.get('[data-test="fullname-input"]').type(fullname)
  cy.get('[data-test="password-input"]').type(password)
  cy.get('[data-test="signup-button"]').click()

  cy.window().then(window => {
    window.localStorage.setItem('signupUsername', username)
  })
})
