USE notifications_db;

TRUNCATE notifications_by_user;
TRUNCATE notifications_by_user_type;
TRUNCATE notifications_by_user_actor;
TRUNCATE notifications_by_user_file;


  INSERT INTO notifications_by_user (user_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', minTimeuuid('2026-11-20 09:00:00+0000'), 'shared', '22222222-2222-2222-2222-222222222222', 'Bob Petrovic', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects');
  INSERT INTO notifications_by_user_type (user_id, n_type, notified_at, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'shared', minTimeuuid('2026-11-20 09:00:00+0000'), '22222222-2222-2222-2222-222222222222', 'Bob Petrovic', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects');
  INSERT INTO notifications_by_user_actor (user_id, by_user_id, notified_at, n_type, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', '22222222-2222-2222-2222-222222222222', minTimeuuid('2026-11-20 09:00:00+0000'), 'shared', 'Bob Petrovic', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects');
  INSERT INTO notifications_by_user_file (user_id, f_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', minTimeuuid('2026-11-20 09:00:00+0000'), 'shared', '22222222-2222-2222-2222-222222222222', 'Bob Petrovic', 'folder', 'Projects');



  INSERT INTO notifications_by_user (user_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name) VALUES ('6b01b85b2cf4749e214d50ee', minTimeuuid('2026-11-20 10:30:00+0000'), 'added', '22222222-2222-2222-2222-222222222222', 'Bob Petrovic', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects', 'cccccccc-cccc-cccc-cccc-ccccccccccc1', 'notes.txt');
  INSERT INTO notifications_by_user_type (user_id, n_type, notified_at, by_user_id, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'added', minTimeuuid('2026-11-20 10:30:00+0000'), '22222222-2222-2222-2222-222222222222', 'Bob Petrovic', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects', 'cccccccc-cccc-cccc-cccc-ccccccccccc1', 'notes.txt');
  INSERT INTO notifications_by_user_actor (user_id, by_user_id, notified_at, n_type, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name) VALUES ('6b01b85b2cf4749e214d50ee', '22222222-2222-2222-2222-222222222222', minTimeuuid('2026-11-20 10:30:00+0000'), 'added', 'Bob Petrovic', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects', 'cccccccc-cccc-cccc-cccc-ccccccccccc1', 'notes.txt');
  INSERT INTO notifications_by_user_file (user_id, f_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_name, added_f_id, added_f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', minTimeuuid('2026-11-20 10:30:00+0000'), 'added', '22222222-2222-2222-2222-222222222222', 'Bob Petrovic', 'folder', 'Projects', 'cccccccc-cccc-cccc-cccc-ccccccccccc1', 'notes.txt');



  INSERT INTO notifications_by_user (user_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', minTimeuuid('2026-11-21 14:15:00+0000'), 'updated', '33333333-3333-3333-3333-333333333333', 'Carol Jovanovic', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'report.pdf');
  INSERT INTO notifications_by_user_type (user_id, n_type, notified_at, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'updated', minTimeuuid('2026-11-21 14:15:00+0000'), '33333333-3333-3333-3333-333333333333', 'Carol Jovanovic', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'report.pdf');
  INSERT INTO notifications_by_user_actor (user_id, by_user_id, notified_at, n_type, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', '33333333-3333-3333-3333-333333333333', minTimeuuid('2026-11-21 14:15:00+0000'), 'updated', 'Carol Jovanovic', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'report.pdf');
  INSERT INTO notifications_by_user_file (user_id, f_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', minTimeuuid('2026-11-21 14:15:00+0000'), 'updated', '33333333-3333-3333-3333-333333333333', 'Carol Jovanovic', 'file', 'report.pdf');



  INSERT INTO notifications_by_user (user_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', minTimeuuid('2026-11-22 08:45:00+0000'), 'shared', '33333333-3333-3333-3333-333333333333', 'Carol Jovanovic', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', 'budget.xlsx');
  INSERT INTO notifications_by_user_type (user_id, n_type, notified_at, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'shared', minTimeuuid('2026-11-22 08:45:00+0000'), '33333333-3333-3333-3333-333333333333', 'Carol Jovanovic', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', 'budget.xlsx');
  INSERT INTO notifications_by_user_actor (user_id, by_user_id, notified_at, n_type, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', '33333333-3333-3333-3333-333333333333', minTimeuuid('2026-11-22 08:45:00+0000'), 'shared', 'Carol Jovanovic', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', 'budget.xlsx');
  INSERT INTO notifications_by_user_file (user_id, f_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb2', minTimeuuid('2026-11-22 08:45:00+0000'), 'shared', '33333333-3333-3333-3333-333333333333', 'Carol Jovanovic', 'file', 'budget.xlsx');



  INSERT INTO notifications_by_user (user_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name) VALUES ('6b01b85b2cf4749e214d50ee', minTimeuuid('2026-11-23 16:20:00+0000'), 'added', '33333333-3333-3333-3333-333333333333', 'Carol Jovanovic', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', 'Photos', 'cccccccc-cccc-cccc-cccc-ccccccccccc2', 'design.png');
  INSERT INTO notifications_by_user_type (user_id, n_type, notified_at, by_user_id, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'added', minTimeuuid('2026-11-23 16:20:00+0000'), '33333333-3333-3333-3333-333333333333', 'Carol Jovanovic', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', 'Photos', 'cccccccc-cccc-cccc-cccc-ccccccccccc2', 'design.png');
  INSERT INTO notifications_by_user_actor (user_id, by_user_id, notified_at, n_type, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name) VALUES ('6b01b85b2cf4749e214d50ee', '33333333-3333-3333-3333-333333333333', minTimeuuid('2026-11-23 16:20:00+0000'), 'added', 'Carol Jovanovic', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', 'Photos', 'cccccccc-cccc-cccc-cccc-ccccccccccc2', 'design.png');
  INSERT INTO notifications_by_user_file (user_id, f_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_name, added_f_id, added_f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa2', minTimeuuid('2026-11-23 16:20:00+0000'), 'added', '33333333-3333-3333-3333-333333333333', 'Carol Jovanovic', 'folder', 'Photos', 'cccccccc-cccc-cccc-cccc-ccccccccccc2', 'design.png');



  INSERT INTO notifications_by_user (user_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', minTimeuuid('2026-11-24 09:10:00+0000'), 'updated', '22222222-2222-2222-2222-222222222222', 'Bob Petrovic', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'report.pdf');
  INSERT INTO notifications_by_user_type (user_id, n_type, notified_at, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'updated', minTimeuuid('2026-11-24 09:10:00+0000'), '22222222-2222-2222-2222-222222222222', 'Bob Petrovic', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'report.pdf');
  INSERT INTO notifications_by_user_actor (user_id, by_user_id, notified_at, n_type, by_user_name, f_type, f_id, f_name) VALUES ('6b01b85b2cf4749e214d50ee', '22222222-2222-2222-2222-222222222222', minTimeuuid('2026-11-24 09:10:00+0000'), 'updated', 'Bob Petrovic', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'report.pdf');
  INSERT INTO notifications_by_user_file (user_id, f_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_name) VALUES ('6b01b85b2cf4749e214d50ee', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', minTimeuuid('2026-11-24 09:10:00+0000'), 'updated', '22222222-2222-2222-2222-222222222222', 'Bob Petrovic', 'file', 'report.pdf');



  INSERT INTO notifications_by_user (user_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('22222222-2222-2222-2222-222222222222', minTimeuuid('2026-11-21 11:00:00+0000'), 'shared', '6b01b85b2cf4749e214d50ee', 'Alice Novak', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'report.pdf');
  INSERT INTO notifications_by_user_type (user_id, n_type, notified_at, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('22222222-2222-2222-2222-222222222222', 'shared', minTimeuuid('2026-11-21 11:00:00+0000'), '6b01b85b2cf4749e214d50ee', 'Alice Novak', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'report.pdf');
  INSERT INTO notifications_by_user_actor (user_id, by_user_id, notified_at, n_type, by_user_name, f_type, f_id, f_name) VALUES ('22222222-2222-2222-2222-222222222222', '6b01b85b2cf4749e214d50ee', minTimeuuid('2026-11-21 11:00:00+0000'), 'shared', 'Alice Novak', 'file', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', 'report.pdf');
  INSERT INTO notifications_by_user_file (user_id, f_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_name) VALUES ('22222222-2222-2222-2222-222222222222', 'bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbb1', minTimeuuid('2026-11-21 11:00:00+0000'), 'shared', '6b01b85b2cf4749e214d50ee', 'Alice Novak', 'file', 'report.pdf');



  INSERT INTO notifications_by_user (user_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('22222222-2222-2222-2222-222222222222', minTimeuuid('2026-11-22 12:00:00+0000'), 'updated', '6b01b85b2cf4749e214d50ee', 'Alice Novak', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects');
  INSERT INTO notifications_by_user_type (user_id, n_type, notified_at, by_user_id, by_user_name, f_type, f_id, f_name) VALUES ('22222222-2222-2222-2222-222222222222', 'updated', minTimeuuid('2026-11-22 12:00:00+0000'), '6b01b85b2cf4749e214d50ee', 'Alice Novak', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects');
  INSERT INTO notifications_by_user_actor (user_id, by_user_id, notified_at, n_type, by_user_name, f_type, f_id, f_name) VALUES ('22222222-2222-2222-2222-222222222222', '6b01b85b2cf4749e214d50ee', minTimeuuid('2026-11-22 12:00:00+0000'), 'updated', 'Alice Novak', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects');
  INSERT INTO notifications_by_user_file (user_id, f_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_name) VALUES ('22222222-2222-2222-2222-222222222222', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', minTimeuuid('2026-11-22 12:00:00+0000'), 'updated', '6b01b85b2cf4749e214d50ee', 'Alice Novak', 'folder', 'Projects');



  INSERT INTO notifications_by_user (user_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name) VALUES ('33333333-3333-3333-3333-333333333333', minTimeuuid('2026-11-23 09:00:00+0000'), 'added', '6b01b85b2cf4749e214d50ee', 'Alice Novak', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects', 'cccccccc-cccc-cccc-cccc-ccccccccccc2', 'design.png');
  INSERT INTO notifications_by_user_type (user_id, n_type, notified_at, by_user_id, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name) VALUES ('33333333-3333-3333-3333-333333333333', 'added', minTimeuuid('2026-11-23 09:00:00+0000'), '6b01b85b2cf4749e214d50ee', 'Alice Novak', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects', 'cccccccc-cccc-cccc-cccc-ccccccccccc2', 'design.png');
  INSERT INTO notifications_by_user_actor (user_id, by_user_id, notified_at, n_type, by_user_name, f_type, f_id, f_name, added_f_id, added_f_name) VALUES ('33333333-3333-3333-3333-333333333333', '6b01b85b2cf4749e214d50ee', minTimeuuid('2026-11-23 09:00:00+0000'), 'added', 'Alice Novak', 'folder', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', 'Projects', 'cccccccc-cccc-cccc-cccc-ccccccccccc2', 'design.png');
  INSERT INTO notifications_by_user_file (user_id, f_id, notified_at, n_type, by_user_id, by_user_name, f_type, f_name, added_f_id, added_f_name) VALUES ('33333333-3333-3333-3333-333333333333', 'aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaa1', minTimeuuid('2026-11-23 09:00:00+0000'), 'added', '6b01b85b2cf4749e214d50ee', 'Alice Novak', 'folder', 'Projects', 'cccccccc-cccc-cccc-cccc-ccccccccccc2', 'design.png');
