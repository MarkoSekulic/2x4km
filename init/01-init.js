const dbFiles = db.getSiblingDB("files_db");

// Create collection if it does not exist
function createCollectionIfMissing(name) {
    const collections = dbFiles.getCollectionNames();

    if (!collections.includes(name)) {
        dbFiles.createCollection(name);
        print(`Created collection: ${name}`);
    } else {
        print(`Collection already exists: ${name}`);
    }
}

// Create collections
createCollectionIfMissing("roots");
createCollectionIfMissing("folders");

// ROOT indexes
dbFiles.roots.createIndex(
    { userId: 1 },
    { unique: true, name: "userId_unique" }
);

// FOLDER indexes
dbFiles.folders.createIndex(
    { parent: 1 },
    { name: "parent_idx" }
);

dbFiles.folders.createIndex(
    { "files.OSID": 1 },
    { name: "files_osid_idx" }
);

print("Files database initialization completed.");
