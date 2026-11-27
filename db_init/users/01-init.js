const usersdb = db.getSiblingDB('users_db');

usersdb.createUser({
  user: 'user',
  pwd: 'pass',
  roles: [{ role: 'readWrite', db: 'users_db' }],
});

usersdb.createCollection('users', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['email', 'passwordHash','firstName', 'lastName', 'createdAt'],
      properties: {
        email: { bsonType: 'string', pattern: '^.+@.+$' },
        firstName: { bsonType: 'string' },
        lastName: { bsonType: 'string' },
        role: { enum: ['admin', 'user', 'guest'] },
        passwordHash: { bsonType: 'string', minLength: 50 },
        dataLimitMB: { bsonType: ['int', 'long', 'double'], minimum: 0 },
        createdAt: { bsonType: 'date' },
      },
    },
  },
});

usersdb.users.createIndex({ email: 1 }, { unique: true });

usersdb.users.insertMany([
  { email: 'admin@example.com', firstName: 'Adam', lastName: "Devic",  role: 'admin', createdAt: new Date(), passwordHash: '$2b$12$b/vqM7vad2v/AfK2d3fK5e4Yx06ElmslYjx9.iBghfljES24KXrLW' }, // adminpass
  { email: 'alice@example.com', firstName: 'Alice',lastName: "Devic", role: 'user',  createdAt: new Date(), passwordHash: '$2b$12$Vd7sASRhp6woLmQ.YqlDieAMLQlRjHYayQ9xbHRsrFBMeqd22V/CO' }, // alicepass
  { email: 'bob@example.com',   firstName: 'Bob',  lastName: "Devic",   role: 'user',  createdAt: new Date(), passwordHash: '$2b$12$KXwLS4D1aw7i.so0TdJviOhioH9nlPEVmfYmPF7vy.Vzzgwtuwda2' }, // bobpass
]);

usersdb.createCollection('verifies')
