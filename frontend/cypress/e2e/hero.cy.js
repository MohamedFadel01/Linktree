describe('Hero Component', () => {
  it('should display correct CTA button when not authenticated', () => {
    cy.visit('/')
    cy.contains('Get Started').should('be.visible')
    cy.contains('Go to Profile').should('not.exist')
  })

  it('should navigate to correct page when clicking CTA', () => {
    cy.visit('/')
    cy.contains('Get Started').click()
    cy.url().should('include', '/signup')
  })
})
