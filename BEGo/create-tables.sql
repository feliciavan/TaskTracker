DROP TABLE IF EXISTS task;
CREATE TABLE task (
  id         INT AUTO_INCREMENT NOT NULL,
  text      VARCHAR(255) NOT NULL,
  day     VARCHAR(255) NOT NULL,
  reminder      boolean NOT NULL,
  PRIMARY KEY (`id`)
);

INSERT INTO task 
  (text, day, reminder) 
VALUES 
  ('buy eggs', 'Sep 1', true),
  ('buy milk', 'Sep 2', true),
  ('buy water', 'Sep 3', true),
 ('buy suger', 'Sep 4', true);