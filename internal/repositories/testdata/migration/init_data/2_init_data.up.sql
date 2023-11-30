INSERT INTO public.users (id,create_at,login,"password",bill) VALUES ('98dcfb07-e16f-4e53-9a28-d2a2e4eed026'::uuid,'2023-01-01 14:00:00.000','TestUser1','password', 0);
INSERT INTO public.users (id,create_at,login,"password",bill) VALUES ('98dcfb07-e16f-4e53-9a28-d2a2e4eed027'::uuid,'2023-01-01 14:00:00.000','TestUser2','password', 100);

INSERT INTO public.orders (id, create_at, "number", accrual, status, user_id) VALUES('334b0360-8222-44fc-bf2e-77ced208f2cd'::uuid, '2023-01-01 14:01:00.000', 4539088167512356, 100.0, 'PROCESSED', '98dcfb07-e16f-4e53-9a28-d2a2e4eed026'::uuid);
INSERT INTO public.orders (id, create_at, "number", accrual, status, user_id) VALUES('334b0360-8222-44fc-bf2e-77ced208f2ce'::uuid, '2023-01-01 14:02:00.000', 3536137811022331, 0, 'NEW', '98dcfb07-e16f-4e53-9a28-d2a2e4eed026'::uuid);
INSERT INTO public.orders (id, create_at, "number", accrual, status, user_id) VALUES('334b0360-8222-44fc-bf2e-77ced208f2cf'::uuid, '2023-01-01 14:03:00.000', 3533841638640315, 0, 'INVALID', '98dcfb07-e16f-4e53-9a28-d2a2e4eed026'::uuid);

INSERT INTO public.withdrawals (id, create_at, order_num, sum, user_id) VALUES('35e1cbd0-c3ba-44eb-8632-0d91c280dee6'::uuid, '2023-11-10 11:00:00.000', 140672056, 12.64, '98dcfb07-e16f-4e53-9a28-d2a2e4eed026'::uuid);
INSERT INTO public.withdrawals (id, create_at, order_num, sum, user_id) VALUES('35e1cbd0-c3ba-44eb-8632-0d91c280dee7'::uuid, '2023-12-10 11:00:00.000', 140672057, 27.385, '98dcfb07-e16f-4e53-9a28-d2a2e4eed026'::uuid);
INSERT INTO public.withdrawals (id, create_at, order_num, sum, user_id) VALUES('35e1cbd0-c3ba-44eb-8632-0d91c280dee8'::uuid, '2023-11-11 11:00:00.000', 140672058, 0.1111111111, '98dcfb07-e16f-4e53-9a28-d2a2e4eed026'::uuid);