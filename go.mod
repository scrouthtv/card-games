module github.com/scrouthtv/card-games

go 1.15

replace (
    github.com/scrouthtv/card-games/doko => ./doko
    github.com/scrouthtv/card-games/logic => ./logic
)

require github.com/gorilla/websocket v1.4.2

require (
    github.com/scrouthtv/card-games/doko v0.0.0-00010101000000-000000000000
    github.com/scrouthtv/card-games/logic v0.0.0-00010101000000-000000000000
)