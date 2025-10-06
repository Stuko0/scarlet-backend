db = db.getSiblingDB('scarlet');

db.createCollection('current_wildfires');
db.createCollection('wildfire_weather');

db.current_wildfires.createIndex({"location": "2dsphere"});

db.current_wildfires.createIndex(
    {"last_updated": 1},
    {expireAfterSeconds:86400}
);

db.runCommand({
    collMod: "current_wildfires",
    validator:{
        $jsonSchema:{
            bsonType: "object",
            required: ["fire_id", "location", "discovery_time"],
            properties:{
                fire_id:{bsonType: "string"},
                location: {
                    bsonType:"object",
                    required: ["type", "coordinates"],
                    properties:{
                        type: {enum: ["Point"]},
                        coordinates:{
                            bsonType: ["array"],
                            items: {bsonType:["double"]}
                        }
                    }
                },
                discovery_time:{bsonType: "date"}
            }
        }
    }
});