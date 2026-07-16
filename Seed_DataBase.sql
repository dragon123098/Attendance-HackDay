IF DB_ID(N'AttendanceHackday') IS NULL
BEGIN
    CREATE DATABASE AttendanceHackday;
END;
GO

USE AttendanceHackday;
GO

-- Clean up existing tables so this file can be re-run safely.
IF OBJECT_ID(N'dbo.ClassroomStudents', N'U') IS NOT NULL DROP TABLE dbo.ClassroomStudents;
IF OBJECT_ID(N'dbo.OwnedShopItems', N'U') IS NOT NULL DROP TABLE dbo.OwnedShopItems;
IF OBJECT_ID(N'dbo.AvatarConfigs', N'U') IS NOT NULL DROP TABLE dbo.AvatarConfigs;
IF OBJECT_ID(N'dbo.Transactions', N'U') IS NOT NULL DROP TABLE dbo.Transactions;
IF OBJECT_ID(N'dbo.AttendanceRecords', N'U') IS NOT NULL DROP TABLE dbo.AttendanceRecords;
IF OBJECT_ID(N'dbo.WeeklyAssignmentTemplates', N'U') IS NOT NULL DROP TABLE dbo.WeeklyAssignmentTemplates;
IF OBJECT_ID(N'dbo.Schedule', N'U') IS NOT NULL DROP TABLE dbo.Schedule;
IF OBJECT_ID(N'dbo.ShopItems', N'U') IS NOT NULL DROP TABLE dbo.ShopItems;
IF OBJECT_ID(N'dbo.Users', N'U') IS NOT NULL DROP TABLE dbo.Users;
IF OBJECT_ID(N'dbo.Classrooms', N'U') IS NOT NULL DROP TABLE dbo.Classrooms;
GO

CREATE TABLE dbo.Users (
    UserID nvarchar(100) NOT NULL PRIMARY KEY,
    Name nvarchar(200) NOT NULL,
    Role nvarchar(50) NOT NULL,
    Email nvarchar(200) NOT NULL,
    PasswordHash nvarchar(300) NOT NULL,
    ClassroomID nvarchar(100) NULL
);

CREATE TABLE dbo.Classrooms (
    ID nvarchar(100) NOT NULL PRIMARY KEY,
    Name nvarchar(200) NOT NULL,
    TeacherID nvarchar(100) NULL
);

CREATE TABLE dbo.ClassroomStudents (
    ClassroomID nvarchar(100) NOT NULL,
    StudentID nvarchar(100) NOT NULL,
    PRIMARY KEY (ClassroomID, StudentID)
);

CREATE TABLE dbo.ShopItems (
    ID nvarchar(100) NOT NULL PRIMARY KEY,
    Name nvarchar(200) NOT NULL,
    Price int NOT NULL,
    Description nvarchar(max) NOT NULL
);

CREATE TABLE dbo.OwnedShopItems (
    UserID nvarchar(100) NOT NULL,
    ShopItemID nvarchar(100) NOT NULL,
    PRIMARY KEY (UserID, ShopItemID)
);

CREATE TABLE dbo.AvatarConfigs (
    UserID nvarchar(100) NOT NULL PRIMARY KEY,
    Base nvarchar(200) NULL,
    HairStyle nvarchar(200) NULL,
    Clothing nvarchar(200) NULL,
    Accessory nvarchar(200) NULL
);

CREATE TABLE dbo.Transactions (
    TransactionID int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    UserID nvarchar(100) NOT NULL,
    Amount int NOT NULL,
    Timestamp datetimeoffset(7) NOT NULL,
    Description nvarchar(max) NOT NULL
);

CREATE TABLE dbo.AttendanceRecords (
    UserID nvarchar(100) NOT NULL,
    ClassroomID nvarchar(100) NOT NULL,
    PresentDates nvarchar(max) NULL,
    AbsentDates nvarchar(max) NULL,
    PRIMARY KEY (UserID, ClassroomID)
);

CREATE TABLE dbo.Schedule (
    ScheduleID int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    ClassroomID nvarchar(100) NOT NULL,
    DayOfWeek nvarchar(50) NOT NULL,
    StartTime nvarchar(20) NOT NULL,
    EndTime nvarchar(20) NOT NULL,
    DoubleDay bit NOT NULL
);

CREATE TABLE dbo.WeeklyAssignmentTemplates (
    WeeklyAssignmentTemplateID int IDENTITY(1,1) NOT NULL PRIMARY KEY,
    ClassroomID nvarchar(100) NOT NULL,
    DueWeekday tinyint NOT NULL,
    Subject nvarchar(100) NOT NULL,
    Title nvarchar(200) NOT NULL,
    DueTime time(0) NOT NULL,
    DisplayOrder int NOT NULL CONSTRAINT DF_WeeklyAssignmentTemplates_DisplayOrder DEFAULT (0),
    CONSTRAINT CK_WeeklyAssignmentTemplates_DueWeekday CHECK (DueWeekday BETWEEN 0 AND 6),
    CONSTRAINT UQ_WeeklyAssignmentTemplates_ClassroomAssignment UNIQUE (ClassroomID, DueWeekday, Subject, Title)
);
GO

