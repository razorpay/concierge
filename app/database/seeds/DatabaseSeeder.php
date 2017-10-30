<?php

class DatabaseSeeder extends Seeder
{
    /**
     * Run the database seeds for setting up inital set of users.
     *
     *
     * @return void
     */
    public function run()
    {
        Eloquent::unguard();

        $users = Config::get('concierge.users');

        foreach ($users as &$user) {
            $user['created_at'] = date('Y-m-d H:i:s');
            $user['updated_at'] = date('Y-m-d H:i:s');
        }
        DB::table('users')->insert($users);
    }
}
