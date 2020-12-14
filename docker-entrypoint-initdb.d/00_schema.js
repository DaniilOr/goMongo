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
db.createCollection('payments', {
    validator: {
        $jsonSchema: {
            bsonType: 'Object',
            required: ['user_id', 'frequent_payments', 'predicted_payments'],
            frequent_payments: {
                payments: {
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
            },
            predicted_payments: {
                payments: {
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
    }
});