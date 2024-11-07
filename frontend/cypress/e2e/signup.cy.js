describe('Signup Flow', () => {
  beforeEach(() => {
    cy.visit('/signup')
  })

  it('should display signup form', () => {
    cy.get('h2').should('contain', 'Sign Up')
    cy.get('[data-test="fullname-input"]').should('exist')
    cy.get('[data-test="username-input"]').should('exist')
    cy.get('[data-test="bio-input"]').should('exist')
    cy.get('[data-test="password-input"]').should('exist')
    cy.get('[data-test="signup-button"]').should('exist')
  })

  it('should show error for invalid signup', () => {
    cy.get('[data-test="fullname-input"]').type('Test User')
    const uniqueUsername = 'testuser'
    cy.get('[data-test="username-input"]').type(uniqueUsername)
    cy.get('[data-test="bio-input"]').type('Test bio')
    cy.get('[data-test="password-input"]').type('short')
    cy.get('[data-test="signup-button"]').click()

    cy.get('.fixed').should('be.visible').and('contain', 'Signup Failed')
  })

  it('should redirect to login after successful signup', () => {
    cy.get('[data-test="fullname-input"]').type('Test User')
    const uniqueUsername = `testuser_${Math.floor(Math.random() * 1000)}`
    cy.get('[data-test="username-input"]').type(uniqueUsername)
    cy.get('[data-test="bio-input"]').type('Test bio')
    cy.get('[data-test="password-input"]').type('validpassword123')
    cy.get('[data-test="signup-button"]').click()

    cy.url().should('include', '/login')
  })
})
