insert into user_profiles(user_id,
                          personal_info,
                          created_at,
                          updated_at)
values (:user_id,
        :personal_info,
        :created_at,
        now())
returning id