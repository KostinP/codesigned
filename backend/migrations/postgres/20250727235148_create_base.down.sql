DROP TABLE IF EXISTS tag_assignments; 
DROP TABLE IF EXISTS tags; 
DROP TABLE IF EXISTS courses;
DROP TABLE IF EXISTS modules;
DROP TABLE IF EXISTS lessons;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS category_assignments;

DROP INDEX IF EXISTS idx_category_assignments_target;
DROP INDEX IF EXISTS idx_category_assignments_category;
DROP INDEX IF EXISTS idx_courses_slug;
DROP INDEX IF EXISTS idx_courses_deleted;
DROP INDEX IF EXISTS idx_modules_course;
DROP INDEX IF EXISTS idx_lessons_module;