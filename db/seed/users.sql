-- Add user with the following credentials
-- username: admin
-- email: admin@admin.com
-- password: adminpass
--
-- Notice that this password only works with bcrypt with cost 12
INSERT INTO users (
	name, 
	email, 
	hashed_password, 
	created, 
	role
) VALUES(
	'admin', 
	'admin@admin.com', 
	'$2a$12$CzMpV/4akHwp7MyAzLLxQ.wZFXREXuWSvHPvApOjr2BVCKiYf7tMi', 
	'2021-10-02 20:07:41', 
	'admin'	
);
