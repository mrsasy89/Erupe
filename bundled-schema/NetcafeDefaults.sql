BEGIN;

INSERT INTO public.cafebonus (time_req, item_type, item_id, quantity)
VALUES
    (1800, 17, 0, 5),
    (3600, 17, 0, 10),
    (7200, 17, 0, 20),
    (10800, 17, 0, 40),
    (18000, 17, 0, 80),
    (28800, 17, 0, 100),
    (43200, 17, 0, 150);

END;
