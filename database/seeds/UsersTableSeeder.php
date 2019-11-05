<?php

use Illuminate\Database\Seeder;
use App\Models\User;

class UsersTableSeeder extends Seeder
{
    /**
     * Run the database seeds.
     *
     * @return void
     */
    public function run()
    {
        User::insert(
            [
                [
                    'name' => 'Admin',
                    'username' => 'admin',
                    'admin'    => 1,
                    'email'    => 'ankit.infra@razorpay.com',
                    'access_token' => '',
                    'google_id' => '',
                ],
            ]
        );
    }
}
