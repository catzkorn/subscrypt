import { cyan } from "color-name";

describe('Navigate to "/"', () => {
  it("visits the homepage", () => {
    cy.visit('http://localhost:5000');


    cy.get('#edit-user-button').click()
    cy.get('#user-name')
      .type('charlotte');
    cy.get('#create-user-button').click()
    cy.get('#user-name')
    .should('have.value', 'charlottecharlotte')
  });
})