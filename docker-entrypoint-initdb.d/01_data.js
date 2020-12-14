db.payments.insertOne(
    {
        user_id: 1,
        frequent_payments: [
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
        predicted_payments: [
            {
                icon: "http://icon1.html",
                name: "Tinkoff",
                link: "http://link1,html"
            },
        ]
    }
)