INSERT INTO dbo.Users (UserID, Name, Role, Email, PasswordHash, ClassroomID) VALUES
('BenJam', 'Benjamin', 'student', 'BenJam@example.com', '$2a$10$CjkYbcz0E4nO5MvxiK4...oJGcCp5ndO0iQ5glmZo2ITGDtMKiZ8K', 'classroom1'),
('DHoney', 'Dylan Steenhoek', 'teacher', 'DHoney@example.com', '$2a$10$lSkGtrwXNpsriRUxhgGor.nLddrL0hiCjGtW/TEQ.NGC99j6xN5DW', ''),
('JRGRUNIG', 'Jed Grunig', 'teacher', 'jed@example.com', '$2a$10$89jnZG5yC7yjXyJEae/8xOiQ3PEP4HAtVc0gxKX6lqw6DCO7KNQqO', ''),
('PeteGrunigi', 'Peter Grunig', 'admin', 'petergrunig@gmail.com', '$2a$10$aGGj2pW8PFXJ/IjO0PESbe0rESTEnAopD9WYt.Qw25N/vg9vFVunq', ''),
('admin1', 'Test Admin', 'admin', 'admin@example.com', '$2a$10$j7TaLACJVyfoNgoiimtyy.b/PR.75ri.RZeEVoU.EBdyFHwjDwWxS', ''),
('jb', 'Joey', 'student', 'jb@example.com', '$2a$10$3fBFhi5ZfE4v7u/yehwR8uNqbqSY35PrYsnXxrYSV35ePJr5MunnC', 'classroom1'),
('sconner1', 'Seth Conner', 'student', 'sconner1@example.com', '$2a$10$ExSwVRKGZf99Q571Qe30Uuh/ecKQyupGFSqWfus2GiE0c8dejGXty', ''),
('student1', 'Test Student', 'student', 'student@example.com', '$2a$10$j7TaLACJVyfoNgoiimtyy.b/PR.75ri.RZeEVoU.EBdyFHwjDwWxS', 'classroom1'),
('teacher1', 'Test Teacher', 'teacher', 'teacher@example.com', '$2a$10$j7TaLACJVyfoNgoiimtyy.b/PR.75ri.RZeEVoU.EBdyFHwjDwWxS', ''),
('test', 'test2', 'teacher', 'hey@gmail.com', '$2a$10$GWPCqQ9BvUTWxFGNMYZhWOBl0v1m/9v0gfpajldzD4XsYhw6LJCYa', '');

INSERT INTO dbo.Classrooms (ID, Name, TeacherID) VALUES
('classroom1', '1st Grade', 'sconner1'),
('classroom2', '5th Grade', 'JRGRUNIG'),
('classroom3', '2nd Grade', 'DHoney');

INSERT INTO dbo.WeeklyAssignmentTemplates (ClassroomID, DueWeekday, Subject, Title, DueTime, DisplayOrder) VALUES
('classroom1', 0, 'Reading', 'Weekend Reading Log', '19:00', 10),
('classroom1', 1, 'Math', 'Addition Practice', '15:30', 10),
('classroom1', 2, 'Reading', 'Story Response', '16:00', 10),
('classroom1', 3, 'Science', 'Weather Journal', '16:30', 10),
('classroom1', 5, 'Spelling', 'Weekly Word Check', '15:00', 10),
('classroom2', 0, 'Reading', 'Novel Reading Log', '19:00', 10),
('classroom2', 1, 'Math', 'Fraction Review', '16:00', 10),
('classroom2', 3, 'Science', 'Ecosystem Notes', '16:30', 10),
('classroom2', 4, 'Writing', 'Persuasive Paragraph', '17:00', 10),
('classroom2', 6, 'Social Studies', 'Map Skills Review', '14:00', 10),
('classroom3', 0, 'Reading', 'Picture Book Response', '18:30', 10),
('classroom3', 2, 'Math', 'Place Value Practice', '15:30', 10),
('classroom3', 3, 'Science', 'Animal Habitat Sketch', '16:00', 10),
('classroom3', 5, 'Spelling', 'Weekly Word Sort', '15:00', 10),
('classroom3', 6, 'Art', 'Color Wheel Practice', '13:00', 10);

INSERT INTO dbo.ClassroomStudents (ClassroomID, StudentID) VALUES
('classroom1', 'student1'),
('classroom1', 'jb'),
('classroom1', 'BenJam'),
('classroom2', 'student3'),
('classroom2', 'student4');

INSERT INTO dbo.ShopItems (ID, Name, Price, Description) VALUES
('cape_gold', 'Golden Cape', 12, 'A shiny cape for extra style.'),
('glasses_rocket', 'Rocket Glasses', 10, 'A bold accessory for your avatar.'),
('hat_star', 'Star Hat', 5, 'A bright hat for a standout student.'),
('trail_rainbow', 'Rainbow Trail', 8, 'A colorful trail effect for your avatar.');

INSERT INTO dbo.OwnedShopItems (UserID, ShopItemID) VALUES
('student1', 'hat_star');

INSERT INTO dbo.Transactions (UserID, Amount, Timestamp, Description) VALUES
('student1', 1, '2026-06-09T15:33:47-06:00', 'Attendance reward'),
('student1', -5, '2026-06-09T16:29:38-06:00', 'Purchased Star Hat'),
('student1', 1, '2026-06-10T16:44:15-06:00', 'Attendance reward for 2026-06-10'),
('student1', 1, '2026-06-15T10:31:04-06:00', 'Attendance reward for 2026-06-15'),
('student1', 1, '2026-06-22T16:25:50-06:00', 'Attendance reward for 2026-06-22'),
('student1', 1, '2026-06-24T16:51:42-06:00', 'Attendance reward for 2026-06-24');

INSERT INTO dbo.AttendanceRecords (UserID, ClassroomID, PresentDates, AbsentDates) VALUES
('student1', 'classroom1', '["2026-06-09","2026-06-10","2026-06-15","2026-06-22","2026-06-24"]', NULL);
GO
