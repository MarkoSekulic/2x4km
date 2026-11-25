const filesdb = db.getSiblingDB('files_db');

filesdb.createUser({
  user: 'user',
  pwd: 'pass',
  roles: [{ role: 'readWrite', db: 'files_db' }],
});

// One root document per user
filesdb.createCollection('root', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['user_id', 'user_first_name', 'user_last_name', 'root_folder_id'],
      properties: {
        user_id: { bsonType: 'string', minLength: 1 },
        user_first_name: { bsonType: 'string' },
        user_last_name: { bsonType: 'string' },
        root_folder_id: { bsonType: 'objectId' },
      },
    },
  },
});

filesdb.root.createIndex({ user_id: 1 }, { unique: true });

filesdb.createCollection('folders', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['name', 'parent', 'root', 'child_folders', 'files', 'created_at'],
      properties: {
        name: { bsonType: 'string', minLength: 1 },
        parent: { bsonType: ['objectId', 'null'] }, // null for the root folder
        root: { bsonType: 'objectId' },
        child_folders: { bsonType: 'array', items: { bsonType: 'objectId' } },
        files: { bsonType: 'array', items: { bsonType: 'objectId' } },
        created_at: { bsonType: 'date' },
      },
    },
  },
});

// Covers queries by parent alone as well as the name check in createNode
filesdb.folders.createIndex({ parent: 1, name: 1 });

filesdb.createCollection('files', {
  validator: {
    $jsonSchema: {
      bsonType: 'object',
      required: ['name', 'parent', 'root', 'OSID', 'created_at'],
      properties: {
        name: { bsonType: 'string', minLength: 1 },
        parent: { bsonType: 'objectId' },
        root: { bsonType: 'objectId' },
        OSID: { bsonType: 'string', minLength: 1 }, // object store id
        created_at: { bsonType: 'date' },
      },
    },
  },
});

filesdb.files.createIndex({ parent: 1, name: 1 });

// ---------------------------------------------------------------------------
// Test data (user ids and names are hardcoded, not linked to users_db)
// ---------------------------------------------------------------------------

// --- user 1: empty root folder ---
const emptyUserRootId = new ObjectId();
const emptyUserRootFolderId = new ObjectId();

filesdb.root.insertOne({
  _id: emptyUserRootId,
  user_id: 'a3f9c1d27b8e4a60',
  user_first_name: 'Mila',
  user_last_name: 'Petrovic',
  root_folder_id: emptyUserRootFolderId,
});

filesdb.folders.insertOne({
  _id: emptyUserRootFolderId,
  name: 'root',
  parent: null,
  root: emptyUserRootId,
  child_folders: [],
  files: [],
  created_at: new Date(),
});

// --- user 2: root -> Documents, Photos ---
const flatUserRootId = new ObjectId();
const flatUserRootFolderId = new ObjectId();
const flatUserDocumentsId = new ObjectId();
const flatUserPhotosId = new ObjectId();
const flatUserReadmeId = new ObjectId();
const flatUserReportId = new ObjectId();
const flatUserNotesId = new ObjectId();
const flatUserHolidayId = new ObjectId();

filesdb.root.insertOne({
  _id: flatUserRootId,
  user_id: '7be04d915c2f8a31',
  user_first_name: 'Luka',
  user_last_name: 'Jovanovic',
  root_folder_id: flatUserRootFolderId,
});

filesdb.folders.insertMany([
  {
    _id: flatUserRootFolderId,
    name: 'root',
    parent: null,
    root: flatUserRootId,
    child_folders: [flatUserDocumentsId, flatUserPhotosId],
    files: [flatUserReadmeId],
    created_at: new Date(),
  },
  {
    _id: flatUserDocumentsId,
    name: 'Documents',
    parent: flatUserRootFolderId,
    root: flatUserRootId,
    child_folders: [],
    files: [flatUserReportId, flatUserNotesId],
    created_at: new Date(),
  },
  {
    _id: flatUserPhotosId,
    name: 'Photos',
    parent: flatUserRootFolderId,
    root: flatUserRootId,
    child_folders: [],
    files: [flatUserHolidayId],
    created_at: new Date(),
  },
]);

filesdb.files.insertMany([
  { _id: flatUserReadmeId,  name: 'readme.txt',  parent: flatUserRootFolderId, root: flatUserRootId, OSID: 'os-4c81e9a2-readme-txt',  created_at: new Date() },
  { _id: flatUserReportId,  name: 'report.pdf',  parent: flatUserDocumentsId,  root: flatUserRootId, OSID: 'os-91d3b7f0-report-pdf',  created_at: new Date() },
  { _id: flatUserNotesId,   name: 'notes.md',    parent: flatUserDocumentsId,  root: flatUserRootId, OSID: 'os-5a6e2c48-notes-md',    created_at: new Date() },
  { _id: flatUserHolidayId, name: 'holiday.jpg', parent: flatUserPhotosId,     root: flatUserRootId, OSID: 'os-e07f19bd-holiday-jpg', created_at: new Date() },
]);

// --- user 3: root -> Work -> Projects ---
const nestedUserRootId = new ObjectId();
const nestedUserRootFolderId = new ObjectId();
const nestedUserWorkId = new ObjectId();
const nestedUserProjectsId = new ObjectId();
const nestedUserTodoId = new ObjectId();
const nestedUserBudgetId = new ObjectId();
const nestedUserPlanId = new ObjectId();

filesdb.root.insertOne({
  _id: nestedUserRootId,
  user_id: 'd25a8f4e0b9c7163',
  user_first_name: 'Sara',
  user_last_name: 'Nikolic',
  root_folder_id: nestedUserRootFolderId,
});

filesdb.folders.insertMany([
  {
    _id: nestedUserRootFolderId,
    name: 'root',
    parent: null,
    root: nestedUserRootId,
    child_folders: [nestedUserWorkId],
    files: [nestedUserTodoId],
    created_at: new Date(),
  },
  {
    _id: nestedUserWorkId,
    name: 'Work',
    parent: nestedUserRootFolderId,
    root: nestedUserRootId,
    child_folders: [nestedUserProjectsId],
    files: [nestedUserBudgetId],
    created_at: new Date(),
  },
  {
    _id: nestedUserProjectsId,
    name: 'Projects',
    parent: nestedUserWorkId,
    root: nestedUserRootId,
    child_folders: [],
    files: [nestedUserPlanId],
    created_at: new Date(),
  },
]);

filesdb.files.insertMany([
  { _id: nestedUserTodoId,   name: 'todo.txt',    parent: nestedUserRootFolderId, root: nestedUserRootId, OSID: 'os-38b0d6c5-todo-txt',    created_at: new Date() },
  { _id: nestedUserBudgetId, name: 'budget.xlsx', parent: nestedUserWorkId,       root: nestedUserRootId, OSID: 'os-b6724e1a-budget-xlsx', created_at: new Date() },
  { _id: nestedUserPlanId,   name: 'plan.docx',   parent: nestedUserProjectsId,   root: nestedUserRootId, OSID: 'os-0fc95a37-plan-docx',   created_at: new Date() },
]);