function searchArtist() {
    var input = document.getElementsByClassName('input')[0].value.toLowerCase();
    var cards = document.querySelectorAll('.flip-card');

    var found = false;
    console.clear();
    console.log(cards[0].getAttribute('data-artist').toLowerCase());
    for (var i = 0; i < cards.length; i++) {
        var artistName = cards[i].getAttribute('data-artist').toLowerCase();
        if (artistName.includes(input)) {
            cards[i].style.display = '';
            found = true;
        } else {
            cards[i].style.display = 'none';
        }
    }
}