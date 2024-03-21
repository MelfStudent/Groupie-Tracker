function searchArtist() {
    var input = document.getElementsByClassName('input')[0].value.toLowerCase();
    var cards = document.querySelectorAll('.flip-card');

    for (var i = 0; i < cards.length; i++) {
        var artistName = cards[i].getAttribute('data-artist').toLowerCase();
        if (artistName.includes(input)) {
            cards[i].style.display = '';
        } else {
            cards[i].style.display = 'none';
        }
    }
}