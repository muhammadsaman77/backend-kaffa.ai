-- Seeding roles
INSERT INTO roles (id, name, description) VALUES ('customer','Customer', 'Regular user role with limited access');
INSERT INTO roles (id, name, description) VALUES ('superadmin','Super Admin', 'Super administrator role with full access to manage the system');
INSERT INTO roles (id, name, description) VALUES ('admin','Administrator', 'Administrator role with access to manage users and stores');

-- Seeding users
INSERT INTO users (id, username, email, password, role_id) VALUES ('01KBHTWATQY5SDQ2JZ7FT91P2G','rinjani_coffee','rinjani.coffee@gmail.com','$2a$12$XEiiqMuUJBTTvpLoxlAjsu4Axw6SuHWK8qQMDKpKER/OG69NjHp16','admin');

-- Seeding stores
INSERT INTO stores (id, user_id, name, description) VALUES ('01KBHV2FX2A0MRX3Q537GNC52C','01KBHTWATQY5SDQ2JZ7FT91P2G','Rinjani Coffee','A cozy coffee shop located in the heart of the city, offering a wide range of coffee blends and pastries. Our mission is to provide a warm and welcoming atmosphere for coffee lovers to enjoy their favorite brews and pastries. We source our coffee beans from local farmers, ensuring freshness and quality in every cup. Whether you prefer a classic espresso or a specialty latte, we have something for everyone. Join us for a cup of coffee and experience the rich flavors and aromas that Rinjani Coffee has to offer. Our friendly staff is always ready to serve you with a smile and make your visit memorable.');