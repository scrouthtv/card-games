/**
 * @class
 */
class CardElement extends HTMLElement {
  constructor() {
    super();

    this.img = document.createElement("img");
    this.img.src = "deck/card-deck-back.png";
    this.img.style.width = "100%";
    this.img.style.height = "100%";

    const root = this.attachShadow({ mode: "closed" });
    root.appendChild(this.img);
  }

  backSide() {
    this.card = null;
    this.back = true;
    this.img.src = "deck/card-deck-back.png";
  }

  /**
   * @param {Card}
   */
  setCard(card) {
    this.card = card;
    this.img.src = "deck/card-deck-" + card.toString() + ".png";
  }

  isBackSide() {
    return this.back;
  }

  getCard() {
    return this.card;
  }
}


export { CardElement };
// customElements.define("spielekiste-card", CardElement);
