db.suggestions.insertOne(
    {
        user_id: 1,
        suggested_payments: [
            {
                icon: "http://icon2.html",
                name: "Yota",
                link: "http://link2,html"
            },
            {
                icon: "http://icon3.html",
                name: "Мегафон",
                link: "http://link3,html"
            },
        ],
    }
)
db.suggestions.insertOne(
    {
        user_id: 2,
        suggested_payments: [
            {
                icon: "http://somepick.html",
                name: "Tele2",
                link: "http://somelink,html"
            },
            {
                icon: "http://newpick.html",
                name: "Beeline",
                link: "http://somelink,html"
            },
        ],
    }
)