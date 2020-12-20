db.createUser({
    user: 'app',
    pwd: 'pass',
    roles: [
        {
            role: 'readWrite',
            db: 'predictions',
        }
    ]
});
db.createCollection('suggestions', {
    validator: {
        $jsonSchema: {
            bsonType: 'Object',
            required: ['user_id', 'suggested_payments'],
            suggested_payments: {
                    bsonType: 'Object',
                    required: ['icon', 'name', 'link'],
                    properties: {
                        icon: {
                            bsonType: "string",
                        },
                        name: {
                            bsonType: "string",
                        },
                        link: {
                            bsonType: "string",
                        }
                    }
            }
        }
    }
});