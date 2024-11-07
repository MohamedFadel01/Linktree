describe('Login Flow', () => {
  let uniqueUsername

  beforeEach(() => {
    const fullname = 'Test User'
    const password = 'validpassword123'

    cy.signup(fullname, password)

    cy.window().then(window => {
      uniqueUsername = window.localStorage.getItem('signupUsername')
    })
  })

  it('should display login form', () => {
    cy.get('h2').should('contain', 'Login')
    cy.get('[data-test="login-username-input"]').should('exist')
    cy.get('[data-test="password-input"]').should('exist')
    cy.get('[data-test="login-button"]').should('exist')
  })

  it('should show error for invalid credentials', () => {
    cy.get('[data-test="login-username-input"]').type('invaliduser')
    cy.get('[data-test="password-input"]').type('invalidpass')
    cy.get('[data-test="login-button"]').click()
    cy.get('div').should('contain', 'Login Failed')
  })

  it('should redirect to profile after successful login', () => {
    cy.login(uniqueUsername, 'validpassword123')

    cy.url().should('include', '/profile')
  })
})
