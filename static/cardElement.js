
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
        
        const root = this.attachShadow({mode: 'closed'});
        root.appendChild(this.img);
    }

    /**
     * @param {Card}
     */
    setCard(card) {
        this.card = card;
        this.img.src = "deck/card-deck-" + card.toString() + ".png";
    }

    getCard() {
        return this.card;
    }
}

customElements.define('spielekiste-card', CardElement);