describe('MessageBox Component', () => {
  it('should display error message correctly', () => {
    cy.visit('/login')
    cy.get('input#username').type('invaliduser')
    cy.get('input#password').type('invalidpass')
    cy.get('button[type="submit"]').click()

    cy.get('div').contains('Login Failed')
  })

  it('should auto-dismiss after duration', () => {
    cy.visit('/login')
    cy.get('input#username').type('invaliduser')
    cy.get('input#password').type('invalidpass')
    cy.get('button[type="submit"]').click()

    cy.get('div').contains('Login Failed').should('be.visible')
    cy.wait(5000)
    cy.get('div').contains('Login Failed').should('not.exist')
  })

  it('should dismiss on click', () => {
    cy.visit('/login')
    cy.get('input#username').type('invaliduser')
    cy.get('input#password').type('invalidpass')
    cy.get('button[type="submit"]').click()

    cy.get('div').contains('Login Failed').click()
    cy.get('div').contains('Login Failed').should('not.exist')
  })
})
