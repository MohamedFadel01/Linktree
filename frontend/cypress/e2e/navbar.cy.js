describe('Navbar Navigation', () => {
  it('should display correct navigation items when not authenticated', () => {
    cy.visit('/')
    cy.get('nav').within(() => {
      cy.contains('Home').should('exist')
      cy.contains('Menu').should('exist')
      cy.get('button').contains('Menu').click()
      cy.contains('Login / Signup').should('be.visible')
      cy.contains('Profile').should('not.exist')
    })
  })

  it('should display correct navigation items when authenticated', () => {
    const fullname = 'Test User'
    const password = 'validpassword123'
    cy.signup(fullname, password)

    cy.window().then(window => {
      const uniqueUsername = window.localStorage.getItem('signupUsername')
      cy.login(uniqueUsername, password)
    })

    cy.url().should('not.include', '/login')

    cy.get('nav').within(() => {
      cy.contains('Home').should('exist')
      cy.contains('Menu').should('exist')
      cy.get('button').contains('Menu').click()
      cy.contains('Profile').should('be.visible')
      cy.contains('Logout').should('be.visible')
    })
  })

  it('should handle user search', () => {
    cy.visit('/')
    cy.get('input[type="text"]').type('testuser{enter}')
    cy.url().should('include', '/testuser')
  })

  it('should show error message for non-existent user search', () => {
    cy.visit('/')
    cy.get('input[type="text"]').type('nonexistentuser{enter}')
    cy.contains('User Not Found').should('be.visible')
  })
})
