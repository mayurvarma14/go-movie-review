// Create app-specific user with read/write access
db.getSiblingDB("admin").auth(
    process.env.MONGO_ROOT_USERNAME,
    process.env.MONGO_ROOT_PASSWORD
  );
  
  db.createUser({
    user: process.env.MONGO_APP_USER,
    pwd: process.env.MONGO_APP_PASSWORD,
    roles: [
      { role: "readWrite", db: process.env.MONGO_INITDB_DATABASE }
    ]
  });